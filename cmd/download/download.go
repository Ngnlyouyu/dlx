package download

import (
	"dlx/internal/download"
	"dlx/internal/download/datafile"
	"dlx/internal/download/extractors"
	"dlx/internal/download/request"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Download(cmd *cobra.Command, args []string) {
	request.DefaultRequest = request.NewRequest(5, "", "", "", true)
	var isErr bool
	for _, videoURL := range args {
		if err := startDownload(videoURL); err != nil {
			fmt.Printf("%v\n", err)
			isErr = true
		}
	}
	if isErr {
		os.Exit(1)
	}
}

func startDownload(videoURL string) error {
	data, err := extractors.Extract(videoURL, datafile.Options{})
	if err != nil {
		return err
	}

	defaultDownloader := download.NewDownloader()
	errs := make([]error, 0)
	for _, item := range data {
		if item.Err != nil {
			// if this error occurs, the preparation step is normal, but the data extraction is wrong.
			// the data is an empty struct.
			errs = append(errs, item.Err)
			continue
		}
		if err = defaultDownloader.Download(item); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return errs[0]
	}
	return nil
}
