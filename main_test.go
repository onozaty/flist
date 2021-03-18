package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestCalcMd5(t *testing.T) {

	tempFile, info, err := craeteTempFile([]byte("ABCDEFG"))

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(tempFile)

	md5, err := calcMd5(tempFile, info)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if md5 != "bb747b3df3130fe1ca4afa93fb7d97c9" {
		t.Fatal("failed test\n", md5)
	}
}

func TestCalcMd5_empty(t *testing.T) {

	tempFile, info, err := craeteTempFile([]byte{})

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(tempFile)

	md5, err := calcMd5(tempFile, info)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if md5 != "d41d8cd98f00b204e9800998ecf8427e" {
		t.Fatal("failed test\n", md5)
	}
}

func TestCalcSize(t *testing.T) {

	tempFile, info, err := craeteTempFile([]byte("ABCDEFG"))

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(tempFile)

	size, err := calcSize(tempFile, info)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if size != "7" {
		t.Fatal("failed test\n", size)
	}
}

func TestPrintDir(t *testing.T) {

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	os.Stdout = w

	tempDir, err := os.MkdirTemp("", "filist")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.RemoveAll(tempDir)

	file1, err := os.Create(filepath.Join(tempDir, "file1"))
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	file1.Write([]byte("a"))

	// サブディレクトリにもファイル配置
	sub1 := filepath.Join(tempDir, "sub1")
	err = os.Mkdir(sub1, 0777)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	file2, err := os.Create(filepath.Join(sub1, "file2"))
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	file2.Write([]byte("abc"))

	err = printDir(tempDir, Option{})
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	os.Stdout = stdout
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	if buf.String() != "file1\nsub1\\file2\n" {
		t.Fatal("failed test\n", buf.String())
	}
}

func TestPrintDir_option(t *testing.T) {

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	os.Stdout = w

	tempDir, err := os.MkdirTemp("", "filist")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.RemoveAll(tempDir)

	file1, err := os.Create(filepath.Join(tempDir, "file1"))
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	file1.Write([]byte("a"))

	err = printDir(tempDir, Option{optionalColumns: []func(string, os.FileInfo) (string, error){calcSize, calcMd5}})
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	os.Stdout = stdout
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	if buf.String() != "file1\t1\t0cc175b9c0f1b6a831c399e269772661\n" {
		t.Fatal("failed test\n", buf.String())
	}
}

func craeteTempFile(contents []byte) (string, os.FileInfo, error) {

	tempFile, err := os.CreateTemp("", "filist")

	if err != nil {
		return "", nil, err
	}
	defer tempFile.Close()

	tempFile.Write(contents)
	info, err := tempFile.Stat()
	if err != nil {
		return "", nil, err
	}

	return tempFile.Name(), info, nil
}
