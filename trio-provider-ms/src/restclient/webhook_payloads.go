package restclient

type OperationType string
type OperationStatus string
type OperationCategory string

const (
	OperationTypeBalances     OperationType   = "fetch_balances"
	OperationTypeTransactions OperationType   = "fetch_transactions"
	OperationTypeAccounts     OperationType   = "fetch_accounts"
	OperationStatusCompleted  OperationStatus = "completed"
	OperationStatusFailed     OperationStatus = "failed"

	OperationCategoryAccounts OperationCategory = "ACCOUNTS"
)

type BalancesWebHookPayload struct {
	Event struct {
		AccountID string        `json:"account_id,omitempty"`
		Category  string        `json:"category,omitempty"`
		Type      OperationType `json:"type,omitempty"`
		Status    string        `json:"status,omitempty"`
	} `json:"event,omitempty"`
	Data *struct {
		Amount struct {
			Amount   float64 `json:"amount,omitempty"`
			Currency string  `json:"currency,omitempty"`
		} `json:"amount,omitempty"`
	} `json:"data,omitempty"`
	Error *struct {
		Message *string `json:"message,omitempty"`
	} `json:"error,omitempty"`
}

type TransactionsWebHookPayload struct {
	Event struct {
		AccountID string        `json:"account_id,omitempty"`
		Category  string        `json:"category,omitempty"`
		Type      OperationType `json:"type,omitempty"`
		Status    string        `json:"status,omitempty"`
	} `json:"event,omitempty"`
	Data *struct {
		TotalTransactions int64 `json:"total_transactions,omitempty"`
	} `json:"data,omitempty"`
	Error *struct {
		Message *string `json:"message,omitempty"`
	} `json:"error,omitempty"`
}

type AccountsWebHookPayload struct {
	Event struct {
		ConnectionID string            `json:"connection_id,omitempty"`
		Category     OperationCategory `json:"category,omitempty"`
		Type         OperationType     `json:"type,omitempty"`
		Status       string            `json:"status,omitempty"`
	} `json:"event,omitempty"`
	Data struct {
		TotalAccounts int64  `json:"total_accounts,omitempty"`
		HolderID      string `json:"holder_id,omitempty"`
	} `json:"data,omitempty"`
	Error struct {
		Message string `json:"message,omitempty"`
	} `json:"error,omitempty"`
}
