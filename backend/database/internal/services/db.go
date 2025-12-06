package dbservice

import (
	"context"
	model "database/internal/models"
	dbrepo "database/internal/repository/database"

	"github.com/sirupsen/logrus"
)

type Database struct {
	log                     *logrus.Logger
	addTransaction          AdderTransaction
	addUser                 AdderUser
	deleteTransaction       DeleterTransaction
	deleteUser              DeleterUser
	editTransaction         EditorTransaction
	requestUserTransactions RequesterUserTransactions
}

func New(log *logrus.Logger, db *dbrepo.DatabaseRepo) *Database {
	return &Database{
		log:                     log,
		addTransaction:          db,
		addUser:                 db,
		deleteTransaction:       db,
		deleteUser:              db,
		editTransaction:         db,
		requestUserTransactions: db,
	}
}

type AdderTransaction interface {
	AddTransaction(ctx context.Context, uid int64, t model.Transaction) (string, error)
}

type AdderUser interface {
	AddUser(ctx context.Context, uid int64) (string, error)
}

type DeleterTransaction interface {
	DeleteTransaction(ctx context.Context, uid, tid int64) (string, error)
}

type DeleterUser interface {
	DeleteUser(ctx context.Context, uid int64) (string, error)
}

type EditorTransaction interface {
	EditTransaction(ctx context.Context, uid int64, t model.Transaction) (string, error)
}

type RequesterUserTransactions interface {
	RequestUserTransactions(ctx context.Context, uid int64) ([]model.Transaction, string, error)
}

func (d *Database) AddTransaction(ctx context.Context, uid int64, t model.Transaction) (string, error) {
	return d.addTransaction.AddTransaction(ctx, uid, t)
}

func (d *Database) AddUser(ctx context.Context, uid int64) (string, error) {
	return d.addUser.AddUser(ctx, uid)
}

func (d *Database) DeleteTransaction(ctx context.Context, uid, tid int64) (string, error) {
	return d.deleteTransaction.DeleteTransaction(ctx, uid, tid)
}

func (d *Database) DeleteUser(ctx context.Context, uid int64) (string, error) {
	return d.deleteUser.DeleteUser(ctx, uid)
}

func (d *Database) EditTransaction(ctx context.Context, uid int64, t model.Transaction) (string, error) {
	return d.editTransaction.EditTransaction(ctx, uid, t)
}

func (d *Database) RequestUserTransactions(ctx context.Context, uid int64) ([]model.Transaction, string, error) {
	return d.requestUserTransactions.RequestUserTransactions(ctx, uid)
}
