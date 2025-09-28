package generator

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/go-pdf/fpdf"
)

// MultiDocument represents a collection of documents to be generated in a single PDF
type MultiDocument struct {
	pdf     *fpdf.Fpdf
	Options *Options
	Header  *HeaderFooter
	Footer  *HeaderFooter
	Docs    []*Document
}

// NewMultiDocument creates a new multi-document generator
func NewMultiDocument(options *Options) *MultiDocument {
	pdf := fpdf.New("P", "mm", "A4", "")

	return &MultiDocument{
		pdf:     pdf,
		Options: options,
		Docs:    make([]*Document, 0),
	}
}

// AddDocument adds a document to the multi-document collection
func (md *MultiDocument) AddDocument(doc *Document) {
	md.Docs = append(md.Docs, doc)
}

// SetHeader sets the header for all documents
func (md *MultiDocument) SetHeader(header *HeaderFooter) {
	md.Header = header
}

// SetFooter sets the footer for all documents
func (md *MultiDocument) SetFooter(footer *HeaderFooter) {
	md.Footer = footer
}

func (md *MultiDocument) GetPdf() *fpdf.Fpdf {
	return md.pdf
}

// Build generates the PDF with all documents
func (md *MultiDocument) Build() (*fpdf.Fpdf, error) {
	// Process each document
	for _, doc := range md.Docs {
		// Add new page for each document
		md.pdf.AddPage()

		// Set up document-specific settings
		doc.pdf = md.pdf
		doc.Options = md.Options

		// Build the document content
		if err := md.buildDocument(doc); err != nil {
			return nil, err
		}
	}

	// Add auto-print if enabled
	if md.Options.AutoPrint {
		md.pdf.SetJavascript("print(true);")
	}

	return md.pdf, nil
}

// getSafeColor returns a safe color array with default values if the input is too short
func (md *MultiDocument) getSafeColor(color []int, defaultColor []int) []int {
	if len(color) >= 3 {
		return color
	}
	return defaultColor
}

// buildDocument builds a single document within the multi-document
func (md *MultiDocument) buildDocument(doc *Document) error {
	// Validate document data
	if err := doc.Validate(); err != nil {
		return err
	}

	// Set header if exists
	if md.Header != nil {
		if err := md.Header.applyHeader(doc); err != nil {
			return err
		}
	}

	// Set footer if exists
	if md.Footer != nil {
		if err := md.Footer.applyFooter(doc); err != nil {
			return err
		}
	}

	// Set position to top of page
	md.pdf.SetXY(10, BaseMarginTop)

	// Load font
	md.pdf.SetFont(md.Options.Font, "", 15)

	// Append document title
	md.appendTitle(doc)

	// Append document metas (ref & version)
	md.appendMetas(doc)

	// Append company contact to doc
	companyBottom := doc.Company.appendCompanyContactToDoc(doc)

	// Append customer contact to doc
	customerBottom := doc.Customer.appendCustomerContactToDoc(doc)

	// Set position to the bottom of the higher contact section
	if customerBottom > companyBottom {
		md.pdf.SetXY(10, customerBottom+5)
	} else {
		md.pdf.SetXY(10, companyBottom+5)
	}

	// Append description
	md.appendDescription(doc)

	// Append items
	md.appendItems(doc)

	// Check page height and add new page if needed
	offset := md.pdf.GetY() + 30
	if doc.Discount != nil {
		offset += 15
	}
	if offset > MaxPageHeight {
		md.pdf.AddPage()
	}

	// Append barcode parallel to total
	md.appendBarcode(doc)

	// Append notes
	md.appendNotes(doc)

	// Append total
	md.appendTotal(doc)

	// Append payment term
	md.appendPaymentTerm(doc)

	return nil
}

