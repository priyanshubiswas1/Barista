package pkg

type Payment struct {
	Id            string  `json:"id"`
	SessionId     string  `json:"sessionId"`
	Amount        float32 `json:"amount"`
	PaymentMethod string  `json:"paymentMethod"`
	Status        string  `json:"status"`
	TransactionId string  `json:"transactionId"`
}
