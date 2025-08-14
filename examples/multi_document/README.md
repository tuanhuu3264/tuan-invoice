# Multi-Document PDF Generator

Ví dụ này hướng dẫn cách tạo một file PDF chứa nhiều đơn (invoices, quotations, delivery notes) trong cùng một file.

## Tính năng

- Tạo nhiều đơn khác nhau trong cùng một file PDF
- Mỗi đơn sẽ được đặt trên một trang riêng biệt
- Header và footer chung cho tất cả các đơn
- Hỗ trợ các loại đơn: Invoice, Quotation, Delivery Note
- Tự động phân trang khi nội dung vượt quá giới hạn

## Cách sử dụng

### 1. Tạo MultiDocument

```go
multiDoc := generator.NewMultiDocument(&generator.Options{
    TextTypeInvoice: "HÓA ĐƠN",
    TextRefTitle:    "Mã số",
    CurrencySymbol:  "₫ ",
    // ... các options khác
})
```

### 2. Thiết lập Header và Footer chung

```go
multiDoc.SetHeader(&generator.HeaderFooter{
    Text:       "<center>CÔNG TY ABC - HỆ THỐNG QUẢN LÝ HÓA ĐƠN</center>",
    Pagination: true,
})

multiDoc.SetFooter(&generator.HeaderFooter{
    Text:       "<center>Cảm ơn quý khách đã sử dụng dịch vụ của chúng tôi!</center>",
    Pagination: true,
})
```

### 3. Tạo các đơn riêng lẻ

```go
// Tạo đơn thứ nhất
doc1, _ := generator.New(generator.Invoice, &generator.Options{
    TextTypeInvoice: "HÓA ĐƠN",
    CurrencySymbol:  "₫ ",
})

doc1.SetRef("INV-2024-001")
doc1.SetCompany(&generator.Contact{...})
doc1.SetCustomer(&generator.Contact{...})
doc1.AppendItem(&generator.Item{...})

// Tạo đơn thứ hai
doc2, _ := generator.New(generator.Quotation, &generator.Options{
    TextTypeQuotation: "BÁO GIÁ",
    CurrencySymbol:    "₫ ",
})

doc2.SetRef("QUO-2024-001")
// ... thiết lập các thông tin khác
```

### 4. Thêm đơn vào MultiDocument

```go
multiDoc.AddDocument(doc1)
multiDoc.AddDocument(doc2)
multiDoc.AddDocument(doc3)
```

### 5. Tạo PDF

```go
pdf, err := multiDoc.Build()
if err != nil {
    log.Fatal(err)
}

err = pdf.OutputFileAndClose("multi_document_example.pdf")
if err != nil {
    log.Fatal(err)
}
```

## Chạy ví dụ

```bash
cd examples/multi_document
go run main.go
```

Kết quả sẽ tạo file `multi_document_example.pdf` chứa:
- Trang 1: Hóa đơn INV-2024-001
- Trang 2: Hóa đơn INV-2024-002  
- Trang 3: Báo giá QUO-2024-001

## Lưu ý

- Mỗi đơn sẽ được đặt trên một trang riêng biệt
- Header và footer sẽ được áp dụng cho tất cả các trang
- Nếu một đơn có nội dung dài, nó sẽ tự động tạo thêm trang mới
- Tất cả các đơn sẽ sử dụng cùng một cấu hình styling (colors, fonts, etc.)
