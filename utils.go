package main

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/himidori/vkapi"
)

func folderExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func downloadFile(url string, name string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func getFileName(url string) string {
	idx := strings.LastIndex(url, "/")
	return url[idx+1:]
}

func getBestLink(photo *vkapi.PhotoAttachment) string {
	if photo.Photo2560 != "" {
		return photo.Photo2560
	} else if photo.Photo1280 != "" {
		return photo.Photo1280
	}

	return photo.Photo604
}

func mkdir(name string) bool {
	err := os.MkdirAll(name, 0755)
	if err != nil {
		return false
	}

	return true
}
