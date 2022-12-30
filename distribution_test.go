// Copyright 2022 E99p1ant. All rights reserved.

package distribution

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

const output = "Hi World"

func goldenFileName() string {
	return fmt.Sprintf("./golden/helloworld_%s_%s", runtime.GOOS, runtime.GOARCH)
}

type embedData struct {
	Host string
	Port uint
}

func Test_NewDistribution(t *testing.T) {
	fileName := goldenFileName()

	f, err := os.CreateTemp(os.TempDir(), "distribution-test-*")
	require.Nil(t, err)

	goldenFile, err := os.Open(fileName)
	require.Nil(t, err)
	defer func() { _ = goldenFile.Close() }()

	_, err = io.Copy(f, goldenFile)
	require.Nil(t, err)
	err = f.Close()
	require.Nil(t, err)

	err = NewDistribution(f.Name(), embedData{
		Host: "localhost",
		Port: 1999,
	})
	require.Nil(t, err)

	err = os.Chmod(f.Name(), 0o755)
	require.Nil(t, err)
	gotOutput, err := exec.Command(f.Name()).Output()
	require.Nil(t, err)
	require.Equal(t, output, string(bytes.TrimSpace(gotOutput)))
}

func Test_ParseFromDistribution(t *testing.T) {
	fileName := goldenFileName()

	f, err := os.CreateTemp(os.TempDir(), "distribution-test-*")
	require.Nil(t, err)

	goldenFile, err := os.Open(fileName)
	require.Nil(t, err)
	defer func() { _ = goldenFile.Close() }()

	_, err = io.Copy(f, goldenFile)
	require.Nil(t, err)
	err = f.Close()
	require.Nil(t, err)

	want := embedData{
		Host: "localhost",
		Port: 1999,
	}
	err = NewDistribution(f.Name(), want)
	require.Nil(t, err)

	var got embedData
	err = ParseFromDistribution(f.Name(), &got)
	require.Nil(t, err)

	require.Equal(t, want, got)
}
