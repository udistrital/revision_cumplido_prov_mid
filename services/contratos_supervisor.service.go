package services

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtenerContratosSupervisor(documento_supervisor string) (contratos_supervisor models.ContratoSupervisor, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ContratosSupervisor", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	dependencias_supervisor, err := ObtenerDependenciasSupervisor(documento_supervisor)
	if err == nil {
		for _, dependencia := range dependencias_supervisor {
			contratos_supervisor.Dependencias_supervisor = append(contratos_supervisor.Dependencias_supervisor, dependencia)
			contratos_dependencia, err := ObtenerContratosDependencia(dependencia.Codigo)
			if err == nil {
				for _, contrato := range contratos_dependencia.Contratos.Contrato {
					informacion_contrato_proveedor, err := helpers.ObtenerInformacionContratoProveedor(contrato.NumeroContrato, contrato.Vigencia)
					if err == nil {
						contratos_supervisor.Contratos = append(contratos_supervisor.Contratos, informacion_contrato_proveedor...)
					} else {
						logs.Error(err)
						continue
					}
				}
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/ContratosSupervisor", "err": err, "status": "404"}
				return contratos_supervisor, outputError
			}
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ContratosSupervisor", "err": err, "status": "404"}
		return contratos_supervisor, outputError
	}
	return contratos_supervisor, nil
}

func ObtenerContratosDependencia(dependencia string) (contratos_dependencia models.ContratoDependencia, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"Success": false,
				"Status":  404,
				"Message": "Error al obtener los contratos de la dependencia: " + dependencia,
				"Error":   err,
			}
		}
	}()

	var respuesta_peticion map[string]interface{}
	//fmt.Println(beego.AppConfig.String("UrlAdministrativaJBPM") + "/contratos_proveedor_dependencia/" + dependencia)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/contratos_proveedor_dependencia/"+dependencia, &respuesta_peticion); err == nil && response == 200 {
		if respuesta_peticion != nil {
			respuesta_json, err_json := json.Marshal(respuesta_peticion)
			if err_json == nil {
				if err := json.Unmarshal(respuesta_json, &contratos_dependencia); err == nil {
					return contratos_dependencia, nil
				} else {
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "/ObtenerContratosDependencia/", "err": err.Error(), "status": "404"}
					return contratos_dependencia, outputError
				}
			} else {
				logs.Error(err_json)
				outputError = map[string]interface{}{"funcion": "/ObtenerContratosDependencia/", "err": err_json.Error(), "status": "404"}
				return contratos_dependencia, outputError
			}
		} else {
			outputError = map[string]interface{}{"funcion": "/ObtenerContratosDependencia/", "err": "No se encontraron contratos para la dependencia: " + dependencia, "status": "404"}
			return contratos_dependencia, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ObtenerContratosDependencia/", "err": err, "status": "404"}
		return contratos_dependencia, outputError
	}

}

func ObtenerDependenciasSupervisor(documento_supervisor string) (dependencias_supervisor []models.Dependencia, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"Success": false,
				"Status":  404,
				"Message": "Error al consultar las dependencias del supervisor identificado con el documento: " + documento_supervisor,
				"Error":   err,
			}
		}
	}()

	var respuesta_peticion map[string]interface{}
	//fmt.Println(beego.AppConfig.String("UrlAdministrativaProduccionJBPM") + "/dependencias_supervisor/" + documento_supervisor)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaProduccionJBPM")+"/dependencias_supervisor/"+documento_supervisor, &respuesta_peticion); err == nil && response == 200 {
		if respuesta_peticion != nil {
			if dependenciasMap, ok := respuesta_peticion["dependencias"].(map[string]interface{}); ok {

				for _, depList := range dependenciasMap {

					if list, ok := depList.([]interface{}); ok {

						for _, dep := range list {

							depMap := dep.(map[string]interface{})
							dependencia := models.Dependencia{

								Codigo: depMap["codigo"].(string),
								Nombre: depMap["nombre"].(string),
							}
							dependencias_supervisor = append(dependencias_supervisor, dependencia)
						}

					} else {
						outputError = map[string]interface{}{"funcion": "/ObtenerDependenciasSupervisor/", "err": "No se encontraron dependencias para el supervisor con documento: " + documento_supervisor, "status": "404"}
						return dependencias_supervisor, outputError
					}
				}
			}
		} else {
			outputError = map[string]interface{}{"funcion": "/ObtenerDependenciasSupervisor/", "err": "No se encontraron dependencias para el supervisor con documento: " + documento_supervisor, "status": "404"}
			return dependencias_supervisor, outputError
		}
	}
	return dependencias_supervisor, nil
}