func (md *MultiDocument) appendTitle(doc *Document) {
	title := doc.typeAsString()

	// Set x y
	md.pdf.SetXY(120, BaseMarginTop)

	// Draw rect with safe color
	darkColor := md.getSafeColor(md.Options.DarkBgColor, []int{0, 0, 0})
	md.pdf.SetFillColor(darkColor[0], darkColor[1], darkColor[2])
	md.pdf.Rect(120, BaseMarginTop, 80, 10, "F")

	// Draw text
	md.pdf.SetFont(md.Options.Font, "", 17)
	md.pdf.CellFormat(80, 10, doc.encodeString(title), "0", 0, "C", false, 0, "")
}

// appendMetas to document
func (md *MultiDocument) appendMetas(doc *Document) {
	// Append ref
	refString := fmt.Sprintf("%s: %s", md.Options.TextRefTitle, doc.Ref)

	md.pdf.SetXY(120, BaseMarginTop+11)
	md.pdf.SetFont(md.Options.Font, "", 9)
	md.pdf.CellFormat(80, 4, doc.encodeString(refString), "0", 0, "R", false, 0, "")

	// Append date
	date := time.Now().Format("02/01/2006")
	if len(doc.Date) > 0 {
		date = doc.Date
	}
	dateString := fmt.Sprintf("%s: %s", md.Options.TextDateTitle, date)
	md.pdf.SetXY(120, BaseMarginTop+15)
	md.pdf.SetFont(md.Options.Font, "", 9)
	md.pdf.CellFormat(80, 4, doc.encodeString(dateString), "0", 0, "R", false, 0, "")
}

// appendDescription to document
func (md *MultiDocument) appendDescription(doc *Document) {
	if len(doc.Description) > 0 {
		md.pdf.SetY(md.pdf.GetY() + 5)
		md.pdf.SetFont(md.Options.Font, "", 13)
		md.pdf.MultiCell(190, 5, doc.encodeString(doc.Description), "B", "L", false)
	}
}

// drawsTableTitles in document
func (md *MultiDocument) drawsTableTitles(doc *Document) {
	// Draw table titles
	md.pdf.SetX(10)
	md.pdf.SetY(md.pdf.GetY() + 5)
	md.pdf.SetFont(md.Options.BoldFont, "B", 9)

	// Draw rect with safe color
	greyColor := md.getSafeColor(md.Options.GreyBgColor, []int{240, 240, 240})
	md.pdf.SetFillColor(greyColor[0], greyColor[1], greyColor[2])
	md.pdf.Rect(10, md.pdf.GetY(), 190, 6, "F")

	// Name
	md.pdf.SetX(ItemColNameOffset)
	md.pdf.CellFormat(
		ItemColUnitPriceOffset-ItemColNameOffset,
		6,
		doc.encodeString(md.Options.TextItemsNameTitle),
		"0",
		0,
		"",
		false,
		0,
		"",
	)

	// Unit price
	md.pdf.SetX(ItemColUnitPriceOffset)
	md.pdf.CellFormat(
		ItemColQuantityOffset-ItemColUnitPriceOffset,
		6,
		doc.encodeString(md.Options.TextItemsUnitCostTitle),
		"0",
		0,
		"",
		false,
		0,
		"",
	)

	// Quantity
	md.pdf.SetX(ItemColQuantityOffset)
	md.pdf.CellFormat(
		ItemColTaxOffset-ItemColQuantityOffset,
		6,
		doc.encodeString(md.Options.TextItemsQuantityTitle),
		"0",
		0,
		"",
		false,
		0,
		"",
	)

	// Total HT
	md.pdf.SetX(ItemColTotalHTOffset)
	md.pdf.CellFormat(
		ItemColTaxOffset-ItemColTotalHTOffset,
		6,
		doc.encodeString(md.Options.TextItemsTotalHTTitle),
		"0",
		0,
		"",
		false,
		0,
		"",
	)

	// Tax
	// md.pdf.SetX(ItemColTaxOffset)
	// md.pdf.CellFormat(
	// 	ItemColDiscountOffset-ItemColTaxOffset,
	// 	6,
	// 	doc.encodeString(md.Options.TextItemsTaxTitle),
	// 	"0",
	// 	0,
	// 	"",
	// 	false,
	// 	0,
	// 	"",
	// )

	// Discount
	md.pdf.SetX(ItemColDiscountOffset)
	md.pdf.CellFormat(
		ItemColTotalTTCOffset-ItemColDiscountOffset,
		6,
		doc.encodeString(md.Options.TextItemsDiscountTitle),
		"0",
		0,
		"",
		false,
		0,
		"",
	)

	// TOTAL TTC
	md.pdf.SetX(ItemColTotalTTCOffset)
	md.pdf.CellFormat(
		190-ItemColTotalTTCOffset,
		6,
		doc.encodeString(md.Options.TextItemsTotalTTCTitle),
		"0",
		0,
		"",
		false,
		0,
		"",
	)
}

