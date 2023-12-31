package pdfgenerator

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/ticket/service"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"io"
	"log"
)

const (
	textWidth            = 40
	textHeight           = 10
	textSize             = 12
	lineBreakAfterHeader = 12
	lineBreak            = 8
)

type Generator struct{}

func (p Generator) GenerateTicket(t service.Ticket, w io.Writer) error {
	pdf := gofpdf.New("P", "mm", "A6", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(textWidth, textHeight, "Ticket Details")
	pdf.Ln(lineBreakAfterHeader)

	pdf.SetFont("Arial", "", textSize)
	pdf.Cell(textWidth, textHeight, fmt.Sprintf("Movie: %s", t.MovieName))
	pdf.Ln(lineBreak)
	pdf.Cell(textWidth, textHeight, fmt.Sprintf("Date: %s", t.Date))
	pdf.Ln(lineBreak)
	pdf.Cell(textWidth, textHeight, fmt.Sprintf("Start time: %s", t.StartTime))
	pdf.Ln(lineBreak)
	pdf.Cell(textWidth, textHeight, fmt.Sprintf("Duration: %d hour(s) %d minute(s)", t.Duration/60, t.Duration%60))
	pdf.Ln(lineBreak)
	pdf.Cell(textWidth, textHeight, fmt.Sprintf("Hall: %d", t.HallId))
	pdf.Ln(lineBreak)
	pdf.Cell(textWidth, textHeight, fmt.Sprintf("Seat Number: %d", t.SeatNumber))

	err := pdf.Output(w)
	if err != nil {
		log.Printf("error while generating PDF file: %v", err)
		return err
	}

	return nil
}
