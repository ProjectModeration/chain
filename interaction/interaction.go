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
type responseImage struct {
	Label string  `json:"prediction"`
	Score float64 `json:"score"`
}

func TextClassification(text string) int {
	encoded := url.QueryEscape(text)
	requestURL := fmt.Sprintf("http://localhost:5000/text?text=%s", encoded)

	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Non-200 response (%d): %s\n", resp.StatusCode, string(body))
		return 0
	}

	var result response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding response:", err)
		return 0
	}

	return result.Label
}

func ImageClassification(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		return ""
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		fmt.Printf("error creating form file: %v\n", err)
		return ""
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Printf("error copying file content: %v\n", err)
		return ""
	}

	if err := writer.Close(); err != nil {
		fmt.Printf("error closing multipart writer: %v\n", err)
		return ""
	}

	resp, err := http.Post("http://localhost:5000/image", writer.FormDataContentType(), body)
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("non-200 response (%d): %s\n", resp.StatusCode, string(body))
		return ""
	}

	var result responseImage
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("error decoding response:", err)
		return ""
	}

	return result.Label
}
