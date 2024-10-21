package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func CambioEstadoCumplido(codigo_abreviacion_cumplido string, cumplido_proveedor_id int) (respuesta_cambio_estado models.CambioEstadoCumplidoResponse, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	// Verificar que se envien todos los datos
	if codigo_abreviacion_cumplido == "" || cumplido_proveedor_id == 0 {
		outputError = fmt.Errorf("Error en los datos de entrada")
		return respuesta_cambio_estado, outputError
	}

	var cumplido_proveedor []models.CumplidoProveedor
	var estado_cumplido []models.EstadoCumplido
	var respuesta_cumplido_proveedor map[string]interface{}
	var respuesta_estado_cumplido map[string]interface{}

	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cumplido_proveedor/?query=Id:"+strconv.Itoa(cumplido_proveedor_id)+"&limit=-1", &respuesta_cumplido_proveedor); (err == nil) && (response == 200) {
		data := respuesta_cumplido_proveedor["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) > 0 {
			helpers.LimpiezaRespuestaRefactor(respuesta_cumplido_proveedor, &cumplido_proveedor)

			if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/estado_cumplido/?query=CodigoAbreviacion:"+codigo_abreviacion_cumplido+"&limit=-1", &respuesta_estado_cumplido); (err == nil) && (response == 200) {
				data := respuesta_estado_cumplido["Data"].([]interface{})
				if len(data[0].(map[string]interface{})) == 0 {
					outputError = fmt.Errorf("El estado del cumplido proveedor ingresado no existe")
					return respuesta_cambio_estado, outputError
				}
				helpers.LimpiezaRespuestaRefactor(respuesta_estado_cumplido, &estado_cumplido)
			} else {
				outputError = fmt.Errorf("Error al consultar el estado del cumplido proveedor")
				return respuesta_cambio_estado, outputError
			}

		} else {
			outputError = fmt.Errorf("El cumplido proveedor ingresado no existe")
			return respuesta_cambio_estado, outputError
		}
	} else {
		outputError = fmt.Errorf("Error al consultar el cumplido proveedor")
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
		outputError = fmt.Errorf("El cumplido ya se encuentra en este estado")
		return respuesta_cambio_estado, outputError
	} else {
		var respuesta_peticion map[string]interface{}
		if err := helpers.SendJson(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido", "POST", &respuesta_peticion, body_cambio_estado); err == nil {
			respuesta_cambio_estado.CumplidoProveedorId = cumplido_proveedor_id
			respuesta_cambio_estado.DocumentoResponsable = body_cambio_estado.DocumentoResponsable
			respuesta_cambio_estado.CargoResponsable = body_cambio_estado.CargoResponsable
			respuesta_cambio_estado.EstadoCumplido = &estado_cumplido[0]

			if codigo_abreviacion_cumplido != "CD" {
				EnviarNotificacionCambioEstado(estado_cumplido[0].Nombre, strconv.Itoa(body_cambio_estado.DocumentoResponsable), strconv.Itoa(ultimo_cambio_estado_cumplido.DocumentoResponsable), cumplido_proveedor[0].NumeroContrato, cumplido_proveedor[0].VigenciaContrato)

			}
			if codigo_abreviacion_cumplido == "AC" {
				respuesta_cambio_estado, outputError = CambioEstadoCumplido("PRO", cumplido_proveedor_id)
			}
			return respuesta_cambio_estado, outputError
		} else {
			outputError = fmt.Errorf("Error al actualizar el cumplido proveedor")
			return respuesta_cambio_estado, outputError
		}
	}
}

func CrearBodyCambioEstadoCumplido(codigo_abreviacion_cumplido string, cumplido_proveedor models.CumplidoProveedor, estado_cumplido models.EstadoCumplido) (body_cambio_estado models.BodyCambioEstadoCumplido, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()
	error_supervisor_contrato := fmt.Errorf("Error al obtener el supervisor del contrato")
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
			outputError = error_supervisor_contrato
			return body_cambio_estado, outputError
		}
	case "PRC":
		body_cambio_estado.EstadoCumplidoId = estado_cumplido
		body_cambio_estado.CumplidoProveedorId = cumplido_proveedor
		body_cambio_estado.DocumentoResponsable = 0
		body_cambio_estado.CargoResponsable = "CONTRATACIÓN"
	case "RC":
		supervisor_contrato, err := ObtenerSupervisorContrato(cumplido_proveedor.NumeroContrato, strconv.Itoa(cumplido_proveedor.VigenciaContrato))
		if err == nil {
			documento_supervisor, _ := strconv.Atoi(supervisor_contrato.Contratos.Supervisor[0].Documento)
			body_cambio_estado.EstadoCumplidoId = estado_cumplido
			body_cambio_estado.CumplidoProveedorId = cumplido_proveedor
			body_cambio_estado.DocumentoResponsable = documento_supervisor
			body_cambio_estado.CargoResponsable = supervisor_contrato.Contratos.Supervisor[0].Cargo
		} else {
			outputError = error_supervisor_contrato
			return body_cambio_estado, outputError
		}
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
			outputError = fmt.Errorf("Error al obtener el ordenador del contrato")
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
			outputError = fmt.Errorf("Error al obtener el ordenador del contrato")
			return body_cambio_estado, outputError
		}
	case "RO":
		supervisor_contrato, err := ObtenerSupervisorContrato(cumplido_proveedor.NumeroContrato, strconv.Itoa(cumplido_proveedor.VigenciaContrato))
		if err == nil {
			documento_supervisor, _ := strconv.Atoi(supervisor_contrato.Contratos.Supervisor[0].Documento)
			body_cambio_estado.EstadoCumplidoId = estado_cumplido
			body_cambio_estado.CumplidoProveedorId = cumplido_proveedor
			body_cambio_estado.DocumentoResponsable = documento_supervisor
			body_cambio_estado.CargoResponsable = supervisor_contrato.Contratos.Supervisor[0].Cargo
		} else {
			outputError = error_supervisor_contrato
			return body_cambio_estado, outputError
		}
	default:
		outputError = fmt.Errorf("El código de abreviación no es válido")
		return body_cambio_estado, outputError
	}
	return body_cambio_estado, outputError

}

