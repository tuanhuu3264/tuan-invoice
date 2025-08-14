package main

import (
	"fmt"

	generator "github.com/angelodlfrtr/go-invoice-generator"
)

func main() {
	doc, _ := generator.New(generator.Invoice, &generator.Options{
		TextTypeInvoice:   "FACTURE",
		TextRefTitle:      "Ref",
		AutoPrint:         true,
		BaseTextColor:     []int{6, 63, 156},
		GreyTextColor:     []int{161, 96, 149},
		GreyBgColor:       []int{171, 240, 129},
		DarkBgColor:       []int{176, 12, 20},
		CurrencyPrecision: 0, // Không có chữ số thập phân - sẽ hiển thị € 100 thay vì € 100.00
		CurrencySymbol:    "€ ",
		CurrencyDecimal:   ".",
		CurrencyThousand:  " ",
	})

	doc.SetHeader(&generator.HeaderFooter{
		Text:       "<center>Invoice with Barcode Example</center>",
		Pagination: true,
	})

	doc.SetFooter(&generator.HeaderFooter{
		Text:       "<center>Thank you for your business!</center>",
		Pagination: true,
	})

	doc.SetRef("INV-2024-001")
	doc.SetVersion("1.0")

	// Set barcode content (Code 128 format)
	doc.SetBarCode("INV-2024-001")

	doc.SetDescription("Sample invoice with barcode")
	doc.SetNotes("This invoice includes a Code 128 barcode for easy scanning and tracking.")

	doc.SetDate("15/01/2024")
	doc.SetPaymentTerm("15/02/2024")

	doc.SetCompany(&generator.Contact{
		Name: "Sample Company Ltd",
		Address: &generator.Address{
			Address:    "123 Business Street",
			Address2:   "Suite 100",
			PostalCode: "12345",
			City:       "Business City",
			Country:    "USA",
		},
	})

	doc.SetCustomer(&generator.Contact{
		Name: "Customer Name",
		Address: &generator.Address{
			Address:    "456 Customer Avenue",
			PostalCode: "67890",
			City:       "Customer City",
			Country:    "USA",
		},
	})

	// Add some sample items with decimal values to see precision effect
	doc.AppendItem(&generator.Item{
		Name:        "Product A",
		Description: "High quality product",
		UnitCost:    "100.75", // Sẽ hiển thị € 101 với precision 0
		Quantity:    "2",
		Tax: &generator.Tax{
			Percent: "10",
		},
	})

	doc.AppendItem(&generator.Item{
		Name:        "Service B",
		Description: "Professional service",
		UnitCost:    "250.25", // Sẽ hiển thị € 250 với precision 0
		Quantity:    "1",
		Tax: &generator.Tax{
			Percent: "10",
		},
		Discount: &generator.Discount{
			Percent: "15",
		},
	})

	pdf, err := doc.Build()
	if err != nil {
		panic(err)
	}

	if err := pdf.OutputFileAndClose("invoice_with_barcode.pdf"); err != nil {
		panic(err)
	}

	fmt.Println("Invoice with barcode generated successfully: invoice_with_barcode.pdf")
}
