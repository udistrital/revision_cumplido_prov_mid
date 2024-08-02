package helpers_contratacion

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers/helpers_ordenador"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func obtenerCumplidos(estado string) (cambios_estado_limpios []models.CambioEstadoCumplido, errorOutput interface{}) {

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

	var cambios_estado map[string]interface{}

	var urlRequest = beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores") + "/cambio_estado_cumplido/?query=EstadoCumplidoId.CodigoAbreviación:" + estado

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
		return nil, errorOutput
	}

	if len(cambios_estado["Data"].([]interface{})[0].(map[string]interface{})) != 0 {
		helpers.LimpiezaRespuestaRefactor(cambios_estado, &cambios_estado_limpios)

	}
	return cambios_estado_limpios, nil
}

func ObtenerCumplidosPendientesContratacion(estado string) (cumplidosInfo []models.SolicituRevisionCumplidoProveedor, errorOutput interface{}) {
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

	cumplidos, e := obtenerCumplidos(estado)

	if e != nil || cumplidos == nil {

		return nil, e
	}

	var info_contrato []models.ContratoGeneral
	for _, cumplido := range cumplidos {

		urlRequest := beego.AppConfig.String("UrlcrudAgora") + "/contrato_general/?query=ContratoSuscrito.NumeroContrato:" + cumplido.CumplidoProveedorId.NumeroContrato
		fmt.Println("Pendiente contratacion: ", urlRequest)
		response, err := helpers.GetJsonWSO2Test(urlRequest, &info_contrato)

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

		var ultimoContrato = len(info_contrato[0].ContratoSuscrito) - 1

		if proveedor, err := helpers_ordenador.ObtenerInfoProveedor(strconv.Itoa(info_contrato[0].Contratista)); err == nil && proveedor != nil {

			if proveedor != nil {
				contrato_disponibilidad, _ := helpers_ordenador.ObtenerContratoDisponiblidad(info_contrato[0].ContratoSuscrito[ultimoContrato].NumeroContrato.Id)

				if contrato_disponibilidad != nil {

					cdprp, _ := helpers_ordenador.ObtenerCrdp(strconv.Itoa(contrato_disponibilidad.NumeroCdp), strconv.Itoa(contrato_disponibilidad.Vigencia))
					if cdprp != nil {

						contrato := models.SolicituRevisionCumplidoProveedor{
							TipoContrato:     info_contrato[0].TipoContrato.TipoContrato,
							NumeroContrato:   info_contrato[0].ContratoSuscrito[ultimoContrato].NumeroContrato.Id,
							VigenciaContrato: info_contrato[0].ContratoSuscrito[ultimoContrato].Vigencia,
							Dependencia:      info_contrato[0].DependenciaSolicitante,
							NombreProveedor:  proveedor.NomProveedor,
							Cdp:              strconv.Itoa(contrato_disponibilidad.NumeroCdp),
							Rp:               cdprp.CDPNumeroDisponibilidad,
							VigenciaRP:       cdprp.RPVigencia,
							Id:               cumplido.CumplidoProveedorId.Id,
							Activo:           cumplido.Activo,
						}
						cumplidosInfo = append(cumplidosInfo, contrato)
					}
				}
			}

		}
	}

	return cumplidosInfo, nil
}
