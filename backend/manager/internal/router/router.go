package router

import (
	"context"
	"errors"
	model "manager/internal/models"
	"net/http"
	"strings"

	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ServiceAPI interface {
	AuthUser(ctx context.Context, initData string) (int64, error)
	AddTransaction(ctx context.Context, uid int64, t model.TransactionMl) error
	DeleteTransaction(ctx context.Context, uid, tid int64) error
	EditTransaction(ctx context.Context, uid int64, t model.Transaction) error
	GetHistory(ctx context.Context, uid int64) ([]model.Transaction, error)
	Chat(ctx context.Context, uid int64, prompt string) (string, error)
	GetFinancialAdvice(ctx context.Context, uid int64) (string, error)
}

type Handler struct {
	log     *logrus.Logger
	service ServiceAPI
}

func New(log *logrus.Logger, service ServiceAPI) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}

func (h *Handler) RouterRegister(r *gin.Engine) {
	api := r.Group("/webapp")
	api.Use(h.authUser)
	{
		api.GET("/datahistory", h.requestDataHistory)
		api.POST("/addt", h.addTransaction)
		api.POST("/deletet", h.deleteTransaction)
		api.POST("/updatet", h.updateTransaction)
		api.POST("/chat", h.chat)
		api.GET("/advice", h.requestAdvice)
	}
}

func (h *Handler) requestDataHistory(c *gin.Context) {
	uid := c.GetInt64("uid")

	history, err := h.service.GetHistory(c.Request.Context(), uid)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": history})
}

func (h *Handler) addTransaction(c *gin.Context) {
	uid := c.GetInt64("uid")

	var t model.TransactionMl
	if err := c.BindJSON(&t); err != nil {
		_ = c.Error(apperror.BadRequestError(err, 4001, "invalid json body"))
		return
	}

	if err := h.service.AddTransaction(c.Request.Context(), uid, t); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (h *Handler) deleteTransaction(c *gin.Context) {
	uid := c.GetInt64("uid")

	var req struct {
		ID int64 `json:"id"`
	}
	if err := c.BindJSON(&req); err != nil {
		_ = c.Error(apperror.BadRequestError(err, 4002, "invalid json body"))
		return
	}

	if err := h.service.DeleteTransaction(c.Request.Context(), uid, req.ID); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
func (h *Handler) updateTransaction(c *gin.Context) {
	uid := c.GetInt64("uid")

	var t model.Transaction
	if err := c.BindJSON(&t); err != nil {
		_ = c.Error(apperror.BadRequestError(err, 4003, "invalid json body"))
		return
	}

	if err := h.service.EditTransaction(c.Request.Context(), uid, t); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *Handler) authUser(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	var authData string
	if authHeader != "" {
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 {
			appErr := apperror.BadRequestError(errors.New("invalid auth header"), 4004, "Invalid authorization header format. Expected 'tma <data>'.")
			_ = c.Error(appErr)
			return
		}
		authData = authParts[1]
	} else {
		raw := c.Request.URL.RawQuery
		const tmaPrefix = "tma"
		idx := strings.Index(raw, tmaPrefix)
		if idx == -1 {
			authData = ""
		} else {
			authData = raw[idx+len(tmaPrefix)+1:]
		}
	}

	if authData == "" {
		err := apperror.BadRequestError(errors.New("missing auth header"), 4005, "Missing authorization header")
		_ = c.Error(err)
		return
	}

	uid, err := h.service.AuthUser(c.Request.Context(), authData)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set("uid", uid)
	c.Next()
}

func (h *Handler) chat(c *gin.Context) {
	uid := c.GetInt64("uid")

	var req struct {
		Prompt string `json:"message"`
	}
	if err := c.BindJSON(&req); err != nil {
		_ = c.Error(apperror.BadRequestError(err, 4006, "invalid json body"))
		return
	}

	resp, err := h.service.Chat(c.Request.Context(), uid, req.Prompt)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": resp})
}

func (h *Handler) requestAdvice(c *gin.Context) {
	uid := c.GetInt64("uid")

	advice, err := h.service.GetFinancialAdvice(c.Request.Context(), uid)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"advice": advice})
}
