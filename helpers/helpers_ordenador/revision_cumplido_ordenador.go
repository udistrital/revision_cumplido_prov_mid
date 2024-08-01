package helpers_ordenador

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
	"strconv"
	"strings"
	"time"
)

func ObtenerNumerosDeContrato(documentoOrdenador string, estado string) (numerosContrato string, errorOutput interface{}) {

	defer func() {
		if err := recover(); err != nil {
			errorMessage := fmt.Sprintf("%v", err)
			errorOutput = map[string]interface{}{
				"Success": false,
				"Status":  400,
				"Message": "Error al consultar las dependencias: " + documentoOrdenador + " - " + errorMessage,
				"Error":   errorMessage,
			}
			fmt.Println("Error recuperado:", errorMessage)
		}
	}()

	var numeros_contrato []string
	var cambios_estado_limpios []models.CambioEstadoCumplido
	var cambios_estado map[string]interface{}
	var urlRequest = beego.AppConfig.String("UrlProveedoresCrud") + "/cambio_estado_cumplido/?query=DocumentoResponsable:" + documentoOrdenador + ",EstadoCumplidoId.Abreviacion:" + estado

	response, err := helpers.GetJsonWSO2Test(urlRequest, &cambios_estado)
	fmt.Println(response)
	if err != nil || response != 200 {
		errorMessage := fmt.Sprintf("%v", err)
		errorOutput = map[string]interface{}{
			"Success": false,
			"Status":  response,
			"Message": "Error al consultar Cumplidos para el proveedor." + errorMessage,
			"Error":   errorMessage,
		}
		logs.Error(err)
		return "", errorOutput
	}

	if len(cambios_estado["Data"].([]interface{})[0].(map[string]interface{})) != 0 {
		helpers.LimpiezaRespuestaRefactor(cambios_estado, &cambios_estado_limpios)

		for _, dependencia := range cambios_estado_limpios {

			numeros_contrato = append(numeros_contrato, dependencia.CumplidoProveedorId.NumeroContrato)
		}

	}
	numero_contrato_unidos := strings.Join(numeros_contrato, "|")
	return numero_contrato_unidos, nil
}

func ObtenerSolicitudesCumplidos(documento string, estado string) (contratos_list []models.Contrato, errorOutput interface{}) {

	var info_contratos []models.ContratoProveedor
	numero_contratos, e := ObtenerNumerosDeContrato(documento, estado)
	urlRequest := beego.AppConfig.String("UrlcrudAgora") + "/contrato_general/?query=ContratoSuscrito.NumeroContrato.in:" + numero_contratos

	defer func() {
		if err := recover(); err != nil {
			errorMessage := fmt.Sprintf("%v", err)
			errorOutput = map[string]interface{}{
				"Success": false,
				"Status":  400,
				"Message": "Error al consultar las dependencias: " + documento + " - " + errorMessage,
				"Error":   errorMessage,
			}
			fmt.Println("Error recuperado:", errorMessage)
		}
	}()

	if e != nil || numero_contratos == "" {
		return nil, e
	}

	response, err := helpers.GetJsonTest(urlRequest, &info_contratos)

	if err != nil || response != 200 {
		errorMessage := fmt.Sprintf("%v", err)
		errorOutput = map[string]interface{}{
			"Success": false,
			"Status":  400,
			"Message": "Error el consultar los contratros para el ordenador en .GetContratos" + errorMessage,
			"Error":   errorMessage,
		}
		logs.Error(err)
		return nil, errorOutput
	}

	for _, contrato := range info_contratos {

		var ultimoContrato = len(contrato.ContratoSuscrito) - 1
		if proveedor, err := ObtenerInfoProveedor(strconv.Itoa(contrato.Contratista)); err == nil && proveedor != nil {

			if proveedor != nil {
				contrato_disponibilidad, _ := ObtenerContratoDisponiblidad(contrato.ContratoSuscrito[ultimoContrato].NumeroContrato.Id)

				if contrato_disponibilidad != nil {

					cdprp, _ := ObtenerCrdp(strconv.Itoa(contrato_disponibilidad.NumeroCdp), strconv.Itoa(contrato_disponibilidad.Vigencia))
					if cdprp != nil {

						contrato := models.Contrato{
							TipoContrato:    contrato.TipoContrato.TipoContrato,
							NumeroContrato:  contrato.ContratoSuscrito[ultimoContrato].NumeroContrato.Id,
							Vigencia:        contrato.ContratoSuscrito[ultimoContrato].Vigencia,
							Dependencia:     contrato.DependenciaSolicitante,
							NombreProveedor: proveedor.NomProveedor,
							Cdp:             strconv.Itoa(contrato_disponibilidad.NumeroCdp),
							Rp:              cdprp.CDPNumeroDisponibilidad,
						}
						contratos_list = append(contratos_list, contrato)
					}
				}
			}

		}
	}

	return contratos_list, nil
}