// appendItems to document
func (md *MultiDocument) appendItems(doc *Document) {
	md.drawsTableTitles(doc)

	md.pdf.SetX(10)
	md.pdf.SetY(md.pdf.GetY() + 8)
	md.pdf.SetFont(md.Options.Font, "", 9)

	for i := 0; i < len(doc.Items); i++ {
		item := doc.Items[i]

		// Check item tax
		if item.Tax == nil {
			item.Tax = doc.DefaultTax
		}

		// Append to pdf
		item.appendColTo(md.Options, doc)

		if md.pdf.GetY() > MaxPageHeight {
			// Add page
			md.pdf.AddPage()
			md.drawsTableTitles(doc)
			md.pdf.SetFont(md.Options.Font, "", 9)
		}

		md.pdf.SetX(10)
		md.pdf.SetY(md.pdf.GetY() + 6)
	}
}

// appendNotes to document
func (md *MultiDocument) appendNotes(doc *Document) {
	if len(doc.Notes) == 0 {
		return
	}

	// Position notes at current Y position
	md.pdf.SetY(md.pdf.GetY() + 40)
	md.pdf.SetFont(md.Options.Font, "", 12)
	md.pdf.SetX(10) // Left side position
	md.pdf.SetRightMargin(100)

	_, lineHt := md.pdf.GetFontSize()
	html := md.pdf.HTMLBasicNew()
	html.Write(lineHt, doc.encodeString(doc.Notes))

	md.pdf.SetRightMargin(BaseMargin)
}

