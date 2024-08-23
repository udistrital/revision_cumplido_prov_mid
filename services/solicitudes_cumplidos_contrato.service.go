package services

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtenerSolicitudesCumplidosContrato(numero_contrato string, vigencia string) (estados_cumplido []models.CumplidosContrato, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ContratosSupervisor", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var respuesta_peticion map[string]interface{}
	var cumplidos_proveedor []models.CumplidoProveedor
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cumplido_proveedor?query=NumeroContrato:"+numero_contrato+",VigenciaContrato:"+vigencia+"&sortby=FechaCreacion&order=desc", &respuesta_peticion); err == nil && response == 200 {
		helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &cumplidos_proveedor)
		for _, cumplido_proveedor := range cumplidos_proveedor {
			estado, err := ObtenerUltimoEstadoCumplidoProveedor(strconv.Itoa(cumplido_proveedor.Id))
			if err == nil {
				estado_cumplido := models.CumplidosContrato{
					NumeroContrato:      estado.CumplidoProveedorId.NumeroContrato,
					FechaCreacion:       estado.FechaCreacion,
					Periodo:             "",
					EstadoCumplido:      estado.EstadoCumplidoId.Nombre,
					CumplidoProveedorId: estado.CumplidoProveedorId,
				}
				estados_cumplido = append(estados_cumplido, estado_cumplido)

			} else {
				outputError = map[string]interface{}{"funcion": "ObtenerSolicitudesCumplidosContrato", "err": err, "status": "502"}
				return nil, outputError
			}
		}
	}
	return estados_cumplido, outputError
}

func ObtenerUltimoEstadoCumplidoProveedor(cumplido_proveedor_id string) (estado_cumplido models.CambioEstadoCumplido, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ContratosSupervisor", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var respuesta_peticion map[string]interface{}
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido?query=CumplidoProveedorId.Id:"+cumplido_proveedor_id+"&sortby=FechaCreacion&order=desc&limit=1", &respuesta_peticion); err == nil && response == 200 {
		estado_josn, err := json.Marshal(respuesta_peticion["Data"].([]interface{})[0])
		if err == nil {
			json.Unmarshal(estado_josn, &estado_cumplido)
			return estado_cumplido, nil
		} else {
			outputError = map[string]interface{}{"funcion": "GetUltimoEstadoCumplidoProveedor", "err": err, "status": "502"}
			return estado_cumplido, outputError
		}
	} else {
		outputError = map[string]interface{}{"funcion": "GetUltimoEstadoCumplidoProveedor", "err": err, "status": "502"}
		return estado_cumplido, outputError
	}
	return estado_cumplido, outputError
}
