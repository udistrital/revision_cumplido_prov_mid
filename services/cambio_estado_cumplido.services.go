package services

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func CambioEstadoCumplido(codigo_abreviacion_cumplido string, cumplido_proveedor_id int) (respuesta_cambio_estado models.CambioEstadoCumplidoResponse, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "CambioEstadoCumplido", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	// Verificar que se envien todos los datos
	if codigo_abreviacion_cumplido == "" || cumplido_proveedor_id == 0 {
		outputError = map[string]interface{}{"funcion": "CambioEstadoCumplido", "status": "404", "mensaje": "Faltan datos por enviar"}
		return respuesta_cambio_estado, outputError
	}

	var cumplido_proveedor []models.CumplidoProveedor
	var estado_cumplido []models.EstadoCumplido
	var respuesta_cumplido_proveedor map[string]interface{}
	var respuesta_estado_cumplido map[string]interface{}

	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cumplido_proveedor/?query=Id:"+strconv.Itoa(cumplido_proveedor_id), &respuesta_cumplido_proveedor); (err == nil) && (response == 200) {
		data := respuesta_cumplido_proveedor["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) > 0 {
			helpers.LimpiezaRespuestaRefactor(respuesta_cumplido_proveedor, &cumplido_proveedor)

			if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/estado_cumplido/?query=CodigoAbreviacion:"+codigo_abreviacion_cumplido, &respuesta_estado_cumplido); (err == nil) && (response == 200) {
				data := respuesta_estado_cumplido["Data"].([]interface{})
				if len(data[0].(map[string]interface{})) == 0 {
					outputError = map[string]interface{}{"funcion": "CambioEstadoCumplido", "status": "404", "mensaje": "El estado del cumplido proveedor ingresado no existe"}
					return respuesta_cambio_estado, outputError
				}
				helpers.LimpiezaRespuestaRefactor(respuesta_estado_cumplido, &estado_cumplido)
			} else {
				outputError = map[string]interface{}{"funcion": "CambioEstadoCumplido", "status": "404", "mensaje": "Error al consultar el estado del cumplido"}
				return respuesta_cambio_estado, outputError
			}

		} else {
			outputError = map[string]interface{}{"funcion": "CambioEstadoCumplido", "status": "404", "mensaje": "El cumplido proveedor ingresado no existe"}
			return respuesta_cambio_estado, outputError
		}
	} else {
		outputError = map[string]interface{}{"funcion": "CambioEstadoCumplido", "status": "404", "mensaje": "Error al consultar el cumplido proveedor"}
		return respuesta_cambio_estado, outputError
	}

	body_cambio_estado, outputError := CrearBodyCambioEstadoCumplido(codigo_abreviacion_cumplido, cumplido_proveedor[0], estado_cumplido[0])
	if outputError != nil {
		return respuesta_cambio_estado, outputError
	}

	ultimo_cambio_estado_cumplido, err := DesactivarCambiosAnterioresCumplido(cumplido_proveedor_id, codigo_abreviacion_cumplido)
	if err != nil {
		outputError = err
		return respuesta_cambio_estado, outputError
	}

	if (ultimo_cambio_estado_cumplido != models.CambioEstadoCumplido{}) && body_cambio_estado.DocumentoResponsable == ultimo_cambio_estado_cumplido.DocumentoResponsable && body_cambio_estado.EstadoCumplidoId.Id == ultimo_cambio_estado_cumplido.EstadoCumplidoId.Id {
		outputError = map[string]interface{}{"funcion": "CambioEstadoCumplido", "status": "404", "mensaje": "El cumplido ya se encuentra en este estado"}
		return respuesta_cambio_estado, outputError
	} else {
		var respuesta_peticion map[string]interface{}
		if err := helpers.SendJson(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido", "POST", &respuesta_peticion, body_cambio_estado); err == nil {
			respuesta_cambio_estado.CumplidoProveedorId = cumplido_proveedor_id
			respuesta_cambio_estado.DocumentoResponsable = body_cambio_estado.DocumentoResponsable
			respuesta_cambio_estado.CargoResponsable = body_cambio_estado.CargoResponsable
			respuesta_cambio_estado.EstadoCumplido = &estado_cumplido[0]

			switch codigo_abreviacion_cumplido {
			case "RC":
				respuesta_cambio_estado, outputError = CambioEstadoCumplido("CD", cumplido_proveedor_id)
			case "AC":
				respuesta_cambio_estado, outputError = CambioEstadoCumplido("PRO", cumplido_proveedor_id)
			case "RO":
				respuesta_cambio_estado, outputError = CambioEstadoCumplido("CD", cumplido_proveedor_id)
			}
			return respuesta_cambio_estado, outputError
		} else {
			outputError = map[string]interface{}{"funcion": "CambioEstadoCumplido", "status": "404", "mensaje": "Error al actualizar el cumplido proveedor"}
			return respuesta_cambio_estado, outputError
		}
	}
}

