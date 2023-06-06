package pdfgenerator

import (
	"fmt"
	"log"

	"github.com/unidoc/unipdf/v3/model"
)

func GeneratePDF(sessionId, userId, seatNumber int, outputPath string) error {
	pdf := model.NewPdfWriter()

	// Создаем страницу PDF-документа.
	page := pdf.P

	text := fmt.Sprintf("Session ID: %d\nUser ID: %d\nSeat Number: %d", sessionId, userId, seatNumber)

	// Добавляем текстовый блок на страницу.
	page.AddText(text, 50, 500, model.NewHelvetica(12))

	// Добавляем страницу в PDF-документ.
	pdf.AddPage(page)

	// Сохраняем PDF-документ в файл.
	err := pdf.WriteToFile("output.pdf")
	if err != nil {
		log.Fatalf("Ошибка при сохранении PDF-файла: %v", err)
	}

	fmt.Println("PDF-файл успешно создан.")
}
