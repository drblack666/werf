package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/gomega"
)

func CopyIn(sourcePath, destinationPath string) {
	err := os.MkdirAll(destinationPath, 0777)
	Ω(err).ShouldNot(HaveOccurred())

	files, err := ioutil.ReadDir(sourcePath)
	Ω(err).ShouldNot(HaveOccurred())

	copyFile := func(srcPath, dstPath string) {
		src, err := os.Open(srcPath)
		Ω(err).ShouldNot(HaveOccurred())
		defer func() { _ = src.Close() }()

		dst, err := os.Create(dstPath)
		Ω(err).ShouldNot(HaveOccurred())
		defer func() { _ = dst.Close() }()

		_, err = io.Copy(dst, src)
		Ω(err).ShouldNot(HaveOccurred())
	}

	for _, f := range files {
		srcPath := filepath.Join(sourcePath, f.Name())
		dstPath := filepath.Join(destinationPath, f.Name())
		if f.IsDir() {
			CopyIn(srcPath, dstPath)
			continue
		}

		copyFile(srcPath, dstPath)
	}
}

func CreateSimpleWerfYaml(path string) {
	err := os.MkdirAll(path, 0777)
	Ω(err).ShouldNot(HaveOccurred())

	werfYamlPath := filepath.Join(path, "werf.yaml")
	data := []byte(`
project: none
configVersion: 1
---
image: ~
from: alpine
shell:
  setup: date
	`)

	err = ioutil.WriteFile(werfYamlPath, bytes.TrimSpace(data), 0644)
	Ω(err).ShouldNot(HaveOccurred())
}
