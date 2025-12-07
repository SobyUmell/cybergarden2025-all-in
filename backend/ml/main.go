package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	OllamaURL    = getEnv("OLLAMA_URL", "http://ollama:11434/api/chat")
	ModelName    = getEnv("MODEL_NAME", "gemma3")
	ContextLimit = 10
	Port         = getEnv("PORT", "8082")
)

var AllowedCategories = []string{
	"Misc", "Food", "Salary", "Shopping", "Electronics", "Restaurants", "Transport",
}

type ContextStore struct {
	sync.RWMutex
	History map[string][]OllamaMessage
}

var store = ContextStore{
	History: make(map[string][]OllamaMessage),
}

func main() {
	r := gin.Default()

	r.POST("/api/categorize", handleCategorize)

	r.POST("/api/chat", handleChat)

	r.POST("/api/advice", handleAdvice)

	r.DELETE("/api/context/:user_id", handleClearContext)

	log.Printf("ML Service started on port %s using model %s", Port, ModelName)
	r.Run(":" + Port)
}

func handleCategorize(c *gin.Context) {
	var req CategorizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	categoriesStr := strings.Join(AllowedCategories, ", ")
	systemPrompt := fmt.Sprintf(`You are a strict data classification machine. 
You will receive transaction details. 
You must return ONLY one word: the category name from the list below that best fits the transaction.
Allowed categories: [%s].
Do NOT write "The category is...", do NOT add punctuation. Return ONLY the category word.
If you cannot decide, return "Misc".`, categoriesStr)
	userPrompt := fmt.Sprintf("Transaction: %s, Amount: %d, Type: %s",
		req.Transaction.Description, req.Transaction.Amount, req.Transaction.Type)

	messages := []OllamaMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}
	category, err := callOllama(messages, 0.0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI Engine error", "details": err.Error()})
		return
	}
	cleanCategory := cleanResponse(category)
	if !isValidCategory(cleanCategory) {
		log.Printf("Model hallucinated: %s. Fallback to Misc.", cleanCategory)
		cleanCategory = "Misc"
	}

	c.JSON(http.StatusOK, CategorizeResponse{
		Kategoria: cleanCategory,
	})
}

func handleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	store.Lock()
	if _, exists := store.History[req.UserID]; !exists {
		store.History[req.UserID] = []OllamaMessage{}
	}
	store.History[req.UserID] = append(store.History[req.UserID], OllamaMessage{Role: "user", Content: req.Prompt})
	if len(store.History[req.UserID]) > ContextLimit {
		store.History[req.UserID] = store.History[req.UserID][len(store.History[req.UserID])-ContextLimit:]
	}

	currentContext := make([]OllamaMessage, len(store.History[req.UserID]))
	copy(currentContext, store.History[req.UserID])
	store.Unlock()
	responseContent, err := callOllama(currentContext, 0.7)
	if err != nil {
		rollbackLastMessage(req.UserID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI Engine error", "details": err.Error()})
		return
	}
	store.Lock()
	store.History[req.UserID] = append(store.History[req.UserID], OllamaMessage{Role: "assistant", Content: responseContent})
	store.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"response": responseContent,
	})
}

func handleAdvice(c *gin.Context) {
	var req AdviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	systemPrompt := `You are a world-class financial advisor. You will receive a JSON string containing a list of user's financial transactions. 
Each transaction has: 'id', 'date' (timestamp), 'kategoria' (category), 'type' (Пополнение/Доход or Списание/Покупка), 'amount' (in RUB), and 'description'.
Your task is to analyze this data and provide constructive, actionable financial advice. 
Focus on identifying spending patterns, suggesting areas for cost reduction, and offering tips on increasing savings or managing debt.
The response must be in Russian and should be formatted as a polite, friendly, and professional text.
Do NOT include the raw JSON data in your final response. Summarize the insights.
If the transaction list is empty, respond with a message about the lack of data and a suggestion to start tracking expenses.
`
	userPrompt := fmt.Sprintf("Analyze the following JSON list of transactions and provide financial advice:\n\n%s", req.Transactions)

	messages := []OllamaMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}
	adviceContent, err := callOllama(messages, 0.8)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI Engine error", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AdviceResponse{
		Advice: adviceContent,
	})
}

func callOllama(messages []OllamaMessage, temp float64) (string, error) {
	reqData := OllamaRequest{
		Model:    ModelName,
		Messages: messages,
		Stream:   false,
		Options:  map[string]interface{}{"temperature": temp},
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Ollama request: %w", err)
	}
	client := &http.Client{Timeout: 120 * time.Second}

	resp, err := client.Post(OllamaURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return "", fmt.Errorf("ollama status %d, failed to read response body: %w", resp.StatusCode, readErr)
		}
		return "", fmt.Errorf("ollama status %d: %s", resp.StatusCode, string(body))
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode Ollama response: %w", err)
	}

	return ollamaResp.Message.Content, nil
}

func handleClearContext(c *gin.Context) {
	userID := c.Param("user_id")
	store.Lock()
	delete(store.History, userID)
	store.Unlock()
	c.JSON(http.StatusOK, gin.H{"status": "cleared"})
}

func rollbackLastMessage(userID string) {
	store.Lock()
	defer store.Unlock()
	if len(store.History[userID]) > 0 {
		store.History[userID] = store.History[userID][:len(store.History[userID])-1]
	}
}

func cleanResponse(input string) string {
	res := strings.TrimSpace(input)
	res = strings.Trim(res, ".")
	res = strings.Trim(res, "\"")
	res = strings.Trim(res, "'")
	parts := strings.Fields(res)
	if len(parts) > 0 {
		return parts[0]
	}
	return "Misc"
}

func isValidCategory(cat string) bool {
	for _, c := range AllowedCategories {
		if strings.EqualFold(c, cat) {
			return true
		}
	}
	return false
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
