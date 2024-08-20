package services

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/jung-kurt/gofpdf"

	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtenerInformacionProveedor(numero_contrato_suscrito string, vigencia string) (informacion_proveedor models.InformacionContratoProveedor, outputError map[string]interface{}) {

	contrato_contratista, err := helpers.ObtenerInformacionContratoContratista(numero_contrato_suscrito, vigencia)
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

	if contratos_persona, outputError := helpers.ObtenerContratosPersona(numero_documento); outputError == nil {
		//fmt.Println("contratos_persona", contratos_persona)
		for _, contrato_persona := range contratos_persona.ContratosPersonas.ContratoPersona {
			var contrato models.InformacionContrato
			contrato, outputError = helpers.ObtenerInformacionContrato(contrato_persona.NumeroContrato, contrato_persona.Vigencia)

			if (contrato == models.InformacionContrato{} || outputError != nil) {
				continue
			}
			var informacion_contrato_contratista models.InformacionContratoContratista
			informacion_contrato_contratista, outputError = helpers.ObtenerInformacionContratoContratista(contrato_persona.NumeroContrato, contrato_persona.Vigencia)
			// se llena el contrato original en el indice 0

			if cdprp, outputError := helpers.ObtenerRP(contrato_persona.NumeroCDP, contrato_persona.Vigencia); outputError == nil {
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

func ObtenerBanco(banco_id int) (banco models.Banco, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			//fmt.Println("error", err)
			outputError = map[string]interface{}{"funcion": "/GetBanco", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var respuesta_peticion map[string]interface{}
	var respuesta_banco models.Banco
	if response, error := helpers.GetJsonTest(beego.AppConfig.String("UrlCoreApi")+"/banco/"+strconv.Itoa(banco_id), &respuesta_peticion); error == nil && response == 200 {
		json_banco, err := json.Marshal(respuesta_peticion)
		if err == nil {
			if err := json.Unmarshal(json_banco, &respuesta_banco); err != nil {
				outputError = map[string]interface{}{"funcion": "/GetBanco", "err": err, "status": "502"}
				return respuesta_banco, outputError
			}
		} else {
			outputError = map[string]interface{}{"funcion": "/GetBanco", "err": err, "status": "502"}
			return respuesta_banco, outputError
		}
	} else {
		outputError = map[string]interface{}{"funcion": "/GetBanco", "err": error, "status": "502"}
		return respuesta_banco, outputError
	}
	return respuesta_banco, nil
}

func CrearInformeSatisfaccion(numero_contrato_suscrito int, vigencia_contrato string, tipo_pago string, periodo_inicio string, periodo_fin string, tipo_factura string, numero_cuenta_factura string, valor_pagar int, tipo_cuenta string, numero_cuenta string, banco_id int) (informe_seguimiento models.InformeSeguimiento, outputError map[string]interface{}) {

	var valor_total_contrato int
	var saldo_contrato int
	balance, err := ObtenerBalanceFinancieroContrato(strconv.Itoa(numero_contrato_suscrito), vigencia_contrato)
	if err == nil {
		valor := strings.Split(balance.TotalContrato, ".")[0]
		valor_total_contrato, _ = strconv.Atoi(valor)
		saldo_contrato, _ = strconv.Atoi(balance.Saldo)
	}

	info_contrato, err := helpers.ObtenerInformacionContrato(strconv.Itoa(numero_contrato_suscrito), vigencia_contrato)
	if err != nil {
		outputError = map[string]interface{}{"funcion": "/CreateInformeSeguimiento", "err": err, "status": "502"}
		return informe_seguimiento, outputError
	}

	supervisor := info_contrato.Contrato.Supervisor.Nombre

	nombre_banco, error := ObtenerBanco(banco_id)
	if error != nil {
		outputError = map[string]interface{}{"funcion": "/CreateInformeSeguimiento", "err": error, "status": "502"}
		return informe_seguimiento, outputError
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(25, 20, 25)
	pdf.SetAutoPageBreak(true, 20)

	pdf.SetHeaderFunc(func() {
		//pdf = headerInformeSatisfaccion(pdf)
		pdf = helpers.Header(pdf, "CUMPLIDO A SATISFACCIÓN POR PARTE DE LA DEPENDENCIA", "Gestión Contractual", "GC-PR-003-FR-012", "05", "13/10/2021")

		pdf.Ln(30)
	})

	pdf.AddPage()

	vigencia, _ := strconv.Atoi(vigencia_contrato)
	contrato, err := ObtenerInformacionProveedor(strconv.Itoa(numero_contrato_suscrito), vigencia_contrato)
	if err != nil {
		outputError = map[string]interface{}{"funcion": "/CreateInformeSeguimiento", "err": err, "status": "502"}
		return informe_seguimiento, outputError
	}

	acta_inicio, err := ObtenerActaInicio(strconv.Itoa(numero_contrato_suscrito), vigencia)
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
			VerificarJefe(info_contrato))
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
			VerificarJefe(info_contrato))
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
		nombre_banco.NombreBanco)

	pdf = footerInformeSeguimiento(pdf,
		contrato.NumeroContratoSuscrito,
		"08/04/2024", //formatear_fecha(acta_inicio.FechaInicio)
		contrato.TipoContrato,
		numero_cuenta_factura,
		supervisor,
		VerificarJefe(info_contrato),
		contrato.Vigencia,
		contrato.NombreDependencia)

	// Crear el PDF

	// Codificar el PDF
	encodedFile := encodePDF(pdf)
	nombre := "prueba"

	informe_seguimiento = models.InformeSeguimiento{File: encodedFile, Archivo: nombre}
	return informe_seguimiento, nil
}

func body_primera_parte(pdf *gofpdf.Fpdf, dependencia string, nombre_proveedor string, numero_nit string, cumplimiento_contrato bool, tipo_contrato string, fecha_inicio string, numero_contrato string, cdp string, vigencia_cdp string, crp string, vigencia_crp string, cargo string) *gofpdf.Fpdf {

	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetFont("Times", "", 10)

	pdf.CellFormat(0, 8, tr("UNIVERSIDAD DISTRITAL FRANCISCO JOSÉ DE CALDAS"), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 8, tr("("+dependencia+")"), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 8, tr(fmt.Sprintf(`En ejercicio de las funciones de (%s)`, cargo)), "", 1, "C", false, 0, "")
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
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`Que el contratista %s identificado con NIT %s cumplió %s a satisfacción con las obligaciones y objeto del %s Nro. %s de fecha %s garantizada y perfeccionada con Certificado de Disponibilidad Presupuestal No. %s de %s y Certificado de Registro Presupuestal No. %s de %s.`, nombre_proveedor, numero_nit, cumplimiento, tipo_contrato, numero_contrato, fecha_inicio, cdp, vigencia_cdp, crp, vigencia_crp)), "", "", false)
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

func body_segunda_parte(pdf *gofpdf.Fpdf, tipo_factura string, numero_cuenta_factura string, valor_total_contrato int, periodo_inicio string, periodo_fin string, saldo_contrato int, fecha_inicio string, fecha_fin string, tipo_cuenta string, numero_cuenta string, nombre_banco string) *gofpdf.Fpdf {

	pdf.SetFont("Times", "", 10)

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetMargins(30, 30, 30)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`Que el valor causado de conformidad con la %s de Venta o Cuenta de cobro No. %s es %s PESOS $%s pesos m/cte.`, tipo_factura, numero_cuenta_factura, strings.ToUpper(helpers.ValorLetras(valor_total_contrato)), helpers.FormatNumber(valor_total_contrato, 0, ".", ","))), "", "J", false)
	pdf.Ln(5)
	pdf.SetMargins(30, 30, 30)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`Que el valor total del contrato corresponde a %s $%s pesos m/cte.`, strings.ToUpper(helpers.ValorLetras(valor_total_contrato)), helpers.FormatNumber(valor_total_contrato, 0, ".", ","))), "", "J", false)
	pdf.Ln(5)
	pdf.SetMargins(30, 30, 30)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`Que el presente pago corresponde al período de %s a %s de ejecución parcial total o único pago del contrato.`, periodo_inicio, periodo_fin)), "", "J", false)
	pdf.Ln(5)
	pdf.SetMargins(30, 30, 30)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`Quedando un saldo correspondiente a $%s pesos m/cte.`, helpers.FormatNumber(saldo_contrato, 0, ".", ","))), "", "", false)
	pdf.Ln(5)
	pdf.SetMargins(30, 30, 30)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`Que el presente pago se encuentra en cumplimiento dentro del tiempo de ejecución del contrato del %s al %s.`, fecha_inicio, fecha_fin)), "", "J", false)

	pdf.SetMargins(30, 30, 30)
	pdf.MultiCell(0, 8, "", "", "J", false)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`Que tal valor debe girarse por petición del contratista a la Cuenta %s  No. %s del Banco %s.`, tipo_cuenta, numero_cuenta, strings.ToUpper(nombre_banco))), "", "J", false)

	return pdf
}

