package mlrepo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	model "manager/internal/models"
	"net/http"
	"time"

	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/sirupsen/logrus"
)

type MLClient struct {
	baseURL string
	client  *http.Client
	log     *logrus.Logger
}

func New(log *logrus.Logger, host string, port int) *MLClient {
	return &MLClient{
		baseURL: fmt.Sprintf("http://%s:%d", host, port),
		client:  &http.Client{Timeout: 120 * time.Second},
		log:     log,
	}
}

type categorizeReq struct {
	UserID      string              `json:"user_id"`
	Transaction model.TransactionMl `json:"transaction"`
}

type categorizeResp struct {
	Kategoria string `json:"kategoria"`
}

type chatReq struct {
	UserID string `json:"user_id"`
	Prompt string `json:"prompt"`
}

type chatResp struct {
	Response string `json:"response"`
}

type adviceReq struct {
	UserID       string `json:"user_id"`
	Transactions string `json:"transactions"`
}

type adviceResp struct {
	Advice string `json:"advice"`
}

func (c *MLClient) CategorizeTransaction(ctx context.Context, uid int64, t model.TransactionMl) (string, error) {
	url := c.baseURL + "/api/categorize"

	body := categorizeReq{
		UserID:      fmt.Sprintf("%d", uid),
		Transaction: t,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", apperror.SystemError(err, 5001, "failed to marshal ml request")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", apperror.SystemError(err, 5002, "failed to create ml request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", apperror.SystemError(err, 5003, "failed to call ml service")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", apperror.SystemError(fmt.Errorf("ml service returned status %d", resp.StatusCode), 5004, "ml service error")
	}

	var result categorizeResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", apperror.SystemError(err, 5005, "failed to decode ml response")
	}

	return result.Kategoria, nil
}

func (c *MLClient) Chat(ctx context.Context, uid int64, prompt string) (string, error) {
	url := c.baseURL + "/api/chat"

	body := chatReq{
		UserID: fmt.Sprintf("%d", uid),
		Prompt: prompt,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", apperror.SystemError(err, 5006, "failed to marshal chat request")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", apperror.SystemError(err, 5007, "failed to create chat request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", apperror.SystemError(err, 5008, "failed to call ml chat service")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", apperror.SystemError(fmt.Errorf("ml chat service returned status %d", resp.StatusCode), 5009, "ml chat service error")
	}

	var result chatResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", apperror.SystemError(err, 5010, "failed to decode chat response")
	}

	return result.Response, nil
}

func (c *MLClient) GetAdvice(ctx context.Context, uid int64, transactions string) (string, error) {
	url := c.baseURL + "/api/advice"

	body := adviceReq{
		UserID:       fmt.Sprintf("%d", uid),
		Transactions: transactions,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", apperror.SystemError(err, 5011, "failed to marshal advice request")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", apperror.SystemError(err, 5012, "failed to create advice request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", apperror.SystemError(err, 5013, "failed to call ml advice service")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", apperror.SystemError(fmt.Errorf("ml advice service returned status %d", resp.StatusCode), 5014, "ml advice service error")
	}

	var result adviceResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", apperror.SystemError(err, 5015, "failed to decode ml advice response")
	}

	return result.Advice, nil
}
