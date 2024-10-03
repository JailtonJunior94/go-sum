package main

import (
	"database/sql"
	"time"
)

type queries struct {
	db *sql.DB
}

func NewQueries(db *sql.DB) *queries {
	return &queries{db: db}
}

func (q *queries) GetInvoiceDates() ([]time.Time, error) {
	query := `SELECT DISTINCT i.[Date] FROM Invoice i ORDER BY i.[Date]`

	rows, err := q.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dates []time.Time
	for rows.Next() {
		var date time.Time
		if err := rows.Scan(&date); err != nil {
			return nil, err
		}
		dates = append(dates, date)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dates, nil
}

func (q *queries) GetInvoices(date time.Time) (*Invoice, error) {
	query := `SELECT
				CAST(i.Id AS CHAR(36)) [InvoiceID],
				i.[Date],
				i.Total,
				CAST(ii.Id AS CHAR(36)) [InvoiceItemID],
				ii.PurchaseDate,
				ii.Description,
				ii.TotalAmount,
				ii.Installment,
				ii.InstallmentValue,
				ii.Tags,
				CAST(c2.Id AS CHAR(36)) [CategoryID],
				c2.Name
			FROM
				Invoice i
				inner join InvoiceItem ii on ii.InvoiceId = i.Id
				inner join Category c2 on c2.Id = ii.CategoryId
			WHERE
				i.Date = @date
				AND ii.Tags != ''
			ORDER BY ii.PurchaseDate`

	rows, err := q.db.Query(query, sql.Named("date", date))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoice Invoice
	var invoiceItem InvoiceItem
	var invoiceItemMap = make(map[string][]InvoiceItem)

	for rows.Next() {
		err = rows.Scan(
			&invoice.ID,
			&invoice.Date,
			&invoice.Total,
			&invoiceItem.ID,
			&invoiceItem.PurchaseDate,
			&invoiceItem.Description,
			&invoiceItem.TotalAmount,
			&invoiceItem.Installment,
			&invoiceItem.InstallmentValue,
			&invoiceItem.Tags,
			&invoiceItem.Category.ID,
			&invoiceItem.Category.Name,
		)
		if err != nil {
			return nil, err
		}

		if _, ok := invoiceItemMap[invoice.ID]; !ok {
			invoiceItemMap[invoice.ID] = []InvoiceItem{invoiceItem}
			continue
		}
		invoiceItemMap[invoice.ID] = append(invoiceItemMap[invoice.ID], invoiceItem)
	}

	invoice.Items = invoiceItemMap[invoice.ID]
	return &invoice, nil
}
