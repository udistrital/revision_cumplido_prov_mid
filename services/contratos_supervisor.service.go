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
			outputError = map[string]interface{}{"funcion": "/ContratosSupervisor", "err": err, "status": "502"}
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
					informacion_contrato_proveedor, err := ObtenerInformacionContratoProveedor(contrato.NumeroContrato, contrato.Vigencia)
					if err == nil {
						for _, contrato_proveedor := range informacion_contrato_proveedor {
							contratos_supervisor.Contratos = append(contratos_supervisor.Contratos, contrato_proveedor)
						}
					} else {
						logs.Error(err)
						outputError = map[string]interface{}{"funcion": "/ContratosSupervisor", "err": err, "status": "502"}
						return contratos_supervisor, outputError
					}
				}
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/ContratosSupervisor", "err": err, "status": "502"}
				return contratos_supervisor, outputError
			}
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ContratosSupervisor", "err": err, "status": "502"}
		return contratos_supervisor, outputError
	}
	return contratos_supervisor, nil
}

func ObtenerInformacionContratoProveedor(numero_contrato_suscrito string, vigencia string) (contratos_proveedor []models.InformacionContratoProveedor, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"Success": false,
				"Status":  502,
				"Message": "Error al obtener la informacion del contrato del proveedor",
				"Error":   err,
			}
		}
	}()

	var respuesta_peticion map[string]interface{}
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/informacion_contrato_proveedor/"+numero_contrato_suscrito+"/"+vigencia, &respuesta_peticion); err == nil && response == 200 {
		if respuesta_peticion != nil {
			if contratosMap, exito := respuesta_peticion["proveedor"].(map[string]interface{}); exito {
				for _, contratoList := range contratosMap {
					if list, exito := contratoList.([]interface{}); exito {
						for _, contrato := range list {
							contratoMap := contrato.(map[string]interface{})
							rp, err := helpers.ObtenerRP(contratoMap["numero_cdp"].(string), contratoMap["vigencia_cdp"].(string))
							if err == nil {
								contrato_proveedor := models.InformacionContratoProveedor{
									TipoContrato:           contratoMap["tipo_contrato"].(string),
									NumeroContratoSuscrito: contratoMap["numero_contrato_suscrito"].(string),
									Vigencia:               contratoMap["vigencia"].(string),
									NumeroRp:               rp.CdpXRp.CdpRp[0].RpNumeroRegistro,
									VigenciaRp:             rp.CdpXRp.CdpRp[0].RpVigencia,
									NombreProveedor:        contratoMap["proveedor"].(string),
									NombreDependencia:      contratoMap["dependencia"].(string),
									NumeroCdp:              contratoMap["numero_cdp"].(string),
									VigenciaCdp:            contratoMap["vigencia_cdp"].(string),
									Rubro:                  contratoMap["rubro"].(string),
								}
								contratos_proveedor = append(contratos_proveedor, contrato_proveedor)
							} else {
								logs.Error(err)
								outputError = map[string]interface{}{"funcion": "/ObtenerInformacionContratoProveedor/", "err": err, "status": "502"}
								return contratos_proveedor, outputError
							}
						}
					} else {
						logs.Error("Error al obtener la informacion del contrato del proveedor")
						outputError = map[string]interface{}{"funcion": "/ObtenerInformacionContratoProveedor/", "err": "Error al obtener la informacion del contrato del proveedor", "status": "502"}
						return contratos_proveedor, outputError
					}
				}
			} else {
				logs.Error("Error al obtener la informacion del contrato del proveedor")
				outputError = map[string]interface{}{"funcion": "/ObtenerInformacionContratoProveedor/", "err": "Error al obtener la informacion del contrato del proveedor", "status": "502"}
				return contratos_proveedor, outputError
			}
		} else {
			logs.Error("Error al obtener la informacion del contrato del proveedor")
			outputError = map[string]interface{}{"funcion": "/ObtenerInformacionContratoProveedor/", "err": "Error al obtener la informacion del contrato del proveedor", "status": "502"}
			return contratos_proveedor, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ObtenerInformacionContratoProveedor/", "err": err, "status": "502"}
		return contratos_proveedor, outputError
	}

	return contratos_proveedor, nil
}

func ObtenerContratosDependencia(dependencia string) (contratos_dependencia models.ContratoDependencia, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"Success": false,
				"Status":  502,
				"Message": "Error al obtener los contratos de la dependencia: " + dependencia,
				"Error":   err,
			}
		}
	}()

	var respuesta_peticion map[string]interface{}
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/contratos_proveedor_dependencia/"+dependencia, &respuesta_peticion); err == nil && response == 200 {
		respuesta_json, err_json := json.Marshal(respuesta_peticion)
		if err_json == nil {
			if err := json.Unmarshal(respuesta_json, &contratos_dependencia); err == nil {
				return contratos_dependencia, nil
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/ObtenerContratosDependencia/", "err": err.Error(), "status": "502"}
				return contratos_dependencia, outputError
			}
		} else {
			logs.Error(err_json)
			outputError = map[string]interface{}{"funcion": "/ObtenerContratosDependencia/", "err": err_json.Error(), "status": "502"}
			return contratos_dependencia, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ObtenerContratosDependencia/", "err": err, "status": "502"}
		return contratos_dependencia, outputError
	}

}

func ObtenerDependenciasSupervisor(documento_supervisor string) (dependencias_supervisor []models.Dependencia, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"Success": false,
				"Status":  502,
				"Message": "Error al consultar las dependencias del supervisor identificado con el documento: " + documento_supervisor,
				"Error":   err,
			}
		}
	}()

	var respuesta_peticion map[string]interface{}
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

					}
				}
			}
		}
	}
	return dependencias_supervisor, nil
}
