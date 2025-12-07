package main

type Transaction struct {
	Date        int64  `json:"date"`
	Type        string `json:"type"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

type CategorizeRequest struct {
	UserID      string      `json:"user_id" binding:"required"`
	Transaction Transaction `json:"transaction" binding:"required"`
}

type ChatRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Prompt string `json:"prompt" binding:"required"`
}

type AdviceRequest struct {
	UserID       string `json:"user_id" binding:"required"`
	Transactions string `json:"transactions" binding:"required"`
}

type CategorizeResponse struct {
	Kategoria string `json:"kategoria"`
}

type AdviceResponse struct {
	Advice string `json:"advice"`
}

type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaRequest struct {
	Model    string                 `json:"model"`
	Messages []OllamaMessage        `json:"messages"`
	Stream   bool                   `json:"stream"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

type OllamaResponse struct {
	Message OllamaMessage `json:"message"`
	Done    bool          `json:"done"`
}
