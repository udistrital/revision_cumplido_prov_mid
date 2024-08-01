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

	var respuesta map[string]interface{}
	var urlRequest = beego.AppConfig.String("UrlProveedoresCrud") + "/cambio_estado_cumplido/?query=DocumentoResponsable:" + documentoOrdenador + ",EstadoCumplidoId.Abreviacion:" + estado
	fmt.Println(urlRequest)
	response, err := helpers.GetJsonWSO2Test(urlRequest, &respuesta)
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

	var dependenciasString []string
	var dependenciasList []models.CambioEstadoCumplido
	if len(respuesta["Data"].([]interface{})[0].(map[string]interface{})) != 0 {
		helpers.LimpiezaRespuestaRefactor(respuesta, &dependenciasList)

		for _, dependencia := range dependenciasList {

			dependenciasString = append(dependenciasString, dependencia.CumplidoProveedorId.NumeroContrato)
		}

	}
	listaString := strings.Join(dependenciasString, "|") + "|"
	return listaString, nil
}

func ObtenerSolicitudesCumplidos(documento string, estado string) (ContratosList []models.Contrato, errorOutput interface{}) {

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

	cumplidosPendientes, e := ObtenerNumerosDeContrato(documento, estado)

	if e != nil || cumplidosPendientes == "" {
		return nil, e
	}

	var respuesta []models.ContratoProveedor
	urlRequest := beego.AppConfig.String("UrlcrudAgora") + "/contrato_general/?query=ContratoSuscrito.NumeroContrato.in:" + cumplidosPendientes
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

	for _, contrato := range respuesta {

		var ultimoContrato = len(contrato.ContratoSuscrito) - 1
		if proveedor, err := ObtenerInfoProveedor(strconv.Itoa(contrato.Contratista)); err == nil && proveedor != nil {

			if proveedor != nil {
				fmt.Println("proverdor no es nulo")
				contratoDisponibilidad, _ := ObtenerContratoDisponiblidad(contrato.ContratoSuscrito[ultimoContrato].NumeroContrato.Id)

				if contratoDisponibilidad != nil {
					fmt.Println("contratoDisponibilidad no es nulo")
					cdprp, _ := ObtenerCrdp(strconv.Itoa(contratoDisponibilidad.NumeroCdp), strconv.Itoa(contratoDisponibilidad.Vigencia))

					if cdprp != nil {
						fmt.Println("proverdor no es nulo")
						contrato := models.Contrato{
							TipoContrato:    contrato.TipoContrato.TipoContrato,
							NumeroContrato:  contrato.ContratoSuscrito[ultimoContrato].NumeroContrato.Id,
							Vigencia:        contrato.ContratoSuscrito[ultimoContrato].Vigencia,
							Dependencia:     contrato.DependenciaSolicitante,
							NombreProveedor: proveedor.NomProveedor,
							Cdp:             strconv.Itoa(contratoDisponibilidad.NumeroCdp),
							Rp:              cdprp.CDPNumeroDisponibilidad,
						}
						ContratosList = append(ContratosList, contrato)
					}
				}
			}

		}
	}

	return ContratosList, nil
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
	cumplidos, e := ObtenerNumerosDeContrato(documentoOrdenador, "AO,FechaCreacion__gt:"+fechaFormateada+",Activo:true")
	fmt.Println(fechaFormateada)
	if e != nil || cumplidos == "" {
		return nil, e
	}

	var respuesta []models.ContratoProveedor
	urlRequest := beego.AppConfig.String("UrlcrudAgora") + "/contrato_general/?query=ContratoSuscrito.NumeroContrato.in:" + cumplidos
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
				contratoDisponibilidad, _ := ObtenerContratoDisponiblidad(contrato.ContratoSuscrito[ultimoContrato].NumeroContrato.Id)

				if contratoDisponibilidad != nil {
					fmt.Println("contratoDisponibilidad no es nulo")
					cdprp, _ := ObtenerCrdp(strconv.Itoa(contratoDisponibilidad.NumeroCdp), strconv.Itoa(contratoDisponibilidad.Vigencia))

					if cdprp != nil {
						fmt.Println("proverdor no es nulo")
						contrato := models.Contrato{
							TipoContrato:       contrato.TipoContrato.TipoContrato,
							NumeroContrato:     contrato.ContratoSuscrito[ultimoContrato].NumeroContrato.Id,
							Vigencia:           contrato.ContratoSuscrito[ultimoContrato].Vigencia,
							Dependencia:        contrato.DependenciaSolicitante,
							NombreProveedor:    proveedor.NomProveedor,
							Cdp:                strconv.Itoa(contratoDisponibilidad.NumeroCdp),
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

func ObtenerInfoProveedor(IdProveedor string) (Provedor *models.Proveedor, errorOutput interface{}) {

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
		Provedor = &respuesta[0]
		return Provedor, nil
	}

	return nil, nil

}

func ObtenerContratoDisponiblidad(NumeroContrato string) (ContratoDisponiblidad *models.ContratoDisponibilidad, errorOutput interface{}) {
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
		ContratoDisponiblidad = &respuesta[len(respuesta)-1]
		return ContratoDisponiblidad, nil

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

		cdprpArrayInterface, ok := respuesta["cdpxrp"].(map[string]interface{})["cdprp"].([]interface{})

		if !ok {
			return nil, errorOutput
		}
		if len(cdprpArrayInterface) > 0 {

			ultimoElemento := cdprpArrayInterface[len(cdprpArrayInterface)-1].(map[string]interface{})
			crdp = &models.CDPRP{
				CDPNumeroDisponibilidad: ultimoElemento["CDP_NUMERO_DISPONIBILIDAD"].(string),
				RPVigencia:              ultimoElemento["RP_VIGENCIA"].(string),
				CDPVigencia:             ultimoElemento["CDP_VIGENCIA"].(string),
			}

			fmt.Println(respuesta)
			return crdp, nil
		}

		return nil, nil

	}

	return nil, nil

}

func GenerarAutorizacion(idSolicitudPago string) (DocumentoGenerado *models.DocuementoAutorizacionPago, errorOutput interface{}) {

	// Obtiene datos de cambio estado
	var respuestaCambioEstado map[string]interface{}
	urlRequest := beego.AppConfig.String("UrlProveedoresCrud") + "/cambio_estado_cumplido/?query=CumplidoProveedorId:" + idSolicitudPago + ",EstadoCumplidoId.Abreviacion:AO,Activo:true"
	response, err := helpers.GetJsonWSO2Test(urlRequest, &respuestaCambioEstado)

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

	var cambioEstado []models.CambioEstadoCumplido
	data := respuestaCambioEstado["Data"]
	if len(data.([]interface{})) > 0 {
		helpers.LimpiezaRespuestaRefactor(respuestaCambioEstado, &cambioEstado)
	}

	if len(cambioEstado) < 1 {
		return nil, nil

	}
	// Obtiene información de los contratos
	var respuesta []models.ContratoProveedor
	urlRequestContrato := beego.AppConfig.String("UrlcrudAgora") + "/contrato_general/?query=ContratoSuscrito.NumeroContrato:" + cambioEstado[0].CumplidoProveedorId.NumeroContrato
	responseContrato, errContrato := helpers.GetJsonWSO2Test(urlRequestContrato, &respuesta)

	if errContrato != nil || responseContrato != 200 {
		errorMessage := fmt.Sprintf("%v", err)
		errorOutput = map[string]interface{}{
			"Success": false,
			"Status":  responseContrato,
			"Message": "Error al consultar Contrato para el proveedor. " + errorMessage,
			"Error":   errorMessage,
		}
		logs.Error(err)
		return nil, errorOutput
	}

	// Obtiene info Proveedor
	if proveedor, err := ObtenerInfoProveedor(strconv.Itoa(respuesta[0].Contratista)); err == nil && proveedor != nil {
		var respuestaOrdenador []models.Proveedor
		urlRequestOrdenador := beego.AppConfig.String("UrlcrudAgora") + "/informacion_proveedor/?query=NumDocumento:" + strconv.Itoa(cambioEstado[0].DocumentoResponsable)
		responseOrdenador, errOrdenador := helpers.GetJsonWSO2Test(urlRequestOrdenador, &respuestaOrdenador)

		if errOrdenador != nil || responseOrdenador != 200 {
			errorMessage := fmt.Sprintf("%v", err)
			errorOutput = map[string]interface{}{
				"Success": false,
				"Status":  responseOrdenador,
				"Message": "Error al consultar Información del proveedor. " + errorMessage,
				"Error":   errorMessage,
			}
			logs.Error(err)
			return nil, errorOutput
		}

		var respuestaDocumentos map[string]interface{}
		urlRequestDocumentos := beego.AppConfig.String("UrlProveedoresCrud") + "/soporte_pago?query=CumplidoProveedorId.id:" + idSolicitudPago
		responseDocuementos, erroDocumentos := helpers.GetJsonWSO2Test(urlRequestDocumentos, &respuestaDocumentos)
		fmt.Println(urlRequestDocumentos)
		var documentosCargados []models.SoportePago
		if len(respuestaDocumentos["Data"].([]interface{})[0].(map[string]interface{})) != 0 {
			helpers.LimpiezaRespuestaRefactor(respuestaDocumentos, &documentosCargados)
		}

		var idDocumentos []string
		for _, documentosCargado := range documentosCargados {
			idDocumentos = append(idDocumentos, strconv.Itoa(documentosCargado.DocumentoId))
		}
		idDocumentosUnidos := strings.Join(idDocumentos, "|")

		if erroDocumentos != nil || responseDocuementos != 200 {

			return nil, errorOutput
		}

		var documentos []models.Documento
		urlRequestDocumentosGestion := "http://pruebasapi.intranetoas.udistrital.edu.co:8094/v1/documento/?query=TipoDocumento.DominioTipoDocumento.CodigoAbreviacion:CUMP_PROV,Id.in:" + idDocumentosUnidos + ",Activo:true"
		responseDocuementosGestion, erroDocumentosGestion := helpers.GetJsonWSO2Test(urlRequestDocumentosGestion, &documentos)

		var lista_documentos_cargados_strings []string
		for _, documento := range documentos {
			lista_documentos_cargados_strings = append(lista_documentos_cargados_strings, documento.TipoDocumento.CodigoAbreviacion)
		}

		indexRespuestaOrdenador := len(respuesta) - 1
		DocumentoGenerado := &models.DocuementoAutorizacionPago{
			NombreOrdenador:    respuestaOrdenador[indexRespuestaOrdenador].NomProveedor,
			DocumentoOrdenador: respuestaOrdenador[indexRespuestaOrdenador].NumDocumento,
			NombreProveedor:    proveedor.NomProveedor,
			DocumentoProveedor: proveedor.NumDocumento,
			DocumentosCargados: lista_documentos_cargados_strings,
		}

		return DocumentoGenerado, nil

		if erroDocumentosGestion != nil || responseDocuementosGestion != 200 {
			return nil, errorOutput
		}

	}
	return nil, nil
}
