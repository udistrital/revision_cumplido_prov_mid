package helper_generar_documento

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"time"
)

func GenerarPdf() {

	now := time.Now()

	formattedDate := now.Format("02/01/2006")
	month := now.Month()
	year := now.Year()
	day := now.Day()
	fmt.Print("Generando documento")
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Times", "", 12)
	pdf.SetMargins(15, 10, 15)
	pdf.AddPage()
	pdf.SetFillColor(240, 240, 240)

	//////Header
	cellX, cellY := 30.0, 10.0
	cellX2, cellY2 := 55.0, 10.0
	cellWidth, cellHeight := 25.0, 25.0
	cellWidth2, cellHeight2 := 55.0, 7.5

	pdf.Rect(cellX, cellY, cellWidth, cellHeight, "")
	logoU := "imgs/image1.png"
	pdf.Image(logoU, cellX+2, cellY, cellWidth-3, cellHeight-5, false, "", 0, "")

	pdf.Rect(cellX+123, cellY, cellWidth+2, cellHeight, "")
	logoSigud := "imgs/image2.png"
	pdf.Image(logoSigud, cellX+127, cellY+8, cellWidth-5, cellHeight-16, false, "", 0, "")

	pdf.Rect(cellX2, cellY2, cellWidth2, cellHeight2, "")
	pdf.SetFont("Times", "B", 10)
	pdf.Text(cellX2+5, cellY2+5, "AUTORIZACION DE GIRO ")

	pdf.SetFont("Times", "", 10)
	pdf.Rect(cellX2, cellY2+7.5, cellWidth2, cellHeight2, "")
	pdf.Text(cellX2+2, cellY2+12, "Macroproceso:  Gestion de Recursos")

	pdf.Rect(cellX2, cellY2+15, cellWidth2, cellHeight2+2.5, "")
	pdf.Text(cellX2+6, cellY2+20, "Proceso: Gestión de Recursos")
	pdf.Text(cellX2+20, cellY2+24, "Financieros")

	pdf.Rect(cellX2+55, cellY2, cellWidth2-12, cellHeight2, "")
	pdf.Text(cellX2+55, cellY2+5, "Codigo: GRF-PR-007-FR-005")

	pdf.SetFont("Times", "", 10)
	pdf.Rect(cellX2+55, cellY2+7.5, cellWidth2-12, cellHeight2, "")
	pdf.Text(cellX2+55, cellY2+13, "Version: 04")

	pdf.Rect(cellX2+55, cellY2+15, cellWidth2-12, cellHeight2+2.5, "")
	pdf.Text(cellX2+55, cellY2+20, "Fecha de Aprobación:")
	pdf.Text(cellX2+55, cellY2+24, formattedDate)

	pdf.Ln(10)
	formattedDate2 := fmt.Sprintf("BOGOTÁ %s %d de %d", month, day, year)
	pdf.Text(cellX, cellY+35, formattedDate2)

	err := pdf.OutputFileAndClose("tabla.pdf")
	if err != nil {
		fmt.Println("Error al guardar el archivo PDF:", err)
	} else {
		fmt.Println("PDF generado exitosamente")

	}

}
