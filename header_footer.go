package generator

import (
	"github.com/creasty/defaults"
	"github.com/go-pdf/fpdf"
)

// HeaderFooter define header or footer informations on document
type HeaderFooter struct {
	UseCustomFunc bool    `json:"-"`
	Text          string  `json:"text,omitempty"`
	FontSize      float64 `json:"font_size,omitempty" default:"10"`
	Pagination    bool    `json:"pagination,omitempty"`
}

type fnc func()

// ApplyFunc allow user to apply custom func
func (hf *HeaderFooter) ApplyFunc(pdf *fpdf.Fpdf, fn fnc) {
	pdf.SetHeaderFunc(fn)
}

// applyHeader apply header to document
func (hf *HeaderFooter) applyHeader(doc *Document) error {
	if err := defaults.Set(hf); err != nil {
		return err
	}

	if !hf.UseCustomFunc {
		doc.pdf.SetHeaderFunc(func() {
			currentY := doc.pdf.GetY()
			currentX := doc.pdf.GetX()

			doc.pdf.SetTopMargin(HeaderMarginTop)
			doc.pdf.SetY(HeaderMarginTop)

			doc.pdf.SetLeftMargin(BaseMargin)
			doc.pdf.SetRightMargin(BaseMargin)

			// Parse Text as html (simple)
			doc.pdf.SetFont(doc.Options.Font, "", hf.FontSize)
			_, lineHt := doc.pdf.GetFontSize()
			html := doc.pdf.HTMLBasicNew()
			html.Write(lineHt, doc.encodeString(hf.Text))

			doc.pdf.SetY(currentY)
			doc.pdf.SetX(currentX)
			doc.pdf.SetMargins(BaseMargin, BaseMarginTop, BaseMargin)
		})
	}

	return nil
}

// applyFooter apply footer to document
func (hf *HeaderFooter) applyFooter(doc *Document) error {
	if err := defaults.Set(hf); err != nil {
		return err
	}

	if !hf.UseCustomFunc {
		doc.pdf.SetFooterFunc(func() {
			currentY := doc.pdf.GetY()
			currentX := doc.pdf.GetX()

			doc.pdf.SetTopMargin(HeaderMarginTop)
			doc.pdf.SetY(287 - HeaderMarginTop)

			// Parse Text as html (simple)
			doc.pdf.SetFont(doc.Options.Font, "", hf.FontSize)
			_, lineHt := doc.pdf.GetFontSize()
			html := doc.pdf.HTMLBasicNew()
			html.Write(lineHt, doc.encodeString(hf.Text))

			doc.pdf.SetY(currentY)
			doc.pdf.SetX(currentX)
			doc.pdf.SetMargins(BaseMargin, BaseMarginTop, BaseMargin)
		})
	}

	return nil
}
