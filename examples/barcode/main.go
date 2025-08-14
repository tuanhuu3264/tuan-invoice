package main

import (
	"fmt"

	generator "github.com/tuanhuu3264/tuan-invoice"
)

func main() {
	doc, _ := generator.New(generator.Invoice, &generator.Options{
		TextTypeInvoice:   "Hóa Đơn",
		AutoPrint:         true,
		CurrencySymbol:    "VND",
		CurrencyThousand:  ".",
		CurrencyDecimal:   ",",
		CurrencyPrecision: 0,
		Format:            "%v %s",
		FormatNegative:    "- %v %s",
		FormatZero:        "0 %s",
		BarCode:           "1234567890",
		TextInvoiceTitle:  "Mã Vận Đơn",

		TextDateTitle:          "Ngày",
		TextRefTitle:           "Mã đơn hàng",
		TextVersionTitle:       "Phiên bản",
		TextTypeQuotation:      "Trích dẫn",
		TextPaymentTermTitle:   "Ngày đặt hàng",
		TextItemsNameTitle:     "Tên sản phẩm",
		TextItemsUnitCostTitle: "Đơn giá",
		TextItemsQuantityTitle: "SL",
		TextItemsDiscountTitle: "Giảm giá",
		TextItemsTaxTitle:      "Thuế",
		TextItemsTotalHTTitle:  "Tổng tiền",
		TextItemsTotalTTCTitle: "Thành tiền",
		TextTotalDiscounted:    "Tổng giảm giá",
		TextTypeDeliveryNote:   "Phiếu giao hàng",
		TextTotalTax:           "Tổng Thuế",
		TextTotalTotal:         "Tổng cộng",
		TextTotalWithTax:       "Thành Tiền",
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
