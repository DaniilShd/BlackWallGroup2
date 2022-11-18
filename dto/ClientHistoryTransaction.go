package dto

type ClientHistoryTransaction struct {
	History []ClientTransaction
}

type ClientTransaction struct {
	TransactionId string `json:"transaction_id"`
	Type          string `json:"type_transaction"`
	Amount        int    `json:"amount"`
}
