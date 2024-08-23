package services

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtenerContratosSupervisor(documento_supervisor string) (contratos_supervisor models.ContratoSupervisor, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ContratosSupervisor", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	if dependencias_supervisor, outputError := ObtenerDependenciasSupervisor(documento_supervisor); outputError == nil {
		for _, dependencia := range dependencias_supervisor {
			contratos_supervisor.Dependencias_supervisor = append(contratos_supervisor.Dependencias_supervisor, dependencia)
			if contratos_dependencia, outputError := ObtenerContratosDependenciaFiltroTemp(dependencia.Codigo, "2024-05", "2024-05"); outputError == nil {
				for _, contrato := range contratos_dependencia.Contratos.Contrato {
					contrato_proveedor, err := helpers.ObtenerInformacionContratoProveedor(contrato.NumeroContrato, contrato.Vigencia)
					if err == nil {
						contratos_supervisor.NombreSupervisor = contrato_proveedor.InformacionContratista.Supervisor.Nombre
						proveedores, err := helpers.ObtenerContratosProveedor(contrato_proveedor.InformacionContratista.Documento.Numero)
						if err == nil {

							for _, proveedor := range proveedores {
								contratos_supervisor.Contratos = append(contratos_supervisor.Contratos, proveedor)
							}
						} else {
							logs.Error(err)
							outputError = map[string]interface{}{"funcion": "/ContratosSupervisor/ObtenerContratosProveedor", "err": err, "status": "502"}
							return contratos_supervisor, outputError
						}
					} else {
						logs.Error(err)
						outputError = map[string]interface{}{"funcion": "/ContratosSupervisor/ObtenerContratosProveedor", "err": err, "status": "502"}
						return contratos_supervisor, outputError
					}
				}
			} else {
				logs.Error(outputError)
				outputError = map[string]interface{}{"funcion": "/ContratosSupervisor/GetContratosDependenciaFiltro", "err": outputError, "status": "502"}
				return contratos_supervisor, outputError
			}
		}
	} else {
		logs.Error(outputError)
		outputError = map[string]interface{}{"funcion": "/ContratosSupervisor/GetDependenciasSupervisor", "err": outputError, "status": "502"}
		return contratos_supervisor, nil

	}
	return contratos_supervisor, outputError
}

func ObtenerDependenciasSupervisor(documento string) (dependenciasList []models.Dependencia, errorOutput interface{}) {
	defer func() {
		if err := recover(); err != nil {
			errorOutput = map[string]interface{}{
				"Success": true,
				"Status":  502,
				"Message": "Error al consultar las dependencias: " + documento,
				"Error":   err,
			}
		}
	}()
	var respuesta map[string]interface{}
	//fmt.Println(beego.AppConfig.String("UrlAdministrativaProduccionJBPM") + "/dependencias_supervisor/" + documento)
	response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaProduccionJBPM")+"/dependencias_supervisor/"+documento, &respuesta)
	if err == nil && response == 200 {
		if respuesta != nil {
			if dependenciasMap, ok := respuesta["dependencias"].(map[string]interface{}); ok {

				for _, depList := range dependenciasMap {

					if list, ok := depList.([]interface{}); ok {

						for _, dep := range list {

							depMap := dep.(map[string]interface{})
							dependencia := models.Dependencia{

								Codigo: depMap["codigo"].(string),
								Nombre: depMap["nombre"].(string),
							}
							dependenciasList = append(dependenciasList, dependencia)
						}

					}
				}
			}
		}
	} else {
		return nil, errorOutput
	}
	if dependenciasList != nil {
		return dependenciasList, nil
	}

	return nil, nil
}

func ObtenerContratosDependenciaFiltro(dependencia string, documento_supervisor string) (contratos_dependencia models.ContratoDependencia, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ObtenerContratosDependenciaFiltro", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var contratos_general []models.ContratoGeneral

	//fmt.Println("URL: ", beego.AppConfig.String("UrlcrudAgora")+"/contrato_general/?query=Supervisor.Documento:"+documento_supervisor+",Supervisor.DependenciaSupervisor:"+dependencia+",TipoContrato.Id.in:9|12|14|15|7|5")
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/contrato_general/?query=Supervisor.Documento:"+documento_supervisor+",Supervisor.DependenciaSupervisor:"+dependencia+",TipoContrato.Id.in:9|12|14|15|7|5", &contratos_general); (err == nil) && (response == 200) {
		if len(contratos_general) > 0 {
			for _, contrato_general := range contratos_general {
				if len(contrato_general.ContratoSuscrito) > 0 {
					var contrato models.ContratoDep
					contrato.Vigencia = strconv.Itoa(contrato_general.ContratoSuscrito[0].Vigencia)
					contrato.NumeroContrato = contrato_general.ContratoSuscrito[0].NumeroContratoSuscrito
					contratos_dependencia.Contratos.Contrato = append(contratos_dependencia.Contratos.Contrato, contrato)
				} else {
					continue
				}
			}
		} else {
			return contratos_dependencia, nil
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ObtenerContratosDependenciaFiltro", "err": err, "status": "502"}
		return contratos_dependencia, outputError
	}

	return

}

func ObtenerContratosDependenciaFiltroTemp(dependencia string, fecha_inicio string, fecha_fin string) (contratos_dependencia models.ContratoDependencia, outputError map[string]interface{}) {
	var temp map[string]interface{}
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlHomologacionDepsJBPM")+"/oikos_argo/"+dependencia, &temp); (err == nil) && (response == 200) {
		json_dep_oikos, error_json := json.Marshal(temp)
		if error_json == nil {
			var depOikos models.HomologacionDepOikos
			if err := json.Unmarshal(json_dep_oikos, &depOikos); err == nil {

				if len(depOikos.Dependencias.Dependencia) != 0 {
					if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/contratos_dependencia_oikos/"+depOikos.Dependencias.Dependencia[0].IDMaster+"/"+fecha_inicio+"/"+fecha_fin, &temp); (err == nil) && (response == 200) {
						json_contrato, error_json := json.Marshal(temp)
						if error_json == nil {
							if err := json.Unmarshal(json_contrato, &contratos_dependencia); err == nil {
								return contratos_dependencia, nil
							} else {
								logs.Error(err)
								outputError = map[string]interface{}{"funcion": "/ObtenerContratosDependenciaFiltroTemp/contratos_dependencia_oikos", "err": err.Error(), "status": "502"}
								return contratos_dependencia, outputError

							}
						} else {
							logs.Error(error_json)
							outputError = map[string]interface{}{"funcion": "/ObtenerContratosDependenciaFiltroTemp/contratos_dependencia_oikos", "err": error_json.Error(), "status": "502"}
							return contratos_dependencia, outputError
						}

					} else {
						logs.Error(err)
						outputError = map[string]interface{}{"funcion": "/ObtenerContratosDependenciaFiltroTemp/contratos_dependencia_oikos", "err": err.Error(), "status": "502"}
						return contratos_dependencia, outputError
					}
				} else {
					outputError = map[string]interface{}{"funcion": "/ObtenerContratosDependenciaFiltroTemp/oikos_argo", "err": "no hay dependencia homologada en oikos", "status": "502"}
					return contratos_dependencia, outputError

				}

			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/GetContratosDependenciaFiltro/oikos_argo", "err": err.Error(), "status": "502"}
				return contratos_dependencia, outputError

			}
		} else {
			logs.Error(error_json)
			outputError = map[string]interface{}{"funcion": "/GetContratosDependenciaFiltro/oikos_argo", "err": error_json.Error(), "status": "502"}
			return contratos_dependencia, outputError
		}
	}
	return
}
