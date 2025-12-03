package interaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

type response struct {
	Content    string `json:"content"`
	Label      int    `json:"label"`
	Prediction string `json:"prediction"`
}

func TextClassification(text string) string {
	encoded := url.QueryEscape(text)
	requestURL := fmt.Sprintf("http://localhost:5000/text?text=%s", encoded)

	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Non-200 response (%d): %s\n", resp.StatusCode, string(body))
		return ""
	}

	var result response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding response:", err)
		return ""
	}

	return result.Prediction
}

func ImageClassification(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return ""
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		fmt.Printf("Error creating form file: %v\n", err)
		return ""
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Printf("Error copying file content: %v\n", err)
		return ""
	}

	err = writer.Close()
	if err != nil {
		fmt.Printf("Error closing multipart writer: %v\n", err)
		return ""
	}

	requestURL := "http://localhost:5000/image"

	resp, err := http.Post(requestURL, writer.FormDataContentType(), body)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Non-200 response (%d): %s\n", resp.StatusCode, string(body))
		return ""
	}

	var result response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding response:", err)
		return ""
	}

	return result.Prediction
}
