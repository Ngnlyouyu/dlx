package utils

import (
	"fmt"
	"os"
)

func MergeFilesWithSameExtension(paths []string, mergedFilePath string) error {
	args := []string{"-y"}
	for _, p := range paths {
		args = append(args, "-i", p)
	}
	args = append(args, "-c:v", "copy", "-c:a", "copy", mergedFilePath)

	ret, stdout, stderr := ExecuteCmd("ffmpeg", args...)
	if ret != 0 {
		return fmt.Errorf("ffmpeg error, stdout: %s, stderr: %s", stdout, stderr)
	}

	for _, p := range paths {
		_ = os.Remove(p)
	}
	return nil
}

func MergeToMP4(paths []string, mergedFilePath, fileName string) error {
	mergeTXTFile := fileName + ".txt"

	mergeFile, err := os.OpenFile(mergeTXTFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	for _, p := range paths {
		mergeFile.WriteString(fmt.Sprintf("file '%s'\n", p))
	}
	mergeFile.Close()

	args := []string{
		"-y",
		"-f",
		"concat",
		"-safe",
		"0",
		"-i",
		mergeTXTFile,
		"-c",
		"copy",
		"-bsf:a",
		"aac_adtstoasc",
		mergedFilePath,
	}
	ret, stdout, stderr := ExecuteCmd("ffmpeg", args...)
	if ret != 0 {
		return fmt.Errorf("ffmpeg error, stdout: %s, stderr: %s", stdout, stderr)
	}

	for _, p := range paths {
		_ = os.Remove(p)
	}
	_ = os.Remove(mergeTXTFile)
	return nil
}
