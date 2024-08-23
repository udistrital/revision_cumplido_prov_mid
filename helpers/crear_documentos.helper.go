package helpers

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func Header(pdf *gofpdf.Fpdf, tipo_documento string, proceso string, codigo string, version string, fecha_aprobacion string) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	path := "static"

	pdf.SetFont("Times", "", 12)
	pdf.SetMargins(10, 10, 10)
	pdf.SetFillColor(240, 240, 240)

	cellX, cellY := 30.0, 10.0
	cellX2, cellY2 := 55.0, 10.0
	cellWidth, cellHeight := 25.0, 25.0
	cellWidth2, cellHeight2 := 55.0, 7.5

	// Primer recuadro con logo
	pdf.Rect(cellX, cellY, cellWidth, cellHeight, "")
	pdf = image(pdf, path+"/img/EscudoUd.png", cellX+1, cellY+1, cellWidth-2, cellHeight-4)

	// Segundo recuadro
	pdf.Rect(cellX2, cellY2, cellWidth2, cellHeight2, "")
	pdf.SetFont("Times", "B", 8)
	textTipoDocumento := tr(tipo_documento)
	textWidthTipoDocumento := pdf.GetStringWidth(textTipoDocumento)
	lineHeightTipoDocumento := pdf.PointConvert(8)
	textLinesTipoDocumento := int(textWidthTipoDocumento/(cellWidth2-6)) + 1
	totalTextHeightTipoDocumento := float64(textLinesTipoDocumento) * lineHeightTipoDocumento
	adjustedYTipoDocumento := (cellHeight2-totalTextHeightTipoDocumento)/2 + cellY2
	pdf.SetXY(cellX2+3, adjustedYTipoDocumento)
	pdf.MultiCell(cellWidth2-6, lineHeightTipoDocumento, textTipoDocumento, "", "C", false)

	// Tercer recuadro
	pdf.SetFont("Times", "", 10)
	pdf.Rect(cellX2, cellY2+7.5, cellWidth2, cellHeight2, "")
	pdf.Text(cellX2+2, cellY2+12, tr("Macroproceso:  Gestión de Recursos"))

	// Cuarto recuadro
	pdf.SetFont("Times", "", 10)
	pdf.Rect(cellX2, cellY2+15, cellWidth2, cellHeight2+2.5, "")
	pdf.SetXY(cellX2+3, cellY2+15)
	textProceso := tr("Proceso: " + proceso)
	textWidthProceso := pdf.GetStringWidth(textProceso)
	lineHeightProceso := pdf.PointConvert(10)
	textLinesProceso := int(textWidthProceso/(cellWidth2-6)) + 1
	totalTextHeightProceso := float64(textLinesProceso) * lineHeightProceso
	adjustedYProceso := (cellHeight2+2.5-totalTextHeightProceso)/2 + cellY2 + 15

	// Centrar el texto dentro de la celda
	pdf.SetXY(cellX2+3, adjustedYProceso)
	pdf.MultiCell(cellWidth2-6, lineHeightProceso, textProceso, "", "C", false)

	// Quinto recuadro
	pdf.SetFont("Times", "", 9)
	pdf.Rect(cellX2+55, cellY2, cellWidth2-12, cellHeight2, "")
	pdf.Text(cellX2+56, cellY2+5, "Codigo: "+codigo)

	// Sexto recuadro
	pdf.SetFont("Times", "", 10)
	pdf.Rect(cellX2+55, cellY2+7.5, cellWidth2-12, cellHeight2, "")
	pdf.Text(cellX2+56, cellY2+12, tr("Versión: "+version))

	// Séptimo recuadro
	pdf.Rect(cellX2+55, cellY2+15, cellWidth2-12, cellHeight2+2.5, "")
	pdf.Text(cellX2+56, cellY2+19, tr("Fecha de Aprobación:"))
	pdf.Text(cellX2+56, cellY2+23, fecha_aprobacion)

	// Último recuadro con logo
	pdf.Rect(cellX+123, cellY, cellWidth+2, cellHeight, "")
	pdf = image(pdf, path+"/img/EscudoSigud.png", cellX+123+1, cellY+10, cellWidth-1, cellHeight-20)

	return pdf
}