func ObtenerSupervisorContrato(numero_contrato_suscrito string, vigencia string) (supervisor_contrato models.SupervisorContratoProveedor, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	error_json := fmt.Errorf("Error al convertir el json")
	var respuesta_peticion map[string]interface{}

	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/informacion_supervisor_contrato/"+numero_contrato_suscrito+"/"+vigencia, &respuesta_peticion); err == nil && response == 200 {
		if respuesta_peticion == nil {
			outputError = fmt.Errorf("No se ha registrado un supervisor para el contrato")
			return supervisor_contrato, outputError
		}
		json_supervisor, err_json := json.Marshal(respuesta_peticion)
		if err_json == nil {
			err := json.Unmarshal(json_supervisor, &supervisor_contrato)
			if err != nil {
				outputError = error_json
				return supervisor_contrato, outputError
			}
		} else {
			outputError = error_json
			return supervisor_contrato, outputError
		}
	} else {
		outputError = fmt.Errorf("Error al consultar el supervisor del contrato")
		return supervisor_contrato, outputError
	}
	return supervisor_contrato, outputError
}

func DesactivarCambiosAnterioresCumplido(cumplido_proveedor_id int, codigo_abreviacion_cumplido string) (ultimo_cambio_cumplido models.CambioEstadoCumplido, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	posibles_estados_siguientes := map[string][]string{
		"CD":  {"PRC"},
		"PRC": {"RC", "AC"},
		"RC":  {"PRC"},
		"AC":  {"PRO"},
		"PRO": {"AO", "RO"},
		"AO":  {"RO"},
		"RO":  {"PRC"},
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
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/?query=CumplidoProveedorId.Id:"+strconv.Itoa(cumplido_proveedor_id)+",Activo:true&sortby=FechaCreacion&order=desc&limit=-1", &respuesta_peticion); err == nil && response == 200 {
		data := respuesta_peticion["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) == 0 {
			outputError = nil
			return
		}
		helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &cambios_anteriores)

		if !contains(posibles_estados_siguientes[cambios_anteriores[0].EstadoCumplidoId.CodigoAbreviacion], codigo_abreviacion_cumplido) {
			outputError = fmt.Errorf("No es posible pasar del estado " + cambios_anteriores[0].EstadoCumplidoId.CodigoAbreviacion + " al estado " + codigo_abreviacion_cumplido)
			return ultimo_cambio_cumplido, outputError
		}

		if (cambios_anteriores[0] != models.CambioEstadoCumplido{}) {
			ultimo_cambio_cumplido = cambios_anteriores[0]
			for _, cambio_anterior := range cambios_anteriores {
				var respuesta_estado_anterior map[string]interface{}
				cambio_anterior.Activo = false
				err := helpers.SendJson(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/"+strconv.Itoa(cambio_anterior.Id), "PUT", &respuesta_estado_anterior, cambio_anterior)
				if err != nil {
					outputError = fmt.Errorf("Error al actualizar el estado de los cambios de estados anteriores")
					return ultimo_cambio_cumplido, outputError
				}
			}
		}

	} else {
		outputError = fmt.Errorf("Error al consultar el estado del cumplido")
		return ultimo_cambio_cumplido, outputError
	}
	return ultimo_cambio_cumplido, outputError
}