func ListaCumplidosReversibles(documentoOrdenador string) (soliciudes []models.Contrato, errorOutput interface{}) {
	defer func() {
		if err := recover(); err != nil {
			errorMessage := fmt.Sprintf("%v", err)
			errorOutput = map[string]interface{}{
				"Success": false,
				"Status":  400,
				"Message": "Error al consultar las dependencias: " + documentoOrdenador + " - " + errorMessage,
				"Error":   errorMessage,
			}
			fmt.Println("Error recuperado:", errorMessage)
		}
	}()

	fechaActual := time.Now()
	fechaMenosQuinceDias := fechaActual.AddDate(0, 0, -15)
	fechaFormateada := fechaMenosQuinceDias.Format("01/02/2006")
	numeros_contrato, e := ObtenerNumerosDeContrato(documentoOrdenador, "AO,FechaCreacion__gt:"+fechaFormateada+",Activo:true")

	if e != nil || numeros_contrato == "" {
		return nil, e
	}

	var respuesta []models.ContratoProveedor
	urlRequest := beego.AppConfig.String("UrlcrudAgora") + "/contrato_general/?query=ContratoSuscrito.NumeroContrato.in:" + numeros_contrato
	println(urlRequest)
	response, err := helpers.GetJsonWSO2Test(urlRequest, &respuesta)

	if err != nil || response != 200 {
		errorMessage := fmt.Sprintf("%v", err)
		errorOutput = map[string]interface{}{
			"Success": false,
			"Status":  400,
			"Message": "Error el consultar los contratros para el ordenador en .GetContratos" + errorMessage,
			"Error":   errorMessage,
		}
		logs.Error(err)
		return nil, errorOutput
	}
	fmt.Println("ssdasdsadadasdasdsad")
	if respuesta == nil {

		return nil, nil

	}

	for _, contrato := range respuesta {

		var ultimoContrato = len(contrato.ContratoSuscrito) - 1
		if proveedor, err := ObtenerInfoProveedor(strconv.Itoa(contrato.Contratista)); err == nil && proveedor != nil {

			if proveedor != nil {
				fmt.Println("proverdor no es nulo")
				contrato_disponibilidad, _ := ObtenerContratoDisponiblidad(contrato.ContratoSuscrito[ultimoContrato].NumeroContrato.Id)

				if contrato_disponibilidad != nil {
					fmt.Println("contratoDisponibilidad no es nulo")
					cdprp, _ := ObtenerCrdp(strconv.Itoa(contrato_disponibilidad.NumeroCdp), strconv.Itoa(contrato_disponibilidad.Vigencia))

					if cdprp != nil {
						fmt.Println("proverdor no es nulo")
						contrato := models.Contrato{
							TipoContrato:       contrato.TipoContrato.TipoContrato,
							NumeroContrato:     contrato.ContratoSuscrito[ultimoContrato].NumeroContrato.Id,
							Vigencia:           contrato.ContratoSuscrito[ultimoContrato].Vigencia,
							Dependencia:        contrato.DependenciaSolicitante,
							NombreProveedor:    proveedor.NomProveedor,
							Cdp:                strconv.Itoa(contrato_disponibilidad.NumeroCdp),
							Rp:                 cdprp.CDPNumeroDisponibilidad,
							DocumentoOrdenador: documentoOrdenador,
						}
						soliciudes = append(soliciudes, contrato)
					}
				}
			}

		}
	}

	return soliciudes, nil

}

