package transaction

import "time"

type CampaignTransactionFormatter struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Amount    int    `json:"amount"`
	CreatedAt string `json:"created_at"`
}

func FormatCampaignTransaction(transaction Transaction) CampaignTransactionFormatter {
	formatter := CampaignTransactionFormatter{
		ID:        transaction.ID,
		Name:      transaction.User.Name,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return formatter
}

func FormatCampaignTransactions(transactions []Transaction) []CampaignTransactionFormatter {
	if len(transactions) == 0 {
		return []CampaignTransactionFormatter{}
	}
	transactionsFormatter := make([]CampaignTransactionFormatter, len(transactions))
	for i, transaction := range transactions {
		transactionsFormatter[i] = FormatCampaignTransaction(transaction)
	}
	return transactionsFormatter
}

type UserTransactionFormatter struct {
	ID        int                                `json:"id"`
	Amount    int                                `json:"amount"`
	Status    string                             `json:"status"`
	CreatedAt string                             `json:"created_at"`
	Campaign  CampaignTransactionDetailFormatter `json:"campaign"`
}

type CampaignTransactionDetailFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func FormatUserTransaction(transaction Transaction) UserTransactionFormatter {
	formatter := UserTransactionFormatter{
		ID:        transaction.ID,
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		CreatedAt: transaction.CreatedAt.Format("2006-01-02 15:04:05"),
		Campaign: CampaignTransactionDetailFormatter{
			Name:     transaction.Campaign.Name,
			ImageURL: "",
		},
	}
	if len(transaction.Campaign.CampaignImages) > 0 {
		formatter.Campaign.ImageURL = transaction.Campaign.CampaignImages[0].FileName
	}
	return formatter
}

func FormatUserTransactions(transactions []Transaction) []UserTransactionFormatter {
	if len(transactions) == 0 {
		return []UserTransactionFormatter{}
	}
	transactionsFormatter := make([]UserTransactionFormatter, len(transactions))
	for i, transaction := range transactions {
		transactionsFormatter[i] = FormatUserTransaction(transaction)
	}
	return transactionsFormatter
}

type TransactionFormatter struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	CampaignID int       `json:"campaign_id"`
	Amount     int       `json:"amount"`
	Status     string    `json:"status"`
	Code       string    `json:"code"`
	PaymentURL string    `json:"payment_url"`
	CreatedAt  time.Time `json:"created_at"`
}

func FormatTransaction(transaction Transaction) TransactionFormatter {
	formatter := TransactionFormatter{
		ID:         transaction.ID,
		UserID:     transaction.UserID,
		CampaignID: transaction.CampaignID,
		Amount:     transaction.Amount,
		Status:     transaction.Status,
		Code:       transaction.Code,
		PaymentURL: transaction.PaymentURL,
		CreatedAt:  transaction.CreatedAt,
	}
	return formatter
}
