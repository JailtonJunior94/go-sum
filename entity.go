package main

import "time"

type (
	Invoice struct {
		ID    string
		Date  time.Time
		Total float64
		Items []InvoiceItem
	}

	InvoiceItem struct {
		ID               string
		PurchaseDate     time.Time
		Description      string
		TotalAmount      float64
		Installment      int
		InstallmentValue float64
		Tags             string
		Category         Category
	}

	Category struct {
		ID   string
		Name string
	}
)