func footerInformeSeguimiento(pdf *gofpdf.Fpdf, contrato_suscrito string, fecha_inicio string, tipo_contrato string, numero_factura string, supervisor string, jefe string, vigencia string, dependencia string) *gofpdf.Fpdf {

	dia := time.Now().Day()
	mes := int(time.Now().Month())
	año := time.Now().Year()

	pdf.SetFont("Times", "", 10)

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetMargins(30, 0, 30)
	pdf.MultiCell(0, 8, "", "", "J", false)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`Con el presente cumplido y de acuerdo a lo establecido en los numerales 32 y 33 del Artículo 18° de la Resolución de Rectoría No. 629 de 2016- Manual de Interventoría y Supervisión certifico que los informes físicos técnicos financieros y administrativos sobre el avance de la ejecución del objeto contractual reposan en el expediente del %s No. %s de %s. De igual forma certifico que se verificaron las condiciones y elementos que hacen parte de la(s) factura(s) No. %s acorde con lo establecido en la ficha técnica del proceso en mención garantizando la calidad del bien o servicio adquirido por la Universidad.`, tipo_contrato, contrato_suscrito, fecha_inicio, numero_factura)), "", "J", false)
	pdf.Ln(15)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`La presente se expide a los %s días del mes %s de %s.`, strconv.Itoa(dia), ObtenerMes(mes), strconv.Itoa(año))), "", "J", false)
	pdf.Ln(5) // Espacio para la firma
	pdf.SetMargins(30, 30, 30)
	pdf.SetFont("Times", "B", 10)
	pdf.MultiCell(0, 8, tr(`__________________________`), "", "", false)
	pdf.MultiCell(0, 8, tr(`NOMBRE`), "", "", false)
	pdf.MultiCell(0, 8, tr(supervisor), "", "", false)
	pdf.MultiCell(0, 8, tr(`C.C ______________ de Bogotá`), "", "", false)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`CARGO %s`, jefe)), "", "", false)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`DEPENDENCIA %s`, dependencia)), "", "", false)
	pdf.MultiCell(0, 8, tr(fmt.Sprintf(`Supervisor %s Contrato/ Contrato de Comisión/Orden de Compra/Orden de Servicio/ Orden de Compra CCE) No. %s de %s.`, supervisor, contrato_suscrito, vigencia)), "", "", false)

	return pdf

}

