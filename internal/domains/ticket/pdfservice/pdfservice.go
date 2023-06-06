package pdfgenerator

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

func GeneratePDF(sessionId, userId, seatNumber int, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Ticket Details")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Session ID: %d", sessionId))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("User ID: %d", userId))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Seat Number: %d", seatNumber))

	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return err
	}

	return nil
}
