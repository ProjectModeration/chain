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
	"strconv"
)

type response struct {
	Content    string `json:"content"`
	Label      int    `json:"label"`
	Prediction string `json:"prediction"`
}
type responseImage struct {
	Label int    `json:"label"`
	Score string `json:"score"`
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

	return strconv.Itoa(result.Label)
}

func ImageClassification(filePath string) int {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		return -1
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		fmt.Printf("error creating form file: %v\n", err)
		return -1
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Printf("error copying file content: %v\n", err)
		return -1
	}

	if err := writer.Close(); err != nil {
		fmt.Printf("error closing multipart writer: %v\n", err)
		return -1
	}

	resp, err := http.Post("http://localhost:5000/image", writer.FormDataContentType(), body)
	if err != nil {
		fmt.Println("error:", err)
		return -1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("non-200 response (%d): %s\n", resp.StatusCode, string(body))
		return -1
	}

	var result responseImage
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("error decoding response:", err)
		return -1
	}

	return result.Label
}