// appendTotal to document
func (md *MultiDocument) appendTotal(doc *Document) {
	md.pdf.SetY(md.pdf.GetY() - 35)
	md.pdf.SetFont(md.Options.Font, "", LargeTextFontSize)
	// Set text color with safe values
	baseTextColor := md.getSafeColor(md.Options.BaseTextColor, []int{35, 35, 35})
	md.pdf.SetTextColor(baseTextColor[0], baseTextColor[1], baseTextColor[2])

	// Draw TOTAL HT title
	md.pdf.SetX(120)
	darkColor := md.getSafeColor(md.Options.DarkBgColor, []int{0, 0, 0})
	md.pdf.SetFillColor(darkColor[0], darkColor[1], darkColor[2])
	md.pdf.Rect(120, md.pdf.GetY(), 40, 10, "F")
	md.pdf.CellFormat(38, 10, doc.encodeString(md.Options.TextTotalTotal), "0", 0, "R", false, 0, "")

	// Draw TOTAL HT amount
	md.pdf.SetX(162)
	greyColor := md.getSafeColor(md.Options.GreyBgColor, []int{240, 240, 240})
	md.pdf.SetFillColor(greyColor[0], greyColor[1], greyColor[2])
	md.pdf.Rect(160, md.pdf.GetY(), 40, 10, "F")
	md.pdf.CellFormat(
		40,
		10,
		doc.encodeString(doc.ac.FormatMoneyDecimal(doc.TotalWithoutTaxAndWithoutDocumentDiscount())),
		"0",
		0,
		"L",
		false,
		0,
		"",
	)

	if doc.Discount != nil {
		baseY := md.pdf.GetY() + 10

		// Draw discounted title
		md.pdf.SetXY(120, baseY)
		darkColor := md.getSafeColor(md.Options.DarkBgColor, []int{0, 0, 0})
		md.pdf.SetFillColor(darkColor[0], darkColor[1], darkColor[2])
		md.pdf.Rect(120, md.pdf.GetY(), 40, 15, "F")

		// title
		md.pdf.CellFormat(38, 7.5, doc.encodeString(md.Options.TextTotalDiscounted), "0", 0, "BR", false, 0, "")

		// description
		md.pdf.SetXY(120, baseY+7.5)
		md.pdf.SetFont(md.Options.Font, "", BaseTextFontSize)
		// Set grey text color with safe values
		greyTextColor := md.getSafeColor(md.Options.GreyTextColor, []int{128, 128, 128})
		md.pdf.SetTextColor(greyTextColor[0], greyTextColor[1], greyTextColor[2])

		var descString bytes.Buffer
		_, discountAmount := doc.Discount.getDiscount()

		md.pdf.CellFormat(38, 7.5, doc.encodeString(descString.String()), "0", 0, "TR", false, 0, "")

		md.pdf.SetFont(md.Options.Font, "", LargeTextFontSize)
		// Set base text color with safe values
		baseTextColor := md.getSafeColor(md.Options.BaseTextColor, []int{35, 35, 35})
		md.pdf.SetTextColor(baseTextColor[0], baseTextColor[1], baseTextColor[2])

		// Draw discount amount
		md.pdf.SetY(baseY)
		md.pdf.SetX(162)
		greyColor := md.getSafeColor(md.Options.GreyBgColor, []int{240, 240, 240})
		md.pdf.SetFillColor(greyColor[0], greyColor[1], greyColor[2])
		md.pdf.Rect(160, md.pdf.GetY(), 40, 15, "F")
		md.pdf.CellFormat(
			40,
			15,
			doc.encodeString(doc.ac.FormatMoneyDecimal(discountAmount)),
			"0",
			0,
			"L",
			false,
			0,
			"",
		)
		md.pdf.SetY(md.pdf.GetY() + 15)
	} else {
		md.pdf.SetY(md.pdf.GetY() + 10)
	}

	// Draw tax title
	md.pdf.SetX(120)
	darkColor = md.getSafeColor(md.Options.DarkBgColor, []int{0, 0, 0})
	md.pdf.SetFillColor(darkColor[0], darkColor[1], darkColor[2])
	md.pdf.Rect(120, md.pdf.GetY(), 40, 10, "F")
	md.pdf.CellFormat(38, 10, doc.encodeString(md.Options.TextTotalTax), "0", 0, "R", false, 0, "")

	// Draw tax amount
	md.pdf.SetX(162)
	greyColor = md.getSafeColor(md.Options.GreyBgColor, []int{240, 240, 240})
	md.pdf.SetFillColor(greyColor[0], greyColor[1], greyColor[2])
	md.pdf.Rect(160, md.pdf.GetY(), 40, 10, "F")
	md.pdf.CellFormat(
		40,
		10,
		doc.encodeString(doc.ac.FormatMoneyDecimal(doc.Tax())),
		"0",
		0,
		"L",
		false,
		0,
		"",
	)

	// Draw total with tax title
	md.pdf.SetY(md.pdf.GetY() + 10)
	md.pdf.SetX(120)
	darkColor = md.getSafeColor(md.Options.DarkBgColor, []int{0, 0, 0})
	md.pdf.SetFillColor(darkColor[0], darkColor[1], darkColor[2])
	md.pdf.Rect(120, md.pdf.GetY(), 40, 10, "F")
	md.pdf.CellFormat(38, 10, doc.encodeString(md.Options.TextTotalWithTax), "0", 0, "R", false, 0, "")

	// Draw total with tax amount
	md.pdf.SetX(162)
	greyColor = md.getSafeColor(md.Options.GreyBgColor, []int{240, 240, 240})
	md.pdf.SetFillColor(greyColor[0], greyColor[1], greyColor[2])
	md.pdf.Rect(160, md.pdf.GetY(), 40, 10, "F")
	md.pdf.CellFormat(
		40,
		10,
		doc.encodeString(doc.ac.FormatMoneyDecimal(doc.TotalWithTax())),
		"0",
		0,
		"L",
		false,
		0,
		"",
	)
}

