package helpers_soporte

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func AgregarComentarioSoporte(soporte_id string, cambio_estado_id string, comentario string) (respuesta models.RespuestaComentarioSoporte, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/AgregarComentarioSoporte", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	var respuesta_peticion map[string]interface{}
	var comentario_soporte models.ComentarioSoporte
	var soporte_pago models.SoportePago
	var cambio_estado_cumplido models.CambioEstadoCumplido

	if soporte_id == "" || cambio_estado_id == "" || comentario == "" {
		outputError = map[string]interface{}{"funcion": "/AgregarComentarioSoporte", "err": "Faltan datos en la solicitud", "status": "400"}
		return respuesta, outputError
	}

	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/soporte_pago/"+soporte_id, &respuesta_peticion); (err == nil) && (response == 200) {
		LimpiezaRespuestaRefactor(respuesta_peticion, &soporte_pago)
		if response, err := getJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/"+cambio_estado_id, &respuesta_peticion); (err == nil) && (response == 200) {
			LimpiezaRespuestaRefactor(respuesta_peticion, &cambio_estado_cumplido)
			comentario_soporte.Comentario = comentario
			comentario_soporte.SoportePagoId = &soporte_pago
			comentario_soporte.CambioEstadoCumplidoId = &cambio_estado_cumplido
			if err := sendJson(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/comentario_soporte", "POST", &respuesta, comentario_soporte); err == nil {
				respuesta.SoportePagoId = soporte_pago.Id
				respuesta.CambioEstadoCumplidoId = cambio_estado_cumplido.Id
				respuesta.Comentario = comentario
				return respuesta, nil
			} else {
				outputError = map[string]interface{}{"funcion": "/AgregarComentarioSoporte/comentario_soporte", "err": err, "status": "502"}
				return respuesta, outputError
			}
		} else {
			outputError = map[string]interface{}{"funcion": "/AgregarComentarioSoporte/cambio_estado_cumplido", "err": err, "status": "502"}
			return respuesta, outputError
		}
	}
	return respuesta, outputError
}
