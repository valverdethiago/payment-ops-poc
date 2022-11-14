package restclient

import "time"

type AccountType string
type BankAccountType string
type FetchRequestStatus string

const (
	BANK_ACCOUNT_TYPE AccountType        = "bank_account"
	CHECKING          BankAccountType    = "checking"
	SAVINGS           BankAccountType    = "savings"
	COMPLETED         FetchRequestStatus = "completed"
	FAILED            FetchRequestStatus = "failed"
)

type ListTransactionsResponse struct {
	Data []struct {
		Account struct {
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
		} `json:"account"`
		AccountingType string `json:"accounting_type"`
		Amount         struct {
			Amount   int    `json:"amount"`
			Currency string `json:"currency"`
		} `json:"amount"`
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
	} `json:"data"`
	Meta struct {
		After  string `json:"after"`
		Before string `json:"before"`
	} `json:"meta"`
}

type ListBalanceResponse struct {
	Data []struct {
		Account struct {
			AccountNumber   string          `json:"account_number,omitempty"`
			AccountType     AccountType     `json:"account_type,omitempty"`
			BankAccountType BankAccountType `json:"bank_account_type,omitempty"`
			BranchNumber    string          `json:"branch_number,omitempty"`
			ConnectionId    string          `json:"connection_id,omitempty"`
			Currency        string          `json:"currency,omitempty"`
			ID              string          `json:"id,omitempty"`
			Institution     struct {
				Code string `json:"code,omitempty"`
				ID   string `json:"id,omitempty"`
				Name string `json:"name,omitempty"`
			}
			VariationCode string `json:"variation_code,omitempty"`
		}
		Amount struct {
			Amount   float64 `json:"amount,omitempty"`
			Currency string  `json:"currency,omitempty"`
		}
		ID         string    `json:"id,omitempty"`
		InsertedAt time.Time `json:"inserted_at,omitempty"`
		Timestamp  time.Time `json:"timestamp,omitempty"`
		UpdatedAt  time.Time `json:"updated_at,omitempty"`
	} `json:"data,omitempty"`
	Meta struct {
		After  string `json:"after,omitempty"`
		Before string `json:"before,omitempty"`
	} `json:"meta,omitempty"`
}

type FetchRequestResponse struct {
	Event struct {
		AccountId string             `json:"account_id"`
		Category  string             `json:"category"`
		Type      string             `json:"type"`
		Status    FetchRequestStatus `json:"status"`
	} `json:"event"`
	Data struct {
		Amount struct {
			Amount   int    `json:"amount"`
			Currency string `json:"currency"`
		} `json:"amount"`
	} `json:"data"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}
