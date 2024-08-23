package helpers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtenerContratosProveedor(numero_documento string) (contrato_proveedor []models.InformacionContratoProveedor, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			//fmt.Println("error", err)
			outputError = map[string]interface{}{"funcion": "/ObtenerContratosProveedor", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	if contratos_persona, outputError := ObtenerContratosPersona(numero_documento); outputError == nil {
		for _, contrato_persona := range contratos_persona.ContratosPersonas.ContratoPersona {
			contrato_persona.FechaFin = time.Date(contrato_persona.FechaFin.Year(), contrato_persona.FechaFin.Month(), contrato_persona.FechaFin.Day(), 0, 0, 0, 0, contrato_persona.FechaFin.Location())
			if time.Now().Before(contrato_persona.FechaFin) {
				var contrato models.InformacionContrato
				contrato, outputError = ObtenerInformacionContrato(contrato_persona.NumeroContrato, contrato_persona.Vigencia)

				if (contrato == models.InformacionContrato{} || outputError != nil) {
					continue
				}
				var informacion_contrato_proveedor models.DatosContratoProveedor
				informacion_contrato_proveedor, outputError = ObtenerInformacionContratoProveedor(contrato_persona.NumeroContrato, contrato_persona.Vigencia)
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
						contrato_proveedor_individual.NombreProveedor = informacion_contrato_proveedor.InformacionContratista.NombreCompleto
						contrato_proveedor_individual.NombreDependencia = informacion_contrato_proveedor.InformacionContratista.Dependencia
						contrato_proveedor_individual.NumDocumentoSupervisor = contrato.Contrato.Supervisor.DocumentoIdentificacion
						contrato_proveedor_individual.NumeroCdp = contrato_persona.NumeroCDP
						contrato_proveedor_individual.VigenciaCdp = contrato_persona.Vigencia
						contrato_proveedor_individual.Rubro = contrato.Contrato.Rubro
						if response, err := GetJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/tipo_contrato/"+contrato.Contrato.TipoContrato, &tipo_contrato); err == nil && response == 200 {
							contrato_proveedor_individual.TipoContrato = tipo_contrato.TipoContrato
						} else {
							logs.Error(err)
							outputError = map[string]interface{}{"funcion": "/ObtenerContratosProveedor/GetContratosPersona", "err": err, "status": "502"}
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
		outputError = map[string]interface{}{"funcion": "/ObtenerContratosProveedor/GetContratosPersona", "err": outputError, "status": "502"}
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

	//var temp map[string]interface{}
	var contratos []models.ContratoGeneral
	var proveedor []models.InformacionProveedor
	var contrato_disponibilidad []models.ContratoDisponibilidad

	if response, err := GetJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/informacion_proveedor/?query=NumDocumento:"+num_documento, &proveedor); err == nil && response == 200 {
		if len(proveedor) > 0 {
			if response, err := GetJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/contrato_general/?query=Contratista:"+strconv.Itoa(proveedor[0].Id)+"&sortby=FechaRegistro&order=desc", &contratos); err == nil && response == 200 {
				acta_inicio, err := ObtenerActaInicio(contratos[0].ContratoSuscrito[0].NumeroContratoSuscrito, contratos[0].ContratoSuscrito[0].Vigencia)
				if err == nil {
					if response, err := GetJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/contrato_disponibilidad/?query=NumeroContrato:"+contratos[0].Id+",Vigencia:"+strconv.Itoa(contratos[0].VigenciaContrato), &contrato_disponibilidad); err == nil && response == 200 {
						contratoPersona := struct {
							NumeroContrato string    `json:"numero_contrato"`
							Vigencia       string    `json:"vigencia"`
							NumeroCDP      string    `json:"cdp"`
							FechaInicio    time.Time `json:"fecha_inicio"`
							FechaFin       time.Time `json:"fecha_fin"`
						}{
							NumeroContrato: contratos[0].ContratoSuscrito[0].NumeroContratoSuscrito,
							Vigencia:       strconv.Itoa(contratos[0].ContratoSuscrito[0].Vigencia),
							NumeroCDP:      strconv.Itoa(contrato_disponibilidad[0].NumeroCdp),
							FechaInicio:    acta_inicio.FechaInicio,
							FechaFin:       acta_inicio.FechaFin,
						}
						contratos_persona.ContratosPersonas.ContratoPersona = append(contratos_persona.ContratosPersonas.ContratoPersona, contratoPersona)
						fmt.Println("contratos_persona", contratos_persona)

					} else {
						outputError = map[string]interface{}{"funcion": "/ObtenerContratosPersona", "message": "Error al obtener la información del contrato de disponibilidad", "err": err, "status": "502"}
						return contratos_persona, outputError
					}

				} else {
					outputError = map[string]interface{}{"funcion": "/ObtenerContratosPersona", "message": "Error al obtener la información del acta de inicio", "err": err, "status": "502"}
					return contratos_persona, outputError
				}
			} else {
				outputError = map[string]interface{}{"funcion": "/ObtenerContratosPersona", "message": "Error al la informacion del contrato del proveedor", "err": err, "status": "502"}
				return contratos_persona, outputError
			}
		}
	} else {
		outputError = map[string]interface{}{"funcion": "/ObtenerContratosPersona", "message": "Error al obtener la información del proveedor", "err": err, "status": "502"}
		return contratos_persona, outputError
	}

	return contratos_persona, nil

}

func ObtenerInformacionContratoProveedor(num_contrato_suscrito string, vigencia string) (informacion_contrato_proveedor models.DatosContratoProveedor, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ObtenerInformacionContratoProveedor", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var temp map[string]interface{}
	//fmt.Println("URL GetInformacionContratoContratista", beego.AppConfig.String("UrlAdministrativaJBPM")+"/"+"informacion_contrato_contratista/"+num_contrato_suscrito+"/"+vigencia)
	if response, err := GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/informacion_contrato_contratista/"+num_contrato_suscrito+"/"+vigencia, &temp); (err == nil) && (response == 200) {
		json_contrato, error_json := json.Marshal(temp)
		if error_json == nil {
			var contrato_proveedor models.DatosContratoProveedor
			if err := json.Unmarshal(json_contrato, &contrato_proveedor); err == nil {
				informacion_contrato_proveedor = contrato_proveedor
				return informacion_contrato_proveedor, nil
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/ObtenerInformacionContratoProveedor", "err": err, "status": "502"}
				return informacion_contrato_proveedor, outputError
			}
		} else {
			logs.Error(error_json.Error())
			outputError = map[string]interface{}{"funcion": "/ObtenerInformacionContratoProveedor", "err": error_json.Error(), "status": "502"}
			return informacion_contrato_proveedor, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ObtenerInformacionContratoProveedor", "err": err, "status": "502"}
		return informacion_contrato_proveedor, outputError
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

func ObtenerActaInicio(numero_contrato_suscrito string, vigencia_contrato int) (acta_inicio models.ActaInicio, outputError map[string]interface{}) {
	var contratos_suscrito []models.ContratoSuscrito
	if response, err := GetJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/contrato_suscrito/?query=NumeroContratoSuscrito:"+numero_contrato_suscrito+",Vigencia:"+strconv.Itoa(vigencia_contrato), &contratos_suscrito); (err == nil) && (response == 200) {
		if len(contratos_suscrito) > 0 {
			var actasInicio []models.ActaInicio
			if response, err := GetJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/acta_inicio/?query=NumeroContrato:"+contratos_suscrito[0].NumeroContrato.Id+",Vigencia:"+strconv.Itoa(contratos_suscrito[0].Vigencia), &actasInicio); (err == nil) && (response == 200) {
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
						outputError = map[string]interface{}{"funcion": "/ObtenerActaInicio", "message": "La unidad de ejecucuion no es un valor de tiempo", "status": "502"}
						return acta_inicio, outputError
					}
				} else {
					acta_inicio = actasInicio[0]
					return acta_inicio, nil
				}
			} else {
				outputError = map[string]interface{}{"funcion": "/ObtenerActaInicio", "err": err, "message": "Error al obtener la acta de inicio", "status": "502"}
				return acta_inicio, outputError
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
