package utils

import (
	"fmt"
	"io"
	"os"
)

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}

func IsDir(filePath string) bool {
	stat, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func CopyFile(srcFile, dstFile string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return fmt.Errorf("failed to open source file(%s): %w", srcFile, err)
	}
	defer src.Close()

	dst, err := os.Create(dstFile)
	if err != nil {
		return fmt.Errorf("failed to create destination file(%s): %w", dstFile, err)
	}
	defer dst.Close()

	srcStat, _ := src.Stat()
	if err := os.Chmod(dstFile, srcStat.Mode()); err != nil {
		return fmt.Errorf("failed to change permissinos for destination file(%s): %w", dstFile, err)
	}

	// Copy the content of srcFile to dstFile
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf(
			"failed copy file content from source(%s) to destination(%s): %w",
			srcFile,
			dstFile,
			err,
		)
	}

	return nil
}
