package dbrepo

import (
	"context"
	"database/internal/config"
	"database/internal/migrator"
	model "database/internal/models"
	"database/sql"
	"fmt"

	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/PrototypeSirius/ruglogger/logger"
	"github.com/sirupsen/logrus"
)

type DatabaseRepo struct {
	db *sql.DB
}

func New(log *logrus.Logger, cfg config.DatabaseConfig) (*DatabaseRepo, func() error) {
	connstr := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.User,
		cfg.Password,
	)
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		appErr := apperror.SystemError(err, 1201, "Connecting to the database")
		logger.FatalOnError(appErr, "Failed to connect to the database", logrus.Fields{
			"host":     cfg.Host,
			"port":     cfg.Port,
			"database": cfg.Database,
		})
	}
	log.Info("The database is connected")
	err = migrator.Run(log, db, cfg.MigrationsPath)
	if err != nil {
		logger.FatalOnError(err, "Failed to migrate the database", logrus.Fields{
			"migrations_path": cfg.MigrationsPath,
		})
	}
	log.Info("Database migration completed")
	return &DatabaseRepo{db: db}, db.Close
}

func (d *DatabaseRepo) AddTransaction(ctx context.Context, uid int64, t model.Transaction) (string, error) {
	query := `INSERT INTO transactions (user_id, date, category, type, amount, description) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := d.db.ExecContext(ctx, query, uid, t.Date, t.Kategoria, t.Type, t.Amount, t.Description)
	if err != nil {
		return "failed to add transaction", apperror.BadRequestError(err, 1081, "Failed to add transaction")
	}
	return "", nil
}

func (d *DatabaseRepo) AddUser(ctx context.Context, uid int64) (string, error) {
	query := `INSERT INTO users (id) VALUES ($1) ON CONFLICT (id) DO NOTHING`
	_, err := d.db.ExecContext(ctx, query, uid)
	if err != nil {
		return "failed to add user", apperror.BadRequestError(err, 1082, "Failed to add user")
	}
	return "success", nil
}

func (d *DatabaseRepo) DeleteTransaction(ctx context.Context, uid, tid int64) (string, error) {
	query := `DELETE FROM transactions WHERE id = $1`
	res, err := d.db.ExecContext(ctx, query, tid)
	if err != nil {
		return "failed to delete transaction", apperror.BadRequestError(err, 1083, "Failed to delete transaction")
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return "failed to delete transaction", apperror.BadRequestError(err, 1084, "Error getting rows affected")
	}
	if rows == 0 {
		return "Transaction not found", apperror.BadRequestError(err, 1085, "Transaction not found")
	}
	return "success", nil
}

func (d *DatabaseRepo) DeleteUser(ctx context.Context, uid int64) (string, error) {
	query := `DELETE FROM users WHERE id = $1`
	res, err := d.db.ExecContext(ctx, query, uid)
	if err != nil {
		return "failed to delete user", apperror.BadRequestError(err, 1086, "Failed to delete user")
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return "failed to delete user", apperror.BadRequestError(err, 1087, "Error getting rows affected")
	}
	if rows == 0 {
		return "User not found", apperror.BadRequestError(err, 1088, "User not found")
	}
	return "success", nil
}

func (d *DatabaseRepo) EditTransaction(ctx context.Context, uid int64, t model.Transaction) (string, error) {
	query := `UPDATE transactions SET date = $1, category = $2, type = $3, amount = $4, description = $5 WHERE id = $6 AND user_id = $7`
	res, err := d.db.ExecContext(ctx, query, t.Date, t.Kategoria, t.Type, t.Amount, t.Description, t.ID, uid)
	if err != nil {
		return "failed to edit transaction", apperror.BadRequestError(err, 1089, "Failed to edit transaction")
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return "failed to edit transaction", apperror.BadRequestError(err, 1090, "Error getting rows affected")
	}
	if rows == 0 {
		return "Transaction not found", apperror.BadRequestError(err, 1091, "Transaction not found")
	}
	return "success", nil
}

func (d *DatabaseRepo) RequestUserTransactions(ctx context.Context, uid int64) ([]model.Transaction, string, error) {
	query := `SELECT id, date, category, type, amount, description FROM transactions WHERE user_id = $1 ORDER BY date DESC`
	rows, err := d.db.QueryContext(ctx, query, uid)
	if err != nil {
		return []model.Transaction{}, "failed to get transactions", apperror.BadRequestError(err, 1092, "Failed to get transactions")
	}
	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		err := rows.Scan(&t.ID, &t.Date, &t.Kategoria, &t.Type, &t.Amount, &t.Description)
		if err != nil {
			return []model.Transaction{}, "failed to get transactions", apperror.BadRequestError(err, 1093, "Error scanning transaction")
		}
		transactions = append(transactions, t)
	}
	if err := rows.Err(); err != nil {
		return []model.Transaction{}, "failed to get transactions", apperror.BadRequestError(err, 1094, "Error getting transactions")
	}
	if len(transactions) == 0 {
		return []model.Transaction{}, "No transactions found", nil
	}
	return transactions, "success", nil
}