// appendPaymentTerm to document
func (md *MultiDocument) appendPaymentTerm(doc *Document) {
	if len(doc.PaymentTerm) > 0 {
		paymentTermString := fmt.Sprintf(
			"%s: %s",
			doc.encodeString(md.Options.TextPaymentTermTitle),
			doc.encodeString(doc.PaymentTerm),
		)
		md.pdf.SetY(md.pdf.GetY() + 15)

		md.pdf.SetX(120)
		md.pdf.SetFont(md.Options.BoldFont, "B", 13)
		md.pdf.CellFormat(80, 4, doc.encodeString(paymentTermString), "0", 0, "R", false, 0, "")
	}
}

// generateBarcode generates a Code 128 barcode image
func (md *MultiDocument) generateBarcode(content string) ([]byte, error) {
	if len(content) == 0 {
		return nil, nil
	}

	// Create Code 128 barcode
	bc, err := code128.Encode(content)
	if err != nil {
		return nil, err
	}

	// Scale the barcode (even larger size)
	scaledBc, err := barcode.Scale(bc, 300, 80)
	if err != nil {
		return nil, err
	}

	// Encode to JPEG
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, scaledBc, &jpeg.Options{Quality: 90})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// appendBarcode to document
func (md *MultiDocument) appendBarcode(doc *Document) {
	if len(doc.BarCode) == 0 {
		return
	}

	// Generate barcode image
	barcodeBytes, err := md.generateBarcode(doc.BarCode)
	if err != nil {
		// If barcode generation fails, just skip it
		return
	}

	// Position barcode on the same row as total section (left side)

	// Create filename for barcode
	fileName := "barcode_" + doc.Ref

	// Create reader from barcode bytes
	ioReader := bytes.NewReader(barcodeBytes)

	// Register image in pdf
	imageInfo := md.pdf.RegisterImageOptionsReader(fileName, fpdf.ImageOptions{
		ImageType: "JPEG",
	}, ioReader)

	if imageInfo != nil {
		// Store current Y position for total section
		currentY := md.pdf.GetY()

		// Position barcode on the left side, same row as total
		x := 10.0
		y := currentY + 10 // Same Y offset as total section

		md.pdf.ImageOptions(fileName, x, y, 0, 20, false, fpdf.ImageOptions{
			ImageType: "PNG",
		}, 0, "")
		// Add barcode text below, perfectly centered within barcode width
		md.pdf.SetY(y + 21)
		md.pdf.SetFont(md.Options.Font, "", 9)

		// Get text width for centering calculation
		textWidth := md.pdf.GetStringWidth(doc.encodeString(doc.BarCode))

		// Use the barcode's actual rendered width (300px scaled to PDF units)
		// The barcode is scaled to 300px width, so we need to account for the PDF scaling
		barcodeWidth := 300.0 * (20.0 / 80.0) // Scale factor: PDF height / original height

		// Calculate perfect center position
		centerX := x + (barcodeWidth-textWidth)/2
		md.pdf.SetX(centerX)

		// Draw the text centered
		md.pdf.CellFormat(textWidth, 4, doc.encodeString(doc.BarCode), "0", 0, "C", false, 0, "")

		// Reset Y position to where total section will start
		md.pdf.SetY(currentY)
	}
}