func ObtenerEstado(estado string) (Estado *models.EstadoCumplidoId, errorOutput interface{}) {

	defer func() {
		if err := recover(); err != nil {
			errorOutput = map[string]interface{}{
				"Success": false,
				"Status":  400,
				"Message": "Error al consultar las ObtenerEstado: " + estado,
				"Error":   err,
			}
		}
	}()

	var respuesta map[string]interface{}
	urlRequest := beego.AppConfig.String("UrlProveedoresCrud") + "/estado_cumplido?query=Abreviacion:" + estado
	println(urlRequest)
	response, err := helpers.GetJsonWSO2Test(urlRequest, &respuesta)

	if err != nil || response != 200 {
		errorMessage := fmt.Sprintf("%v", err)
		errorOutput = map[string]interface{}{
			"Success": false,
			"Status":  400,
			"Message": "Error el consultar los contratros para el ordenador en .GetContratos" + errorMessage,
			"Error":   errorMessage,
		}
		logs.Error(err)
		return nil, errorOutput
	}
	var estado_list []models.EstadoCumplidoId

	helpers.LimpiezaRespuestaRefactor(respuesta, &estado_list)
	if respuesta != nil {
		Estado = &estado_list[0]
		return Estado, nil

	}

	return nil, nil

}

func ObtenerInfoProveedor(IdProveedor string) (provedor *models.Proveedor, errorOutput interface{}) {

	defer func() {
		if err := recover(); err != nil {
			errorOutput = map[string]interface{}{
				"Success": false,
				"Status":  400,
				"Message": "Error al consultar las dependencias: " + IdProveedor,
				"Error":   err,
			}
		}
	}()

	var respuesta []models.Proveedor
	urlRequest := beego.AppConfig.String("UrlcrudAgora") + "/informacion_proveedor/?query=id:" + IdProveedor
	println(urlRequest)
	response, err := helpers.GetJsonWSO2Test(urlRequest, &respuesta)

	if err != nil || response != 200 {
		errorMessage := fmt.Sprintf("%v", err)
		errorOutput = map[string]interface{}{
			"Success": false,
			"Status":  400,
			"Message": "Error el consultar los contratros para el ordenador en .GetContratos" + errorMessage,
			"Error":   errorMessage,
		}
		logs.Error(err)
		return nil, errorOutput
	}

	if respuesta != nil {
		provedor = &respuesta[0]
		return provedor, nil
	}

	return nil, nil

}

func ObtenerContratoDisponiblidad(NumeroContrato string) (contrato_disponibilidad *models.ContratoDisponibilidad, errorOutput interface{}) {
	defer func() {
		if err := recover(); err != nil {
			errorOutput = map[string]interface{}{
				"Success": false,
				"Status":  400,
				"Message": "Error al consultar las dependencias: " + NumeroContrato,
				"Error":   err,
			}
		}
	}()

	var respuesta []models.ContratoDisponibilidad
	urlRequest := beego.AppConfig.String("UrlcrudAgora") + "/contrato_disponibilidad?query=NumeroContrato:" + NumeroContrato
	fmt.Println(urlRequest)
	response, err := helpers.GetJsonWSO2Test(urlRequest, &respuesta)

	if err != nil || response != 200 {
		errorMessage := fmt.Sprintf("%v", err)
		errorOutput = map[string]interface{}{
			"Success": false,
			"Status":  400,
			"Message": "Error el consultar los contratros para el ordenador en .GetContratos" + errorMessage,
			"Error":   errorMessage,
		}
		logs.Error(err)
		return nil, errorOutput
	}

	if respuesta != nil {
		println("respuesta", respuesta)
		contrato_disponibilidad = &respuesta[len(respuesta)-1]
		return contrato_disponibilidad, nil

	}
	return nil, nil

}

