package helpers_supervisor

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/jung-kurt/gofpdf"

	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func CreateInformeSeguimiento(numero_contrato_suscrito string, vigencia_contrato string, tipo_pago string, periodo_inicio string, periodo_fin string, tipo_soporte_pagar string, numero_cuenta_factura string, valor_pagar string, tipo_cuenta string, numero_cuenta string, banco string) (informe_seguimiento models.InformeSeguimiento, outputError map[string]interface{}) {

	// Se obtiene la información del proveedor
	informacion_proveedor, err := GetInformacionProveedor(numero_contrato_suscrito, vigencia_contrato)
	if err == nil {
		fmt.Println(informacion_proveedor)
		// Se crea el PDF

		pdf := gofpdf.New("P", "mm", "Letter", "")
		pdf.AddPage()
		pdf.SetMargins(7, 7, 7)
		pdf.SetAutoPageBreak(true, 7) // margen inferior
		pdf.SetY(pdf.GetCellMargin())
		pdf.SetX(pdf.GetCellMargin())

		pdf = header(pdf, true)
		pdf = footer(pdf, "-COPIA ESTUDIANTE-")
		pdf = separador(pdf)

		// Se codifica el PDF
		encodedFile := encodePDF(pdf)
		nombre := "InformeSeguimiento_" + informacion_proveedor.NumeroContratoSuscrito + "_" + informacion_proveedor.Vigencia + ".pdf"

		informe_seguimiento = models.InformeSeguimiento{File: encodedFile, Archivo: nombre}
		return informe_seguimiento, nil
	}
	return
}

type styling struct {
	mL float64 // margen izq
	mT float64 // margen sup
	mR float64 // margen der
	mB float64 // margen inf
	wW float64 // ancho area trabajo
	hW float64 // alto area trabajo
	//hH    float64 // alto header
	hF float64 // alto footer
	//lh int     // alto linea común
	//brdrs string  // estilo border común
}

func encodePDF(pdf *gofpdf.Fpdf) string {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	//pdf.OutputFileAndClose("../docs/recibo.pdf") // para guardar el archivo localmente
	pdf.Output(writer)
	writer.Flush()
	encodedFile := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return encodedFile
}

func image(pdf *gofpdf.Fpdf, image string, x, y, w, h float64) *gofpdf.Fpdf {
	//The ImageOptions method takes a file path, x, y, width, and height parameters, and an ImageOptions struct to specify a couple of options.
	pdf.ImageOptions(image, x, y, w, h, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	return pdf
}

func fontStyle(pdf *gofpdf.Fpdf, style string, size float64, bw int) {
	pdf.SetTextColor(bw, bw, bw)
	pdf.SetFont("Helvetica", style, size)
}

// Description: genera el encabezado reutilizable del recibo de pago
func header(pdf *gofpdf.Fpdf, banco bool) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	path := beego.AppConfig.String("StaticPath")
	pdf = image(pdf, path+"/img/UDEscudo2.png", 7, pdf.GetY(), 0, 17.5)

	if banco {
		pdf = image(pdf, path+"/img/banco.PNG", 198, pdf.GetY(), 0, 12.5)
	}

	pdf.SetXY(7, pdf.GetY())
	fontStyle(pdf, "B", 10, 0)
	pdf.Cell(13, 10, "")
	pdf.Cell(140, 10, "UNIVERSIDAD DISTRITAL")
	if banco {
		fontStyle(pdf, "B", 8, 0)
		pdf.Cell(50, 10, "PAGUE UNICAMENTE EN")
		fontStyle(pdf, "B", 10, 0)
	}
	pdf.Ln(4)
	pdf.Cell(13, 10, "")
	pdf.Cell(60, 10, tr("Francisco José de Caldas"))
	pdf.Cell(80, 10, "COMPROBANTE DE PAGO No ")

	if banco {
		fontStyle(pdf, "B", 8, 0)
		pdf.Cell(50, 10, "BANCO DE OCCIDENTE")
	} /* else {
		fontStyle(pdf, "", 8, 70)
		pdf.Cell(50, 10, "espacio para serial")
	} */

	pdf.Ln(4)
	fontStyle(pdf, "", 8, 0)
	pdf.Cell(13, 10, "")
	pdf.Cell(50, 10, "NIT 899.999.230-7")
	pdf.Ln(10)
	return pdf
}

// Description: genera el pie de paǵina reutilizable del recibo de pago
func footer(pdf *gofpdf.Fpdf, copiaPara string) *gofpdf.Fpdf {
	fontStyle(pdf, "", 8, 70)
	pdf.CellFormat(134, 5, copiaPara, "", 0, "C", false, 0, "")
	pdf.SetXY(142.9, pdf.GetY())
	pdf.CellFormat(66, 5, "-Espacio para timbre o sello Banco-", "", 0, "C", false, 0, "")
	fontStyle(pdf, "", 8, 0)
	pdf.Ln(5)

	return pdf
}

// Description: genera linea de corte reutilizable del recibo de pago
func separador(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	fontStyle(pdf, "", 8, 70)
	pdf.CellFormat(201.9, 5, "...........................................................................................................................Doblar...........................................................................................................................", "", 0, "TC", false, 0, "")
	fontStyle(pdf, "", 8, 0)
	pdf.Ln(5)
	return pdf
}

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
