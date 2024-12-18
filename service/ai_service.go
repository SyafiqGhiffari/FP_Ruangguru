package service

import (
	"a21hc3NpZ25tZW50/model"
	"net/http"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AIService struct {
	Client HTTPClient
}

func (s *AIService) AnalyzeData(table map[string][]string, token, query string) (string, error) {
	if len(table) == 0 {
		return "", errors.New("tabel kosong")
	}

	payload := map[string]interface{}{
		"inputs": map[string]interface{}{
			"table": table,
			"query": query,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", errors.New("gagal membuat payload JSON")
	}

	request, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/google/tapas-base-finetuned-wtq", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", errors.New("gagal membuat permintaan HTTP")
	}

	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")

	response, err := s.Client.Do(request)
	if err != nil {
		return "", errors.New("gagal mengirim permintaan ke API")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return "", fmt.Errorf("API gagal: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", errors.New("gagal mendekode respons JSON")
	}

	answer, ok := result["answer"].(string)
	if !ok {
		log.Printf("Tidak ditemukan 'answer' dalam respons: %v", result)
		return "", errors.New("respon tidak valid atau 'answer' tidak ditemukan")
	}
	confidence, _ := result["confidence"].(float64)
	if confidence < 0.5 {
		log.Printf("confidence rate: %.2f", confidence)
	}

	return answer, nil // TODO: replace this
}

func (s *AIService) ChatWithAI(context, query, token string) (model.ChatResponse, error) {
	// Prepare payload
	payload := map[string]string{
		"context": context,
		"query":   query,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return model.ChatResponse{}, errors.New("gagal membuat payload JSON")
	}

	// Create HTTP request
	request, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/microsoft/Phi-3.5-mini-instruct/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return model.ChatResponse{}, errors.New("gagal membuat permintaan HTTP")
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")

	// Send request
	response, err := s.Client.Do(request)
	if err != nil {
		return model.ChatResponse{}, err
	}
	defer response.Body.Close()

	// Check successful response
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return model.ChatResponse{}, fmt.Errorf("failed to chat with AI: %s, response: %s", response.Status, string(body))
	}

	// Decode the response into a slice of ChatResponse
	var chatResponses []model.ChatResponse
	if err := json.NewDecoder(response.Body).Decode(&chatResponses); err != nil {
		return model.ChatResponse{}, err
	}

	// Return the first response (or handle it as needed)
	if len(chatResponses) > 0 {
		return chatResponses[0], nil
	}

	return model.ChatResponse{}, errors.New("no responses found")
	// TODO: answer here
}