func CrearBodyCambioEstadoCumplido(codigo_abreviacion_cumplido string, cumplido_proveedor models.CumplidoProveedor, estado_cumplido models.EstadoCumplido) (body_cambio_estado models.BodyCambioEstadoCumplido, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "CrearBodyCambioEstadoCumplido", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	switch codigo_abreviacion_cumplido {
	case "CD":
		supervisor_contrato, err := ObtenerSupervisorContrato(cumplido_proveedor.NumeroContrato, strconv.Itoa(cumplido_proveedor.VigenciaContrato))
		if err == nil {
			documento_supervisor, _ := strconv.Atoi(supervisor_contrato.Contratos.Supervisor[0].Documento)
			body_cambio_estado.EstadoCumplidoId = estado_cumplido
			body_cambio_estado.CumplidoProveedorId = cumplido_proveedor
			body_cambio_estado.DocumentoResponsable = documento_supervisor
			body_cambio_estado.CargoResponsable = supervisor_contrato.Contratos.Supervisor[0].Cargo
		} else {
			outputError = map[string]interface{}{"funcion": "CambioEstadoCumplido", "status": "404", "mensaje": "Error al obtener el supervisor del contrato"}
			return body_cambio_estado, outputError
		}
	case "PRC":
		body_cambio_estado.EstadoCumplidoId = estado_cumplido
		body_cambio_estado.CumplidoProveedorId = cumplido_proveedor
		body_cambio_estado.DocumentoResponsable = 0
		body_cambio_estado.CargoResponsable = "CONTRATACIÓN"
	case "RC":
		body_cambio_estado.EstadoCumplidoId = estado_cumplido
		body_cambio_estado.CumplidoProveedorId = cumplido_proveedor
		body_cambio_estado.DocumentoResponsable = 0
		body_cambio_estado.CargoResponsable = "CONTRATACIÓN"
	case "AC":
		body_cambio_estado.EstadoCumplidoId = estado_cumplido
		body_cambio_estado.CumplidoProveedorId = cumplido_proveedor
		body_cambio_estado.DocumentoResponsable = 0
		body_cambio_estado.CargoResponsable = "CONTRATACIÓN"
	case "PRO":
		ordenador_contrato, err := helpers.ObtenerOrdenadorContrato(cumplido_proveedor.NumeroContrato, strconv.Itoa(cumplido_proveedor.VigenciaContrato))
		if err == nil {
			documento_ordenador, _ := strconv.Atoi(ordenador_contrato.Contratos.Ordenador[0].Documento)
			body_cambio_estado.EstadoCumplidoId = estado_cumplido
			body_cambio_estado.CumplidoProveedorId = cumplido_proveedor
			body_cambio_estado.DocumentoResponsable = documento_ordenador
			body_cambio_estado.CargoResponsable = ordenador_contrato.Contratos.Ordenador[0].RolOrdenador
		} else {
			outputError = map[string]interface{}{"funcion": "CambioEstadoCumplido", "status": "404", "mensaje": "Error al obtener el ordenador del contrato"}
			return body_cambio_estado, outputError
		}
	case "AO":
		ordenador_contrato, err := helpers.ObtenerOrdenadorContrato(cumplido_proveedor.NumeroContrato, strconv.Itoa(cumplido_proveedor.VigenciaContrato))
		if err == nil {
			documento_ordenador, _ := strconv.Atoi(ordenador_contrato.Contratos.Ordenador[0].Documento)
			body_cambio_estado.EstadoCumplidoId = estado_cumplido
			body_cambio_estado.CumplidoProveedorId = cumplido_proveedor
			body_cambio_estado.DocumentoResponsable = documento_ordenador
			body_cambio_estado.CargoResponsable = ordenador_contrato.Contratos.Ordenador[0].RolOrdenador
		} else {
			outputError = map[string]interface{}{"funcion": "CambioEstadoCumplido", "status": "404", "mensaje": "Error al obtener el ordenador del contrato"}
			return body_cambio_estado, outputError
		}
	case "RO":
		ordenador_contrato, err := helpers.ObtenerOrdenadorContrato(cumplido_proveedor.NumeroContrato, strconv.Itoa(cumplido_proveedor.VigenciaContrato))
		if err == nil {
			documento_ordenador, _ := strconv.Atoi(ordenador_contrato.Contratos.Ordenador[0].Documento)
			body_cambio_estado.EstadoCumplidoId = estado_cumplido
			body_cambio_estado.CumplidoProveedorId = cumplido_proveedor
			body_cambio_estado.DocumentoResponsable = documento_ordenador
			body_cambio_estado.CargoResponsable = ordenador_contrato.Contratos.Ordenador[0].RolOrdenador
		} else {
			outputError = map[string]interface{}{"funcion": "CambioEstadoCumplido", "status": "404", "mensaje": "Error al obtener el ordenador del contrato"}
			return body_cambio_estado, outputError
		}
	default:
		outputError = map[string]interface{}{"funcion": "CambioEstadoCumplido", "status": "404", "mensaje": "El código de abreviación no es válido"}
		return body_cambio_estado, outputError
	}
	return body_cambio_estado, outputError

}

