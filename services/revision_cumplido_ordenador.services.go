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

func ObtenerCumplidosPendientesOrdenador(documento_ordenador string) (cambios_estado []models.CambioEstadoCumplido, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ObtenerCumplidosPendientesOrdenador", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var respuesta_peticion map[string]interface{}
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/?query=DocumentoResponsable:"+documento_ordenador+",EstadoCumplidoId.CodigoAbreviacion:PRO,Activo:true", &respuesta_peticion); (err == nil) && (response == 200) {
		helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &cambios_estado)
		return cambios_estado, nil
	} else {
		outputError = map[string]interface{}{"funcion": "/ObtenerCumplidosPendientesOrdenador", "message": "Error al consultar los cumplidos pendeinetes de revision por el ordenador", "err": err, "status": "502"}
		return nil, outputError
	}
	return cambios_estado, nil
}

func ObtenerSolicitudesCumplidos(documento_ordenador string) (cumplidosInfo []models.SolicituRevisionCumplidoProveedor, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ObtenerSolicitudesCumplidos", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	cumplidos, err := ObtenerCumplidosPendientesOrdenador(documento_ordenador)
	if err == nil {
		for _, cumplido := range cumplidos {
			informacion_contrato_proveedor, err := helpers.ObtenerInformacionContratoProveedor(cumplido.CumplidoProveedorId.NumeroContrato, strconv.Itoa(cumplido.CumplidoProveedorId.VigenciaContrato))
			if err == nil {
				vigencia, _ := strconv.Atoi(informacion_contrato_proveedor[0].Vigencia)
				solicitudes_cumplido := models.SolicituRevisionCumplidoProveedor{
					TipoContrato:      informacion_contrato_proveedor[0].TipoContrato,
					NumeroContrato:    informacion_contrato_proveedor[0].NumeroContratoSuscrito,
					VigenciaContrato:  vigencia,
					Dependencia:       informacion_contrato_proveedor[0].NombreDependencia,
					NombreProveedor:   informacion_contrato_proveedor[0].NombreProveedor,
					Cdp:               informacion_contrato_proveedor[0].NumeroCdp,
					Rp:                informacion_contrato_proveedor[0].NumeroRp,
					VigenciaRP:        informacion_contrato_proveedor[0].VigenciaRp,
					CumplidoId:        cumplido.CumplidoProveedorId.Id,
					Activo:            cumplido.Activo,
					FechaCreacion:     cumplido.CumplidoProveedorId.FechaCreacion,
					FechaModificacion: cumplido.CumplidoProveedorId.FechaModificacion,
				}
				cumplidosInfo = append(cumplidosInfo, solicitudes_cumplido)
			} else {
				outputError = map[string]interface{}{"funcion": "/ObtenerCumplidosPendientesOrdenador", "message": "Error al consultar los cumplidos pendeinetes de revision por el ordenador", "err": err, "status": "502"}
				return nil, outputError
			}
		}
	} else {
		outputError = map[string]interface{}{"funcion": "/ObtenerCumplidosPendientesOrdenador", "message": "Error al consultar los cumplidos pendeinetes de revision por el ordenador", "err": err, "status": "502"}
		return nil, outputError
	}

	return cumplidosInfo, nil
}

func ListaCumplidosReversibles(documento_ordenador string) (soliciudes_revertibles []models.SolicituRevisionCumplidoProveedor, outputError interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ListaCumplidosReversibles", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	fechaActual := time.Now()
	fechaMenosQuinceDias := fechaActual.AddDate(0, 0, -15)
	fechaFormateada := fechaMenosQuinceDias.Format("01/02/2006")

	var respuesta_peticion map[string]interface{}
	var cumplidos []models.CambioEstadoCumplido

	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/?query=DocumentoResponsable:"+documento_ordenador+",EstadoCumplidoId.CodigoAbreviacion:AO,Activo:true,FechaModificacion__gte:"+fechaFormateada, &respuesta_peticion); err == nil && response == 200 {
		if respuesta_peticion["Data"] != nil {
			helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &cumplidos)
			if len(cumplidos) > 0 {
				for _, cumplido := range cumplidos {
					informacion_contrato_proveedor, err := helpers.ObtenerInformacionContratoProveedor(cumplido.CumplidoProveedorId.NumeroContrato, strconv.Itoa(cumplido.CumplidoProveedorId.VigenciaContrato))
					if err == nil {
						vigencia, _ := strconv.Atoi(informacion_contrato_proveedor[0].Vigencia)
						solicitudes_cumplido := models.SolicituRevisionCumplidoProveedor{
							TipoContrato:      informacion_contrato_proveedor[0].TipoContrato,
							NumeroContrato:    informacion_contrato_proveedor[0].NumeroContratoSuscrito,
							VigenciaContrato:  vigencia,
							Dependencia:       informacion_contrato_proveedor[0].NombreDependencia,
							NombreProveedor:   informacion_contrato_proveedor[0].NombreProveedor,
							Cdp:               informacion_contrato_proveedor[0].NumeroCdp,
							Rp:                informacion_contrato_proveedor[0].NumeroRp,
							VigenciaRP:        informacion_contrato_proveedor[0].VigenciaRp,
							CumplidoId:        cumplido.CumplidoProveedorId.Id,
							Activo:            cumplido.Activo,
							FechaCreacion:     cumplido.CumplidoProveedorId.FechaCreacion,
							FechaModificacion: cumplido.CumplidoProveedorId.FechaModificacion,
						}
						soliciudes_revertibles = append(soliciudes_revertibles, solicitudes_cumplido)
					} else {
						outputError = map[string]interface{}{"funcion": "/ListaCumplidosReversibles", "message": "Error al consultar los cumplidos pendeinetes de revision por el ordenador", "err": err, "status": "502"}
						return nil, outputError
					}
				}
			} else {
				outputError = map[string]interface{}{"funcion": "/ListaCumplidosReversibles", "message": "No hay cumplidos que se puedan revertir", "status": "502"}
				return nil, outputError
			}
		} else {
			outputError = map[string]interface{}{"funcion": "/ListaCumplidosReversibles", "message": "No hay cumplidos que se puedan revertir", "err": err, "status": "502"}
			return nil, outputError
		}

	}

	return soliciudes_revertibles, nil
}

func GenerarAutorizacionPago(id_solicitud_pago string) (datos_documento *models.DocuementoAutorizacionPago, errorOutput interface{}) {

	var respuesta_cambioEstado map[string]interface{}

	url_request := beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores") + "/cambio_estado_cumplido/?query=CumplidoProveedorId:" + id_solicitud_pago + ",EstadoCumplidoId.CodigoAbreviación:PRO,Activo:true"
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

		data := respuesta_cambioEstado["Data"]
		if len(data.([]interface{})) > 0 {
			helpers.LimpiezaRespuestaRefactor(respuesta_cambioEstado, &cambio_estado)
			if cambio_estado[0].EstadoCumplidoId == nil {
				return nil, nil
			}
		}

		if len(cambio_estado) < 0 {
			return nil, nil
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
