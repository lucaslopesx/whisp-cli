package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("failed to load .env")
		os.Exit(1)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("error missing OPENAI_API_KEY")
		os.Exit(1)
	}

	var filePath string
	flag.StringVar(&filePath, "path", "./input.mp3", "path to the audio file to be transcribed (e.g., ./input.mp3)")
	flag.Parse()

	if filePath == "" {
		fmt.Println("error: missing required -path flag")
		flag.Usage()
		os.Exit(1)
	}

	res, err := transcribe(apiKey, filePath)
	if err != nil {
		fmt.Printf("Error: %v", err)
		fmt.Println()
		os.Exit(1)
	}

	fmt.Println("Transcription:")
	fmt.Println(res)
}

func transcribe(apiKey, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy audio to form: %v", err)
	}

	writer.WriteField("model", "whisper-1")
	writer.Close()

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/transcriptions", body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return string(respBody), nil
}
