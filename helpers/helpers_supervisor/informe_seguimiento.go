package helpers_supervisor

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/jung-kurt/gofpdf"

	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func GetInformacionProveedor(numero_contrato_suscrito string, vigencia string) (informacion_proveedor models.InformacionContratoProveedor, outputError map[string]interface{}) {

	contrato_contratista, err := GetInformacionContratoContratista(numero_contrato_suscrito, vigencia)
	if err == nil {
		contratistas, err := ContratosContratistaTemp(contrato_contratista.InformacionContratista.Documento.Numero)
		if err == nil {
			return contratistas[0], nil
		} else {
			return informacion_proveedor, map[string]interface{}{"funcion": "/GetInformacionProveedor", "err": err, "status": "502"}
		}
	} else {
		return informacion_proveedor, map[string]interface{}{"funcion": "/GetInformacionProveedor", "err": err, "status": "502"}
	}
	return
}

func ContratosContratistaTemp(numero_documento string) (contrato_proveedor []models.InformacionContratoProveedor, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			//fmt.Println("error", err)
			outputError = map[string]interface{}{"funcion": "/ContratosContratista", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	if contratos_persona, outputError := GetContratosPersona(numero_documento); outputError == nil {
		//fmt.Println("contratos_persona", contratos_persona)
		for _, contrato_persona := range contratos_persona.ContratosPersonas.ContratoPersona {
			var contrato models.InformacionContrato
			contrato, outputError = GetContrato(contrato_persona.NumeroContrato, contrato_persona.Vigencia)

			if (contrato == models.InformacionContrato{} || outputError != nil) {
				continue
			}
			var informacion_contrato_contratista models.InformacionContratoContratista
			informacion_contrato_contratista, outputError = GetInformacionContratoContratista(contrato_persona.NumeroContrato, contrato_persona.Vigencia)
			// se llena el contrato original en el indice 0

			if cdprp, outputError := GetRP(contrato_persona.NumeroCDP, contrato_persona.Vigencia); outputError == nil {
				for _, rp := range cdprp.CdpXRp.CdpRp {
					var contrato_proveedor_individual models.InformacionContratoProveedor
					contrato_proveedor_individual.TipoContrato = contrato.Contrato.TipoContrato
					contrato_proveedor_individual.NumeroContratoSuscrito = contrato_persona.NumeroContrato
					contrato_proveedor_individual.Vigencia = contrato_persona.Vigencia
					contrato_proveedor_individual.NumeroRp = rp.RpNumeroRegistro
					contrato_proveedor_individual.VigenciaRp = rp.RpVigencia
					contrato_proveedor_individual.NombreProveedor = informacion_contrato_contratista.InformacionContratista.NombreCompleto
					contrato_proveedor_individual.NombreDependencia = informacion_contrato_contratista.InformacionContratista.Dependencia
					contrato_proveedor_individual.NumDocumentoSupervisor = contrato.Contrato.Supervisor.DocumentoIdentificacion
					contrato_proveedor_individual.NumeroCdp = contrato_persona.NumeroCDP
					contrato_proveedor_individual.VigenciaCdp = contrato_persona.Vigencia
					contrato_proveedor_individual.Rubro = contrato.Contrato.Rubro
					contrato_proveedor = append(contrato_proveedor, contrato_proveedor_individual)
				}

			} else {
				logs.Error(outputError)
				continue
			}

		}
	} else {
		logs.Error(outputError)
		outputError = map[string]interface{}{"funcion": "/contratosContratista/GetContratosPersona", "err": outputError, "status": "502"}
		return nil, outputError
	}
	return contrato_proveedor, nil
}

func CreateInformeSeguimiento(numero_contrato_suscrito int, vigencia_contrato string, tipo_pago string, periodo_inicio string, periodo_fin string, tipo_factura string, numero_cuenta_factura string, valor_pagar int, tipo_cuenta string, numero_cuenta string, banco string) (informe_seguimiento models.InformeSeguimiento, outputError map[string]interface{}) {

	var valor_total_contrato int
	var saldo_contrato int
	balance, err := GetBalanceFinancieroContrato(strconv.Itoa(numero_contrato_suscrito), vigencia_contrato)
	if err == nil {
		valor := strings.Split(balance.TotalContrato, ".")[0]
		valor_total_contrato, _ = strconv.Atoi(valor)
		saldo_contrato, _ = strconv.Atoi(balance.Saldo)
	}
	fmt.Println("valor_total_contrato", valor_total_contrato)
	fmt.Println("saldo_contrato", saldo_contrato)
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetMargins(7, 7, 7)
	pdf.SetAutoPageBreak(true, 7)

	pdf.SetHeaderFunc(func() {
		pdf = header(pdf)
		pdf.Ln(35)
	})
	//body_primera_parte(pdf *gofpdf.Fpdf, nombre_proveedor string, numero_nit string, cumplimiento_contrato bool, tipo_contrato string, fecha_inicio string, numero_contrato string, cdp string, vigencia_cdp string, crp string, vigencia_crp string)
	pdf.AddPage()

	vigencia, _ := strconv.Atoi(vigencia_contrato)
	contrato, err := GetInformacionProveedor(strconv.Itoa(numero_contrato_suscrito), vigencia_contrato)
	if err != nil {
		outputError = map[string]interface{}{"funcion": "/CreateInformeSeguimiento", "err": err, "status": "502"}
		return informe_seguimiento, outputError
	}

	acta_inicio, err := GetActaInicio(strconv.Itoa(numero_contrato_suscrito), vigencia)
	if err == nil {
		pdf = body_primera_parte(
			pdf,
			contrato.NombreDependencia,
			contrato.NombreProveedor,
			"830006800-3",
			true,
			"Contrato de comisión",
			formatear_fecha(acta_inicio.FechaInicio),
			contrato.NumeroContratoSuscrito,
			contrato.NumeroCdp,
			contrato.VigenciaCdp,
			contrato.NumeroRp,
			contrato.VigenciaRp,
		)
	} else {

		pdf = body_primera_parte(
			pdf,
			contrato.NombreDependencia,
			contrato.NombreProveedor,
			"830006800-3",
			true,
			"Contrato de comisión",
			"08/04/2024",
			contrato.NumeroContratoSuscrito,
			contrato.NumeroCdp,
			contrato.VigenciaCdp,
			contrato.NumeroRp,
			contrato.VigenciaRp,
		)
	}

	pdf = body_segunda_parte(
		pdf,
		tipo_factura,
		numero_cuenta_factura,
		valor_total_contrato,
		periodo_inicio,
		periodo_fin,
		saldo_contrato,
		formatear_fecha(acta_inicio.FechaInicio),
		formatear_fecha(acta_inicio.FechaFin),
		tipo_cuenta,
		numero_cuenta,
		banco)

	// Crear el PDF

	// Codificar el PDF
	encodedFile := encodePDF(pdf)
	nombre := "prueba"

	informe_seguimiento = models.InformeSeguimiento{File: encodedFile, Archivo: nombre}
	return informe_seguimiento, nil
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

func fontStyle(pdf *gofpdf.Fpdf, style string, size float64, bw int) {
	pdf.SetTextColor(bw, bw, bw)
	pdf.SetFont("Helvetica", style, size)
}

func header(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	path := beego.AppConfig.String("StaticPath")

	now := time.Now()
	formattedDate := now.Format("02/01/2006")

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
	pdf.Text(cellX2+3, cellY2+3, tr("CUMPLIDO A SATISFACCIÓN POR"))
	pdf.Text(cellX2+7, cellY2+6, tr("PARTE DE LA DEPENDENCIA"))

	// Tercer recuadro
	pdf.SetFont("Times", "", 10)
	pdf.Rect(cellX2, cellY2+7.5, cellWidth2, cellHeight2, "")
	pdf.Text(cellX2+2, cellY2+12, tr("Macroproceso:  Gestión de Recursos"))

	// Cuarto recuadro
	pdf.SetFont("Times", "", 10)
	pdf.Rect(cellX2, cellY2+15, cellWidth2, cellHeight2+2.5, "")
	pdf.Text(cellX2+6, cellY2+20, tr("Proceso: Gestión Contractual"))

	// Quinto recuadro
	pdf.SetFont("Times", "", 9)
	pdf.Rect(cellX2+55, cellY2, cellWidth2-12, cellHeight2, "")
	pdf.Text(cellX2+56, cellY2+5, "Codigo: GC-PR-003-FR-012")

	// Sexto recuadro
	pdf.SetFont("Times", "", 10)
	pdf.Rect(cellX2+55, cellY2+7.5, cellWidth2-12, cellHeight2, "")
	pdf.Text(cellX2+56, cellY2+12, tr("Versión: 05"))

	// Séptimo recuadro
	pdf.Rect(cellX2+55, cellY2+15, cellWidth2-12, cellHeight2+2.5, "")
	pdf.Text(cellX2+56, cellY2+20, tr("Fecha de Aprobación:"))
	pdf.Text(cellX2+56, cellY2+24, formattedDate)

	// Último recuadro con logo
	pdf.Rect(cellX+123, cellY, cellWidth+2, cellHeight, "")
	pdf = image(pdf, path+"/img/EscudoSigud.png", cellX+123+1, cellY+10, cellWidth-1, cellHeight-20)

	return pdf
}

func textoNegrilla(pdf *gofpdf.Fpdf, texto string) string {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetFont("Times", "B", 12)
	pdf.CellFormat(0, 10, tr(texto), "", 1, "C", false, 0, "")
	return texto
}

func textoMayuscula(pdf *gofpdf.Fpdf, texto string) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetFont("Times", "", 12)
	pdf.CellFormat(0, 10, tr(strings.ToUpper(texto)), "", 1, "C", false, 0, "")
	return pdf
}

func body_primera_parte(pdf *gofpdf.Fpdf, dependencia string, nombre_proveedor string, numero_nit string, cumplimiento_contrato bool, tipo_contrato string, fecha_inicio string, numero_contrato string, cdp string, vigencia_cdp string, crp string, vigencia_crp string) *gofpdf.Fpdf {

	var cumplimiento string
	if cumplimiento_contrato {
		cumplimiento = "totalmente"
	} else {
		cumplimiento = "parcialmente"
	}

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetFont("Times", "", 12)

	pdf.CellFormat(0, 10, tr("UNIVERSIDAD DISTRITAL FRANCISCO JOSÉ DE CALDAS"), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 10, tr("("+dependencia+")"), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 10, tr("En ejercicio de las funciones de (jefe de la dependencia o supervisor del contrato)"), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 10, tr("CERTIFICA"), "", 1, "C", false, 0, "")
	pdf.Ln(10) // Espacio después de la certificación

	/**
	pdf.SetFont("Times", "B", 12)
	nombre_proveedor = strings.ToUpper(nombre_proveedor)
	**/
	// Contenido principal
	pdf.SetFont("Times", "", 12)
	pdf.MultiCell(0, 10, tr(fmt.Sprintf(`Que el contratista %s identificado con NIT %s cumplió %s a satisfacción con las obligaciones y objeto del %s Nro. %s de fecha %s garantizada y perfeccionada con Certificado de Disponibilidad Presupuestal No. %s de %s y Certificado de Registro Presupuestal No. %s de %s.`, nombre_proveedor, numero_nit, cumplimiento, tipo_contrato, numero_contrato, fecha_inicio, cdp, vigencia_cdp, crp, vigencia_crp)), "", "", false)

	pdf.MultiCell(0, 10, tr(`Que conforme con los documentos aportados el contratista cumple con la afiliación y pagos al Sistema General de Seguridad Social de Salud y Pensiones Riesgos Laborales y las obligaciones parafiscales por el período y desembolso aquí causados y autorizados. Así mismo los documentos requeridos (RUT con impresión actualizada Certificado de Cámara de Comercio “no mayor a 90 días” cuenta bancaria fotocopia de la Cédula Actas de Entrega de Elementos o Remisiones Informes de Seguimiento de Supervisión Evaluación del Proveedor y Actas de Liquidación “si se requiere”) para el giro respectivo.`), "", "", false)

	return pdf
}

func body_segunda_parte(pdf *gofpdf.Fpdf, tipo_factura string, numero_cuenta_factura string, valor_total_contrato int, periodo_inicio, periodo_fin string, saldo_contrato int, fecha_inicio string, fecha_fin string, tipo_cuenta string, numero_cuenta string, banco string) *gofpdf.Fpdf {

	pdf.SetFont("Times", "", 12)

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.MultiCell(0, 10, tr(fmt.Sprintf(`Que el valor causado de conformidad con la %s de Venta o Cuenta de cobro No. %s es %s PESOS $%s pesos m/cte.`, tipo_factura, numero_cuenta_factura, strings.ToUpper(ValorLetras(valor_total_contrato)), FormatNumber(valor_total_contrato, 0, ".", ","))), "", "", false)

	pdf.MultiCell(0, 10, tr(fmt.Sprintf(`Que el valor total del contrato corresponde a %s $%s pesos m/cte.`, strings.ToUpper(ValorLetras(valor_total_contrato)), FormatNumber(valor_total_contrato, 0, ".", ","))), "", "", false)

	pdf.MultiCell(0, 10, tr(fmt.Sprintf(`Que el presente pago corresponde al período de %s a %s de ejecución parcial total o único pago del contrato.`, periodo_inicio, periodo_fin)), "", "", false)

	pdf.MultiCell(0, 10, tr(fmt.Sprintf(`Quedando un saldo correspondiente a $%s pesos m/cte.`, FormatNumber(saldo_contrato, 0, ".", ","))), "", "", false)

	pdf.MultiCell(0, 10, tr(fmt.Sprintf(`Que el presente pago se encuentra en cumplimiento dentro del tiempo de ejecución del contrato del %s al %s.`, fecha_inicio, fecha_fin)), "", "", false)

	pdf.MultiCell(0, 10, tr(fmt.Sprintf(`Que tal valor debe girarse por petición del contratista a la Cuenta %s  No. %s del Banco %s.`, tipo_cuenta, numero_cuenta, strings.ToUpper(banco))), "", "", false)

	pdf.MultiCell(0, 10, tr(`Con el presente cumplido y de acuerdo a lo establecido en los numerales 32 y 33 del Artículo 18° de la Resolución de Rectoría No. 629 de 2016- Manual de Interventoría y Supervisión certifico que los informes físicos técnicos financieros y administrativos sobre el avance de la ejecución del objeto contractual reposan en el expediente del ______________ (Contrato/ Contrato de Comisión/Orden de Compra/Orden de Servicio/ Orden de Compra CCE) No.________ de________. De igual forma certifico que se verificaron las condiciones y elementos que hacen parte de la(s) factura(s) No. ___________________ acorde con lo establecido en la ficha técnica del proceso en mención garantizando la calidad del bien o servicio adquirido por la Universidad.`), "", "", false)

	pdf.MultiCell(0, 10, tr(`La presente se expide a los ____ días del mes _____ de ______.`), "", "", false)

	pdf.Ln(10) // Espacio para la firma

	pdf.SetFont("Times", "B", 12)
	pdf.MultiCell(0, 10, tr(`__________________________`), "", "", false)
	pdf.MultiCell(0, 10, tr(`NOMBRE`), "", "", false)
	pdf.MultiCell(0, 10, tr(`C.C ______________ de Bogotá`), "", "", false)
	pdf.MultiCell(0, 10, tr(`CARGO`), "", "", false)
	pdf.MultiCell(0, 10, tr(`DEPENDENCIA`), "", "", false)
	pdf.MultiCell(0, 10, tr(`Supervisor ______________________Contrato/ Contrato de Comisión/Orden de Compra/Orden de Servicio/ Orden de Compra CCE) No. ________ de________.`), "", "", false)

	// Tabla de elaboración, revisión y aprobación
	pdf.Ln(10)
	pdf.SetFont("Times", "B", 10)
	pdf.CellFormat(0, 10, tr(`Elaboró`), "1", 0, "C", false, 0, "")
	pdf.CellFormat(0, 10, tr(`Revisó`), "1", 0, "C", false, 0, "")
	pdf.CellFormat(0, 10, tr(`Aprobó`), "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("Times", "", 10)
	pdf.CellFormat(0, 10, tr(`NOMBRE`), "1", 0, "C", false, 0, "")
	pdf.CellFormat(0, 10, tr(`NOMBRE`), "1", 0, "C", false, 0, "")
	pdf.CellFormat(0, 10, tr(`NOMBRE`), "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	pdf.CellFormat(0, 10, tr(`CARGO`), "1", 0, "C", false, 0, "")
	pdf.CellFormat(0, 10, tr(`CARGO`), "1", 0, "C", false, 0, "")
	pdf.CellFormat(0, 10, tr(`CARGO`), "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	pdf.CellFormat(0, 10, tr(`FIRMA`), "1", 0, "C", false, 0, "")
	pdf.CellFormat(0, 10, tr(`FIRMA`), "1", 0, "C", false, 0, "")
	pdf.CellFormat(0, 10, tr(`FIRMA`), "1", 0, "C", false, 0, "")

	return pdf
}
