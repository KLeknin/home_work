package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrFileExist             = errors.New("file exist")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {

	sourceFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	fi, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return ErrUnsupportedFile
	}

	if fi.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	newOffset, err := sourceFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	if newOffset < offset {
		return ErrOffsetExceedsFileSize
	}

	destFile, err := os.Open(toPath)
	if err == nil || destFile != nil {
		destFile.Close()
		return ErrFileExist
	}

	destFile, err = os.Create(toPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	max := limit
	if max == 0 || offset+max > fi.Size() {
		max = fi.Size() - offset
	}
	bufSize := max / 100
	if bufSize < 1 {
		bufSize = 1
	}
	buf := make([]byte, bufSize)
	var progress int64

	pBar := pb.StartNew(int(max))
	defer pBar.Finish()

	for progress < max {
		if int64(len(buf)) > max-progress {
			buf = buf[0 : max-progress]
		}
		readBites, err := sourceFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			progress = bufSize - int64(readBites)
		}

		_, err = destFile.Write(buf[0:readBites])
		if err != nil {
			return err
		}
		progress += int64(readBites)
		pBar.SetCurrent(progress)
	}

	return nil
}
