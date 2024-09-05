package services

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtenerCumplidosPendientesContratacion() (cumplidosInfo []models.SolicituRevisionCumplidoProveedor, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var cumplidos []models.CambioEstadoCumplido
	var respuesta_peticion map[string]interface{}

	//fmt.Println("UrlCrudRevisionCumplidosProveedores", beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/?query=EstadoCumplidoId.CodigoAbreviación:PRC,Activo:true")
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/?query=EstadoCumplidoId.CodigoAbreviacion:PRC,Activo:true", &respuesta_peticion); err == nil && response == 200 {
		data := respuesta_peticion["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) > 0 {
			helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &cumplidos)
			if len(cumplidos) > 0 {
				for _, cumplido := range cumplidos {
					informacion_contrato, err := helpers.ObtenerInformacionContratoProveedor(cumplido.CumplidoProveedorId.NumeroContrato, strconv.Itoa(cumplido.CumplidoProveedorId.VigenciaContrato))
					if err == nil {
						contrato := models.SolicituRevisionCumplidoProveedor{
							TipoContrato:     informacion_contrato[0].TipoContrato,
							NumeroContrato:   informacion_contrato[0].NumeroContratoSuscrito,
							VigenciaContrato: cumplido.CumplidoProveedorId.VigenciaContrato,
							Dependencia:      informacion_contrato[0].NombreDependencia,
							NombreProveedor:  informacion_contrato[0].NombreProveedor,
							Cdp:              informacion_contrato[0].NumeroCdp,
							Rp:               informacion_contrato[0].NumeroRp,
							VigenciaRP:       informacion_contrato[0].VigenciaRp,
							CumplidoId:       cumplido.CumplidoProveedorId.Id,
							Activo:           cumplido.Activo,
						}
						cumplidosInfo = append(cumplidosInfo, contrato)
					} else {
						logs.Error(err)
						outputError = fmt.Errorf("Error al obtener la informacion del contrato")
						return cumplidosInfo, outputError
					}
				}
			} else {
				outputError = fmt.Errorf("No se encontraron cumplidos pendientes")
			}
		} else {
			outputError = fmt.Errorf("No se encontraron cumplidos pendientes por revision de contratación")
			return cumplidosInfo, outputError
		}

	} else {
		logs.Error(err)
		outputError = fmt.Errorf("Error al obtener los cumplidos pendientes")
		return cumplidosInfo, outputError
	}
	return cumplidosInfo, outputError
}
