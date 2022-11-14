package restclient

import "time"

type AccountType string
type BankAccountType string
type FetchRequestStatus string

const (
	CheckingBankAccountType    BankAccountType    = "checking"
	SavingsBankAccountType     BankAccountType    = "savings"
	CompleteFetchRequestStatus FetchRequestStatus = "completed"
	FailedFetchRequestStatus   FetchRequestStatus = "failed"
)

type Account struct {
	AccountNumber   string `json:"account_number"`
	AccountType     string `json:"account_type"`
	BankAccountType string `json:"bank_account_type"`
	BranchNumber    string `json:"branch_number"`
	Currency        string `json:"currency"`
	Id              string `json:"id"`
	Institution     struct {
		Code string `json:"code"`
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"institution"`
	VariationCode interface{} `json:"variation_code"`
}

type Transaction struct {
	Account         Account     `json:"account"`
	AccountingType  string      `json:"accounting_type"`
	Amount          Amount      `json:"amount"`
	Category        interface{} `json:"category"`
	CreditCardData  interface{} `json:"credit_card_data"`
	Description     string      `json:"description"`
	DescriptionType string      `json:"description_type"`
	Id              string      `json:"id"`
	Identification  string      `json:"identification"`
	InsertedAt      string      `json:"inserted_at"`
	Status          string      `json:"status"`
	Timestamp       time.Time   `json:"timestamp"`
	Type            interface{} `json:"type"`
	UpdatedAt       string      `json:"updated_at"`
}

type Meta struct {
	After  string `json:"after"`
	Before string `json:"before"`
}
type Error struct {
	Message string `json:"message"`
}

type ListTransactionsResponse struct {
	Data []Transaction `json:"data"`
	Meta Meta          `json:"meta"`
}

type Amount struct {
	Amount   float64 `json:"amount,omitempty"`
	Currency string  `json:"currency,omitempty"`
}

type Event struct {
	AccountId string             `json:"account_id"`
	Category  string             `json:"category"`
	Type      string             `json:"type"`
	Status    FetchRequestStatus `json:"status"`
}

type ListBalanceResponse struct {
	Data []struct {
		Account    Account   `json:"Account"`
		Amount     Amount    `json:"Amount"`
		ID         string    `json:"id,omitempty"`
		InsertedAt time.Time `json:"inserted_at,omitempty"`
		Timestamp  time.Time `json:"timestamp,omitempty"`
		UpdatedAt  time.Time `json:"updated_at,omitempty"`
	} `json:"data,omitempty"`
	Meta Meta `json:"meta,omitempty"`
}

type FetchRequestResponse struct {
	Data struct {
		AccountId string `json:"account_id"`
		Status    string `json:"status"`
		Type      string `json:"type"`
	} `json:"data"`
}
