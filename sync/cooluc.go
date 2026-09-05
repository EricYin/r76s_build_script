package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const baseURL = "https://init.cooluc.com/"
const outputDir = "./download_dir/"

func main() {
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating root dir:", err)
		return
	}

	downloadDirectory(baseURL, outputDir)
	fmt.Println("Done")
}

func downloadDirectory(currentURL string, currentLocalDir string) {
	fmt.Println("Scanning:", currentURL)

	resp, err := http.Get(currentURL)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return
	}
	html := string(bodyBytes)

	re := regexp.MustCompile(`href="([^"]+)"`)
	matches := re.FindAllStringSubmatch(html, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		link := match[1]
		fullURL := currentURL + link
		// if link == "/" || link == "//" {
		// 	continue
		// }
		if strings.HasPrefix(link, "/") {
			fmt.Println("ignore link :" + link + " ==> " + fullURL)
			continue
		}
		if strings.HasPrefix(link, "?") {
			fmt.Println("ignore link :" + link + " ==> " + fullURL)
			continue
		}
		if strings.HasPrefix(link, "..") {
			fmt.Println("ignore link :" + link + " ==> " + fullURL)
			continue
		}
		
		if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
			fmt.Println("ignore link :" + link + " ==> " + fullURL)
			continue
		}

		if strings.HasSuffix(link, "/") {
			subDirName := strings.TrimSuffix(link, "/")
			newLocalDir := filepath.Join(currentLocalDir, subDirName)

			err := os.MkdirAll(newLocalDir, os.ModePerm)
			if err != nil {
				fmt.Println("Error creating sub dir:", err)
				continue
			}

			downloadDirectory(fullURL, newLocalDir)
		} else {
			localFilePath := filepath.Join(currentLocalDir, link)
			fmt.Println("Downloading:", link)

			err := downloadFile(fullURL, localFilePath)
			if err != nil {
				fmt.Println("Error downloading file:", err)
			}
		}
	}
}

func downloadFile(url string, localPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
