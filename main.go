package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	COPY_BUFFER_SIZE = 4096
)

type DownloadOptions struct {
	Filepath string
}

type DownloadStatus struct {
	Error           error
	IsComplete      bool
	BytesDownloaded int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [options] URL\n", os.Args[0])
		os.Exit(1)
	}

	url := os.Args[1]
	filename := getFilenameFromUrl(url)
	statusChannel := make(chan DownloadStatus)
	ticker := time.NewTicker(time.Millisecond * 500)
	lastStatus := DownloadStatus{}

	go downloadUrl(url, &DownloadOptions{Filepath: filename}, statusChannel)
	done := false
	for done != true {
		select {
		case <-ticker.C:
			n, unit := getFormattedSize(lastStatus.BytesDownloaded)
			if lastStatus.IsComplete == true {
				fmt.Printf("\x1B[1K\rcompleted: %.2f %s\n", n, unit)
				done = true
				break
			}
			if lastStatus.Error != nil {
				fmt.Fprintf(
					os.Stderr,
					"download failed: %s\n",
					lastStatus.Error)
				done = true
				break
			}
			fmt.Printf("\x1B[1K\r%.2f %s", n, unit)
		case status := <-statusChannel:
			lastStatus = status
		}
	}
}

func getFormattedSize(size int64) (float64, string) {
	formatUnits := []string{"B", "KiB", "MiB", "GiB"}
	var formattedSize float64 = float64(size)
	var formattedUnit string
	for i := 0; i < 4; i++ {
		if formattedSize < 1000 {
			formattedUnit = formatUnits[i]
			break
		}
		formattedSize /= float64(1000)
	}
	return formattedSize, formattedUnit
}

func downloadUrl(url string, options *DownloadOptions, statusChannel chan DownloadStatus) {
	resp, err := http.Get(url)
	if err != nil {
		statusChannel <- DownloadStatus{
			Error: err,
		}
		return
	}
	if resp.StatusCode != 200 {
		statusChannel <- DownloadStatus{
			Error: fmt.Errorf("server returned: %s", resp.Status),
		}
		return
	}
	file, err := os.OpenFile(
		options.Filepath,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0o600)
	if err != nil {
		statusChannel <- DownloadStatus{
			Error: fmt.Errorf("open file %v failed: %v", options.Filepath, err),
		}
		return
	}
	defer file.Close()
	copyVerbose(file, resp.Body, statusChannel)
}

func getFilenameFromUrl(url string) string {
	filenameIndex := strings.LastIndex(url, "/")
	if filenameIndex == -1 {
		return "download"
	}
	return url[filenameIndex+1:]
}

func copyVerbose(dest io.Writer, src io.Reader, statusChannel chan DownloadStatus) {
	buf := make([]byte, COPY_BUFFER_SIZE)
	status := DownloadStatus{}

	for {
		nread, err := src.Read(buf)
		if nread > 0 {
			nwritten, err := dest.Write(buf[0:nread])
			if err != nil {
				status.Error = err
				break
			}
			status.BytesDownloaded += int64(nwritten)
		}
		if err != nil {
			if err == io.EOF {
				status.IsComplete = true
			} else {
				status.Error = err
			}
		}
		statusChannel <- status
	}
}
