package helpers_estado_pago

import (
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func CambioEstadoPago(estado_cumplido_id string, cumplido_proveedor_id string, documento_responsable string, cargo_responsable string) (respuesta_cambio_estado models.CambioEstadoCumplidoResponse, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "CambioEstadoPago", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	// Verificar que se envien todos los datos
	if estado_cumplido_id == "" || cumplido_proveedor_id == "" || documento_responsable == "" || cargo_responsable == "" {
		outputError = map[string]interface{}{"funcion": "CambioEstadoPago", "status": "502", "mensaje": "Faltan datos por enviar"}
		return respuesta_cambio_estado, outputError
	}

	var respuesta_peticion map[string]interface{}
	var respuesta map[string]interface{}
	var estado_cumplido []models.EstadoCumplido
	var cambios_anteriores []models.CambioEstadoCumplido
	id_estado, _ := strconv.Atoi(estado_cumplido_id)
	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/?query=CumplidoProveedorId.Id:"+cumplido_proveedor_id+",Activo:true&sortby=FechaCreacion&order=desc", &respuesta_peticion); err == nil && response == 200 {
		LimpiezaRespuestaRefactor(respuesta_peticion, &cambios_anteriores)
		if len(cambios_anteriores) > 0 && cambios_anteriores[0].EstadoCumplidoId.Id == id_estado && cambios_anteriores[0].Activo == true {
			outputError = map[string]interface{}{"funcion": "CambioEstadoPago", "status": "200", "mensaje": "El cumplido ya se encuentra en este estado"}
			return respuesta_cambio_estado, outputError
		} else if len(cambios_anteriores) > 0 {
			for _, cambio_anterior := range cambios_anteriores {
				var respuesta_estado_anterior map[string]interface{}
				cambio_anterior.Activo = false
				cambio_anterior.FechaModificacion = time.Now()
				err := sendJson(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/"+strconv.Itoa(cambio_anterior.Id), "PUT", &respuesta_estado_anterior, cambio_anterior)
				if err != nil {
					outputError = map[string]interface{}{"funcion": "CambioEstadoPago", "status": "502", "mensaje": "Error al actualizar el cumplido proveedor"}
					return respuesta_cambio_estado, outputError
				}
			}
		}
	} else {
		outputError = map[string]interface{}{"funcion": "CambioEstadoPago", "status": "502", "mensaje": "Error al consultar el estado del cumplido"}
		return respuesta_cambio_estado, outputError
	}

	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/estado_cumplido/?query=Id:"+estado_cumplido_id, &respuesta); (err == nil) && (response == 200) {
		LimpiezaRespuestaRefactor(respuesta, &estado_cumplido)
		if estado_cumplido != nil {
			var cumplido_proveedor []models.CumplidoProveedor
			if response, err := getJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cumplido_proveedor/?query=Id:"+cumplido_proveedor_id, &respuesta); (err == nil) && (response == 200) {
				LimpiezaRespuestaRefactor(respuesta, &cumplido_proveedor)
				if respuesta != nil {
					documento, _ := strconv.Atoi(documento_responsable)
					var body = models.BodyCambioEstadoCumplido{
						EstadoCumplidoId:     estado_cumplido[0],
						CumplidoProveedorId:  cumplido_proveedor[0],
						DocumentoResponsable: documento,
						CargoResponsable:     cargo_responsable,
						Activo:               true,
						FechaCreacion:        time.Now(),
						FechaModificacion:    time.Now(),
					}
					if err := sendJson(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido", "POST", &respuesta, body); err == nil {
						var respuesta_cambio_estado models.CambioEstadoCumplidoResponse
						respuesta_cambio_estado.CumplidoProveedorId = cumplido_proveedor[0].Id
						respuesta_cambio_estado.DocumentoResponsable = documento
						respuesta_cambio_estado.EstadoCumplido = &estado_cumplido[0]
						return respuesta_cambio_estado, nil
					} else {
						outputError = map[string]interface{}{"funcion": "CambioEstadoPago", "status": "502", "mensaje": "Error al actualizar el cumplido proveedor"}
						return respuesta_cambio_estado, outputError
					}
				} else {
					outputError = map[string]interface{}{"funcion": "CambioEstadoPago", "status": "502", "mensaje": "No se encontró el cumplido proveedor"}
					return respuesta_cambio_estado, outputError
				}
			} else {
				outputError = map[string]interface{}{"funcion": "CambioEstadoPago", "status": "502", "mensaje": "Error al consultar el cumplido proveedor"}
				return respuesta_cambio_estado, outputError
			}
		}
	}
	return respuesta_cambio_estado, outputError
}
