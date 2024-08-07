package helpers_contratacion

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers/helpers_ordenador"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
	"strconv"
	"strings"
)

func ObteberNumeroDeContrato() (dependencias string, errorOutput interface{}) {

	defer func() {
		if err := recover(); err != nil {
			errorMessage := fmt.Sprintf("%v", err)
			errorOutput = map[string]interface{}{
				"Success": false,
				"Status":  400,
				"Message": "Error al consultar las dependencias: " + " - " + errorMessage,
				"Error":   errorMessage,
			}
			fmt.Println("Error recuperado:", errorMessage)
		}
	}()

	var respuesta map[string]interface{}
	var urlRequest = beego.AppConfig.String("UrlProveedoresCrud") + "/cambio_estado_cumplido/?query=EstadoCumplidoId.CodigoAbreviación:PRC,Activo:true"
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
	if len(dependenciasString) > 0 {
		listaString := strings.Join(dependenciasString, "|") + "|"
		return listaString, nil
	}
	return "0", nil
}

func ObternerCumplidosPendientesContratacion() (solicitudes []models.Contrato, errorOutput interface{}) {
	defer func() {
		if err := recover(); err != nil {
			errorMessage := fmt.Sprintf("%v", err)
			errorOutput = map[string]interface{}{
				"Success": false,
				"Status":  400,
				"Message": "Error al consultar las dependencias: " + " - " + errorMessage,
				"Error":   errorMessage,
			}
			fmt.Println("Error recuperado:", errorMessage)
		}
	}()

	cumplidos, e := ObteberNumeroDeContrato()

	if e != nil || cumplidos == "" {

		return nil, e
	}

	var respuesta []models.ContratoProveedor
	urlRequest := beego.AppConfig.String("UrlcrudAgora") + "/contrato_general/?query=ContratoSuscrito.NumeroContrato.in:" + cumplidos
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
	for _, contrato := range respuesta {

		var ultimoContrato = len(contrato.ContratoSuscrito) - 1
		if proveedor, err := helpers_ordenador.ObtenerInfoProveedor(strconv.Itoa(contrato.Contratista)); err == nil && proveedor != nil {

			if proveedor != nil {
				fmt.Println("proverdor no es nulo")
				contratoDisponibilidad, _ := helpers_ordenador.ObtenerContratoDisponiblidad(contrato.ContratoSuscrito[ultimoContrato].NumeroContrato.Id)

				if contratoDisponibilidad != nil {
					fmt.Println("contratoDisponibilidad no es nulo")
					cdprp, _ := helpers_ordenador.ObtenerCrdp(strconv.Itoa(contratoDisponibilidad.NumeroCdp), strconv.Itoa(contratoDisponibilidad.Vigencia))

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
						solicitudes = append(solicitudes, contrato)
					}
				}
			}

		}
	}

	return solicitudes, nil
}
