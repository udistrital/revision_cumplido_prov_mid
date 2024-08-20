package helpers

import (
	"encoding/json"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtenerContratosContratista(numero_documento string) (contrato_proveedor []models.InformacionContratoProveedor, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			//fmt.Println("error", err)
			outputError = map[string]interface{}{"funcion": "/ContratosContratista", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	if contratos_persona, outputError := ObtenerContratosPersona(numero_documento); outputError == nil {
		for _, contrato_persona := range contratos_persona.ContratosPersonas.ContratoPersona {
			/*
				vigencia, _ := strconv.Atoi(contrato_persona.Vigencia)
				fechaFin, err := ObtenerActaInicio(contrato_persona.NumeroContrato, vigencia)
				if err == nil {
				 *Toda la asignacion de datos del contrato
				}
			*/
			contrato_persona.FechaInicio = time.Date(contrato_persona.FechaInicio.Year(), contrato_persona.FechaInicio.Month(), contrato_persona.FechaInicio.Day(), 0, 0, 0, 0, contrato_persona.FechaInicio.Location())
			contrato_persona.FechaFin = time.Date(contrato_persona.FechaFin.Year(), contrato_persona.FechaFin.Month(), contrato_persona.FechaFin.Day(), 0, 0, 0, 0, contrato_persona.FechaFin.Location())
			if time.Now().Before(contrato_persona.FechaFin) {
				var contrato models.InformacionContrato
				contrato, outputError = ObtenerInformacionContrato(contrato_persona.NumeroContrato, contrato_persona.Vigencia)

				if (contrato == models.InformacionContrato{} || outputError != nil) {
					continue
				}
				var informacion_contrato_contratista models.InformacionContratoContratista
				informacion_contrato_contratista, outputError = ObtenerInformacionContratoContratista(contrato_persona.NumeroContrato, contrato_persona.Vigencia)
				// se llena el contrato original en el indice 0

				if cdprp, outputError := ObtenerRP(contrato_persona.NumeroCDP, contrato_persona.Vigencia); outputError == nil {
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
						if response, err := GetJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/tipo_contrato/"+contrato.Contrato.TipoContrato, &tipo_contrato); err == nil && response == 200 {
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

func ObtenerInformacionContrato(num_contrato_suscrito string, vigencia string) (informacion_contrato models.InformacionContrato, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetContrato", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var temp map[string]interface{}
	//fmt.Println(beego.AppConfig.String("UrlAdministrativaJBPM") + "/" + "contrato/" + num_contrato_suscrito + "/" + vigencia)
	if response, err := GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/contrato/"+num_contrato_suscrito+"/"+vigencia, &temp); (err == nil) && (response == 200) {
		json_contrato, error_json := json.Marshal(temp)
		if error_json == nil {
			var contrato models.InformacionContrato
			if err := json.Unmarshal(json_contrato, &contrato); err == nil {
				informacion_contrato = contrato
				//Se valida si esta vacio el objeto
				if informacion_contrato == (models.InformacionContrato{}) {
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "/ObtenerInformacionContrato/EmptyResponse", "err": err, "status": "502"}
					return informacion_contrato, outputError
				}
				return informacion_contrato, nil
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/ObtenerInformacionContrato", "err": err, "status": "502"}
				return informacion_contrato, outputError
			}
		} else {
			logs.Error(error_json.Error())
			outputError = map[string]interface{}{"funcion": "/ObtenerInformacionContrato", "err": error_json.Error(), "status": "502"}
			return informacion_contrato, outputError
		}

	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ObtenerInformacionContrato", "err": err, "status": "502"}
		return informacion_contrato, outputError
	}

	return informacion_contrato, nil
}

func ObtenerContratosPersona(num_documento string) (contratos_persona models.InformacionContratosPersona, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ObtenerContratosPersona", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var temp map[string]interface{}
	var contratos models.InformacionContratosPersona
	//fmt.Println(beego.AppConfig.String("UrlAdministrativaJBPM") + "/contratos_contratista/" + num_documento)
	if response, err := GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/contratos_contratista/"+num_documento, &temp); (err == nil) && (response == 200) {
		json_contratos, error_json := json.Marshal(temp)
		if error_json == nil {
			err := json.Unmarshal(json_contratos, &contratos)
			if err == nil {
				contratos_persona = contratos
				//fmt.Println("Contratos personas", contratos_persona)
				return contratos_persona, nil
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/ObtenerContratosPersona", "err": err, "status": "502"}
				return contratos_persona, outputError
			}

		} else {
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/ObtenerContratosPersona", "err": error_json.Error(), "status": "502"}
			return contratos_persona, outputError
		}

	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ObtenerContratosPersona", "err": err, "status": "502"}
		return contratos_persona, outputError
	}

	return contratos_persona, nil

}

func ObtenerInformacionContratoContratista(num_contrato_suscrito string, vigencia string) (informacion_contrato_contratista models.InformacionContratoContratista, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetInformacionContratoContratista", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var temp map[string]interface{}
	//fmt.Println("URL GetInformacionContratoContratista", beego.AppConfig.String("UrlAdministrativaJBPM")+"/"+"informacion_contrato_contratista/"+num_contrato_suscrito+"/"+vigencia)
	if response, err := GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/informacion_contrato_contratista/"+num_contrato_suscrito+"/"+vigencia, &temp); (err == nil) && (response == 200) {
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

func ObtenerRP(numero_cdp string, vigencia_cdp string) (rp models.InformacionCdpRp, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetRP0", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var temp map[string]interface{}
	var temp_cdp_rp models.InformacionCdpRp
	//fmt.Println(beego.AppConfig.String("UrlFinancieraJBPM") + "/" + "cdprp/" + numero_cdp + "/" + vigencia_cdp + "/01")
	if response, err := GetJsonWSO2Test(beego.AppConfig.String("UrlFinancieraJBPM")+"/cdprp/"+numero_cdp+"/"+vigencia_cdp+"/01", &temp); (err == nil) && (response == 200) {
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
