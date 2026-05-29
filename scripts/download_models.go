package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// WriteCounter counts the number of bytes written to it and prints progress.
type WriteCounter struct {
	Total         uint64
	ContentLength uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc *WriteCounter) PrintProgress() {
	// Print the progress on the same line using \r
	if wc.ContentLength > 0 {
		percentage := float64(wc.Total) / float64(wc.ContentLength) * 100
		fmt.Printf("\rDownloading... %.2f%% (%d/%d bytes)", percentage, wc.Total, wc.ContentLength)
	} else {
		fmt.Printf("\rDownloading... %d bytes", wc.Total)
	}
}

// DownloadFile downloads a file from a URL to a local path with progress tracking.
func DownloadFile(destPath string, url string) error {
	// Create the destination directory if it doesn't exist
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("could not create directory %s: %v", destDir, err)
	}

	// Create a temporary file for the download to avoid partial files on failure
	tmpPath := destPath + ".tmp"
	out, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Perform the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the server returned a successful status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Initialize the progress counter
	counter := &WriteCounter{
		ContentLength: uint64(resp.ContentLength),
	}

	// Copy the data from the response body to the file, using TeeReader to count bytes
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		return err
	}

	// Clear the progress line
	fmt.Print("\n")

	// Close the file before renaming
	out.Close()

	// Rename the temporary file to the final destination path
	if err = os.Rename(tmpPath, destPath); err != nil {
		return fmt.Errorf("could not rename temporary file: %v", err)
	}

	return nil
}

func main() {
	fmt.Println("--------------------------------------------------")
	fmt.Println("Speaking_Hearts - AI Model Downloader")
	fmt.Println("--------------------------------------------------")

	// Define the models and their corresponding HuggingFace resolve URLs
	// Note: These are placeholder URLs for the skeleton. 
	// Real URLs for faster-whisper and NLLB would be used here.
	models := []struct {
		Dir  string
		File string
		URL  string
	}{
		{
			Dir:  "models/whisper/base",
			File: "model.bin",
			URL:  "https://huggingface.co/Systran/faster-whisper-base/resolve/main/model.bin",
		},
		{
			Dir:  "models/whisper/base",
			File: "config.json",
			URL:  "https://huggingface.co/Systran/faster-whisper-base/resolve/main/config.json",
		},
		{
			Dir:  "models/nllb/distilled-600M",
			File: "sentencepiece.bpe.model",
			URL:  "https://huggingface.co/facebook/nllb-200-distilled-600M/resolve/main/sentencepiece.bpe.model",
		},
	}

	for _, m := range models {
		destPath := filepath.Join(m.Dir, m.File)
		
		// Check if the file already exists to avoid redundant downloads
		if _, err := os.Stat(destPath); err == nil {
			fmt.Printf("[-] %s already exists, skipping.\n", destPath)
			continue
		}

		fmt.Printf("[+] Downloading %s to %s...\n", m.File, m.Dir)
		if err := DownloadFile(destPath, m.URL); err != nil {
			fmt.Printf("[!] Error downloading %s: %v\n", m.File, err)
			
			// If it's a 404, it might be the placeholder URL failing
			if strings.Contains(err.Error(), "404") {
				fmt.Println("    (Check if the HuggingFace URLs are correct or if you have internet access)")
			}
		} else {
			fmt.Printf("[✓] Successfully downloaded %s\n", m.File)
		}
	}

	fmt.Println("--------------------------------------------------")
	fmt.Println("Model download process completed.")
	fmt.Println("--------------------------------------------------")
}
