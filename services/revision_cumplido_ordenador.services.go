package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtenerCumplidosOrdenador(docuemento_ordenador string, estado string) (cambios_estado_limpios []models.CambioEstadoCumplido, errorOutput interface{}) {

	defer func() {
		if err := recover(); err != nil {
			errorMessage := fmt.Sprintf("%v", err)
			errorOutput = map[string]interface{}{
				"Success": false,
				"Status":  400,
				"Message": "Error al consultar las dependencias: " + docuemento_ordenador + " - " + errorMessage,
				"Error":   errorMessage,
			}
			fmt.Println("Error recuperado:", errorMessage)
		}
	}()

	var cambios_estado map[string]interface{}

	var urlRequest = beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores") + "/cambio_estado_cumplido/?query=DocumentoResponsable:" + docuemento_ordenador + ",EstadoCumplidoId.CodigoAbreviación:" + estado

	response, err := helpers.GetJsonWSO2Test(urlRequest, &cambios_estado)
	//fmt.Println(response)
	if err != nil || response != 200 {
		errorMessage := fmt.Sprintf("%v", err)
		errorOutput = map[string]interface{}{
			"Success": false,
			"Status":  response,
			"Message": "Error al consultar Cumplidos para el proveedor." + errorMessage,
			"Error":   errorMessage,
		}
		logs.Error(err)
		return nil, errorOutput
	}

	if len(cambios_estado["Data"].([]interface{})[0].(map[string]interface{})) != 0 {
		helpers.LimpiezaRespuestaRefactor(cambios_estado, &cambios_estado_limpios)

	}
	return cambios_estado_limpios, nil
}

