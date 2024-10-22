package helpers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtenerInformacionContratoProveedor(numero_contrato_suscrito string, vigencia string) (contratos_proveedor []models.InformacionContratoProveedor, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
		}
	}()

	var respuesta_peticion map[string]interface{}
	//fmt.Println(beego.AppConfig.String("UrlAdministrativaJBPM") + "/informacion_contrato_proveedor/" + numero_contrato_suscrito + "/" + vigencia)
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
									RPFechaRegistro:        rp.CdpXRp.CdpRp[0].RpFechaRegistro,
									NombreProveedor:        contratoMap["proveedor"].(string),
									NombreDependencia:      contratoMap["dependencia"].(string),
									NumeroCdp:              contratoMap["numero_cdp"].(string),
									VigenciaCdp:            contratoMap["vigencia_cdp"].(string),
									CDPFechaExpedicion:     rp.CdpXRp.CdpRp[0].CdpFechaExpedicion,
									Rubro:                  contratoMap["rubro"].(string),
								}
								contratos_proveedor = append(contratos_proveedor, contrato_proveedor)
							} else {
								logs.Error(err)
								outputError = fmt.Errorf("Error al obtener la información del RP asociado al contrato con número %s y vigencia %s", contratoMap["numero_contrato_suscrito"], contratoMap["vigencia"])
								return contratos_proveedor, outputError
							}
						}
					} else {
						logs.Error("Error al procesar la lista de contratos del proveedor")
						outputError = fmt.Errorf("Error al procesar la lista de contratos del proveedor")
						return contratos_proveedor, outputError
					}
				}
			} else {
				logs.Error("Error al extraer los datos del proveedor en la respuesta")
				outputError = fmt.Errorf("Error al extraer los datos del proveedor en la respuesta")
				return contratos_proveedor, outputError
			}
		} else {
			logs.Error("No se encontro información del contrato del proveedor")
			outputError = fmt.Errorf("No se encontro información del contrato del proveedor")
			return contratos_proveedor, outputError
		}
	} else {
		logs.Error(err)
		outputError = fmt.Errorf("Error al obtener los datos del servicio para el contrato número %s con vigencia %s", numero_contrato_suscrito, vigencia)
		return contratos_proveedor, outputError
	}

	return contratos_proveedor, nil
}

func ObtenerOrdenadorContrato(numero_contrato_suscrito string, vigencia string) (ordenador_contrato models.OrdenadorContratoProveedor, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var respuesta_peticion map[string]interface{}

	fmt.Println("URL ordenador contrato: ", beego.AppConfig.String("UrlAdministrativaJBPM")+"/informacion_ordenador_contrato/"+numero_contrato_suscrito+"/"+vigencia)
	if response, err := GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/informacion_ordenador_contrato/"+numero_contrato_suscrito+"/"+vigencia, &respuesta_peticion); err == nil && response == 200 {
		if respuesta_peticion == nil {
			outputError = fmt.Errorf("No se encontro el ordenador del contrato")
			return ordenador_contrato, outputError
		}
		json_ordenador, err_json := json.Marshal(respuesta_peticion)
		if err_json == nil {
			err := json.Unmarshal(json_ordenador, &ordenador_contrato)
			if err != nil {
				outputError = fmt.Errorf("Error al convertir el json")
				return ordenador_contrato, outputError
			}
		} else {
			outputError = fmt.Errorf("Error al convertir el json")
			return ordenador_contrato, outputError
		}
	} else {
		outputError = fmt.Errorf("Error al consultar el ordenador del contrato")
		return ordenador_contrato, outputError
	}
	return ordenador_contrato, outputError
}

func ObtenerRP(numero_cdp string, vigencia_cdp string) (rp models.InformacionCdpRp, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var temp map[string]interface{}
	var temp_cdp_rp models.InformacionCdpRp
	//fmt.Println(beego.AppConfig.String("UrlFinancieraJBPM") + "/" + "cdprp/" + numero_cdp + "/" + vigencia_cdp + "/01")
	if response, err := GetJsonWSO2Test(beego.AppConfig.String("UrlFinancieraJBPM")+"/cdprp/"+numero_cdp+"/"+vigencia_cdp+"/01", &temp); (err == nil) && (response == 200) {
		if temp == nil {
			outputError = fmt.Errorf("No se encontro el RP")
			return rp, outputError
		}
		json_cdp_rp, error_json := json.Marshal(temp)

		if error_json == nil {
			if err := json.Unmarshal(json_cdp_rp, &temp_cdp_rp); err == nil {
				rp = temp_cdp_rp
				return rp, nil
			} else {
				logs.Error(err)
				outputError = fmt.Errorf("Error al convertir a Json")
				return rp, outputError
			}
		} else {
			logs.Error(error_json)
			outputError = fmt.Errorf("Error al convertir a Json")
			return rp, outputError
		}
	} else {
		logs.Error(err)
		outputError = fmt.Errorf("Error al consultar el RP del contrato")
		return rp, outputError
	}

}

func ObtenerActaInicio(numero_contrato_suscrito string, vigencia_contrato int) (acta_inicio models.ActaInicio, outputError error) {
	var contratos_suscrito []models.ContratoSuscrito
	if response, err := GetJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/contrato_suscrito/?query=NumeroContratoSuscrito:"+numero_contrato_suscrito+",Vigencia:"+strconv.Itoa(vigencia_contrato)+"&limit=-1", &contratos_suscrito); (err == nil) && (response == 200) {
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
								outputError = fmt.Errorf("La unidad de ejecucion no es un valor de tiempo")
								return acta_inicio, outputError
							}
							return acta_inicio, nil
						} else {
							outputError = fmt.Errorf("No existe el contrato general ingresado")
							return acta_inicio, outputError
						}
					} else {
						outputError = fmt.Errorf("Error al obtener el contrato general")
						return acta_inicio, outputError
					}
				} else {
					acta_inicio = actasInicio[0]
				}
			} else {
				outputError = fmt.Errorf("Error al obtener la acta de inicio")
			}
		} else {
			outputError = fmt.Errorf("No existe el contrato suscrito ingresado")
			return acta_inicio, outputError
		}
	} else {
		outputError = fmt.Errorf("Error al obtener el contrato suscrito")
		return acta_inicio, outputError
	}

	return acta_inicio, outputError
}

func ObtenerContratoGeneralProveedor(numero_contrato_suscrito string, vigencia_contrato string) (contrato_general models.ContratoGeneral, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var contrato []models.ContratoGeneral
	//fmt.Println(beego.AppConfig.String("UrlcrudAgora") + "/contrato_general/?query=ContratoSuscrito.NumeroContratoSuscrito:" + numero_contrato_suscrito + ",VigenciaContrato:" + vigencia_contrato)
	if response, err := GetJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/contrato_general/?query=ContratoSuscrito.NumeroContratoSuscrito:"+numero_contrato_suscrito+",VigenciaContrato:"+vigencia_contrato, &contrato); (err == nil) && (response == 200) {
		if len(contrato) > 0 {
			return contrato[0], nil
		} else {
			outputError = fmt.Errorf("No se encontró contrato")
			return contrato[0], outputError
		}
	} else {
		outputError = fmt.Errorf("Error al obtener el contrato general")
		return contrato_general, outputError
	}
}
