package transaction

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
	var transactionsFormatter []CampaignTransactionFormatter
	for _, transaction := range transactions {
		transactionFormatter := FormatCampaignTransaction(transaction)
		transactionsFormatter = append(transactionsFormatter, transactionFormatter)
	}
	return transactionsFormatter
}
