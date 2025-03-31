package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jailtonjunior94/go-sum/pkg/excel"

	"github.com/shopspring/decimal"
)

func main() {
	db, err := NewSQLServerDatabase("sqlserver", "")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	queries := NewQueries(db)

	ctx := context.Background()
	provider := excel.NewProvider()
	xls := provider.NewFile(ctx)

	date := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	invoices, err := queries.GetInvoices(date)
	if err != nil {
		log.Fatal(err)
	}

	groupedByTag := GroupBy(invoices.Items, func(i InvoiceItem) string {
		return i.Tags
	})

	tagSheet := xls.NewSheet(ctx, date.Format("2006-01"))
	for tag, group := range groupedByTag {
		sum := Sum(group, func(i InvoiceItem) float64 {
			return i.InstallmentValue
		})

		sumDecimal := decimal.NewFromFloat(sum).Round(2)
		fmt.Printf("Total de gastos na categoria %s: R$ %s\n", tag, sumDecimal)

		row := NewBudgetExportRow(tag, sumDecimal.String())
		if err := tagSheet.Write(ctx, row); err != nil {
			log.Fatal(err)
		}
	}

	err = xls.SaveAs(ctx, fmt.Sprintf("./files/orcamento-domestico-%s.xlsx", date.Format("2006-01")))
	if err != nil {
		log.Fatal(err)
	}
}

func GroupBy[T any, K comparable](items []T, fn func(T) K) map[K][]T {
	grouped := make(map[K][]T)
	for _, item := range items {
		key := fn(item)
		grouped[key] = append(grouped[key], item)
	}
	return grouped
}

func Sum[T any](items []T, fn func(T) float64) float64 {
	var sum float64
	for _, item := range items {
		sum += fn(item)
	}
	return sum
}

type ItemsExportRow struct {
	Category  interface{} `column:"A" header:"Categoria"`
	Price     interface{} `column:"B" header:"Preço"`
	Item      interface{} `column:"C" header:"Itens"`
	PriceItem interface{} `column:"D" header:"Valor do Item"`
}

func NewItemsExportRow(category, price, item, priceItem string) ItemsExportRow {
	return ItemsExportRow{
		Category:  category,
		Price:     price,
		Item:      item,
		PriceItem: priceItem,
	}
}

type BudgetExportRow struct {
	Budget interface{} `column:"A" header:"Orçamento"`
	Price  interface{} `column:"B" header:"Preço"`
}

func NewBudgetExportRow(budget, price string) BudgetExportRow {
	return BudgetExportRow{
		Budget: budget,
		Price:  price,
	}
}

type FilterFunc[T any] func(T) bool

func Filter[T any](items []T, fn FilterFunc[T]) []T {
	newSlice := make([]T, len(items))
	for index, item := range items {
		if fn(item) {
			newSlice[index] = item
		}
	}
	return newSlice
}
