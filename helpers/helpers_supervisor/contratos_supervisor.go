package helpers_supervisor

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ContratosSupervisor(documento_supervisor string) (contratos_supervisor models.ContratoSupervisor, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ContratosSupervisor", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	if dependencias_supervisor, outputError := GetDependenciasSupervisor(documento_supervisor); outputError == nil {
		for _, dependencia := range dependencias_supervisor {
			contratos_supervisor.Dependencias_supervisor = append(contratos_supervisor.Dependencias_supervisor, dependencia)
			if contratos_dependencia, outputError := GetContratosDependenciaFiltroTemp(dependencia.Codigo, "2024-05", "2024-05"); outputError == nil {
				for _, contrato := range contratos_dependencia.Contratos.Contrato {
					contrato_contratista, err := GetInformacionContratoContratista(contrato.NumeroContrato, contrato.Vigencia)
					if err == nil {
						contratos_supervisor.NombreSupervisor = contrato_contratista.InformacionContratista.Supervisor.Nombre
						contratistas, err := ContratosContratista(contrato_contratista.InformacionContratista.Documento.Numero)
						if err == nil {

							for _, contratista := range contratistas {
								contratos_supervisor.Contratos = append(contratos_supervisor.Contratos, contratista)
							}
						} else {
							logs.Error(err)
							outputError = map[string]interface{}{"funcion": "/ContratosSupervisor/ContratosContratista", "err": err, "status": "502"}
							return contratos_supervisor, outputError
						}
					} else {
						logs.Error(err)
						outputError = map[string]interface{}{"funcion": "/ContratosSupervisor/ContratosContratista", "err": err, "status": "502"}
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

func GetDependenciasSupervisor(documento string) (dependenciasList []models.Dependencia, errorOutput interface{}) {
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
	response, err := getJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaProduccionJBPM")+"/dependencias_supervisor/"+documento, &respuesta)
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

func GetContrato(num_contrato_suscrito string, vigencia string) (informacion_contrato models.InformacionContrato, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetContrato", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var temp map[string]interface{}
	//fmt.Println(beego.AppConfig.String("UrlAdministrativaJBPM") + "/" + "contrato/" + num_contrato_suscrito + "/" + vigencia)
	if response, err := getJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/"+"contrato/"+num_contrato_suscrito+"/"+vigencia, &temp); (err == nil) && (response == 200) {
		json_contrato, error_json := json.Marshal(temp)
		if error_json == nil {
			var contrato models.InformacionContrato
			if err := json.Unmarshal(json_contrato, &contrato); err == nil {
				informacion_contrato = contrato
				//Se valida si esta vacio el objeto
				if informacion_contrato == (models.InformacionContrato{}) {
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "/GetContrato/EmptyResponse", "err": err, "status": "502"}
					return informacion_contrato, outputError
				}
				return informacion_contrato, nil
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/GetContrato", "err": err, "status": "502"}
				return informacion_contrato, outputError
			}
		} else {
			logs.Error(error_json.Error())
			outputError = map[string]interface{}{"funcion": "/GetContrato", "err": error_json.Error(), "status": "502"}
			return informacion_contrato, outputError
		}

	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/GetContrato", "err": err, "status": "502"}
		return informacion_contrato, outputError
	}

	return informacion_contrato, nil
}

func GetContratosPersona(num_documento string) (contratos_persona models.InformacionContratosPersona, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetContratosPersona", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var temp map[string]interface{}
	var contratos models.InformacionContratosPersona
	fmt.Println(beego.AppConfig.String("UrlAdministrativaJBPM") + "/contratos_contratista/" + num_documento)
	if response, err := getJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/contratos_contratista/"+num_documento, &temp); (err == nil) && (response == 200) {
		json_contratos, error_json := json.Marshal(temp)
		if error_json == nil {
			err := json.Unmarshal(json_contratos, &contratos)
			if err == nil {
				contratos_persona = contratos
				//fmt.Println("Contratos personas", contratos_persona)
				return contratos_persona, nil
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/GetContratosPersona", "err": err, "status": "502"}
				return contratos_persona, outputError
			}

		} else {
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/GetContratosPersona", "err": error_json.Error(), "status": "502"}
			return contratos_persona, outputError
		}

	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/GetContratosPersona", "err": err, "status": "502"}
		return contratos_persona, outputError
	}

	return contratos_persona, nil

}

func GetInformacionContrato(numContrato string, vigencia string) (informacion_contrato models.InformacionContrato, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetContrato", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var temp map[string]interface{}

	if response, err := getJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaProduccionJBPM")+"/"+"contrato/"+numContrato+"/"+vigencia, &temp); (err == nil) && (response == 200) {
		json_contrato, error_json := json.Marshal(temp)
		if error_json == nil {
			var contrato models.InformacionContrato
			if err := json.Unmarshal(json_contrato, &contrato); err == nil {
				informacion_contrato = contrato
				//Se valida si esta vacio el objeto
				if informacion_contrato == (models.InformacionContrato{}) {
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "/GetContrato/EmptyResponse", "err": err, "status": "502"}
					return informacion_contrato, outputError
				}
				return informacion_contrato, nil
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/GetContrato", "err": err, "status": "502"}
				return informacion_contrato, outputError
			}
		} else {
			logs.Error(error_json.Error())
			outputError = map[string]interface{}{"funcion": "/GetContrato", "err": error_json.Error(), "status": "502"}
			return informacion_contrato, outputError
		}

	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/GetContrato", "err": err, "status": "502"}
		return informacion_contrato, outputError
	}

	return informacion_contrato, nil
}

func GetContratosDependenciaFiltro(dependencia string, documento_supervisor string) (contratos_dependencia models.ContratoDependencia, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetContratosDependenciaFiltro", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var contratos_general []models.ContratoGeneral

	//fmt.Println("URL: ", beego.AppConfig.String("UrlcrudAgora")+"/contrato_general/?query=Supervisor.Documento:"+documento_supervisor+",Supervisor.DependenciaSupervisor:"+dependencia+",TipoContrato.Id.in:9|12|14|15|7|5")
	if response, err := getJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/contrato_general/?query=Supervisor.Documento:"+documento_supervisor+",Supervisor.DependenciaSupervisor:"+dependencia+",TipoContrato.Id.in:9|12|14|15|7|5", &contratos_general); (err == nil) && (response == 200) {
		if len(contratos_general) > 0 {
			for _, contrato_general := range contratos_general {
				if len(contrato_general.ContratoSuscrito) > 0 {
					var contrato models.Contrato
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
		outputError = map[string]interface{}{"funcion": "/GetContratosDependenciaFiltro", "err": err, "status": "502"}
		return contratos_dependencia, outputError
	}

	return

}

func GetContratosDependenciaFiltroTemp(dependencia string, fecha_inicio string, fecha_fin string) (contratos_dependencia models.ContratoDependencia, outputError map[string]interface{}) {
	var temp map[string]interface{}
	if response, err := getJsonWSO2Test(beego.AppConfig.String("UrlHomologacionDepsJBPM")+"/oikos_argo/"+dependencia, &temp); (err == nil) && (response == 200) {
		json_dep_oikos, error_json := json.Marshal(temp)
		if error_json == nil {
			var depOikos models.HomologacionDepOikos
			if err := json.Unmarshal(json_dep_oikos, &depOikos); err == nil {

				if len(depOikos.Dependencias.Dependencia) != 0 {
					if response, err := getJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/contratos_dependencia_oikos/"+depOikos.Dependencias.Dependencia[0].IDMaster+"/"+fecha_inicio+"/"+fecha_fin, &temp); (err == nil) && (response == 200) {
						json_contrato, error_json := json.Marshal(temp)
						if error_json == nil {
							if err := json.Unmarshal(json_contrato, &contratos_dependencia); err == nil {
								return contratos_dependencia, nil
							} else {
								logs.Error(err)
								outputError = map[string]interface{}{"funcion": "/GetContratosDependenciaFiltro/contratos_dependencia_oikos", "err": err.Error(), "status": "502"}
								return contratos_dependencia, outputError

							}
						} else {
							logs.Error(error_json)
							outputError = map[string]interface{}{"funcion": "/GetContratosDependenciaFiltro/contratos_dependencia_oikos", "err": error_json.Error(), "status": "502"}
							return contratos_dependencia, outputError
						}

					} else {
						logs.Error(err)
						outputError = map[string]interface{}{"funcion": "/GetContratosDependenciaFiltro/contratos_dependencia_oikos", "err": err.Error(), "status": "502"}
						return contratos_dependencia, outputError
					}
				} else {
					outputError = map[string]interface{}{"funcion": "/GetContratosDependenciaFiltro/oikos_argo", "err": "no hay dependencia homologada en oikos", "status": "502"}
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

func GetRP(numero_cdp string, vigencia_cdp string) (rp models.InformacionCdpRp, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetRP0", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var temp map[string]interface{}
	var temp_cdp_rp models.InformacionCdpRp
	//fmt.Println(beego.AppConfig.String("UrlFinancieraJBPM") + "/" + "cdprp/" + numero_cdp + "/" + vigencia_cdp + "/01")
	if response, err := getJsonWSO2Test(beego.AppConfig.String("UrlFinancieraJBPM")+"/"+"cdprp/"+numero_cdp+"/"+vigencia_cdp+"/01", &temp); (err == nil) && (response == 200) {
		json_cdp_rp, error_json := json.Marshal(temp)

		if error_json == nil {
			if err := json.Unmarshal(json_cdp_rp, &temp_cdp_rp); err == nil {
				rp = temp_cdp_rp
				return rp, nil
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/GetRP1", "err": err, "status": "502"}
				return rp, outputError
			}
		} else {
			logs.Error(error_json)
			outputError = map[string]interface{}{"funcion": "/GetRP2", "err": error_json, "status": "502"}
			return rp, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/GetRP3", "err": err, "status": "502"}
		return rp, outputError
	}
	return rp, outputError
}

func GetInformacionContratoContratista(num_contrato_suscrito string, vigencia string) (informacion_contrato_contratista models.InformacionContratoContratista, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetInformacionContratoContratista", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var temp map[string]interface{}
	//fmt.Println("URL GetInformacionContratoContratista", beego.AppConfig.String("UrlAdministrativaJBPM")+"/"+"informacion_contrato_contratista/"+num_contrato_suscrito+"/"+vigencia)
	if response, err := getJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/"+"informacion_contrato_contratista/"+num_contrato_suscrito+"/"+vigencia, &temp); (err == nil) && (response == 200) {
		json_contrato, error_json := json.Marshal(temp)
		if error_json == nil {
			var contrato_contratista models.InformacionContratoContratista
			if err := json.Unmarshal(json_contrato, &contrato_contratista); err == nil {
				informacion_contrato_contratista = contrato_contratista
				return informacion_contrato_contratista, nil
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/GetInformacionContratoContratista", "err": err, "status": "502"}
				return informacion_contrato_contratista, outputError
			}
		} else {
			logs.Error(error_json.Error())
			outputError = map[string]interface{}{"funcion": "/GetInformacionContratoContratista", "err": error_json.Error(), "status": "502"}
			return informacion_contrato_contratista, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/getInformacionContratosContratista", "err": err, "status": "502"}
		return informacion_contrato_contratista, outputError
	}
}

func ContratosContratista(numero_documento string) (contrato_proveedor []models.InformacionContratoProveedor, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			//fmt.Println("error", err)
			outputError = map[string]interface{}{"funcion": "/ContratosContratista", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	if contratos_persona, outputError := GetContratosPersona(numero_documento); outputError == nil {
		for _, contrato_persona := range contratos_persona.ContratosPersonas.ContratoPersona {
			contrato_persona.FechaInicio = time.Date(contrato_persona.FechaInicio.Year(), contrato_persona.FechaInicio.Month(), contrato_persona.FechaInicio.Day(), 0, 0, 0, 0, contrato_persona.FechaInicio.Location())
			contrato_persona.FechaFin = time.Date(contrato_persona.FechaFin.Year(), contrato_persona.FechaFin.Month(), contrato_persona.FechaFin.Day(), 0, 0, 0, 0, contrato_persona.FechaFin.Location())
			if time.Now().Before(contrato_persona.FechaFin) {
				var contrato models.InformacionContrato
				contrato, outputError = GetContrato(contrato_persona.NumeroContrato, contrato_persona.Vigencia)

				if (contrato == models.InformacionContrato{} || outputError != nil) {
					continue
				}
				var informacion_contrato_contratista models.InformacionContratoContratista
				informacion_contrato_contratista, outputError = GetInformacionContratoContratista(contrato_persona.NumeroContrato, contrato_persona.Vigencia)
				// se llena el contrato original en el indice 0

				if cdprp, outputError := GetRP(contrato_persona.NumeroCDP, contrato_persona.Vigencia); outputError == nil {
					for _, rp := range cdprp.CdpXRp.CdpRp {
						var tipo_contrato models.TipoContrato
						var contrato_proveedor_individual models.InformacionContratoProveedor
						contrato_proveedor_individual.TipoContrato = contrato.Contrato.TipoContrato
						contrato_proveedor_individual.NumeroContratoSuscrito = contrato_persona.NumeroContrato
						contrato_proveedor_individual.Vigencia = contrato_persona.Vigencia
						contrato_proveedor_individual.NumeroRp = rp.RpNumeroRegistro
						contrato_proveedor_individual.VigenciaRp = rp.RpVigencia
						contrato_proveedor_individual.NombreProveedor = informacion_contrato_contratista.InformacionContratista.NombreCompleto
						contrato_proveedor_individual.NombreDependencia = informacion_contrato_contratista.InformacionContratista.Dependencia
						contrato_proveedor_individual.NumDocumentoSupervisor = contrato.Contrato.Supervisor.DocumentoIdentificacion
						contrato_proveedor_individual.NumeroCdp = contrato_persona.NumeroCDP
						contrato_proveedor_individual.VigenciaCdp = contrato_persona.Vigencia
						contrato_proveedor_individual.Rubro = contrato.Contrato.Rubro
						if response, err := getJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/tipo_contrato/"+contrato.Contrato.TipoContrato, &tipo_contrato); err == nil && response == 200 {
							contrato_proveedor_individual.TipoContrato = tipo_contrato.TipoContrato
						} else {
							logs.Error(err)
							outputError = map[string]interface{}{"funcion": "/contratosContratista/GetContratosPersona", "err": err, "status": "502"}
							return nil, outputError

						}
						contrato_proveedor = append(contrato_proveedor, contrato_proveedor_individual)
					}

				} else {
					logs.Error(outputError)
					continue
				}

			}

		}
	} else {
		logs.Error(outputError)
		outputError = map[string]interface{}{"funcion": "/contratosContratista/GetContratosPersona", "err": outputError, "status": "502"}
		return nil, outputError
	}
	return contrato_proveedor, nil
}

func GetActaInicio(numero_contrato_suscrito string, vigencia_contrato int) (acta_inicio models.ActaInicio, outputError map[string]interface{}) {
	var contratos_suscrito []models.ContratoSuscrito
	if response, err := getJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/contrato_suscrito/?query=NumeroContratoSuscrito:"+numero_contrato_suscrito+",Vigencia:"+strconv.Itoa(vigencia_contrato), &contratos_suscrito); (err == nil) && (response == 200) {
		if len(contratos_suscrito) > 0 {
			var actasInicio []models.ActaInicio
			if response, err := getJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/acta_inicio/?query=NumeroContrato:"+contratos_suscrito[0].NumeroContrato.Id+",Vigencia:"+strconv.Itoa(contratos_suscrito[0].Vigencia), &actasInicio); (err == nil) && (response == 200) {
				if len(actasInicio) == 0 {
					acta_inicio.Id = 0
					acta_inicio.NumeroContrato = contratos_suscrito[0].NumeroContratoSuscrito
					acta_inicio.Vigencia = contratos_suscrito[0].Vigencia
					acta_inicio.FechaInicio = contratos_suscrito[0].FechaSuscripcion
					acta_inicio.Descripcion = "No se ha registrado acta de inicio"

					switch contratos_suscrito[0].NumeroContrato.UnidadEjecucion.Id {
					case 205:
						acta_inicio.FechaFin = acta_inicio.FechaInicio.AddDate(0, 0, contratos_suscrito[0].NumeroContrato.PlazoEjecucion)
					case 206:
						acta_inicio.FechaFin = acta_inicio.FechaInicio.AddDate(0, contratos_suscrito[0].NumeroContrato.PlazoEjecucion, 0)
					case 207:
						acta_inicio.FechaFin = acta_inicio.FechaInicio.AddDate(contratos_suscrito[0].NumeroContrato.PlazoEjecucion, 0, 0)
					default:
						outputError = map[string]interface{}{"funcion": "/getInformacionContratosContratista", "message": "La unidad de ejecucuion no es un valor de tiempo", "status": "502"}
						return acta_inicio, outputError
					}
				} else {
					acta_inicio = actasInicio[0]
					return acta_inicio, nil
				}
			} else {
				outputError = map[string]interface{}{"funcion": "/GetActaInicio", "err": err, "message": "Error al obtener la acta de inicio", "status": "502"}
				return acta_inicio, outputError
			}
		} else {
			outputError = map[string]interface{}{"funcion": "/GetActaInicio", "err": err, "message": "No existe el contrato suscrito ingresado", "status": "502"}
			return acta_inicio, outputError
		}
	} else {
		outputError = map[string]interface{}{"funcion": "/GetActaInicio", "err": err, "message": "Error al obtener el contrato suscrito", "status": "502"}
		return acta_inicio, outputError
	}

	return acta_inicio, outputError
}

func FechasContratoConNovedades(numero_contrato string, vigencia_contrato string, numero_cdp string, num_doc string) (fechas models.FechasConNovedades, outputError map[string]interface{}) {

	if contratos_persona, err := GetContratosPersona(num_doc); err == nil {
		for _, contrato := range contratos_persona.ContratosPersonas.ContratoPersona {
			if contrato.NumeroContrato == numero_contrato && contrato.Vigencia == vigencia_contrato && contrato.NumeroCDP == numero_cdp {
				fechas.FechaInicio = contrato.FechaInicio
				fechas.FechaFin = contrato.FechaFin
				return fechas, nil
			}
		}
		outputError = map[string]interface{}{"funcion": "/FechasContratoConNovedades", "err": "No se encontro el contrato", "status": "502"}
		return fechas, outputError
	} else {
		outputError = map[string]interface{}{"funcion": "/FechasContratoConNovedades/GetContratoPersona", "err": err, "status": "502"}
		return fechas, outputError
	}

	return fechas, outputError
}
