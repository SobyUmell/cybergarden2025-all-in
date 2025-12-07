package model

type Transaction struct {
	ID          int64  `json:"id"`
	Date        int64  `json:"date"`
	Kategoria   string `json:"kategoria"`
	Type        string `json:"type"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

type TransactionMl struct {
	Date        int64  `json:"date"`
	Type        string `json:"type"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

type UserID struct {
	Uid int64 `json:"uid"`
}
