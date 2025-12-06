package databaserepo

import (
	"context"
	"fmt"
	model "manager/internal/models"

	cyberdatabase "github.com/PrototypeSirius/protos_service/gen/cybergarden/database"
	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DBClient struct {
	db  cyberdatabase.DatabaseClient
	log *logrus.Logger
}

func New(log *logrus.Logger, host string, port int) (*DBClient, error) {
	dbaddr := fmt.Sprintf("%s:%d", host, port)
	dcc, err := grpc.NewClient(dbaddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, apperror.SystemError(err, 1021, "error load database service")
	}
	return &DBClient{db: cyberdatabase.NewDatabaseClient(dcc), log: logrus.New()}, nil
}

func (c *DBClient) AddUser(ctx context.Context, uid int64) error {
	resp, err := c.db.AddUser(ctx, &cyberdatabase.AddUserRequest{UserID: uid})
	if err != nil {
		return apperror.SystemError(err, 1031, "grpc add user failed")
	}
	if resp.ErrorMes != "" && resp.ErrorMes != "success" {
		c.log.Warnf("AddUser response: %s", resp.ErrorMes)
	}
	return nil
}

func (c *DBClient) AddTransaction(ctx context.Context, uid int64, t model.Transaction) error {
	_, err := c.db.AddTransaction(ctx, &cyberdatabase.AddTransactionRequest{
		UserID:      uid,
		Transaction: mapModelToProto(t),
	})
	if err != nil {
		return apperror.SystemError(err, 1032, "grpc add transaction failed")
	}
	return nil
}

func mapModelToProto(t model.Transaction) *cyberdatabase.Transaction {
	return &cyberdatabase.Transaction{
		ID:          t.ID,
		Date:        t.Date,
		Kategoria:   t.Kategoria,
		Type:        t.Type,
		Amount:      t.Amount,
		Description: t.Description,
	}
}

func (c *DBClient) DeleteTransaction(ctx context.Context, uid, tid int64) error {
	_, err := c.db.DeleteTransaction(ctx, &cyberdatabase.DeleteTransactionRequest{
		UserID: uid,
		ID:     tid,
	})
	if err != nil {
		return apperror.SystemError(err, 1033, "grpc delete transaction failed")
	}
	return nil
}

func (c *DBClient) EditTransaction(ctx context.Context, uid int64, t model.Transaction) error {
	_, err := c.db.EditTransaction(ctx, &cyberdatabase.EditTransactionRequest{
		UserID:      uid,
		Transaction: mapModelToProto(t),
	})
	if err != nil {
		return apperror.SystemError(err, 1034, "grpc edit transaction failed")
	}
	return nil
}

func (c *DBClient) RequestUserTransactions(ctx context.Context, uid int64) ([]model.Transaction, error) {
	resp, err := c.db.RequestUserTransactions(ctx, &cyberdatabase.RequestUserTransactionsRequest{
		UserID: uid,
	})
	if err != nil {
		return nil, apperror.SystemError(err, 1035, "grpc request transactions failed")
	}
	transactions := make([]model.Transaction, 0)
	for _, t := range resp.GetTransactions() {
		transactions = append(transactions, model.Transaction{
			ID:          t.GetID(),
			Date:        t.GetDate(),
			Kategoria:   t.GetKategoria(),
			Type:        t.GetType(),
			Amount:      t.GetAmount(),
			Description: t.GetDescription(),
		})
	}

	return transactions, nil
}
