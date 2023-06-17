package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileNotExist          = errors.New("file does not exist")
	ErrNoPermission          = errors.New("no permission to the file")
	ErrNoPermissionCreatFile = errors.New("no permission to create the file")
	ErrNegativeLimit         = errors.New("limit cannot be negative number")
	ErrNegativeOffset        = errors.New("offset cannot be negative number")
	ErrEqualPaths            = errors.New("paths are equal")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if limit < 0 {
		return ErrNegativeLimit
	}

	if offset < 0 {
		return ErrNegativeOffset
	}

	fileFrom, err := os.Open(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotExist
		}
		if os.IsPermission(err) {
			return ErrNoPermission
		}
		return err
	}
	defer fileFrom.Close()

	fileFromInfo, err := fileFrom.Stat()
	if err != nil {
		return err
	}

	// create directories in the toPath if they don't exist
	err = os.MkdirAll(filepath.Dir(toPath), os.ModePerm)
	if err != nil {
		return err
	}

	// create toPath file
	fileTo, err := os.Create(toPath)
	if err != nil {
		if os.IsPermission(err) {
			return ErrNoPermissionCreatFile
		}
		return err
	}
	defer fileTo.Close()

	if err := comparePaths(fromPath, toPath); err != nil {
		return err
	}

	// check on error that offset exceeds file size
	if fileFromInfo.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	_, err = fileFrom.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	if limit == 0 {
		limit = fileFromInfo.Size() - offset
	}

	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(fileFrom)

	_, err = io.CopyN(fileTo, barReader, limit)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	bar.Finish()

	return nil
}

func comparePaths(path1, path2 string) error {
	normalizedPath1, err := filepath.Abs(filepath.Clean(path1))
	if err != nil {
		return fmt.Errorf("error normalize path1")
	}

	normalizedPath2, err := filepath.Abs(filepath.Clean(path2))
	if err != nil {
		return fmt.Errorf("error normalize path2")
	}

	file1Info, err := os.Stat(normalizedPath1)
	if err != nil {
		return err
	}

	file2Info, err := os.Stat(normalizedPath2)
	if err != nil {
		return err
	}

	if ok := os.SameFile(file1Info, file2Info); ok {
		return ErrEqualPaths
	}

	return nil
}
