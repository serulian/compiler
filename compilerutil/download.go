// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilerutil

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// ProgressCallback is a callback for the progress of a file download operation.
type ProgressCallback func(bytesDownloaded int64, totalBytes int64, startCall bool)

// DownloadURLToFileWithProgress downloads the contents of the given URL to the given file path on disk,
// invoking the callback function to report progress.
func DownloadURLToFileWithProgress(url *url.URL, filePath string, callback ProgressCallback) error {
	// Retrieve the size of the file by making a HEAD request with a Content-Length.
	headResp, err := http.Head(url.String())
	if err != nil {
		return err
	}
	defer headResp.Body.Close()

	size, err := strconv.ParseInt(headResp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return err
	}

	callback(0, int64(size), true)

	// Open the file on disk in which we will download.
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Kick off a goroutine to monitor the size of the file and call the progress callback.
	doneCh := make(chan bool, 1)
	monitorFileSize(filePath, size, callback, doneCh)

	// Start the download.
	getResp, err := http.Get(url.String())
	if err != nil {
		doneCh <- true
		return err
	}

	defer getResp.Body.Close()

	_, err = io.Copy(file, getResp.Body)
	doneCh <- true
	return err
}

func monitorFileSize(filePath string, totalBytes int64, callback ProgressCallback, doneCh chan bool) {
	for {
		select {
		case <-doneCh:
			return

		default:
			file, err := os.Open(filePath)
			if err != nil {
				return
			}

			stat, err := file.Stat()
			if err != nil {
				return
			}

			callback(stat.Size(), totalBytes, false)
		}

		time.Sleep(time.Second)
	}
}