func GenerarPdfAutorizacionPago(autorizacion *models.DocuementoAutorizacionPago) string {

	if autorizacion == nil {
		fmt.Println("Error al generar el documento")
		return ""

	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	now := time.Now()
	//////Header
	cellX, cellY := 30.0, 10.0
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
		pdf = Header(pdf, "AUTORIZACIÓN DE GIRO", "Gestión de Recursos Financieros", "GRF-PR-007-FR-005", "04", "14/01/2022")
	})
	pdf.AddPage()

	pdf = body(pdf, cellX, cellY, month, day, year, autorizacion)

	file_name := fmt.Sprintf("AutorizacionPago_%s.pdf", autorizacion.NombreProveedor+"-"+autorizacion.DocumentoProveedor)
	err := pdf.OutputFileAndClose(file_name + ".pdf")
	if err != nil {
		fmt.Println("Error al guardar el archivo PDF:", err)
	} else {
		fmt.Println("PDF generado exitosamente")

	}

	fileData, err := ioutil.ReadFile(file_name + ".pdf")
	if err != nil {
		fmt.Println("Error al leer el archivo PDF:", err)
		return ""
	}

	base64Data := base64.StdEncoding.EncodeToString(fileData)

	err = os.Remove(file_name)
	if err != nil {
		fmt.Println("Error al eliminar el archivo temporal:", err)
	}
	return base64Data
}

func header(pdf *gofpdf.Fpdf, cellX float64, cellY float64, cellX2 float64, cellY2 float64, cellWidth float64, cellHeight float64, cellWidth2 float64, cellHeight2 float64, formattedDate string) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.Rect(cellX, cellY, cellWidth, cellHeight, "")
	logoU := "static/img/EscudoUd.png"
	pdf.Image(logoU, cellX+2, cellY, cellWidth-3, cellHeight-5, false, "", 0, "")

	pdf.Rect(cellX+123, cellY, cellWidth+2, cellHeight, "")
	logoSigud := "static/img/EscudoSigud.png"
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
	pdf = CrearTablaDocumentos(pdf, cellx, cellY-10, autorizacion)

	return pdf
}

