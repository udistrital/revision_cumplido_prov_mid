package helper_generar_documento

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
	"time"
)

func GenerarPdf(atorizacion *models.DocuementoAutorizacionPago) {

	if atorizacion == nil {
		fmt.Println("Error al generar el documento")
		return

	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	now := time.Now()
	//////Header
	cellX, cellY := 30.0, 10.0
	cellX2, cellY2 := 55.0, 10.0
	cellWidth, cellHeight := 25.0, 25.0
	cellWidth2, cellHeight2 := 55.0, 7.5
	formattedDate := now.Format("02/01/2006")
	month := int(now.Month())
	year := now.Year()
	day := now.Day()
	fmt.Print("Generando documento")
	pdf.SetFont("Times", "", 12)
	pdf.SetMargins(15, 10, 15)

	pdf.SetFillColor(240, 240, 240)

	pdf.SetFooterFunc(func() {
		pdf = footer(pdf, cellX, cellY)
	})
	pdf = footer(pdf, cellX, cellY)
	pdf.SetHeaderFunc(func() {
		pdf = header(pdf, cellX, cellY, cellX2, cellY2, cellWidth, cellHeight, cellWidth2, cellHeight2, formattedDate)
	})
	pdf.AddPage()
	pdf = header(pdf, cellX, cellY, cellX2, cellY2, cellWidth, cellHeight, cellWidth2, cellHeight2, formattedDate)

	pdf = body(pdf, cellX, cellY, month, day, year, atorizacion)

	err := pdf.OutputFileAndClose("tabla.pdf")
	if err != nil {
		fmt.Println("Error al guardar el archivo PDF:", err)
	} else {
		fmt.Println("PDF generado exitosamente")

	}

}

func header(pdf *gofpdf.Fpdf, cellX float64, cellY float64, cellX2 float64, cellY2 float64, cellWidth float64, cellHeight float64, cellWidth2 float64, cellHeight2 float64, formattedDate string) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.Rect(cellX, cellY, cellWidth, cellHeight, "")
	logoU := "imgs/image1.png"
	pdf.Image(logoU, cellX+2, cellY, cellWidth-3, cellHeight-5, false, "", 0, "")

	pdf.Rect(cellX+123, cellY, cellWidth+2, cellHeight, "")
	logoSigud := "imgs/image2.png"
	pdf.Image(logoSigud, cellX+127, cellY+8, cellWidth-5, cellHeight-16, false, "", 0, "")

	pdf.Rect(cellX2, cellY2, cellWidth2, cellHeight2, "")
	pdf.SetFont("Times", "B", 10)
	pdf.Text(cellX2+5, cellY2+5, tr("AUTORIZACIÓN DE GIRO"))

	pdf.SetFont("Times", "", 10)
	pdf.Rect(cellX2, cellY2+7.5, cellWidth2, cellHeight2, "")
	pdf.Text(cellX2+2, cellY2+12, "Macroproceso:  Gestion de Recursos")

	pdf.Rect(cellX2, cellY2+15, cellWidth2, cellHeight2+2.5, "")
	pdf.Text(cellX2+6, cellY2+20, tr("Proceso: Gestión de Recursos"))
	pdf.Text(cellX2+20, cellY2+24, "Financieros")

	pdf.Rect(cellX2+55, cellY2, cellWidth2-12, cellHeight2, "")
	pdf.Text(cellX2+55, cellY2+5, "Codigo: GRF-PR-007-FR-005")

	pdf.SetFont("Times", "", 10)
	pdf.Rect(cellX2+55, cellY2+7.5, cellWidth2-12, cellHeight2, "")
	pdf.Text(cellX2+55, cellY2+13, "Version: 04")

	pdf.Rect(cellX2+55, cellY2+15, cellWidth2-12, cellHeight2+2.5, "")
	pdf.Text(cellX2+55, cellY2+20, tr("Fecha de Aprobación:"))
	pdf.Text(cellX2+55, cellY2+24, "02/08/2024")

	return pdf
}

func body(pdf *gofpdf.Fpdf, cellx float64, cellY float64, month int, day int, year int, autorizacion *models.DocuementoAutorizacionPago) *gofpdf.Fpdf {

	tr := pdf.UnicodeTranslatorFromDescriptor("")
	formattedDate2 := tr(fmt.Sprintf("BOGOTÁ %s %d de %d", obtenerMes(int(month)), day, year))
	pdf.Text(cellx, cellY+35, formattedDate2)

	pdf.SetXY(29, 55)
	pdf.MultiCell(153, 5, tr(fmt.Sprintf("Yo %s  en calidad de Ordenador del Gasto del (los) Rubro(s) __________________________, anexo los documentos detallados en la presente, como soporte a la orden de pago correspondiente.", autorizacion.NombreOrdenador)), "", "", false)
	println(cellY, "y")
	pdf = tablaDocumentos(pdf, cellx, cellY-10, autorizacion)

	return pdf
}