func ObtenerCrdp(Cdp string, Vigencia string) (crdp *models.CDPRP, errorOutput interface{}) {
	defer func() {
		if err := recover(); err != nil {
			errorMessage := fmt.Sprintf("%v", err)
			errorOutput = map[string]interface{}{
				"Success": false,
				"Status":  400,
				"Message": "Error al consultar las dependencias: " + Cdp + " - " + errorMessage,
				"Error":   errorMessage,
			}
			fmt.Println("Error recuperado:", errorMessage)
		}
	}()

	var respuesta map[string]interface{}
	urlRequest := beego.AppConfig.String("UrlFinancieraJBPM") + "/cdprp/" + Cdp + "/" + Vigencia + "/01"
	println(urlRequest)

	response, err := helpers.GetJsonWSO2Test(urlRequest, &respuesta)

	if err != nil || response != 200 {
		errorMessage := fmt.Sprintf("%v", err)
		errorOutput = map[string]interface{}{
			"Success": false,
			"Status":  400,
			"Message": "Error el consultar los contratros para el ordenador en .GetContratos" + errorMessage,
			"Error":   errorMessage,
		}
		logs.Error(err)
		return nil, errorOutput
	}

	if respuesta != nil {

		cdprp_array, ok := respuesta["cdpxrp"].(map[string]interface{})["cdprp"].([]interface{})

		if !ok {
			return nil, errorOutput
		}
		if len(cdprp_array) > 0 {

			ultimoElemento := cdprp_array[len(cdprp_array)-1].(map[string]interface{})
			crdp = &models.CDPRP{
				CDPNumeroDisponibilidad: ultimoElemento["CDP_NUMERO_DISPONIBILIDAD"].(string),
				RPVigencia:              ultimoElemento["RP_VIGENCIA"].(string),
				CDPVigencia:             ultimoElemento["CDP_VIGENCIA"].(string),
			}

			fmt.Println(respuesta)
			return crdp, nil
		}

	}

	return nil, nil

}

