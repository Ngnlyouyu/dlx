package download

import (
	"dlx/internal/download/datafile"
	"dlx/internal/download/request"
	"dlx/log"
	"dlx/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
)

type Downloader struct {
	bar *pb.ProgressBar
}

func NewDownloader() *Downloader {
	return &Downloader{}
}

func (d *Downloader) Download(data *datafile.Data) error {
	if len(data.Streams) == 0 {
		err := fmt.Errorf("no streams in title: %s", data.Title)
		log.Sugar.Error(err)
		return err
	}

	sortedStreams := genSortedStreams(data.Streams)

	title := data.Title

	streamName := sortedStreams[0].ID
	stream, ok := data.Streams[streamName]
	if !ok {
		err := fmt.Errorf("no stream named: %s", streamName)
		log.Sugar.Error(err)
		return err
	}

	// TODO: download caption

	mergedFilePath := fmt.Sprintf("%s.%s", title, stream.Ext)
	_, mergedFileExists, err := utils.FileSize(mergedFilePath)
	if err != nil {
		return err
	}
	if mergedFileExists {
		fmt.Printf("%s: file already exists, skipping\n", mergedFilePath)
	}

	d.bar = utils.ProgressBar(stream.Size)
	d.bar.Start()

	if len(stream.Parts) == 1 {
		if err := d.save(stream.Parts[0], data.URL, title); err != nil {
			return err
		}
		d.bar.Finish()
		return nil
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	errs := make([]error, 0)
	parts := make([]string, len(stream.Parts))
	for idx, part := range stream.Parts {
		if len(errs) > 0 {
			break
		}

		partFileName := fmt.Sprintf("%s[%d]", title, idx)
		partFilePath := fmt.Sprintf("%s.%s", partFileName, part.Ext)
		parts[idx] = partFilePath

		wg.Add(1)
		go func(part *datafile.Part, fileName string) {
			defer wg.Done()
			if err := d.save(part, data.URL, fileName); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}(part, partFileName)
	}
	wg.Wait()
	if len(errs) > 0 {
		return errs[0]
	}
	d.bar.Finish()

	if data.Type != datafile.DataTypeVideo {
		return nil
	}

	fmt.Printf("Merging video parts into %s\n", mergedFilePath)
	if stream.Ext != datafile.ExtensionMP4 || stream.NeedMux {
		return utils.MergeFilesWithSameExtension(parts, mergedFilePath)
	}
	return utils.MergeToMP4(parts, mergedFilePath, title)
}

func (d *Downloader) save(part *datafile.Part, refer, fileName string) error {
	filePath := fmt.Sprintf("%s.%s", fileName, part.Ext)
	fileSize, exists, err := utils.FileSize(filePath)
	if err != nil {
		return err
	}
	// Skip segment file
	// TODO: Live video URLs will not return the size
	if exists && fileSize == part.Size {
		d.bar.Add64(fileSize)
		return nil
	}

	tempFilePath := filePath + ".download"
	tempFileSize, _, err := utils.FileSize(tempFilePath)
	if err != nil {
		return err
	}
	headers := map[string]string{
		"Referer": refer,
	}
	var (
		file      *os.File
		fileError error
	)
	if tempFileSize > 0 {
		// range start from 0, 0-1023 means the first 1024 bytes of the file
		headers["Range"] = fmt.Sprintf("bytes=%d-", tempFileSize)
		file, fileError = os.OpenFile(tempFilePath, os.O_APPEND|os.O_WRONLY, 0644)
		d.bar.Add64(tempFileSize)
	} else {
		file, fileError = os.Create(tempFilePath)
	}
	if fileError != nil {
		return fileError
	}

	// close and rename temp file at the end of this function
	defer func() {
		// must close the file before rename or it will cause
		// `The process cannot access the file because it is being used by another process.` error.
		file.Close() // nolint
		if err == nil {
			os.Rename(tempFilePath, filePath) // nolint
		}
	}()

	temp := tempFileSize
	for i := 0; ; i++ {
		written, err := d.writeFile(part.URL, file, headers)
		if err == nil {
			break
		} else if i+1 >= 5 {
			return err
		}
		temp += written
		headers["Range"] = fmt.Sprintf("bytes=%d-", temp)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (d *Downloader) writeFile(url string, file *os.File, headers map[string]string) (int64, error) {
	res, err := request.Request(http.MethodGet, url, nil, headers)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	barWriter := d.bar.NewProxyWriter(file)
	// Note that io.Copy reads 32kb(maximum) from input and writes them to output, then repeats.
	// So don't worry about memory.
	written, copyErr := io.Copy(barWriter, res.Body)
	if copyErr != nil && copyErr != io.EOF {
		return written, fmt.Errorf("file copy error: %s", copyErr)
	}
	return written, nil
}
