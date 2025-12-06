package model

type Transaction struct {
	ID          int64
	Date        int64
	Kategoria   string
	Type        string
	Amount      int64
	Description string
}
