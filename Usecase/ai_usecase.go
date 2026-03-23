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

func (u *AIUsecase) GenerateResponse(noteID, userID primitive.ObjectID, message string) (string, error) {
	note, err := u.Repo.GetNoteByID(noteID, userID)
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`You are a helpful AI assistant for the Notebook app.

Current note title: "%s"

Full note content:
%s

User request: %s

Answer clearly, use markdown when helpful, keep it concise and useful.`, note.Title, note.Content, message)

	return u.callOpenAI(prompt)
}

func (u *AIUsecase) callOpenAI(prompt string) (string, error) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		return "", errors.New("OpenAI API key is not configured")
	}

	reqBody := map[string]any{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a helpful assistant specialized in note-taking, summarizing, explaining, and expanding ideas."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
		"max_tokens":  1000,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
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
		return "", fmt.Errorf("OpenAI error (%d): %s", resp.StatusCode, string(body))
	}

	var openaiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&openaiResp); err != nil {
		return "", err
	}

	if len(openaiResp.Choices) == 0 || openaiResp.Choices[0].Message.Content == "" {
		return "", errors.New("no response from OpenAI")
	}

	return openaiResp.Choices[0].Message.Content, nil
}