func ObtenerSolicitudesCumplidos(documento string, estado string) (cumplidosInfo []models.SolicituRevisionCumplidoProveedor, errorOutput interface{}) {
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

	cumplidos, err := ObtenerCumplidosOrdenador(documento, estado)
	if err != nil {
		errorOutput = map[string]interface{}{
			"Success": false,
			"Status":  400,
			"Message": "Error al consultar los cumplidos para el proveedor en .ObtenerCumplidosPorEstado",
			"Error":   err,
		}
		return nil, errorOutput
	}

	for _, cumplido := range cumplidos {
		info_contrato, err := helpers.ObtenerInformacionContratoProveedor(cumplido.CumplidoProveedorId.NumeroContrato, strconv.Itoa(cumplido.CumplidoProveedorId.VigenciaContrato))
		if err != nil {

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

		if proveedor, err := helpers.ObtenerContratosProveedor(info_contrato.InformacionContratista.Documento.Numero); err == nil && proveedor != nil {

			var vigencia, _ = strconv.Atoi(proveedor[0].Vigencia)
			contrato := models.SolicituRevisionCumplidoProveedor{
				TipoContrato:     proveedor[0].TipoContrato,
				NumeroContrato:   proveedor[0].NumeroContratoSuscrito,
				VigenciaContrato: vigencia,
				Dependencia:      proveedor[0].NombreDependencia,
				NombreProveedor:  proveedor[0].NombreProveedor,
				Cdp:              proveedor[0].NumeroCdp,
				Rp:               proveedor[0].NumeroRp,
				VigenciaRP:       proveedor[0].VigenciaRp,
				Id:               cumplido.CumplidoProveedorId.Id,
				Activo:           cumplido.Activo,
			}
			cumplidosInfo = append(cumplidosInfo, contrato)
		}
	}
	return cumplidosInfo, nil
}

func ListaCumplidosReversibles(docuemento_ordenador string) (soliciudes []models.SolicituRevisionCumplidoProveedor, errorOutput interface{}) {
	defer func() {
		if err := recover(); err != nil {
			errorMessage := fmt.Sprintf("%v", err)
			errorOutput = map[string]interface{}{
				"Success": false,
				"Status":  400,
				"Message": "Error al consultar las dependencias: " + docuemento_ordenador + " - " + errorMessage,
				"Error":   errorMessage,
			}
			fmt.Println("Error recuperado:", errorMessage)
		}
	}()

	fechaActual := time.Now()
	fechaMenosQuinceDias := fechaActual.AddDate(0, 0, -15)
	fechaFormateada := fechaMenosQuinceDias.Format("01/02/2006")
	cumplidos, e := ObtenerCumplidosOrdenador(docuemento_ordenador, "AO,FechaCreacion__gt:"+fechaFormateada+",Activo:true")

	if e != nil || cumplidos == nil {
		return nil, e
	}
	for _, cumplido := range cumplidos {
		info_contrato, err := helpers.ObtenerInformacionContratoProveedor(cumplido.CumplidoProveedorId.NumeroContrato, strconv.Itoa(cumplido.CumplidoProveedorId.VigenciaContrato))
		if err != nil {

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

		if err != nil {
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

		if proveedor, err := helpers.ObtenerContratosProveedor(info_contrato.InformacionContratista.Documento.Numero); err == nil && proveedor != nil {

			var vigencia, _ = strconv.Atoi(proveedor[0].Vigencia)
			contrato := models.SolicituRevisionCumplidoProveedor{
				TipoContrato:     proveedor[0].TipoContrato,
				NumeroContrato:   proveedor[0].NumeroContratoSuscrito,
				VigenciaContrato: vigencia,
				Dependencia:      proveedor[0].NombreDependencia,
				NombreProveedor:  proveedor[0].NombreProveedor,
				Cdp:              proveedor[0].NumeroCdp,
				Rp:               proveedor[0].NumeroRp,
				VigenciaRP:       proveedor[0].VigenciaRp,
				Id:               cumplido.CumplidoProveedorId.Id,
				Activo:           cumplido.Activo,
			}
			soliciudes = append(soliciudes, contrato)
		}
	}
	return soliciudes, nil
}

func ObtenerEstadoCumplido(estado string) (Estado *models.EstadoCumplidoId, errorOutput interface{}) {

	defer func() {
		if err := recover(); err != nil {
			errorOutput = map[string]interface{}{
				"Success": false,
				"Status":  400,
				"Message": "Error al consultar las ObtenerEstadoCumplido: " + estado,
				"Error":   err,
			}
		}
	}()

	var respuesta map[string]interface{}
	urlRequest := beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores") + "/estado_cumplido?query=CodigoAbreviación:" + estado
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

func GenerarAutorizacionPago(id_solicitud_pago string) (datos_documento *models.DocuementoAutorizacionPago, errorOutput interface{}) {

	// Obtiene datos de cambio estado
	var respuesta_cambioEstado map[string]interface{}

	url_request := beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores") + "/cambio_estado_cumplido/?query=CumplidoProveedorId:" + id_solicitud_pago + ",EstadoCumplidoId.CodigoAbreviación:PRO,Activo:true"
	response, err := helpers.GetJsonWSO2Test(url_request, &respuesta_cambioEstado)
	var cambio_estado []models.CambioEstadoCumplido
	//fmt.Println(url_request)
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

		//fmt.Println(len(cambio_estado))

		data := respuesta_cambioEstado["Data"]
		if len(data.([]interface{})) > 0 {
			helpers.LimpiezaRespuestaRefactor(respuesta_cambioEstado, &cambio_estado)
			if cambio_estado[0].EstadoCumplidoId == nil {
				return nil, nil
			}
		}

		if len(cambio_estado) < 0 {
			return nil, nil
			//fmt.Println("entro por que esta vacio?")
		}

		// Obtiene información de los contratos
		var respuesta []models.ContratoGeneral
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

			url_request_documentos := beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores") + "/soporte_cumplido?query=CumplidoProveedorId.id:" + id_solicitud_pago
			responseDocuementos, error_documentos := helpers.GetJsonWSO2Test(url_request_documentos, &respuesta_documentos)
			//fmt.Println(url_request_documentos)
			var documentosCargados []models.SoporteCumplido
			if respuesta_documentos["Data"] != nil {
				if len(respuesta_documentos["Data"].([]interface{})[0].(map[string]interface{})) != 0 {
					helpers.LimpiezaRespuestaRefactor(respuesta_documentos, &documentosCargados)
				}

			} else {
				return nil, errorOutput
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

			//fmt.Println("documentos", lista_documentos_cargados_strings)
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
	//println(urlRequest)
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
