# Barcode Feature Example

This example demonstrates how to add barcode functionality to your invoices, delivery notes, and quotations.

## Features Added

- **BarCode Field**: New field in Document struct to store barcode content
- **SetBarCode()**: Setter method to set barcode content
- **Code 128 Support**: Generates Code 128 barcodes (most common format)
- **Automatic Positioning**: Barcode is positioned at the bottom right of the document
- **Text Display**: Shows the barcode text below the barcode image

## Usage

```go
// Create a new document
doc, _ := generator.New(generator.Invoice, &generator.Options{
    TextTypeInvoice: "FACTURE",
    AutoPrint:       true,
})

// Set barcode content (Code 128 format)
doc.SetBarCode("INV-2024-001")

// Build the document
pdf, err := doc.Build()
if err != nil {
    panic(err)
}

// Save to file
pdf.OutputFileAndClose("invoice_with_barcode.pdf")
```

## Barcode Format

The barcode feature supports **Code 128** format, which can encode:
- Numbers (0-9)
- Uppercase letters (A-Z)
- Lowercase letters (a-z)
- Special characters

## Dependencies

The barcode feature requires the following dependency:
```
github.com/boombuler/barcode v1.0.1
```

## Running the Example

1. Make sure you have the barcode dependency installed:
   ```bash
   go mod tidy
   ```

2. Run the example:
   ```bash
   cd examples/barcode
   go run main.go
   ```

3. Check the generated PDF file: `invoice_with_barcode.pdf`

## Customization

You can customize the barcode by modifying the `generateBarcode()` function in `build.go`:

- **Size**: Change the scale parameters (currently 200x50)
- **Position**: Modify the x, y coordinates in `appendBarcode()`
- **Format**: Support other barcode types by importing different packages

## Notes

- If barcode generation fails, the document will still be generated without the barcode
- The barcode is positioned at the bottom right of the document
- The barcode text is displayed below the barcode image for easy reading
