package services

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtenerSolicitudesCumplidosContrato(numero_contrato string, vigencia string) (solicitudes_cumplido []models.CumplidosContrato, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var respuesta_peticion map[string]interface{}
	var cumplidos_proveedor []models.CumplidoProveedor
	//fmt.Println(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores") + "/cumplido_proveedor/?query=NumeroContrato:" + numero_contrato + ",VigenciaContrato:" + vigencia + "&sortby=FechaCreacion&order=desc")
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cumplido_proveedor/?query=NumeroContrato:"+numero_contrato+",VigenciaContrato:"+vigencia+"&sortby=FechaCreacion&order=desc", &respuesta_peticion); err == nil && response == 200 {
		data := respuesta_peticion["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) == 0 {
			outputError = fmt.Errorf("No se encontraron cumplidos para el contrato")
			return nil, outputError
		}
		helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &cumplidos_proveedor)
		for i, cumplido_proveedor := range cumplidos_proveedor {
			estado, err := ObtenerUltimoEstadoCumplidoProveedor(strconv.Itoa(cumplido_proveedor.Id))
			if err == nil {
				solicitud_cumplido := models.CumplidosContrato{
					ConsecutivoCumplido:             i + 1,
					NumeroContrato:                  estado.CumplidoProveedorId.NumeroContrato,
					FechaCreacion:                   estado.FechaCreacion,
					Periodo:                         ObtenerPeriodoInformacionPago(cumplido_proveedor.Id),
					EstadoCumplido:                  estado.EstadoCumplidoId.Nombre,
					CodigoAbreviacionEstadoCumplido: estado.EstadoCumplidoId.CodigoAbreviacion,
					CumplidoProveedorId:             estado.CumplidoProveedorId,
				}
				solicitudes_cumplido = append(solicitudes_cumplido, solicitud_cumplido)

			} else {
				outputError = fmt.Errorf("Error al obtener el ultimo estado del cumplido")
				return nil, outputError
			}
		}
	}
	return solicitudes_cumplido, outputError
}

func ObtenerPeriodoInformacionPago(cumplido_proveedor_id int) (periodo_pago string) {

	periodo_pago = ""
	var respuesta_peticion map[string]interface{}
	var informacion_pago_proveedor []models.InformacionPago
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/informacion_pago/?query=CumplidoProveedorId.Id:"+strconv.Itoa(cumplido_proveedor_id), &respuesta_peticion); err == nil && response == 200 {
		data := respuesta_peticion["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) == 0 {
			return periodo_pago
		}
		helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &informacion_pago_proveedor)
		fecha_inicio := informacion_pago_proveedor[0].FechaInicial.Format("2006-01-02")
		fecha_fin := informacion_pago_proveedor[0].FechaFinal.Format("2006-01-02")

		periodo_pago = fecha_inicio + " - " + fecha_fin

		return periodo_pago

	} else {
		return periodo_pago
	}
}

func ObtenerUltimoEstadoCumplidoProveedor(cumplido_proveedor_id string) (estado_cumplido models.CambioEstadoCumplido, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var respuesta_peticion map[string]interface{}
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/?query=Activo:true,CumplidoProveedorId.Id:"+cumplido_proveedor_id+"&sortby=FechaCreacion&order=desc&limit=1", &respuesta_peticion); err == nil && response == 200 {
		data := respuesta_peticion["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) > 0 {
			estado_json, err := json.Marshal(respuesta_peticion["Data"].([]interface{})[0])
			if err == nil {
				json.Unmarshal(estado_json, &estado_cumplido)
				return estado_cumplido, nil
			} else {
				outputError = fmt.Errorf("Error al convertir el Json")
				return estado_cumplido, outputError
			}
		} else {
			outputError = fmt.Errorf("No se encontraron estados para el cumplido")
			return estado_cumplido, outputError
		}
	} else {
		outputError = fmt.Errorf("Error al obtener el ultimo estado del cumplido")
		return estado_cumplido, outputError
	}
}
