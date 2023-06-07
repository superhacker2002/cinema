package pdfgenerator

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/ticket/service"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"log"
)

type PDFGenerator struct{}

func New() PDFGenerator {
	return PDFGenerator{}
}

func (p PDFGenerator) GenerateTicket(t service.Ticket, outputPath string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A6", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Ticket Details")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Movie: %s", t.MovieName))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Start time: %s", t.StartTime))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Date: %s", t.Date))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Start time: %s", t.StartTime))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Duration: %d hour(s) %d minute(s)", t.Duration/60, t.Duration%60))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Hall: %d", t.HallId))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Seat Number: %d", t.SeatNumber))

	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return outputPath, nil
}
