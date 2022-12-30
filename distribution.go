// Copyright 2022 E99p1ant. All rights reserved.

package distribution

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

const delimiter = 0xe9

func NewDistribution(binaryPath string, embedData interface{}) error {
	binaryFile, err := os.OpenFile(binaryPath, os.O_RDWR|os.O_APPEND, 0)
	if err != nil {
		return errors.Wrap(err, "open")
	}
	defer func() { _ = binaryFile.Close() }()

	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(embedData); err != nil {
		return errors.Wrap(err, "gob encode")
	}
	size := buffer.Len()

	_, err = binaryFile.Write(buffer.Bytes())
	if err != nil {
		return errors.Wrap(err, "write embed data")
	}
	strBytes := []byte(fmt.Sprintf("%d", size))
	// Add delimiter.
	strBytes = append([]byte{delimiter}, strBytes...)

	// Write the delimiter and the embed data size.
	_, err = binaryFile.Write(strBytes)
	if err != nil {
		return errors.Wrap(err, "write delimiter and size")
	}
	return nil
}

func ParseFromDistribution(binaryPath string, v interface{}) error {
	binaryFile, err := os.Open(binaryPath)
	if err != nil {
		return errors.Wrap(err, "open")
	}
	defer func() { _ = binaryFile.Close() }()

	fileInfo, err := binaryFile.Stat()
	if err != nil {
		return errors.Wrap(err, "stat")
	}
	fileSize := fileInfo.Size()

	// Find the delimiter and read the embed data size.
	var sizeOffset int
	for i := 0; i < 256; i++ {
		_, err := binaryFile.Seek(fileSize-int64(i), 0)
		if err != nil {
			return errors.Wrap(err, "seek")
		}

		var b [1]byte
		_, err = binaryFile.Read(b[:])
		if err != nil && err != io.EOF {
			return errors.Wrap(err, "read")
		}

		if b[0] == delimiter {
			sizeOffset = i
			break
		}
	}

	strBytes := make([]byte, sizeOffset-1)
	_, err = binaryFile.Read(strBytes)
	if err != nil {
		return errors.Wrap(err, "read size")
	}

	size, err := strconv.Atoi(string(strBytes))
	if err != nil {
		return errors.Wrap(err, "parse size")
	}

	// Read the embed data.
	_, err = binaryFile.Seek(fileSize-int64(sizeOffset)-int64(size), 0)
	if err != nil {
		return errors.Wrap(err, "seek")
	}

	embedData := make([]byte, size)
	_, err = binaryFile.Read(embedData)
	if err != nil {
		return errors.Wrap(err, "read embed data")
	}

	if err := gob.NewDecoder(bytes.NewReader(embedData)).Decode(v); err != nil {
		return errors.Wrap(err, "gob decode")
	}
	return nil
}