func CrearTablaDocumentos(pdf *gofpdf.Fpdf, cellx float64, cellY float64, autorizacion *models.DocuementoAutorizacionPago) *gofpdf.Fpdf {
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

		if DocumentoEnLista(autorizacion.DocumentosCargados, key) {
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

func DocumentoEnLista(documentos []string, documento string) bool {
	for _, value := range documentos {
		if value == documento {
			return true
		}
	}
	return false

}

func CrearPdfInformeSatisfaccion(dependencia string, nombre_proveedor string, numero_nit string, cumplimiento_contrato bool, tipo_contrato string, fecha_inicio time.Time, numero_contrato string, cdp string, vigencia_cdp string, rp string, vigencia_rp string, cargo string, tipo_factura string, numero_cuenta_factura string, valor_total_contrato int, periodo_inicio string, periodo_fin string, saldo_contrato int, fecha_fin time.Time, tipo_cuenta string, numero_cuenta string, nombre_banco string, supervisor string, vigencia string) (informe_satisfaccion models.InformeSatisfaccion, outputError interface{}) {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(25, 20, 25)
	pdf.SetAutoPageBreak(true, 20)

	pdf.SetHeaderFunc(func() {
		pdf = Header(pdf, "CUMPLIDO A SATISFACCIÓN POR PARTE DE LA DEPENDENCIA", "Gestión Contractual", "GC-PR-003-FR-012", "05", "13/10/2021")
		pdf.Ln(20)
	})

	pdf.AddPage()

	pdf = body_primera_parte(
		pdf,
		dependencia,
		nombre_proveedor,
		numero_nit,
		cumplimiento_contrato,
		tipo_contrato,
		fecha_inicio,
		numero_contrato,
		cdp,
		vigencia_cdp,
		rp,
		vigencia_rp,
		cargo)

	pdf = body_segunda_parte(
		pdf,
		tipo_factura,
		numero_cuenta_factura,
		valor_total_contrato,
		periodo_inicio,
		periodo_fin,
		saldo_contrato,
		fecha_inicio,
		fecha_fin,
		tipo_cuenta,
		numero_cuenta,
		nombre_banco)

	pdf = footerInformeSeguimiento(pdf,
		numero_contrato,
		fecha_inicio,
		tipo_contrato,
		numero_cuenta_factura,
		supervisor,
		cargo,
		vigencia,
		dependencia)

	encodedFile := encodePDF(pdf)
	nombre := "prueba"
	informe_satisfaccion = models.InformeSatisfaccion{File: encodedFile, Archivo: nombre}
	return informe_satisfaccion, nil
}

func body_primera_parte(pdf *gofpdf.Fpdf, dependencia string, nombre_proveedor string, numero_nit string, cumplimiento_contrato bool, tipo_contrato string, fecha_inicio time.Time, numero_contrato string, cdp string, vigencia_cdp string, crp string, vigencia_crp string, cargo string) *gofpdf.Fpdf {

	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetFont("Times", "", 10)

	pdf.CellFormat(0, 8, tr("UNIVERSIDAD DISTRITAL FRANCISCO JOSÉ DE CALDAS"), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 8, tr("("+dependencia+")"), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 8, tr(fmt.Sprintf(`En ejercicio de las funciones de (%s)`, VerificarJefe(cargo))), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 8, tr("CERTIFICA"), "", 1, "C", false, 0, "")

	var cumplimiento string
	if cumplimiento_contrato {
		cumplimiento = "totalmente"
	} else {
		cumplimiento = "parcialmente"
	}

	// Espacio después de la certificación

	// Contenido principal
	pdf.SetFont("Times", "", 10)
	pdf.SetMargins(30, 30, 30)
	pdf.MultiCell(0, 8, "", "", "J", false)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`Que el contratista %s identificado con NIT %s cumplió %s a satisfacción con las obligaciones y objeto del %s Nro. %s de fecha %s garantizada y perfeccionada con Certificado de Disponibilidad Presupuestal No. %s de %s y Certificado de Registro Presupuestal No. %s de %s.`, nombre_proveedor, numero_nit, cumplimiento, tipo_contrato, numero_contrato, formatear_fecha(fecha_inicio), cdp, vigencia_cdp, crp, vigencia_crp)), "", "", false)
	pdf.Ln(5)
	pdf.SetMargins(30, 30, 30)
	pdf.MultiCell(0, 8, tr(`Que conforme con los documentos aportados el contratista cumple con la afiliación y pagos al Sistema General de Seguridad Social de Salud y Pensiones Riesgos Laborales y las obligaciones parafiscales por el período y desembolso aquí causados y autorizados. Así mismo los documentos requeridos (RUT con impresión actualizada Certificado de Cámara de Comercio “no mayor a 90 días” cuenta bancaria fotocopia de la Cédula Actas de Entrega de Elementos o Remisiones Informes de Seguimiento de Supervisión Evaluación del Proveedor y Actas de Liquidación “si se requiere”) para el giro respectivo.`), "", "", false)
	pdf.Ln(5)
	return pdf
}

func ObtenerMes(mes int) string {
	meses := []string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}
	return meses[mes-1]
}

func body_segunda_parte(pdf *gofpdf.Fpdf, tipo_factura string, numero_cuenta_factura string, valor_total_contrato int, periodo_inicio string, periodo_fin string, saldo_contrato int, fecha_inicio time.Time, fecha_fin time.Time, tipo_cuenta string, numero_cuenta string, nombre_banco string) *gofpdf.Fpdf {

	fmt.Println("Periodo inicio", periodo_inicio)
	pdf.SetFont("Times", "", 10)

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetMargins(30, 30, 30)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`Que el valor causado de conformidad con la %s de Venta o Cuenta de cobro No. %s es %s PESOS $%s pesos m/cte.`, tipo_factura, numero_cuenta_factura, strings.ToUpper(ValorLetras(valor_total_contrato)), FormatNumber(valor_total_contrato, 0, ".", ","))), "", "J", false)
	pdf.Ln(5)
	pdf.SetMargins(30, 30, 30)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`Que el valor total del contrato corresponde a %s $%s pesos m/cte.`, strings.ToUpper(ValorLetras(valor_total_contrato)), FormatNumber(valor_total_contrato, 0, ".", ","))), "", "J", false)
	pdf.Ln(5)
	pdf.SetMargins(30, 30, 30)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`Que el presente pago corresponde al período de %s a %s de ejecución parcial total o único pago del contrato.`, periodo_inicio, periodo_fin)), "", "J", false)
	pdf.Ln(5)
	pdf.SetMargins(30, 30, 30)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`Quedando un saldo correspondiente a $%s pesos m/cte.`, FormatNumber(saldo_contrato, 0, ".", ","))), "", "", false)
	pdf.Ln(5)
	pdf.SetMargins(30, 30, 30)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`Que el presente pago se encuentra en cumplimiento dentro del tiempo de ejecución del contrato del %s al %s.`, formatear_fecha(fecha_inicio), formatear_fecha(fecha_fin))), "", "J", false)

	pdf.SetMargins(30, 30, 30)
	pdf.MultiCell(0, 8, "", "", "J", false)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`Que tal valor debe girarse por petición del contratista a la Cuenta %s  No. %s del Banco %s.`, tipo_cuenta, numero_cuenta, strings.ToUpper(nombre_banco))), "", "J", false)

	return pdf
}