func tablaDocumentos(pdf *gofpdf.Fpdf, cellx float64, cellY float64, autorizacion *models.DocuementoAutorizacionPago) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	col1Width := 120.0
	col2Width := 20.0
	totalWidth := col1Width + col2Width
	pageWidth, _ := pdf.GetPageSize()
	startX := (pageWidth - totalWidth) / 2
	cellY += 75

	pdf.SetXY(startX, cellY)
	pdf.SetFont("Times", "B", 10)
	pdf.CellFormat(col1Width, 7, "DOCUMENTO", "1", 0, "C", false, 0, "")
	pdf.CellFormat(col2Width, 7, "X", "1", 0, "C", false, 0, "")

	pdf.SetFont("Times", "", 10)
	for key, value := range docuementos() {
		cellY += 7
		fmt.Println("yyyy", cellY)
		pdf.SetXY(startX, cellY)
		pdf.CellFormat(col1Width, 7, tr(value), "1", 0, "L", false, 0, "")

		if contains(autorizacion.DocumentosCargados, key) {
			pdf.CellFormat(col2Width, 7, "X", "1", 0, "C", false, 0, "")
		} else {
			pdf.CellFormat(col2Width, 7, "", "1", 0, "C", false, 0, "")
		}
		println(cellY)
		if cellY > 250 {
			pdf.AddPage()
			cellY = 35
		}
	}
	pdf.SetFont("Times", "B", 10)
	cellY += 7
	pdf.SetXY(startX, cellY)
	pdf.SetFont("Times", "B", 10)
	pdf.CellFormat(col1Width, 7, "OTROS DOCUMENTOS (DETALLAR)", "1", 0, "C", false, 0, "")
	pdf.CellFormat(col2Width, 7, "X", "1", 0, "C", false, 0, "")

	println("guia", cellY)
	if cellY > 235 {
		pdf.AddPage()
		cellY = 35
		pdf = body2(pdf, cellx, cellY, 0, 0, 0, autorizacion)

	} else {
		pdf = body2(pdf, cellx, cellY, 0, 0, 0, autorizacion)
	}
	return pdf
}

func body2(pdf *gofpdf.Fpdf, cellx float64, cellY float64, month int, day int, year int, autorizacion *models.DocuementoAutorizacionPago) *gofpdf.Fpdf {
	pdf.SetFont("Times", "", 10)
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	cellY += 10
	pdf.SetXY(29, cellY)
	fmt.Println(autorizacion.DocumentoProveedor)
	pdf.MultiCell(153, 5, tr(fmt.Sprintf(`Autorizo a la Tesorería General a girar a favor de %s con C.C., NIT, TI, OTROS Nº %s para realizar el giro una vez sean deducidos los descuentos de Ley correspondientes. El valor bruto de la presente autorización es de _______________ pesos m/cte. ($____________). `, tr(autorizacion.NombreProveedor), tr(autorizacion.DocumentoProveedor))), "", "", false)
	cellY += 35
	pdf.SetXY(29, cellY)
	pdf.MultiCell(153, 5, tr("___________________________ "), "", "C", false)
	cellY += 5
	pdf.SetXY(29, cellY)
	pdf.MultiCell(153, 5, tr("Firma"), "", "C", false)
	cellY += 10
	pdf.SetFont("Times", "B", 10)
	pdf.SetXY(29, cellY)
	pdf.MultiCell(153, 5, tr("NOTA. "), "", "", false)

	pdf.SetFont("Times", "", 10)
	pdf.SetXY(29, cellY)
	pdf.MultiCell(153, 5, tr("NOTA. De igual forma, se deben reservar presupuestalmente aquellos saldos de órdenes de compra o serviciosen que no se utilizo la totalidad del registro presupuestal "), "", "", false)

	return pdf
}

func obtenerMes(numero_mes int) (mes string) {
	meses := []string{
		"Enero",
		"Febrero",
		"Marzo",
		"Abril",
		"Mayo",
		"Junio",
		"Julio",
		"Agosto",
		"Septiembre",
		"Octubre",
		"Noviembre",
		"Diciembre",
	}
	return meses[numero_mes-1]
}

