package helpers_ordenador

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
	"strconv"
	"strings"
)

func ObteberNumeroDeContrato(documentoOrdenador string, estado string) (dependencias string, errorOutput interface{}) {

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
	response, err := helpers.GetJsonWSO2Test(urlRequest, &respuesta)

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

func ObternerContratos(documento string, estado string) (ContratosList []models.Contrato, errorOutput interface{}) {

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

	cumplidos, e := ObteberNumeroDeContrato(documento, estado)

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

	for _, contrato := range respuesta {
		var ultimoContrato = len(contrato.ContratoSuscrito) - 1
		if proveedor, err := ObtenerInfoProveedor(strconv.Itoa(contrato.Contratista)); err == nil && proveedor != nil {
			/*contratoDisponibilidad, _ := ObtenerContratoDisponiblidad(contrato.ContratoSuscrito[ultimoContrato].NumeroContrato.Id)*/
			contrato := models.Contrato{
				TipoContrato:    contrato.TipoContrato.TipoContrato,
				NumeroContrato:  contrato.ContratoSuscrito[ultimoContrato].NumeroContrato.Id,
				Vigencia:        contrato.ContratoSuscrito[ultimoContrato].Vigencia,
				Dependencia:     contrato.DependenciaSolicitante,
				NombreProveedor: proveedor.NomProveedor,
			}
			ContratosList = append(ContratosList, contrato)

		}
	}

	return ContratosList, nil
}

func ObternerAprobadoSupervisor(idPago string) (CumplidoFinal *models.CambioEstadoCumplido, errorOutput interface{}) {

	defer func() {
		if err := recover(); err != nil {
			errorMessage := fmt.Sprintf("%v", err)
			errorOutput = map[string]interface{}{
				"Success": false,
				"Status":  400,
				"Message": "Error al consultar las dependencias: " + idPago + " - " + errorMessage,
				"Error":   errorMessage,
			}
			fmt.Println("Error recuperado:", errorMessage)
		}
	}()

	var respuesta map[string]interface{}
	urlRequest := beego.AppConfig.String("UrlProveedoresCrud") + "/cambio_estado_cumplido/?query=CumplidoProveedorId:" + idPago + ",EstadoCumplidoId.Abreviacion:AO"
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
	var Cumplidos []models.CambioEstadoCumplido
	if len(respuesta["Data"].([]interface{})[0].(map[string]interface{})) != 0 {
		helpers.LimpiezaRespuestaRefactor(respuesta, &Cumplidos)

	}
	if respuesta != nil {
		CumplidoFinal = &Cumplidos[0]
	}
	print(CumplidoFinal)
	return CumplidoFinal, nil
}

func RevertirAprobadoSupervisor(Id string) (DependenciasList *models.CambioEstadoCumplido, errorOutput interface{}) {

	defer func() {
		if err := recover(); err != nil {
			errorMessage := fmt.Sprintf("%v", err)
			errorOutput = map[string]interface{}{
				"Success": false,
				"Status":  400,
				"Message": "Error al consultar el pago " + Id + " - " + errorMessage,
				"Error":   errorMessage,
			}
			fmt.Println("Error recuperado:", errorMessage)
		}
	}()

	var CumplidoARevertir, err = ObternerAprobadoSupervisor(Id)

	if err != nil {
		errorMessage := fmt.Sprintf("%v", err)
		errorOutput = map[string]interface{}{
			"Success": false,
			"Status":  400,
			"Message": "Error el cambio de de estado para : " + Id + "--" + errorMessage,
			"Error":   err,
		}
	} else if CumplidoARevertir.CumplidoProveedorId != nil {

		estado, errorEstado := ObtenerEstado("PRO")
		fmt.Println("estado")
		fmt.Println(estado)
		if errorEstado == nil && estado != nil {

			updateCamioEstado := &models.CambioEstadoCumplido{
				Id:                  CumplidoARevertir.Id,
				EstadoCumplidoId:    (*models.EstadoCumplido)(estado),
				CumplidoProveedorId: CumplidoARevertir.CumplidoProveedorId,
			}
			var respuesta map[string]interface{}
			helpers.SendJson(beego.AppConfig.String("UrlProveedoresCrud")+"/cambio_estado_cumplido/"+Id, "PUT", &respuesta, updateCamioEstado)

			helpers.LimpiezaRespuestaRefactor(respuesta, &updateCamioEstado)
			return updateCamioEstado, nil
		}

	}
	return nil, nil
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
	var estado_List []models.EstadoCumplidoId

	helpers.LimpiezaRespuestaRefactor(respuesta, &estado_List)
	if respuesta != nil {
		Estado = &estado_List[0]
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
		ContratoDisponiblidad = &respuesta[0]
		return ContratoDisponiblidad, nil

	}
	return nil, nil

}

func ObtenerCumplido(IdPago string) (cumplido *models.CumplidoProveedorId, errorOutput interface{}) {

	defer func() {
		if err := recover(); err != nil {
			errorMessage := fmt.Sprintf("%v", err)
			errorOutput = map[string]interface{}{
				"Success": false,
				"Status":  400,
				"Message": "Error al consultar la solicitud de : " + IdPago + " - " + errorMessage,
				"Error":   errorMessage,
			}
			fmt.Println("Error recuperado:", errorMessage)
		}
	}()

	var respuesta map[string]interface{}
	var urlRequest = beego.AppConfig.String("UrlProveedoresCrud") + "/cumplido_proveedor/?query=Id:" + IdPago
	response, err := helpers.GetJsonWSO2Test(urlRequest, &respuesta)
	fmt.Println("respuestaPeticion", respuesta)
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

	var cumplidoResponse []models.CumplidoProveedorId
	if len(respuesta["Data"].([]interface{})[0].(map[string]interface{})) != 0 {
		helpers.LimpiezaRespuestaRefactor(respuesta, &cumplidoResponse)
		cumplido = &cumplidoResponse[0]

	}

	return cumplido, nil
}

func ObtenerInfoContratoPorId(IdContrato string) (contratoRespuesta *models.Contrato, errorOutput interface{}) {

	defer func() {
		if err := recover(); err != nil {
			errorMessage := fmt.Sprintf("%v", err)
			errorOutput = map[string]interface{}{
				"Success": false,
				"Status":  500, // Cambiado a 500 para errores internos
				"Message": "Error al consultar las dependencias: " + IdContrato + " - " + errorMessage,
				"Error":   errorMessage,
			}
			fmt.Println("Error recuperado:", errorMessage)
		}
	}()

	var respuesta []models.ContratoProveedor
	urlRequest := beego.AppConfig.String("UrlcrudAgora") + "/contrato_general/?query=id:" + IdContrato
	println(urlRequest)
	response, err := helpers.GetJsonWSO2Test(urlRequest, &respuesta)

	if err != nil {
		errorMessage := fmt.Sprintf("Error al consultar los contratos para el ordenador en .GetContratos: %v", err)
		errorOutput = map[string]interface{}{
			"Success": false,
			"Status":  400,
			"Message": errorMessage,
			"Error":   err.Error(),
		}
		logs.Error(err)
		return nil, errorOutput
	}

	if response != 200 {
		errorMessage := fmt.Sprintf("Error al consultar los contratos para el ordenador en .GetContratos: código de estado %d", response)
		errorOutput = map[string]interface{}{
			"Success": false,
			"Status":  400,
			"Message": errorMessage,
			"Error":   errorMessage,
		}
		logs.Error(errorMessage)
		return nil, errorOutput
	}

	for _, contrato := range respuesta {
		ultimoContrato := len(contrato.ContratoSuscrito) - 1
		proveedor, err := ObtenerInfoProveedor(strconv.Itoa(contrato.Contratista))
		if err == nil && proveedor != nil {
			contratoRespuesta = &models.Contrato{
				TipoContrato:    contrato.TipoContrato.TipoContrato,
				NumeroContrato:  contrato.ContratoSuscrito[ultimoContrato].NumeroContrato.Id,
				Vigencia:        contrato.ContratoSuscrito[ultimoContrato].Vigencia,
				Dependencia:     contrato.DependenciaSolicitante,
				NombreProveedor: proveedor.NomProveedor,
			}
			return contratoRespuesta, nil
		}
	}

	errorOutput = map[string]interface{}{
		"Success": false,
		"Status":  404,
		"Message": "No se encontró ningún contrato válido para el ID proporcionado: " + IdContrato,
		"Error":   "Contrato no encontrado",
	}
	return nil, errorOutput
}
