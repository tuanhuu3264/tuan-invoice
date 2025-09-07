package generator

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
	"image"

	"github.com/go-pdf/fpdf"
)

// Contact contact a company informations
type Contact struct {
	Name           string   `json:"name,omitempty" validate:"required,min=1,max=256"`
	Logo           []byte   `json:"logo,omitempty"` // Logo byte array
	Address        *Address `json:"address,omitempty"`
	Phone          string   `json:"phone,omitempty"`
	AddtionnalInfo []string `json:"additional_info,omitempty"`
}

// appendContactTODoc append the contact to the document
func (c *Contact) appendContactTODoc(
	x float64,
	y float64,
	fill bool,
	logoAlign string,
	doc *Document,
) float64 {
	doc.pdf.SetXY(x, y)

	// Logo
	if c.Logo != nil {
		// Create filename
		fileName := b64.StdEncoding.EncodeToString([]byte(c.Name))

		// Create reader from logo bytes
		ioReader := bytes.NewReader(c.Logo)

		// Get image format
		_, format, _ := image.DecodeConfig(bytes.NewReader(c.Logo))

		// Register image in pdf
		imageInfo := doc.pdf.RegisterImageOptionsReader(fileName, fpdf.ImageOptions{
			ImageType: format,
		}, ioReader)

		if imageInfo != nil {
			var imageOpt fpdf.ImageOptions
			imageOpt.ImageType = format
			doc.pdf.ImageOptions(fileName, x, y-30, 20, 30, false, imageOpt, 0, "")
		}
	}

	// Name
	if fill {
		doc.pdf.SetFillColor(
			doc.Options.GreyBgColor[0],
			doc.Options.GreyBgColor[1],
			doc.Options.GreyBgColor[2],
		)
	} else {
		doc.pdf.SetFillColor(255, 255, 255)
	}

	// Reset x
	doc.pdf.SetX(x)

	// Calculate total height for unified block
	var totalHeight float64 = 10 // Name height

	if c.Phone != "" {
		totalHeight += 5 // Phone height
	}

	if c.Address != nil {
		addrHeight := 17
		if len(c.Address.Address2) > 0 {
			addrHeight += 5
		}
		if len(c.Address.Country) == 0 {
			addrHeight -= 5
		}
		totalHeight += float64(addrHeight)
	}

	if c.AddtionnalInfo != nil {
		totalHeight += float64(len(c.AddtionnalInfo))*3 + 2 // Additional info height
	}

	// Create unified background rectangle for all contact info
	doc.pdf.Rect(x, doc.pdf.GetY(), 80, totalHeight, "F")

	// Set name - match Title Invoice styling
	doc.pdf.SetFont(doc.Options.Font, "B", 10)
	doc.pdf.CellFormat(80, 10, doc.encodeString(c.Name), "0", 0, "L", false, 0, "")

	if c.Phone != "" {
		doc.pdf.SetXY(x, doc.pdf.GetY()+10)
		doc.pdf.SetFont(doc.Options.Font, "", 10)
		doc.pdf.CellFormat(80, 5, doc.encodeString(fmt.Sprintf("%s: %s", doc.Options.TextPhoneTitle, c.Phone)), "0", 0, "L", false, 0, "")
	}

	if c.Address != nil {
		// Set address - match Title Invoice width
		doc.pdf.SetFont(doc.Options.Font, "", 10)
		doc.pdf.SetXY(x, doc.pdf.GetY()+5)
		doc.pdf.MultiCell(80, 5, doc.encodeString(c.Address.ToString()), "0", "L", false)
	}

	// Addtionnal info
	if c.AddtionnalInfo != nil {
		doc.pdf.SetXY(x, doc.pdf.GetY())
		doc.pdf.SetFontSize(SmallTextFontSize)
		doc.pdf.SetXY(x, doc.pdf.GetY()+2)

		for _, line := range c.AddtionnalInfo {
			doc.pdf.SetXY(x, doc.pdf.GetY())
			doc.pdf.MultiCell(80, 3, doc.encodeString(line), "0", "L", false)
		}

		doc.pdf.SetXY(x, doc.pdf.GetY())
		doc.pdf.SetFontSize(BaseTextFontSize)
	}

	return doc.pdf.GetY()
}

// appendCompanyContactToDoc append the company contact to the document
func (c *Contact) appendCompanyContactToDoc(doc *Document) float64 {
	// Always start at the same Y position regardless of logo
	return c.appendContactTODoc(10, BaseMarginTop+28, true, "L", doc)
}

// appendCustomerContactToDoc append the customer contact to the document
func (c *Contact) appendCustomerContactToDoc(doc *Document) float64 {
	// Always start at the same Y position regardless of logo
	return c.appendContactTODoc(120, BaseMarginTop+28, true, "R", doc)
}