func VerificarJefe(info_contrato models.InformacionContrato) string {
	cargo := strings.Fields(strings.ToLower(info_contrato.Contrato.Supervisor.Cargo))
	if cargo[0] == "jefe" {
		return info_contrato.Contrato.Supervisor.Cargo
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

func ObtenerBalanceFinancieroContrato(numero_contrato_suscrito string, vigencia_contrato string) (balance_contrato models.BalanceContrato, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetBalanceContrato", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	informacion_contrato, err := helpers.ObtenerInformacionContratoContratista(numero_contrato_suscrito, vigencia_contrato)
	if err == nil {
		informacion_contratista, err := helpers.ObtenerContratosContratista(informacion_contrato.InformacionContratista.Documento.Numero)
		if err == nil {
			contrato, err := helpers.ObtenerInformacionContrato(numero_contrato_suscrito, vigencia_contrato)
			if err == nil {
				valor_girado, err := ObtenerValorGiradoPorCdp(informacion_contratista[0].NumeroCdp, informacion_contratista[0].VigenciaCdp, contrato.Contrato.UnidadEjecutora)
				if err == nil {
					total_contrato, err := strconv.ParseFloat(informacion_contrato.InformacionContratista.ValorContrato, 64)
					if err == nil {
						saldo_contrato := int(total_contrato) - valor_girado
						balance_contrato.TotalContrato = informacion_contrato.InformacionContratista.ValorContrato
						balance_contrato.Saldo = strconv.Itoa(saldo_contrato)
						return balance_contrato, nil
					} else {
						outputError = map[string]interface{}{"funcion": "/GetBalanceContrato", "err": err, "status": "502"}
						return balance_contrato, outputError
					}

				} else {
					outputError = map[string]interface{}{"funcion": "/GetBalanceContrato", "err": err, "status": "502"}
					return balance_contrato, outputError
				}
			}
		} else {
			outputError = map[string]interface{}{"funcion": "/GetBalanceContrato", "err": err, "status": "502"}
			return balance_contrato, outputError
		}
	} else {
		outputError = map[string]interface{}{"funcion": "/GetBalanceContrato", "err": err, "status": "502"}
		return balance_contrato, outputError
	}

	return balance_contrato, outputError
}

func ObtenerValorGiradoPorCdp(cdp string, vigencia_cdp string, unidad_ejecucion string) (valor_girado int, err error) {
	var temp_giros_tercero map[string]interface{}
	var giros_tercero models.GirosTercero
	valor_girado = 0
	fmt.Println(beego.AppConfig.String("UrlFinancieraJBPM") + "/" + "giros_tercero/" + cdp + "/" + vigencia_cdp + "/" + unidad_ejecucion)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlFinancieraJBPM")+"/giros_tercero/"+cdp+"/"+vigencia_cdp+"/"+unidad_ejecucion, &temp_giros_tercero); (err == nil) && (response == 200) {
		json_giros_tercero, error_json := json.Marshal(temp_giros_tercero)
		if error_json == nil {
			if err := json.Unmarshal(json_giros_tercero, &giros_tercero); err == nil {
				//fmt.Println("giros "+cdp, giros_tercero)
				for _, giro := range giros_tercero.Giros.Tercero {
					total_girado, err := strconv.Atoi(giro.ValorBrutoGirado)
					//fmt.Println(total_girado)
					if err == nil {
						valor_girado = valor_girado + total_girado
					}
				}
				//fmt.Println(valor_girado)
				return valor_girado, nil

			} else {
				err = errors.New("Error Unmarshal giros_tercero")
				return valor_girado, err
			}

		} else {
			err = errors.New("Error Marshal giros_tercero")
			return valor_girado, err
		}

	} else {
		return valor_girado, err
	}
	return
}
