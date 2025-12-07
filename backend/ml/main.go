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

// Конфигурация
var (
	OllamaURL    = getEnv("OLLAMA_URL", "http://ollama:11434/api/chat")
	ModelName    = getEnv("MODEL_NAME", "gemma3")
	ContextLimit = 10
	Port         = getEnv("PORT", "8082")
)

// Разрешенные категории
var AllowedCategories = []string{
	"Misc", "Food", "Salary", "Shopping", "Electronics", "Restaurants", "Transport",
}

// Хранилище контекста для чата
type ContextStore struct {
	sync.RWMutex
	History map[string][]OllamaMessage
}

var store = ContextStore{
	History: make(map[string][]OllamaMessage),
}

func main() {
	r := gin.Default()

	// 1. Хендлер для строгой классификации
	r.POST("/api/categorize", handleCategorize)

	// 2. Хендлер для чата с контекстом
	r.POST("/api/chat", handleChat)

	// Сброс контекста
	r.DELETE("/api/context/:user_id", handleClearContext)

	log.Printf("ML Service started on port %s using model %s", Port, ModelName)
	r.Run(":" + Port)
}

// --- ЛОГИКА КЛАССИФИКАЦИИ ---

func handleCategorize(c *gin.Context) {
	var req CategorizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Формируем строгий системный промпт
	categoriesStr := strings.Join(AllowedCategories, ", ")
	systemPrompt := fmt.Sprintf(`You are a strict data classification machine. 
You will receive transaction details. 
You must return ONLY one word: the category name from the list below that best fits the transaction.
Allowed categories: [%s].
Do NOT write "The category is...", do NOT add punctuation. Return ONLY the category word.
If you cannot decide, return "Misc".`, categoriesStr)

	// Формируем описание транзакции для модели
	userPrompt := fmt.Sprintf("Transaction: %s, Amount: %d, Type: %s",
		req.Transaction.Description, req.Transaction.Amount, req.Transaction.Type)

	messages := []OllamaMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	// Отправляем в Ollama с температурой 0 (максимальная точность)
	category, err := callOllama(messages, 0.0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI Engine error", "details": err.Error()})
		return
	}

	// Чистим ответ (убираем пробелы, точки)
	cleanCategory := cleanResponse(category)

	// Валидация (если модель все-таки сошла с ума, ставим Misc)
	if !isValidCategory(cleanCategory) {
		log.Printf("Model hallucinated: %s. Fallback to Misc.", cleanCategory)
		cleanCategory = "Misc"
	}

	c.JSON(http.StatusOK, CategorizeResponse{
		Kategoria: cleanCategory,
	})
}

// --- ЛОГИКА ЧАТА ---

func handleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Работа с историей
	store.Lock()
	if _, exists := store.History[req.UserID]; !exists {
		store.History[req.UserID] = []OllamaMessage{}
	}
	store.History[req.UserID] = append(store.History[req.UserID], OllamaMessage{Role: "user", Content: req.Prompt})

	// Обрезка истории
	if len(store.History[req.UserID]) > ContextLimit {
		store.History[req.UserID] = store.History[req.UserID][len(store.History[req.UserID])-ContextLimit:]
	}

	currentContext := make([]OllamaMessage, len(store.History[req.UserID]))
	copy(currentContext, store.History[req.UserID])
	store.Unlock()

	// Отправляем в Ollama с температурой 0.7 (креативность для чата)
	responseContent, err := callOllama(currentContext, 0.7)
	if err != nil {
		// Откат истории при ошибке
		rollbackLastMessage(req.UserID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI Engine error", "details": err.Error()})
		return
	}

	// Сохраняем ответ
	store.Lock()
	store.History[req.UserID] = append(store.History[req.UserID], OllamaMessage{Role: "assistant", Content: responseContent})
	store.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"response": responseContent,
	})
}

// --- ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ---

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
	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Post(OllamaURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama status %d: %s", resp.StatusCode, string(body))
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", err
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
	// Удаляем лишние пробелы, точки в конце, переносы строк
	res := strings.TrimSpace(input)
	res = strings.Trim(res, ".")
	res = strings.Trim(res, "\"")
	res = strings.Trim(res, "'")
	// Берем первое слово, если модель вдруг выдала "Food category"
	parts := strings.Fields(res)
	if len(parts) > 0 {
		return parts[0]
	}
	return "Misc"
}

func isValidCategory(cat string) bool {
	for _, c := range AllowedCategories {
		if strings.EqualFold(c, cat) { // Регистронезависимое сравнение
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
