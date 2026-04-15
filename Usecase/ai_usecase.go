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

	prompt := fmt.Sprintf(`
You are an AI assistant inside a Notebook application.

Your job is to help users understand, summarize, and expand their notes.

-----------------------------------
NOTE TITLE:
%s

NOTE CONTENT:
%s
-----------------------------------

USER REQUEST:
%s

-----------------------------------
INSTRUCTIONS:

1. Detect the user's intent:

- If the user asks to SUMMARIZE:
  • Return ONLY bullet points
  • Maximum 3–5 bullet points
  • Each bullet must be short (1 sentence)
  • Do NOT explain anything

- If the user asks to EXPLAIN:
  • Give a clear explanation in simple terms
  • MUST include at least one real-world example
  • Use short paragraphs or bullet points

- If the user asks to EXPAND:
  • Add new ideas or deeper details
  • Include examples if possible

- If the request is unclear:
  • Ask a clarifying question to better understand the user's intent

2. Response rules:
   - Be clear, concise, and helpful
   - Use markdown formatting (headings, bullet points) when useful
   - Do NOT repeat the full note unless necessary
   - Focus only on relevant parts of the note

3. Output format:
   - Use headings if needed
   - Use bullet points for lists
   - Keep response structured and easy to read

4. Keep response length moderate (not too long, not too short)

-----------------------------------
FINAL ANSWER:
`, note.Title, note.Content, message)

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
