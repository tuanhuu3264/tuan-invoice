package main

import (
	"fmt"
	"log"

	generator "github.com/tuanhuu3264/tuan-invoice"
)

func main() {
	// Tạo MultiDocument với options chung
	multiDoc := generator.NewMultiDocument(&generator.Options{
		TextTypeInvoice: "HÓA ĐƠN",
		TextRefTitle:    "Mã số",
		AutoPrint:       false,
		BaseTextColor:   []int{6, 63, 156},
		GreyTextColor:   []int{161, 96, 149},
		GreyBgColor:     []int{171, 240, 129},
		DarkBgColor:     []int{176, 12, 20},
		Font:            "Helvetica",
		CurrencySymbol:  "₫ ",
	})
	multiDoc.SetHeader(&generator.HeaderFooter{
		Text: "HÓA ĐƠN",
	})
	multiDoc.SetFooter(&generator.HeaderFooter{
		Text: "HÓA ĐƠN",
	})

	// Tạo đơn thứ nhất
	doc1, _ := generator.New(generator.Invoice, &generator.Options{
		TextTypeInvoice: "HÓA ĐƠN",
		TextRefTitle:    "Mã số",
		CurrencySymbol:  "₫ ",
	})

	doc1.SetRef("INV-2024-001")
	doc1.SetVersion("1.0")
	doc1.SetDate("15/01/2024")
	doc1.SetPaymentTerm("15/02/2024")
	doc1.SetDescription("Hóa đơn dịch vụ tháng 1/2024")
	doc1.SetNotes("Thanh toán qua chuyển khoản ngân hàng")

	doc1.SetCompany(&generator.Contact{
		Name: "CÔNG TY ABC",
		Address: &generator.Address{
			Address:    "123 Đường ABC",
			Address2:   "Tầng 5, Tòa nhà XYZ",
			PostalCode: "70000",
			City:       "TP. Hồ Chí Minh",
			Country:    "Việt Nam",
		},
	})

	doc1.SetCustomer(&generator.Contact{
		Name: "KHÁCH HÀNG A",
		Address: &generator.Address{
			Address:    "456 Đường DEF",
			PostalCode: "10000",
			City:       "Hà Nội",
			Country:    "Việt Nam",
		},
	})

	// Thêm items cho đơn thứ nhất
	doc1.AppendItem(&generator.Item{
		Name:        "Dịch vụ tư vấn",
		Description: "Tư vấn thiết kế website",
		UnitCost:    "5000000",
		Quantity:    "1",
		Tax: &generator.Tax{
			Percent: "10",
		},
	})

	// Tạo đơn thứ hai
	doc2, _ := generator.New(generator.Invoice, &generator.Options{
		TextTypeInvoice: "HÓA ĐƠN",
		TextRefTitle:    "Mã số",
		CurrencySymbol:  "₫ ",
	})

	doc2.SetRef("INV-2024-002")
	doc2.SetVersion("1.0")
	doc2.SetDate("20/01/2024")
	doc2.SetPaymentTerm("20/02/2024")
	doc2.SetDescription("Hóa đơn dịch vụ tháng 1/2024")
	doc2.SetNotes("Thanh toán tiền mặt")

	doc2.SetCompany(&generator.Contact{
		Name: "CÔNG TY ABC",
		Address: &generator.Address{
			Address:    "123 Đường ABC",
			Address2:   "Tầng 5, Tòa nhà XYZ",
			PostalCode: "70000",
			City:       "TP. Hồ Chí Minh",
			Country:    "Việt Nam",
		},
	})

	doc2.SetCustomer(&generator.Contact{
		Name: "KHÁCH HÀNG B",
		Address: &generator.Address{
			Address:    "789 Đường GHI",
			PostalCode: "50000",
			City:       "Đà Nẵng",
			Country:    "Việt Nam",
		},
	})

	// Thêm items cho đơn thứ hai
	doc2.AppendItem(&generator.Item{
		Name:        "Dịch vụ bảo trì",
		Description: "Bảo trì hệ thống tháng 1",
		UnitCost:    "3000000",
		Quantity:    "1",
		Tax: &generator.Tax{
			Percent: "10",
		},
	})

	// Thêm tất cả đơn vào MultiDocument
	multiDoc.AddDocument(doc1)
	multiDoc.AddDocument(doc2)

	// Tạo PDF
	pdf, err := multiDoc.Build()
	if err != nil {
		log.Fatal(err)
	}

	// Lưu file
	err = pdf.OutputFileAndClose("basic_multi_document_example.pdf")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Đã tạo thành công file PDF với nhiều đơn: basic_multi_document_example.pdf")
	fmt.Println("File này chứa:")
	fmt.Println("- Hóa đơn INV-2024-001 (trang 1)")
	fmt.Println("- Hóa đơn INV-2024-002 (trang 2)")
}