func ObtenerSupervisorContrato(numero_contrato_suscrito string, vigencia string) (supervisor_contrato models.SupervisorContratoProveedor, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ObtenerSupervisorContrato", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	var respuesta_peticion map[string]interface{}

	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/informacion_supervisor_contrato/"+numero_contrato_suscrito+"/"+vigencia, &respuesta_peticion); err == nil && response == 200 {
		if respuesta_peticion == nil {
			outputError = map[string]interface{}{"funcion": "ObtenerSupervisorContrato", "status": "404", "mensaje": "No se ha registrado un supervisor para el contrato"}
			return supervisor_contrato, outputError
		}
		json_supervisor, err_json := json.Marshal(respuesta_peticion)
		if err_json == nil {
			err := json.Unmarshal(json_supervisor, &supervisor_contrato)
			if err != nil {
				outputError = map[string]interface{}{"funcion": "ObtenerSupervisorContrato", "status": "404", "mensaje": "Error al convertir el json"}
				return supervisor_contrato, outputError
			}
		} else {
			outputError = map[string]interface{}{"funcion": "ObtenerSupervisorContrato", "status": "404", "mensaje": "Error al convertir el json"}
			return supervisor_contrato, outputError
		}
	} else {
		outputError = map[string]interface{}{"funcion": "ObtenerSupervisorContrato", "status": "404", "mensaje": "Error al consultar el supervisor del contrato"}
		return supervisor_contrato, outputError
	}
	return supervisor_contrato, outputError
}

func DesactivarCambiosAnterioresCumplido(cumplido_proveedor_id int, codigo_abreviacion_cumplido string) (ultimo_cambio_cumplido models.CambioEstadoCumplido, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "DesactivarCambiosAnterioresCumplido", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	posibles_estados_siguientes := map[string][]string{
		"CD":  {"PRC"},
		"PRC": {"RC", "AC"},
		"RC":  {"CD"},
		"AC":  {"PRO"},
		"PRO": {"AO", "RO"},
		"AO":  {"RO"},
		"RO":  {"CD"},
	}

	contains := func(posibles_estados []string, estado string) bool {
		for _, a := range posibles_estados {
			if a == estado {
				return true
			}
		}
		return false
	}

	var respuesta_peticion map[string]interface{}
	var cambios_anteriores []models.CambioEstadoCumplido
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/?query=CumplidoProveedorId.Id:"+strconv.Itoa(cumplido_proveedor_id)+",Activo:true&sortby=FechaCreacion&order=desc", &respuesta_peticion); err == nil && response == 200 {
		data := respuesta_peticion["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) == 0 {
			outputError = nil
			return
		}
		helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &cambios_anteriores)

		if !contains(posibles_estados_siguientes[cambios_anteriores[0].EstadoCumplidoId.CodigoAbreviacion], codigo_abreviacion_cumplido) {
			outputError = map[string]interface{}{"funcion": "DesactivarCambiosAnterioresCumplido", "status": "404", "mensaje": "No es posible pasar del estado " + cambios_anteriores[0].EstadoCumplidoId.CodigoAbreviacion + " al estado " + codigo_abreviacion_cumplido}
			return ultimo_cambio_cumplido, outputError
		}

		if (cambios_anteriores[0] != models.CambioEstadoCumplido{}) {
			ultimo_cambio_cumplido = cambios_anteriores[0]
			for _, cambio_anterior := range cambios_anteriores {
				var respuesta_estado_anterior map[string]interface{}
				cambio_anterior.Activo = false
				err := helpers.SendJson(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/"+strconv.Itoa(cambio_anterior.Id), "PUT", &respuesta_estado_anterior, cambio_anterior)
				if err != nil {
					outputError = map[string]interface{}{"funcion": "DesactivarCambiosAnterioresCumplido", "status": "404", "mensaje": "Error al actualizar el estado de los cambios de estados anteriores"}
					return ultimo_cambio_cumplido, outputError
				}
			}
		}

	} else {
		outputError = map[string]interface{}{"funcion": "DesactivarCambiosAnterioresCumplido", "status": "404", "mensaje": "Error al consultar el estado del cumplido"}
		return ultimo_cambio_cumplido, outputError
	}
	return ultimo_cambio_cumplido, outputError
}
