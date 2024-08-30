package helpers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

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
	fmt.Println(beego.AppConfig.String("UrlAdministrativaJBPM") + "/informacion_contrato_proveedor/" + numero_contrato_suscrito + "/" + vigencia)
	if response, err := GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/informacion_contrato_proveedor/"+numero_contrato_suscrito+"/"+vigencia, &respuesta_peticion); err == nil && response == 200 {
		if respuesta_peticion != nil {
			if contratosMap, exito := respuesta_peticion["proveedor"].(map[string]interface{}); exito {
				for _, contratoList := range contratosMap {
					if list, exito := contratoList.([]interface{}); exito {
						for _, contrato := range list {
							contratoMap := contrato.(map[string]interface{})
							rp, err := ObtenerRP(contratoMap["numero_cdp"].(string), contratoMap["vigencia_cdp"].(string))
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

func ObtenerOrdenadorContrato(numero_contrato_suscrito string, vigencia string) (ordenador_contrato models.OrdenadorContratoProveedor, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ObtenerOrdenadorContrato", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var respuesta_peticion map[string]interface{}

	if response, err := GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/informacion_ordenador_contrato/"+numero_contrato_suscrito+"/"+vigencia, &respuesta_peticion); err == nil && response == 200 {
		json_ordenador, err_json := json.Marshal(respuesta_peticion)
		if err_json == nil {
			err := json.Unmarshal(json_ordenador, &ordenador_contrato)
			if err != nil {
				outputError = map[string]interface{}{"funcion": "ObtenerOrdenadorContrato", "status": "502", "mensaje": "Error al convertir el json"}
				return ordenador_contrato, outputError
			}
		} else {
			outputError = map[string]interface{}{"funcion": "ObtenerOrdenadorContrato", "status": "502", "mensaje": "Error al convertir el json"}
			return ordenador_contrato, outputError
		}
	} else {
		outputError = map[string]interface{}{"funcion": "ObtenerOrdenadorContrato", "status": "502", "mensaje": "Error al consultar el ordenador del contrato"}
		return ordenador_contrato, outputError
	}
	return ordenador_contrato, outputError
}

func ObtenerRP(numero_cdp string, vigencia_cdp string) (rp models.InformacionCdpRp, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetRP0", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var temp map[string]interface{}
	var temp_cdp_rp models.InformacionCdpRp
	fmt.Println(beego.AppConfig.String("UrlFinancieraJBPM") + "/" + "cdprp/" + numero_cdp + "/" + vigencia_cdp + "/01")
	if response, err := GetJsonWSO2Test(beego.AppConfig.String("UrlFinancieraJBPM")+"/cdprp/"+numero_cdp+"/"+vigencia_cdp+"/01", &temp); (err == nil) && (response == 200) {
		if temp == nil {
			outputError = map[string]interface{}{"funcion": "/GetRP", "err": "No se encontro el RP", "status": "404"}
			return rp, outputError
		}
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

}

func ObtenerActaInicio(numero_contrato_suscrito string, vigencia_contrato int) (acta_inicio models.ActaInicio, outputError map[string]interface{}) {
	var contratos_suscrito []models.ContratoSuscrito
	if response, err := GetJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/contrato_suscrito/?query=NumeroContratoSuscrito:"+numero_contrato_suscrito+",Vigencia:"+strconv.Itoa(vigencia_contrato), &contratos_suscrito); (err == nil) && (response == 200) {
		if len(contratos_suscrito) > 0 {
			var actasInicio []models.ActaInicio
			//fmt.Println("URL acta inicio: ", beego.AppConfig.String("UrlcrudAgora")+"/acta_inicio/?query=NumeroContrato:"+contratos_suscrito[0].NumeroContrato.Id+",Vigencia:"+strconv.Itoa(contratos_suscrito[0].Vigencia))
			if response, err := GetJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/acta_inicio/?query=NumeroContrato:"+contratos_suscrito[0].NumeroContrato.Id+",Vigencia:"+strconv.Itoa(contratos_suscrito[0].Vigencia), &actasInicio); (err == nil) && (response == 200) {
				if len(actasInicio) == 0 {
					var contratos_generales []models.ContratoGeneral
					if respuesta, err := GetJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/contrato_general/?query=Id:"+contratos_suscrito[0].NumeroContrato.Id+",VigenciaContrato:"+strconv.Itoa(contratos_suscrito[0].Vigencia), &contratos_generales); (err == nil) && (respuesta == 200) {
						//fmt.Println("Contrato general: ", contratos_generales)
						if len(contratos_generales) > 0 {
							acta_inicio.Id = 0
							acta_inicio.NumeroContrato = contratos_suscrito[0].NumeroContratoSuscrito
							acta_inicio.Vigencia = contratos_suscrito[0].Vigencia
							acta_inicio.FechaInicio = contratos_suscrito[0].FechaSuscripcion
							acta_inicio.Descripcion = "No se ha registrado acta de inicio"

							switch contratos_generales[0].UnidadEjecucion.Id {
							case 205:
								acta_inicio.FechaFin = acta_inicio.FechaInicio.AddDate(0, 0, contratos_generales[0].PlazoEjecucion)
							case 206:
								acta_inicio.FechaFin = acta_inicio.FechaInicio.AddDate(0, contratos_generales[0].PlazoEjecucion, 0)
							case 207:
								acta_inicio.FechaFin = acta_inicio.FechaInicio.AddDate(contratos_generales[0].PlazoEjecucion, 0, 0)
							default:
								outputError = map[string]interface{}{"funcion": "/ObtenerActaInicio", "message": "La unidad de ejecucuion no es un valor de tiempo", "status": "502"}
								return acta_inicio, outputError
							}
							return acta_inicio, nil
						} else {
							outputError = map[string]interface{}{"funcion": "/ObtenerActaInicio", "message": "No existe el contrato general ingresado", "status": "502"}
							return acta_inicio, outputError
						}
					} else {
						outputError = map[string]interface{}{"funcion": "/ObtenerActaInicio", "message": "Error al obtener el contrato general", "status": "502"}
						return acta_inicio, outputError
					}
				} else {
					acta_inicio = actasInicio[0]
				}
			} else {
				outputError = map[string]interface{}{"funcion": "/ObtenerActaInicio", "err": err, "message": "Error al obtener la acta de inicio", "status": "502"}
			}
		} else {
			outputError = map[string]interface{}{"funcion": "/ObtenerActaInicio", "err": err, "message": "No existe el contrato suscrito ingresado", "status": "502"}
			return acta_inicio, outputError
		}
	} else {
		outputError = map[string]interface{}{"funcion": "/ObtenerActaInicio", "err": err, "message": "Error al obtener el contrato suscrito", "status": "502"}
		return acta_inicio, outputError
	}

	return acta_inicio, outputError
}

func ObtenerContratoGeneralProveedor(numero_contrato_suscrito string, vigencia_contrato string) (contrato_general models.ContratoGeneral, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ObtenerContratoGeneralProveedor", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var contrato []models.ContratoGeneral
	if response, err := GetJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/contrato_general/?query=ContratoSuscrito.NumeroContratoSuscrito:"+numero_contrato_suscrito+",VigenciaContrato:"+vigencia_contrato, &contrato); (err == nil) && (response == 200) {
		if len(contrato) > 0 {
			return contrato[0], nil
		} else {
			outputError = map[string]interface{}{"funcion": "/ObtenerContratoGeneralProveedor", "err": "No se encontró contrato", "status": "404"}
			return contrato[0], outputError
		}
	}
	return contrato_general, outputError
}
