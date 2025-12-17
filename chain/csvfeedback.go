package chain

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"strconv"
)

func fileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil // File exists
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil // File does not exist
	}
	return false, err // Other error occurred
}

func addNewDataToFeedBack(text string, textRes int) {
	feedbackFile := "./feedback.csv"
	textResStr := strconv.Itoa(textRes)

	// check if the file exists to determine if we need to write the header
	exists, err := fileExists(feedbackFile)
	if err != nil {
		log.Fatal(err)
	}

	// open the file the flags mean:
	// os.O_APPEND: append to the file (write to the end).
	// os.O_CREATE: create the file if it does not exist.
	// os.O_WRONLY: open the file for writing only.
	file, err := os.OpenFile(feedbackFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// if the file did not exist, write the header row first
	if !exists {
		header := []string{"Text", "Result"}
		if err := writer.Write(header); err != nil {
			log.Fatal("Error writing header to CSV:", err)
		}
	}

	// write the new data record
	data := []string{text, textResStr}
	if err := writer.Write(data); err != nil {
		log.Fatal("Error writing data to CSV:", err)
	}
}
