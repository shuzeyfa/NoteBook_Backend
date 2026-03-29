package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	domain "taskmanagement/Domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AIUsecase struct {
	Repo domain.NoteRepository
}

// GenerateResponse is called from controller
func (u *AIUsecase) GenerateResponse(noteID, userID primitive.ObjectID, message string) (string, error) {
	note, err := u.Repo.GetNoteByID(noteID, userID)
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`You are a helpful AI assistant inside a Notebook app.

Current note title: "%s"

Full note content:
%s

User's request: %s

Answer in a clear, friendly and useful way. Use markdown when helpful. Keep responses concise but complete.`,
		note.Title, note.Content, message)

	return u.callGroq(prompt)
}

func (u *AIUsecase) callGroq(prompt string) (string, error) {
	key := os.Getenv("GROQ_API_KEY")
	if key == "" {
		return "", errors.New("GROQ_API_KEY is not set in .env file")
	}

	reqBody := map[string]any{
		"model": "llama-3.3-70b-versatile", // Best free model on Groq (very strong)
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a smart, helpful, and concise assistant specialized in note-taking, summarizing, explaining, and expanding ideas.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
		"max_tokens":  1200,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Groq API Error (%d): %s", resp.StatusCode, string(body))
	}

	var groqResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&groqResp); err != nil {
		return "", err
	}

	if len(groqResp.Choices) == 0 || groqResp.Choices[0].Message.Content == "" {
		return "", errors.New("no response received from Groq")
	}

	return groqResp.Choices[0].Message.Content, nil
}
