package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"a21hc3NpZ25tZW50/service"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// Initialize the services
var fileService = &service.FileService{}
var aiService = &service.AIService{Client: &http.Client{}}
var store = sessions.NewCookieStore([]byte("my-key"))

func getSession(r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "chat-session")
	return session
}

type Body struct {
	File     string `json:"file"`
	Question string `json:"question"`
}

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the Hugging Face token from the environment variables
	token := os.Getenv("HUGGINGFACE_TOKEN")
	if token == "" {
		log.Fatal("HUGGINGFACE_TOKEN is not set in the .env file")
	}

	// Set up the router
	router := mux.NewRouter()

	// File upload endpoint
	router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20) // Limit to 10MB
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		// Retrieve the file
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Read the file content
		fileContent, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		// Retrieve the question from the form
		question := r.FormValue("question")
		if question == "" {
			http.Error(w, "Question is required", http.StatusBadRequest)
			return
		}

		// Process the file content
		dataMap, err := fileService.ProcessFile(string(fileContent))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Analyze the data using AI
		analysisResult, err := aiService.AnalyzeData(dataMap, token, question)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the analysis result
		response := map[string]string{
			"status": "success",
			"answer": analysisResult,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		// TODO: answer here
	}).Methods("POST")

	// Chat endpoint
	router.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		// Parse the JSON request body
		var request struct {
			Query string `json:"query"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Prepare the request to the Hugging Face API
		apiURL := "https://api-inference.huggingface.co/models/microsoft/Phi-3.5-mini-instruct/v1/chat/completions"
		reqBody := map[string]interface{}{
			"messages": []map[string]string{
				{"role": "system", "content": "You are a helpful assistant."},
				{"role": "user", "content": request.Query},
			},
			"max_tokens": 1000, // Set max tokens (jawaban Ai)
		}

		jsonReqBody, err := json.Marshal(reqBody)
		if err != nil {
			http.Error(w, "Error creating request body", http.StatusInternalServerError)
			return
		}

		// Create and send the HTTP request
		httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonReqBody))
		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return
		}
		httpReq.Header.Set("Authorization", "Bearer "+os.Getenv("HUGGINGFACE_TOKEN"))
		httpReq.Header.Set("Content-Type", "application/json")

		responseAI, err := aiService.Client.Do(httpReq)
		if err != nil {
			http.Error(w, "Error communicating with AI service", http.StatusInternalServerError)
			return
		}
		defer responseAI.Body.Close()

		// Handle the AI response
		if responseAI.StatusCode != http.StatusOK {
			http.Error(w, "AI service returned an error", responseAI.StatusCode)
			return
		}

		var aiResponse struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}

		if err := json.NewDecoder(responseAI.Body).Decode(&aiResponse); err != nil {
			http.Error(w, "Error decoding AI response", http.StatusInternalServerError)
			return
		}

		// Retrieve the AI-generated text
		answer := ""
		if len(aiResponse.Choices) > 0 {
			answer = aiResponse.Choices[0].Message.Content
		}

		// Respond with the AI-generated text
		response := map[string]string{
			"status": "success",
			"answer": answer,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		// TODO: answer here
	}).Methods("POST")

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Allow your React app's origin
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler(router)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
