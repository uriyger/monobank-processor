package model

// MonoStatement represents mono request body
type MonoStatement struct {
	Type string
	Data struct {
		StatementItem StatementItem `json:"statementItem"`
		Account       string        `json:"account"`
	}
}

// StatementItem represents statement from bank
type StatementItem struct {
	ID              string `json:"id"`
	Time            int64  `json:"time"`
	Description     string `json:"description"`
	Mcc             int    `json:"mcc"`
	Comment         string `json:"comment"`
	Hold            bool   `json:"hold"`
	Amount          int    `json:"amount"`
	OperationAmount int    `json:"operationAmount"`
	CurrencyCode    int    `json:"currencyCode"`
	CommissionRate  int    `json:"commissionRate"`
	CashbackAmount  int    `json:"cashbackAmount"`
	Balance         int    `json:"balance"`
}