func footer(pdf *gofpdf.Fpdf, cellx float64, cellY float64) *gofpdf.Fpdf {

	col1Width := 60.0
	col2Width := 20.0
	totalWidth := col1Width*2 + col2Width*2
	pageWidth, _ := pdf.GetPageSize()
	startX := (pageWidth - totalWidth) / 2
	cellY += 150
	pdf.SetXY(startX, cellY)

	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetFont("Times", "", 10)
	pdf.SetXY(29, 265)
	pdf.MultiCell(153, 5, tr("Este documento es propiedad de la Universidad Distrital Francisco José de Caldas. Prohibida su reproducción por cualquier medio, sin previa autorización. "), "", "C", false)

	return pdf
}

func tablaFooter(pdf *gofpdf.Fpdf, cellx float64, cellY float64) *gofpdf.Fpdf {

	col1Width := 60.0
	col2Width := 20.0
	totalWidth := col1Width*2 + col2Width*2
	pageWidth, _ := pdf.GetPageSize()
	startX := (pageWidth - totalWidth) / 2
	cellY += 150
	pdf.SetXY(startX, cellY)
	pdf.SetFont("Times", "B", 8)
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	//fila 1
	pdf.CellFormat(col2Width, 5, "               ", "1", 0, "C", false, 0, "")
	pdf.CellFormat(col1Width, 5, "NOMBRE", "1", 0, "C", false, 0, "")
	pdf.CellFormat(col1Width, 5, "CARGO/TIPO CONTRATO", "1", 0, "C", false, 0, "")
	pdf.CellFormat(col2Width, 5, "FIRMA", "1", 0, "C", false, 0, "")

	//fila 2
	pdf.SetXY(startX, cellY+5)
	pdf.CellFormat(col2Width, 5, tr("PROYECTÓ"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(col1Width, 5, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(col1Width, 5, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(col2Width, 5, "", "1", 0, "C", false, 0, "")

	//fila 2
	pdf.SetXY(startX, cellY+10)
	pdf.CellFormat(col2Width, 5, tr("REVISOÓ"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(col1Width, 5, "", "1", 0, "l", false, 0, "")
	pdf.CellFormat(col1Width, 5, "", "1", 0, "l", false, 0, "")
	pdf.CellFormat(col2Width, 5, "", "1", 0, "l", false, 0, "")

	//fila 2
	pdf.SetXY(startX, cellY+15)
	pdf.CellFormat(col2Width, 5, tr("APROBÓ"), "1", 0, "l", false, 0, "")
	pdf.CellFormat(col1Width, 5, "", "1", 0, "l", false, 0, "")
	pdf.CellFormat(col1Width, 5, "", "1", 0, "l", false, 0, "")
	pdf.CellFormat(col2Width, 5, "", "1", 0, "l", false, 0, "")

	return pdf
}

func docuementos() (docuementos map[string]string) {
	docuementos = make(map[string]string)
	docuementos["COPS1"] = "Copia de la Orden o Contrato de Prestación de Servicio (para el primer pago"
	docuementos["COPS2"] = "Copia de la Orden de Compra o Suministro (para el primer pago). "
	docuementos["COPS3"] = "Copia de Certificados de Disponibilidad y Reserva Presupuestal. "
	docuementos["COPS4"] = "Acta de inicio (únicamente para el primer pago). "
	docuementos["COPS5"] = "Certificado de cumplido o recibo a satisfacción por parte del supervisor."
	docuementos["FACT6"] = "Factura original. "
	docuementos["COPS"] = "Entrada a almacén de bienes devolutivos. "
	docuementos["ARL"] = "Copia del pago del aporte al Sistema General de Seguridad Social en Salud. "
	docuementos["SSP"] = "Copia del pago del aporte al Sistema General de Seguridad Social en Pensión. "
	docuementos["RUT"] = "Copia del Registro Único Tributario (RUT). "
	docuementos["MMEN"] = "Medio magnético y escrito de la nomina. "
	docuementos["POL"] = "Copia de las pólizas necesarias. "
	docuementos["AL"] = "Acta de liquidación o Acta de Terminación. "
	docuementos["PTA"] = "Pasabordo en tiquetes aéreos. "
	docuementos["RMAAV"] = "Resolución motivada autorizando avance y /o viáticos. "
	docuementos["POLA"] = "Pólizas acordadas en el contrato y en el presupuesto para los anticipos. "

	return docuementos

}

func contains(documentos []string, documento string) bool {
	for _, value := range documentos {
		if value == documento {
			return true
		}
	}
	return false

}