func EnviarNotificacionCambioEstado(nombre_estado string, documento_responsable string, documento_responsable_anterior string, num_contrato_suscrito string, vigencia int) (outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()
	error_json := fmt.Errorf("Error al convertir el json")
	var responsable_anterior string
	var email string

	if documento_responsable == "0" {
		email = "fctrujilloo@udistrital.edu.co"
	} else {
		var autenticacion_persona models.AutenticacionPersona
		body_autenticacion := map[string]interface{}{
			"numero": documento_responsable,
		}
		var respuesta_peticion map[string]interface{}
		// Se busca con el numero de documento del responsable los datos para recuperar el correo electronico
		if err := helpers.SendJson(beego.AppConfig.String("UrlAutenticacionMid")+"/token/documentoToken", "POST", &respuesta_peticion, body_autenticacion); err == nil {
			json_autenticacion, err := json.Marshal(respuesta_peticion)
			if err != nil {
				outputError = error_json
				return outputError
			}
			err = json.Unmarshal(json_autenticacion, &autenticacion_persona)
			if err != nil {
				outputError = error_json
				return outputError
			}

			if autenticacion_persona.Email == "" {
				outputError = fmt.Errorf("No se encontró el correo electrónico del responsable")
				return outputError
			}
			//email = autenticacion_persona.Email
			email = "fctrujilloo@udistrital.edu.co"
		} else {
			outputError = fmt.Errorf("Error al consultar el responsable")
			return outputError
		}
	}

	if documento_responsable_anterior != "0" {
		// Se busca con el numero de documento del responsable anterior los datos para recuperar el nombre
		var informacion_personal []models.InformacionProveedor
		//fmt.Println(beego.AppConfig.String("UrlcrudAgora") + "/informacion_proveedor/?query=NumDocumento:" + documento_responsable_anterior)
		if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlcrudAgora")+"/informacion_proveedor/?query=NumDocumento:"+documento_responsable_anterior, &informacion_personal); err == nil && response == 200 {
			if len(informacion_personal) == 0 {
				outputError = fmt.Errorf("No se encontraron datos del responsable anterior")
				return outputError
			}
			responsable_anterior = informacion_personal[0].NomProveedor
		} else {
			outputError = fmt.Errorf("Error al consultar el responsable anterior")
			return outputError
		}
	} else {
		responsable_anterior = "Contratación"
	}

	// Se busca con el numero de contrato suscrito y la vigencia los datos para recuperar el nombre del proveedor
	informacion_contrato, outputError := helpers.ObtenerInformacionContratoProveedor(num_contrato_suscrito, strconv.Itoa(vigencia))
	if outputError != nil {
		outputError = fmt.Errorf("Error al consultar el proveedor")
		return outputError
	}

	// Se crea el body y se hace el Post para enviar un correo electronico
	var body_enviar_notificacion models.NotificacionEmail
	body_enviar_notificacion.Source = "notificacionescumplidosproveedores@udistrital.edu.co"
	body_enviar_notificacion.Template = "REVISION_CUMPLIDOS_PROVEEDORES_PLANTILLA"
	body_enviar_notificacion.Destinations = make([]models.DestinationItem, 1)
	body_enviar_notificacion.Destinations[0].Destination.ToAddresses = append(body_enviar_notificacion.Destinations[0].Destination.ToAddresses, email)
	body_enviar_notificacion.Destinations[0].ReplacementTemplateData.EstadoCumplido = nombre_estado
	body_enviar_notificacion.Destinations[0].ReplacementTemplateData.NombreEmisorCumplido = responsable_anterior
	body_enviar_notificacion.Destinations[0].ReplacementTemplateData.NombreProveedor = informacion_contrato[0].NombreProveedor
	body_enviar_notificacion.Destinations[0].ReplacementTemplateData.RevisionCumplidosProveedoresUrl = beego.AppConfig.String("UrlRevisionCumplidosProveedoresCliente")
	body_enviar_notificacion.Destinations[0].Attachments = []string{}
	body_enviar_notificacion.Destinations[0].Destination.BccAddresses = []string{}
	body_enviar_notificacion.Destinations[0].Destination.CcAddresses = []string{}
	body_enviar_notificacion.DefaultTemplateData.EstadoCumplido = nombre_estado
	body_enviar_notificacion.DefaultTemplateData.NombreEmisorCumplido = responsable_anterior
	body_enviar_notificacion.DefaultTemplateData.NombreProveedor = informacion_contrato[0].NombreProveedor
	body_enviar_notificacion.DefaultTemplateData.RevisionCumplidosProveedoresUrl = beego.AppConfig.String("UrlRevisionCumplidosProveedoresCliente")

	var respuesta map[string]interface{}
	fmt.Println(beego.AppConfig.String("UrlNotificacionesMid") + "/email/enviar_templated_email")
	if err := helpers.SendJsonTls(beego.AppConfig.String("UrlNotificacionesMid")+"/email/enviar_templated_email", "POST", &respuesta, body_enviar_notificacion); err != nil {
		jsonData, err := json.MarshalIndent(body_enviar_notificacion, "", "    ")
		if err != nil {
			log.Fatalf("Error al convertir a JSON: %s", err)
		}
		fmt.Println(string(jsonData))
		outputError = fmt.Errorf("Error al enviar la notificacion al correo")
		return outputError
	}

	return
}