func footerInformeSeguimiento(pdf *gofpdf.Fpdf, contrato_suscrito string, fecha_inicio time.Time, tipo_contrato string, numero_factura string, supervisor string, cargo string, vigencia string, dependencia string) *gofpdf.Fpdf {

	dia := time.Now().Day()
	mes := int(time.Now().Month())
	año := time.Now().Year()

	pdf.SetFont("Times", "", 10)

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetMargins(30, 0, 30)
	pdf.MultiCell(0, 8, "", "", "J", false)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`Con el presente cumplido y de acuerdo a lo establecido en los numerales 32 y 33 del Artículo 18° de la Resolución de Rectoría No. 629 de 2016- Manual de Interventoría y Supervisión certifico que los informes físicos técnicos financieros y administrativos sobre el avance de la ejecución del objeto contractual reposan en el expediente del %s No. %s de %s. De igual forma certifico que se verificaron las condiciones y elementos que hacen parte de la(s) factura(s) No. %s acorde con lo establecido en la ficha técnica del proceso en mención garantizando la calidad del bien o servicio adquirido por la Universidad.`, tipo_contrato, contrato_suscrito, formatear_fecha(fecha_inicio), numero_factura)), "", "J", false)
	pdf.Ln(15)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`La presente se expide a los %s días del mes %s de %s.`, strconv.Itoa(dia), ObtenerMes(mes), strconv.Itoa(año))), "", "J", false)
	pdf.Ln(5) // Espacio para la firma
	pdf.SetMargins(30, 30, 30)
	pdf.SetFont("Times", "B", 10)
	pdf.MultiCell(0, 8, tr(`__________________________`), "", "", false)
	pdf.MultiCell(0, 8, tr(`NOMBRE`), "", "", false)
	pdf.MultiCell(0, 8, tr(supervisor), "", "", false)
	pdf.MultiCell(0, 8, tr(`C.C ______________ de Bogotá`), "", "", false)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`CARGO %s`, VerificarJefe(cargo))), "", "", false)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`DEPENDENCIA %s`, dependencia)), "", "", false)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`Supervisor %s Contrato/ Contrato de Comisión/Orden de Compra/Orden de Servicio/ Orden de Compra CCE) No. %s de %s.`, supervisor, contrato_suscrito, vigencia)), "", "", false)

	return pdf

}

func VerificarJefe(cargo string) string {
	palabras := strings.Fields(strings.ToLower(cargo))
	if len(palabras) > 0 && palabras[0] == "jefe" {
		return cargo
	} else {
		return "Supervisor del Contrato"
	}
}

func formatear_fecha(fecha time.Time) (fecha_formateada string) {
	layout := "02/01/2006"
	return fecha.Format(layout)
}

func encodePDF(pdf *gofpdf.Fpdf) string {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	pdf.OutputFileAndClose("/home/faidercamilo/go/src/github.com/udistrital/prueba.pdf") // para guardar el archivo localmente
	pdf.Output(writer)
	writer.Flush()
	encodedFile := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return encodedFile
}

func image(pdf *gofpdf.Fpdf, image string, x, y, w, h float64) *gofpdf.Fpdf {
	pdf.ImageOptions(image, x, y, w, h, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	return pdf
}
