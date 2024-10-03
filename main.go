package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
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

	// invoiceDates, err := queries.GetInvoiceDates()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	ctx := context.Background()
	provider := excel.NewProvider()
	xls := provider.NewFile(ctx)

	date := time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)
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
		fmt.Printf("Total gasto da tag %s: R$ %s\n", tag, sumDecimal)

		row := NewBudgetExportRow(tag, sumDecimal.String())
		if err := tagSheet.Write(ctx, row); err != nil {
			log.Fatal(err)
		}
	}

	// categorySheet := xls.NewSheet(ctx, date.Format("2006-01"))
	// for category, group := range groupedByCategory {
	// 	sum := Sum(group, func(i InvoiceItem) float64 {
	// 		return i.InstallmentValue
	// 	})

	// 	sumDecimal := decimal.NewFromFloat(sum).Round(2)
	// 	fmt.Printf("Total gasto da categoria %s: R$ %s\n", category, sumDecimal)
	// 	row := NewItemsExportRow(category, sumDecimal.String(), "", "")
	// 	if err := categorySheet.Write(ctx, row); err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	for _, item := range group {
	// 		fmt.Printf("  descrição: %s\n", item.Description)
	// 		row := NewItemsExportRow("", "", item.Description, decimal.NewFromFloat(item.InstallmentValue).Round(2).String())
	// 		if err := categorySheet.Write(ctx, row); err != nil {
	// 			log.Fatal(err)
	// 		}
	// 	}
	// }

	err = xls.SaveAs(ctx, fmt.Sprintf("./files/orcamento-domestico-%s.xlsx", uuid.New().String()))
	if err != nil {
		log.Fatal(err)
	}

	// percentage calculation
	value1, err := decimal.NewFromString("2400")
	if err != nil {
		panic(err)
	}

	value2, err := decimal.NewFromString("4350")
	if err != nil {
		panic(err)
	}

	hundred, err := decimal.NewFromString("100")
	if err != nil {
		panic(err)
	}

	percentage := value1.Div(value2).Mul(hundred)
	fmt.Println(percentage)
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

// func regraDeTres(valor1 float64, valor2 float64, valor3 float64) float64 {
// 	return (valor2 * valor3) / valor1
// }

// func main() {
// 	fmt.Println(regraDeTres(2, 4, 6)) // Deve imprimir 12

// 	// Se 5 livros custam 125 reais, quanto custam 8 livros?
// 	fmt.Println(regraDeTres(5, 125, 8)) // Deve imprimir 200

// 	// Se 30 metros de tecido custam 450 reais, quanto custam 50 metros?
// 	fmt.Println(regraDeTres(30, 450, 50)) // Deve imprimir 750

// 	// Se uma viagem de 120 km é feita em 2 horas, quanto tempo levaria para fazer uma viagem de 200 km?
// 	fmt.Println(regraDeTres(120, 2, 200)) // Deve imprimir 3.33
// }

// func main() {
// 	value, err := decimal.NewFromString("17175")
// 	if err != nil {
// 		panic(err)
// 	}

// 	percentage, err := decimal.NewFromString("6")
// 	if err != nil {
// 		panic(err)
// 	}

// 	hundred, err := decimal.NewFromString("100")
// 	if err != nil {
// 		panic(err)
// 	}

// 	percentageValue := value.Mul(percentage).Div(hundred)
// 	newValue := value.Sub(percentageValue)
// 	fmt.Println(newValue)
// }

// func main() {
//     value, err := decimal.NewFromString("50")
//     if err != nil {
//         panic(err)
//     }

//     percentage, err := decimal.NewFromString("1")
//     if err != nil {
//         panic(err)
//     }

//     hundred, err := decimal.NewFromString("100")
//     if err != nil {
//         panic(err)
//     }

//     percentageValue := value.Mul(percentage).Div(hundred)
//     newValue := value.Add(percentageValue)
//     fmt.Println(newValue)
// }

// func main() {
// 	value, err := decimal.NewFromString("15000")
// 	if err != nil {
// 		panic(err)
// 	}

// 	percentage, err := decimal.NewFromString("30")
// 	if err != nil {
// 		panic(err)
// 	}

// 	hundred, err := decimal.NewFromString("100")
// 	if err != nil {
// 		panic(err)
// 	}

// 	percentageValue := value.Mul(percentage).Div(hundred)
// 	fmt.Println(percentageValue)
// 	newValue := value.Add(percentageValue)
// 	fmt.Println(newValue)
// }
