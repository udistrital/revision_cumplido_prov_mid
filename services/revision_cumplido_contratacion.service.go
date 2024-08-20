package services

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtenerCumplidosPorEstado(estado string) (cambios_estado_limpios []models.CambioEstadoCumplido, errorOutput interface{}) {

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

	cumplidos, error := ObtenerCumplidosPorEstado(estado)
	if error != nil || cumplidos == nil {
		errorOutput = map[string]interface{}{
			"Success": false,
			"Status":  400,
			"Message": "Error al consultar los cumplidos para el proveedor en .ObtenerCumplidosPorEstado",
			"Error":   error,
		}
		return nil, error
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
