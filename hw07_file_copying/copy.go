package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileNotExist          = errors.New("file does not exist")
	ErrNoPermission          = errors.New("no permission to the file")
	ErrNoReadPermission      = errors.New("no permission to read the file")
	ErrNoPermissionCreatFile = errors.New("no permission to create the file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
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

	// check on a read permission of the file
	if fileFromInfo.Mode().Perm()&0400 == 0 { //nolint:gofumpt
		return ErrNoReadPermission
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

	// check on error that offset exceeds file size
	if fileFromInfo.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	_, err = fileFrom.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	// use func CopyN if limit more than zero
	if limit > 0 {
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

	bar := pb.Full.Start64(fileFromInfo.Size() - offset)
	barReader := bar.NewProxyReader(fileFrom)

	_, err = io.Copy(fileTo, barReader)
	if err != nil {
		return err
	}

	bar.Finish()

	return nil
}
