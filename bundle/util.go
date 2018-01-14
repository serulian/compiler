// Copyright 2018 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bundle

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

type byteFile struct {
	filename string
	kind     FileKind
	content  []byte
}

func (bf byteFile) Filename() string {
	return bf.filename
}

func (bf byteFile) Kind() FileKind {
	return bf.kind
}

func (bf byteFile) Reader() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(bf.content))
}

type inmemoryBundle struct {
	files map[string]BundledFile
}

func (imb inmemoryBundle) LookupFile(name string) (BundledFile, bool) {
	file, ok := imb.files[name]
	return file, ok
}

func (imb inmemoryBundle) Files() []BundledFile {
	files := make([]BundledFile, 0, len(imb.files))
	for _, file := range imb.files {
		files = append(files, file)
	}
	return files
}

type bundledFile struct {
	baseBundle Bundle
	file       BundledFile
}

func (bf bundledFile) LookupFile(name string) (BundledFile, bool) {
	if name == bf.file.Filename() {
		return bf.file, true
	}
	return bf.baseBundle.LookupFile(name)
}

func (bf bundledFile) Files() []BundledFile {
	files := bf.baseBundle.Files()
	return append(files, bf.file)
}

// WithFile returns a new file bundle that includes the given file.
func WithFile(bundle Bundle, file BundledFile) Bundle {
	return bundledFile{bundle, file}
}

// InMemoryBundle returns an in-memory bundle of files.
func InMemoryBundle(files map[string]BundledFile) Bundle {
	return inmemoryBundle{files}
}

// EmptyBundle returns an empty bundle of files.
func EmptyBundle() Bundle {
	return inmemoryBundle{map[string]BundledFile{}}
}

// FileFromBytes returns a BundledFile from byte data.
func FileFromBytes(filename string, kind FileKind, content []byte) BundledFile {
	return byteFile{filename, kind, content}
}

// FileFromString returns a BundledFile from byte data.
func FileFromString(filename string, kind FileKind, content string) BundledFile {
	return byteFile{filename, kind, []byte(content)}
}

// DetectContentType returns the content type of the given bundled file. Always returns a content type
// if no error occurred.
func DetectContentType(file BundledFile) (string, error) {
	reader := file.Reader()
	defer reader.Close()

	buffer := make([]byte, 512)
	n, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	return http.DetectContentType(buffer[:n]), nil
}

// WriteToFileSystem writes all the bundled files to the file system, placing them under the given
// directory.
func WriteToFileSystem(bundle Bundle, dir string) error {
	for _, bundledFile := range bundle.Files() {
		filePath := path.Join(dir, bundledFile.Filename())

		diskFile, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer diskFile.Close()

		reader := bundledFile.Reader()
		defer reader.Close()

		if _, err := io.Copy(diskFile, reader); err != nil {
			return err
		}
	}

	return nil
}
