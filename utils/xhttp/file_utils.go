package xhttp

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadImageAndSave 下载图片
// @Description: 返回下载的图片扩展名
// @param url
// @param path
// @return string
// @return error
func DownloadImageAndSave(url string, path string, httpClient *http.Client) (string, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image, status code: %d", resp.StatusCode)
	}

	ext := filepath.Ext(url)
	file, err := os.Create(path + ext)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return ext, nil
}

// DownloadFileAndSave
// @Description: 返回下载的文件扩展名
// @param url
// @param path
// @return string
// @return error
func DownloadFileAndSave(url string, path string, httpClient *http.Client) (string, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file, status code: %d", resp.StatusCode)
	}

	ext := filepath.Ext(url)
	file, err := os.Create(path + ext)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return ext, nil
}

// DownloadImageGetBytes 下载图片
// @Description: 返回下载的图片文件字节切片和对应的扩展名
// @param url
// @return []byte
// @return string
// @return error
func DownloadImageGetBytes(url string, httpClient *http.Client) ([]byte, string, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to download image, status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	ext := filepath.Ext(url)

	return data, ext, nil
}

// DownloadFileGetBytes 下载文件
// @Description: 返回下载的文件字节切片和对应的扩展名
// @param url
// @return []byte
// @return string
// @return error
func DownloadFileGetBytes(url string, httpClient *http.Client) ([]byte, string, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to download file, status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	ext := filepath.Ext(url)

	return data, ext, nil
}