func GenerarAutorizacion(id_solicitud_pago string) (datos_documento *models.DocuementoAutorizacionPago, errorOutput interface{}) {

	// Obtiene datos de cambio estado
	var respuesta_cambioEstado map[string]interface{}
	url_request := beego.AppConfig.String("UrlProveedoresCrud") + "/cambio_estado_cumplido/?query=CumplidoProveedorId:" + id_solicitud_pago + ",EstadoCumplidoId.Abreviacion:AO,Activo:true"
	response, err := helpers.GetJsonWSO2Test(url_request, &respuesta_cambioEstado)
	var cambio_estado []models.CambioEstadoCumplido

	if err != nil || response != 200 {
		errorMessage := fmt.Sprintf("%v", err)
		errorOutput = map[string]interface{}{
			"Success": false,
			"Status":  response,
			"Message": "Error al consultar Cumplidos para el proveedor. " + errorMessage,
			"Error":   errorMessage,
		}
		logs.Error(err)
		return nil, errorOutput
	}

	if respuesta_cambioEstado["Data"] != nil {

		fmt.Println(len(cambio_estado))

		data := respuesta_cambioEstado["Data"]
		if len(data.([]interface{})) > 0 {
			helpers.LimpiezaRespuestaRefactor(respuesta_cambioEstado, &cambio_estado)
			if cambio_estado[0].EstadoCumplidoId == nil {
				return nil, nil
			}
		}

		if len(cambio_estado) < 0 {
			return nil, nil
			fmt.Println("entro por que esta vacio?")
		}

		// Obtiene información de los contratos
		var respuesta []models.ContratoProveedor
		url_request_contrato := beego.AppConfig.String("UrlcrudAgora") + "/contrato_general/?query=ContratoSuscrito.NumeroContrato:" + cambio_estado[0].CumplidoProveedorId.NumeroContrato
		resonse_contrato, err_contrato := helpers.GetJsonWSO2Test(url_request_contrato, &respuesta)

		if err_contrato != nil || resonse_contrato != 200 {
			errorMessage := fmt.Sprintf("%v", err)
			errorOutput = map[string]interface{}{
				"Success": false,
				"Status":  resonse_contrato,
				"Message": "Error al consultar Contrato para el proveedor. " + errorMessage,
				"Error":   errorMessage,
			}
			logs.Error(err)
			return nil, errorOutput
		}

		// Obtiene info Proveedor
		if proveedor, err := ObtenerInfoProveedor(strconv.Itoa(respuesta[0].Contratista)); err == nil && proveedor != nil {
			var info_ordenador []models.Proveedor
			url_request_ordenador := beego.AppConfig.String("UrlcrudAgora") + "/informacion_proveedor/?query=NumDocumento:" + strconv.Itoa(cambio_estado[0].DocumentoResponsable)
			response_ordenador, errOrdenador := helpers.GetJsonWSO2Test(url_request_ordenador, &info_ordenador)

			if errOrdenador != nil || response_ordenador != 200 {
				errorMessage := fmt.Sprintf("%v", err)
				errorOutput = map[string]interface{}{
					"Success": false,
					"Status":  response_ordenador,
					"Message": "Error al consultar Información del proveedor. " + errorMessage,
					"Error":   errorMessage,
				}
				logs.Error(err)
				return nil, errorOutput
			}

			var respuesta_documentos map[string]interface{}
			url_request_documentos := beego.AppConfig.String("UrlProveedoresCrud") + "/soporte_pago?query=CumplidoProveedorId.id:" + id_solicitud_pago
			responseDocuementos, error_documentos := helpers.GetJsonWSO2Test(url_request_documentos, &respuesta_documentos)
			fmt.Println(url_request_documentos)
			var documentosCargados []models.SoportePago
			if len(respuesta_documentos["Data"].([]interface{})[0].(map[string]interface{})) != 0 {
				helpers.LimpiezaRespuestaRefactor(respuesta_documentos, &documentosCargados)
			}

			var id_documentos []string
			for _, documentosCargado := range documentosCargados {
				id_documentos = append(id_documentos, strconv.Itoa(documentosCargado.DocumentoId))
			}
			id_documentos_unidos := strings.Join(id_documentos, "|")

			if error_documentos != nil || responseDocuementos != 200 {

				return nil, errorOutput
			}

			var documentos []models.Documento

			url_request_documentos_destion := "http://pruebasapi.intranetoas.udistrital.edu.co:8094/v1/documento/?query=Id.in:" + id_documentos_unidos + ",Activo:true&limit=-1"
			response_docuementos_gestion, erro_documentos_gestion := helpers.GetJsonTest(url_request_documentos_destion, &documentos)

			if erro_documentos_gestion != nil || response_docuementos_gestion != 200 {
				return nil, errorOutput
			}

			var lista_documentos_cargados_strings []string
			for _, documento := range documentos {
				lista_documentos_cargados_strings = append(lista_documentos_cargados_strings, documento.TipoDocumento.CodigoAbreviacion)
			}

			indexRespuestaOrdenador := len(respuesta) - 1
			datos_documento := &models.DocuementoAutorizacionPago{
				NombreOrdenador:    info_ordenador[indexRespuestaOrdenador].NomProveedor,
				DocumentoOrdenador: info_ordenador[indexRespuestaOrdenador].NumDocumento,
				NombreProveedor:    proveedor.NomProveedor,
				DocumentoProveedor: proveedor.NumDocumento,
				DocumentosCargados: lista_documentos_cargados_strings,
			}

			return datos_documento, nil

		}
		return nil, nil
	} else {
		return nil, nil
	}
}
