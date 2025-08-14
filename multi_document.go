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
	for i, doc := range md.Docs {
		// Add new page for each document (except the first one)
		if i > 0 {
			md.pdf.AddPage()
		} else {
			md.pdf.AddPage()
		}

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

// buildDocument builds a single document within the multi-document
func (md *MultiDocument) buildDocument(doc *Document) error {
	// Validate document data
	if err := doc.Validate(); err != nil {
		return err
	}

	// Set position to top of page
	md.pdf.SetXY(10, BaseMarginTop)

	// Load font
	md.pdf.SetFont(md.Options.Font, "", 12)

	// Append document title
	md.appendTitle(doc)
	md.appendBarcode(doc)

	// Append document metas (ref & version)
	md.appendMetas(doc)

	// Append company contact to doc
	companyBottom := doc.Company.appendCompanyContactToDoc(doc)

	// Append customer contact to doc
	customerBottom := doc.Customer.appendCustomerContactToDoc(doc)

	if customerBottom > companyBottom {
		md.pdf.SetXY(10, customerBottom)
	} else {
		md.pdf.SetXY(10, companyBottom)
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

	// Append notes
	md.appendNotes(doc)

	// Append total
	md.appendTotal(doc)

	// Append payment term
	md.appendPaymentTerm(doc)

	return nil
}

// appendTitle appends the title for a specific document
func (md *MultiDocument) appendTitle(doc *Document) {
	title := doc.typeAsString()

	// Set x y
	md.pdf.SetXY(120, BaseMarginTop)

	// Draw rect with safe color
	md.pdf.SetFillColor(md.Options.DarkBgColor[0], md.Options.DarkBgColor[1], md.Options.DarkBgColor[2])
	md.pdf.Rect(120, BaseMarginTop, 80, 10, "F")

	// Draw text
	md.pdf.SetFont(md.Options.Font, "", 14)
	md.pdf.CellFormat(80, 10, doc.encodeString(title), "0", 0, "C", false, 0, "")
}

// appendBarcode appends barcode for a specific document
func (md *MultiDocument) appendBarcode(doc *Document) {
	if len(doc.BarCode) > 0 {
		// Generate barcode
		code, err := code128.Encode(doc.BarCode)
		if err != nil {
			// Skip barcode if encoding fails
			return
		}

		scaledCode, err := barcode.Scale(code, 80, 20)
		if err != nil {
			// Skip barcode if scaling fails
			return
		}

		// Convert to image
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, scaledCode, nil); err != nil {
			// Skip barcode if encoding fails
			return
		}

		// Add to PDF
		md.pdf.SetXY(120, BaseMarginTop+11)
		md.pdf.Image("", 120, BaseMarginTop+11, 80, 20, false, "JPEG", 0, "")
	}
}

// appendMetas appends metadata for a specific document
func (md *MultiDocument) appendMetas(doc *Document) {
	// Append ref
	refString := fmt.Sprintf("%s: %s", md.Options.TextRefTitle, doc.Ref)

	md.pdf.SetXY(120, BaseMarginTop+11)
	md.pdf.SetFont(md.Options.Font, "", 8)
	md.pdf.CellFormat(80, 4, doc.encodeString(refString), "0", 0, "R", false, 0, "")

	// Append version
	if len(doc.Version) > 0 {
		versionString := fmt.Sprintf("%s: %s", md.Options.TextVersionTitle, doc.Version)
		md.pdf.SetXY(120, BaseMarginTop+15)
		md.pdf.SetFont(md.Options.Font, "", 8)
		md.pdf.CellFormat(80, 4, doc.encodeString(versionString), "0", 0, "R", false, 0, "")
	}

	// Append date
	date := time.Now().Format("02/01/2006")
	if len(doc.Date) > 0 {
		date = doc.Date
	}
	dateString := fmt.Sprintf("%s: %s", md.Options.TextDateTitle, date)
	md.pdf.SetXY(120, BaseMarginTop+19)
	md.pdf.SetFont(md.Options.Font, "", 8)
	md.pdf.CellFormat(80, 4, doc.encodeString(dateString), "0", 0, "R", false, 0, "")
}

// appendDescription appends description for a specific document
func (md *MultiDocument) appendDescription(doc *Document) {
	if len(doc.Description) > 0 {
		md.pdf.SetY(md.pdf.GetY() + 10)
		md.pdf.SetFont(md.Options.Font, "", 10)
		md.pdf.MultiCell(190, 5, doc.encodeString(doc.Description), "B", "L", false)
	}
}

// appendItems appends items for a specific document
func (md *MultiDocument) appendItems(doc *Document) {
	if len(doc.Items) > 0 {
		// Draw table titles
		md.drawsTableTitles(doc)

		// Append each item
		for _, item := range doc.Items {
			item.appendColTo(doc.Options, doc)
		}
	}
}

// drawsTableTitles draws table titles for a specific document
func (md *MultiDocument) drawsTableTitles(doc *Document) {
	// Draw table titles
	md.pdf.SetX(10)
	md.pdf.SetY(md.pdf.GetY() + 5)
	md.pdf.SetFont(md.Options.BoldFont, "B", 8)

	// Draw rect with safe color
	md.pdf.SetFillColor(md.Options.GreyBgColor[0], md.Options.GreyBgColor[1], md.Options.GreyBgColor[2])
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
	md.pdf.SetX(ItemColTaxOffset)
	md.pdf.CellFormat(
		ItemColDiscountOffset-ItemColTaxOffset,
		6,
		doc.encodeString(md.Options.TextItemsTaxTitle),
		"0",
		0,
		"",
		false,
		0,
		"",
	)

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

// appendNotes appends notes for a specific document
func (md *MultiDocument) appendNotes(doc *Document) {
	if len(doc.Notes) > 0 {
		md.pdf.SetY(md.pdf.GetY() + 10)
		md.pdf.SetFont(md.Options.Font, "", 10)
		md.pdf.MultiCell(190, 5, doc.encodeString(doc.Notes), "0", "L", false)
	}
}

// appendTotal appends totals for a specific document
func (md *MultiDocument) appendTotal(doc *Document) {
	md.pdf.SetY(md.pdf.GetY() + 10)
	md.pdf.SetFont(md.Options.Font, "", LargeTextFontSize)

	md.pdf.SetTextColor(
		md.Options.BaseTextColor[0],
		md.Options.BaseTextColor[1],
		md.Options.BaseTextColor[2],
	)

	// Draw TOTAL HT title
	md.pdf.SetX(120)
	md.pdf.SetFillColor(md.Options.DarkBgColor[0], md.Options.DarkBgColor[1], md.Options.DarkBgColor[2])
	md.pdf.Rect(120, md.pdf.GetY(), 40, 10, "F")
	md.pdf.CellFormat(38, 10, doc.encodeString(md.Options.TextTotalTotal), "0", 0, "R", false, 0, "")

	// Draw TOTAL HT amount
	md.pdf.SetX(162)
	md.pdf.SetFillColor(md.Options.GreyBgColor[0], md.Options.GreyBgColor[1], md.Options.GreyBgColor[2])
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
		// Draw discount title
		md.pdf.SetX(120)
		md.pdf.SetFillColor(md.Options.DarkBgColor[0], md.Options.DarkBgColor[1], md.Options.DarkBgColor[2])
		md.pdf.Rect(120, md.pdf.GetY(), 40, 15, "F")
		md.pdf.CellFormat(38, 15, doc.encodeString(md.Options.TextTotalDiscounted), "0", 0, "R", false, 0, "")

		// Draw discount amount
		md.pdf.SetX(162)
		md.pdf.SetFillColor(md.Options.GreyBgColor[0], md.Options.GreyBgColor[1], md.Options.GreyBgColor[2])
		md.pdf.Rect(160, md.pdf.GetY(), 40, 15, "F")
		md.pdf.CellFormat(
			40,
			15,
			doc.encodeString(doc.ac.FormatMoneyDecimal(doc.TotalWithoutTax())),
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
	md.pdf.SetFillColor(md.Options.DarkBgColor[0], md.Options.DarkBgColor[1], md.Options.DarkBgColor[2])
	md.pdf.Rect(120, md.pdf.GetY(), 40, 10, "F")
	md.pdf.CellFormat(38, 10, doc.encodeString(md.Options.TextTotalTax), "0", 0, "R", false, 0, "")

	// Draw tax amount
	md.pdf.SetX(162)
	md.pdf.SetFillColor(md.Options.GreyBgColor[0], md.Options.GreyBgColor[1], md.Options.GreyBgColor[2])
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
	md.pdf.SetFillColor(md.Options.DarkBgColor[0], md.Options.DarkBgColor[1], md.Options.DarkBgColor[2])
	md.pdf.Rect(120, md.pdf.GetY(), 40, 10, "F")
	md.pdf.CellFormat(38, 10, doc.encodeString(md.Options.TextTotalWithTax), "0", 0, "R", false, 0, "")

	// Draw total with tax amount
	md.pdf.SetX(162)
	md.pdf.SetFillColor(md.Options.GreyBgColor[0], md.Options.GreyBgColor[1], md.Options.GreyBgColor[2])
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

// appendPaymentTerm appends payment term for a specific document
func (md *MultiDocument) appendPaymentTerm(doc *Document) {
	if len(doc.PaymentTerm) > 0 {
		md.pdf.SetY(md.pdf.GetY() + 10)
		md.pdf.SetFont(md.Options.Font, "", 10)
		md.pdf.CellFormat(190, 10, doc.encodeString(doc.PaymentTerm), "0", 0, "R", false, 0, "")
	}
}
