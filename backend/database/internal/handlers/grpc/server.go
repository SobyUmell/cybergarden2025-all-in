package grpchandler

import (
	"context"
	model "database/internal/models"
	"errors"

	database "github.com/PrototypeSirius/protos_service/gen/cybergarden/database"
	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/PrototypeSirius/ruglogger/logger"
	"google.golang.org/grpc"
)

type DBService interface {
	AddTransaction(ctx context.Context, uid int64, t model.Transaction) (string, error)
	AddUser(ctx context.Context, uid int64) (string, error)
	DeleteTransaction(ctx context.Context, uid, tid int64) (string, error)
	DeleteUser(ctx context.Context, uid int64) (string, error)
	EditTransaction(ctx context.Context, uid int64, t model.Transaction) (string, error)
	RequestUserTransactions(ctx context.Context, uid int64) ([]model.Transaction, string, error)
}

type serverAPI struct {
	database.UnimplementedDatabaseServer
	db DBService
}

func Register(s *grpc.Server, db DBService) {
	database.RegisterDatabaseServer(s, &serverAPI{db: db})
}

func (s *serverAPI) AddTransaction(ctx context.Context, req *database.AddTransactionRequest) (*database.AddTransactionResponse, error) {
	if req.GetUserID() == 0 {
		appErr := apperror.BadRequestError(errors.New("empty user id"), 1050, "Error while adding transaction")
		logger.LogOnError(appErr, "Error in adding transaction")
		return nil, appErr
	}
	if req.GetTransaction() == nil {
		appErr := apperror.BadRequestError(errors.New("empty transaction"), 1051, "Error while adding transaction")
		logger.LogOnError(appErr, "Error in adding transaction")
		return nil, appErr
	}
	mes, err := s.db.AddTransaction(ctx, req.GetUserID(), interpretatorTransactionAdd(req))
	if err != nil {
		logger.LogOnError(err, "Error in adding transaction")
		return &database.AddTransactionResponse{ErrorMes: mes}, err
	}
	return &database.AddTransactionResponse{ErrorMes: mes}, nil
}

func (s *serverAPI) AddUser(ctx context.Context, req *database.AddUserRequest) (*database.AddUserResponse, error) {
	if req.GetUserID() == 0 {
		appErr := apperror.BadRequestError(errors.New("empty user id"), 1052, "Error while adding user")
		logger.LogOnError(appErr, "Error in adding user")
		return nil, appErr
	}
	mes, err := s.db.AddUser(ctx, req.GetUserID())
	if err != nil {
		logger.LogOnError(err, "Error in adding user")
		return &database.AddUserResponse{ErrorMes: mes}, err
	}
	return &database.AddUserResponse{ErrorMes: mes}, nil
}

func (s *serverAPI) DeleteTransaction(ctx context.Context, req *database.DeleteTransactionRequest) (*database.DeleteTransactionResponse, error) {
	if req.GetUserID() == 0 {
		appErr := apperror.BadRequestError(errors.New("empty user id"), 1053, "Error while deleting transaction")
		logger.LogOnError(appErr, "Error in deleting transaction")
		return nil, appErr
	}
	if req.GetID() == 0 {
		appErr := apperror.BadRequestError(errors.New("empty transaction id"), 1054, "Error while deleting transaction")
		logger.LogOnError(appErr, "Error in deleting transaction")
		return nil, appErr
	}
	mes, err := s.db.DeleteTransaction(ctx, req.GetUserID(), req.GetID())
	if err != nil {
		logger.LogOnError(err, "Error in deleting transaction")
		return &database.DeleteTransactionResponse{ErrorMes: mes}, err
	}
	return &database.DeleteTransactionResponse{ErrorMes: mes}, nil
}

func (s *serverAPI) DeleteUser(ctx context.Context, req *database.DeleteUserRequest) (*database.DeleteUserResponse, error) {
	if req.GetUserID() == 0 {
		appErr := apperror.BadRequestError(errors.New("empty user id"), 1055, "Error while deleting user")
		logger.LogOnError(appErr, "Error in deleting user")
		return nil, appErr
	}
	mes, err := s.db.DeleteUser(ctx, req.GetUserID())
	if err != nil {
		logger.LogOnError(err, "Error in deleting user")
		return &database.DeleteUserResponse{ErrorMes: mes}, err
	}
	return &database.DeleteUserResponse{ErrorMes: mes}, nil
}

func (s *serverAPI) EditTransaction(ctx context.Context, req *database.EditTransactionRequest) (*database.EditTransactionResponse, error) {
	if req.GetUserID() == 0 {
		appErr := apperror.BadRequestError(errors.New("empty user id"), 1056, "Error while editing transaction")
		logger.LogOnError(appErr, "Error in editing transaction")
		return nil, appErr
	}
	if req.GetTransaction() == nil {
		appErr := apperror.BadRequestError(errors.New("empty transaction"), 1057, "Error while editing transaction")
		logger.LogOnError(appErr, "Error in editing transaction")
		return nil, appErr
	}
	mes, err := s.db.EditTransaction(ctx, req.GetUserID(), interpretatorTransactionEdit(req))
	if err != nil {
		logger.LogOnError(err, "Error in editing transaction")
		return &database.EditTransactionResponse{ErrorMes: mes}, err
	}
	return &database.EditTransactionResponse{}, nil
}

func (s *serverAPI) RequestUserTransactions(ctx context.Context, req *database.RequestUserTransactionsRequest) (*database.RequestUserTransactionsResponse, error) {
	if req.GetUserID() == 0 {
		appErr := apperror.BadRequestError(errors.New("empty user id"), 1058, "Error while requesting user transactions")
		logger.LogOnError(appErr, "Error in requesting user transactions")
		return nil, appErr
	}
	transactions, mes, err := s.db.RequestUserTransactions(ctx, req.GetUserID())
	if err != nil {
		logger.LogOnError(err, "Error in requesting user transactions")
		return &database.RequestUserTransactionsResponse{ErrorMes: mes}, err
	}
	return interpretatorTransactionResponse(transactions), nil
}

func interpretatorTransactionAdd(req *database.AddTransactionRequest) model.Transaction {
	t := req.GetTransaction()
	return model.Transaction{
		ID:          t.ID,
		Date:        t.Date,
		Kategoria:   t.Kategoria,
		Type:        t.Type,
		Amount:      t.Amount,
		Description: t.Description,
	}
}

func interpretatorTransactionEdit(req *database.EditTransactionRequest) model.Transaction {
	t := req.GetTransaction()
	return model.Transaction{
		ID:          t.ID,
		Date:        t.Date,
		Kategoria:   t.Kategoria,
		Type:        t.Type,
		Amount:      t.Amount,
		Description: t.Description,
	}
}

func interpretatorTransactionResponse(transactions []model.Transaction) *database.RequestUserTransactionsResponse {
	var protoTransactions []*database.Transaction
	for _, t := range transactions {
		protoTransactions = append(protoTransactions, &database.Transaction{
			ID:          t.ID,
			Date:        t.Date,
			Kategoria:   t.Kategoria,
			Type:        t.Type,
			Amount:      t.Amount,
			Description: t.Description,
		})
	}
	return &database.RequestUserTransactionsResponse{
		Transactions: protoTransactions,
	}
}
