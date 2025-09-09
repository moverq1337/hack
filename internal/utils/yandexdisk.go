package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func UploadToYandexDisk(filePath, fileName string) (string, error) {
	oauthToken := os.Getenv("YANDEX_DISK_TOKEN")
	if oauthToken == "" {
		return "", fmt.Errorf("YANDEX_DISK_TOKEN environment variable not set")
	}

	client := &http.Client{Timeout: 30 * time.Second}

	getUploadURL := fmt.Sprintf("https:

	req, err := http.NewRequest("GET", getUploadURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "OAuth "+oauthToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get upload URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get upload URL, status: %s, response: %s", resp.Status, string(body))
	}

	var uploadResponse struct {
		Href string `json:"href"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&uploadResponse); err != nil {
		return "", fmt.Errorf("failed to parse upload URL response: %v", err)
	}

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	uploadReq, err := http.NewRequest("PUT", uploadResponse.Href, bytes.NewReader(fileData))
	if err != nil {
		return "", fmt.Errorf("failed to create upload request: %v", err)
	}
	uploadReq.Header.Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")

	uploadResp, err := client.Do(uploadReq)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}
	defer uploadResp.Body.Close()

	if uploadResp.StatusCode != http.StatusCreated && uploadResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(uploadResp.Body)
		return "", fmt.Errorf("upload failed, status: %s, response: %s", uploadResp.Status, string(body))
	}

	publishURL := fmt.Sprintf("https:
	publishReq, err := http.NewRequest("PUT", publishURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create publish request: %v", err)
	}
	publishReq.Header.Set("Authorization", "OAuth "+oauthToken)

	publishResp, err := client.Do(publishReq)
	if err != nil {
		return "", fmt.Errorf("failed to publish file: %v", err)
	}
	defer publishResp.Body.Close()

	if publishResp.StatusCode != http.StatusOK && publishResp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(publishResp.Body)
		return "", fmt.Errorf("publish failed, status: %s, response: %s", publishResp.Status, string(body))
	}

	publicURL := fmt.Sprintf("https:
	publicReq, err := http.NewRequest("GET", publicURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create public URL request: %v", err)
	}
	publicReq.Header.Set("Authorization", "OAuth "+oauthToken)

	publicResp, err := client.Do(publicReq)
	if err != nil {
		return "", fmt.Errorf("failed to get public URL: %v", err)
	}
	defer publicResp.Body.Close()

	if publicResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(publicResp.Body)
		return "", fmt.Errorf("failed to get public URL, status: %s, response: %s", publicResp.Status, string(body))
	}

	var publicResponse struct {
		PublicURL string `json:"public_url"`
	}
	if err := json.NewDecoder(publicResp.Body).Decode(&publicResponse); err != nil {
		return "", fmt.Errorf("failed to parse public URL response: %v", err)
	}

	return publicResponse.PublicURL, nil
}
