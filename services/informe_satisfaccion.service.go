package services

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/astaxie/beego"

	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtenerBalanceFinancieroContrato(numero_contrato_suscrito string, vigencia_contrato string) (balance_contrato models.BalanceContrato, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetBalanceContrato", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	informacion_contrato, err := helpers.ObtenerInformacionContratoProveedor(numero_contrato_suscrito, vigencia_contrato)
	if err == nil {
		informacion_contratista, err := helpers.ObtenerContratosProveedor(informacion_contrato.InformacionContratista.Documento.Numero)
		if err == nil {
			contrato, err := helpers.ObtenerInformacionContrato(numero_contrato_suscrito, vigencia_contrato)
			if err == nil {
				valor_girado, err := ObtenerValorGiradoPorCdp(informacion_contratista[0].NumeroCdp, informacion_contratista[0].VigenciaCdp, contrato.Contrato.UnidadEjecutora)
				if err == nil {
					total_contrato, err := strconv.ParseFloat(informacion_contrato.InformacionContratista.ValorContrato, 64)
					if err == nil {
						saldo_contrato := int(total_contrato) - valor_girado
						balance_contrato.TotalContrato = strings.Split(informacion_contrato.InformacionContratista.ValorContrato, ".")[0]
						balance_contrato.Saldo = strconv.Itoa(saldo_contrato)
						return balance_contrato, nil
					} else {
						outputError = map[string]interface{}{"funcion": "/GetBalanceContrato", "err": err, "status": "502"}
						return balance_contrato, outputError
					}

				} else {
					outputError = map[string]interface{}{"funcion": "/GetBalanceContrato", "err": err, "status": "502"}
					return balance_contrato, outputError
				}
			}
		} else {
			outputError = map[string]interface{}{"funcion": "/GetBalanceContrato", "err": err, "status": "502"}
			return balance_contrato, outputError
		}
	} else {
		outputError = map[string]interface{}{"funcion": "/GetBalanceContrato", "err": err, "status": "502"}
		return balance_contrato, outputError
	}

	return balance_contrato, outputError
}

func ObtenerValorGiradoPorCdp(cdp string, vigencia_cdp string, unidad_ejecucion string) (valor_girado int, err error) {
	var temp_giros_tercero map[string]interface{}
	var giros_tercero models.GirosTercero
	valor_girado = 0
	//fmt.Println(beego.AppConfig.String("UrlFinancieraJBPM") + "/" + "giros_tercero/" + cdp + "/" + vigencia_cdp + "/" + unidad_ejecucion)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlFinancieraJBPM")+"/giros_tercero/"+cdp+"/"+vigencia_cdp+"/"+unidad_ejecucion, &temp_giros_tercero); (err == nil) && (response == 200) {
		json_giros_tercero, error_json := json.Marshal(temp_giros_tercero)
		if error_json == nil {
			if err := json.Unmarshal(json_giros_tercero, &giros_tercero); err == nil {
				//fmt.Println("giros "+cdp, giros_tercero)
				for _, giro := range giros_tercero.Giros.Tercero {
					total_girado, err := strconv.Atoi(giro.ValorBrutoGirado)
					//fmt.Println(total_girado)
					if err == nil {
						valor_girado = valor_girado + total_girado
					}
				}
				//fmt.Println(valor_girado)
				return valor_girado, nil

			} else {
				err = errors.New("Error Unmarshal giros_tercero")
				return valor_girado, err
			}

		} else {
			err = errors.New("Error Marshal giros_tercero")
			return valor_girado, err
		}

	} else {
		return valor_girado, err
	}
	return
}

func ObtenerInformacionInformeSatisfaccion(numero_contrato_suscrito string, vigencia string) (informacion_informe models.InformacionInformeSatisfaccion, outputError map[string]interface{}) {

	contrato_contratista, err := helpers.ObtenerInformacionContratoProveedor(numero_contrato_suscrito, vigencia)
	if err == nil {
		contratistas, err := helpers.ObtenerContratosProveedor(contrato_contratista.InformacionContratista.Documento.Numero)
		if err == nil {
			informacion_informe.Dependencia = contratistas[0].NombreDependencia
			informacion_informe.NombreProveedor = contratistas[0].NombreProveedor
			informacion_informe.DocumentoProveedor = contrato_contratista.InformacionContratista.Documento.Numero
			informacion_informe.TipoContrato = contratistas[0].TipoContrato
			informacion_informe.Cdp = contratistas[0].NumeroCdp
			informacion_informe.VigenciaCdp = contratistas[0].VigenciaCdp
			informacion_informe.Rp = contratistas[0].NumeroRp
			informacion_informe.VigenciaRp = contratistas[0].VigenciaRp
			contratos_persona, err := helpers.ObtenerContratosPersona(contrato_contratista.InformacionContratista.Documento.Numero)
			if err == nil {
				informacion_informe.FechaInicio = contratos_persona.ContratosPersonas.ContratoPersona[0].FechaInicio
				informacion_informe.FechaFin = contratos_persona.ContratosPersonas.ContratoPersona[0].FechaFin
			} else {
				outputError = map[string]interface{}{"funcion": "/ObtenerInformacionInformeSatisfaccion", "err": err, "status": "502"}
				return informacion_informe, outputError
			}
			informacion_informe.Supervisor = contrato_contratista.InformacionContratista.Supervisor.Nombre
			informacion_informe.CargoSupervisor = contrato_contratista.InformacionContratista.Supervisor.Cargo
		} else {
			outputError = map[string]interface{}{"funcion": "/ObtenerInformacionInformeSatisfaccion", "err": err, "status": "502"}
			return informacion_informe, outputError
		}
	}
	return informacion_informe, nil
}

func ObtenerBanco(banco_id int) (banco models.Banco, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			//fmt.Println("error", err)
			outputError = map[string]interface{}{"funcion": "/GetBanco", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var respuesta_peticion map[string]interface{}
	var respuesta_banco models.Banco
	if response, error := helpers.GetJsonTest(beego.AppConfig.String("UrlCoreApi")+"/banco/"+strconv.Itoa(banco_id), &respuesta_peticion); error == nil && response == 200 {
		json_banco, err := json.Marshal(respuesta_peticion)
		if err == nil {
			if err := json.Unmarshal(json_banco, &respuesta_banco); err != nil {
				outputError = map[string]interface{}{"funcion": "/GetBanco", "err": err, "status": "502"}
				return respuesta_banco, outputError
			}
		} else {
			outputError = map[string]interface{}{"funcion": "/GetBanco", "err": err, "status": "502"}
			return respuesta_banco, outputError
		}
	} else {
		outputError = map[string]interface{}{"funcion": "/GetBanco", "err": error, "status": "502"}
		return respuesta_banco, outputError
	}
	return respuesta_banco, nil
}

func CrearInformeSatisfaccion(numero_contrato_suscrito int, vigencia_contrato string, tipo_pago string, periodo_inicio string, periodo_fin string, tipo_factura string, numero_cuenta_factura string, valor_pagar int, tipo_cuenta string, numero_cuenta string, banco_id int) (informe_satisfaccion models.InformeSatisfaccion, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			//fmt.Println("error", err)
			outputError = map[string]interface{}{"funcion": "/CrearInformeSatisfaccion", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	nombre_banco := ""

	informacion_informe_satisfaccion, err := ObtenerInformacionInformeSatisfaccion(strconv.Itoa(numero_contrato_suscrito), vigencia_contrato)
	if err == nil {
		banco, err := ObtenerBanco(banco_id)
		if err == nil {
			nombre_banco = banco.NombreBanco
		} else {
			outputError = map[string]interface{}{"funcion": "/CrearInformeSatisfaccion", "err": err, "status": "502"}
			return informe_satisfaccion, outputError
		}

		balance_contrato, err := ObtenerBalanceFinancieroContrato(strconv.Itoa(numero_contrato_suscrito), vigencia_contrato)
		if err != nil {
			outputError = map[string]interface{}{"funcion": "/CrearInformeSatisfaccion", "err": err, "status": "502"}
			return informe_satisfaccion, outputError
		}

		total_contrato, _ := strconv.Atoi(balance_contrato.TotalContrato)
		saldo_contrato, _ := strconv.Atoi(balance_contrato.Saldo)

		informe_satisfaccion, _ = helpers.CrearPdfInformeSatisfaccion(informacion_informe_satisfaccion.Dependencia,
			informacion_informe_satisfaccion.NombreProveedor,
			informacion_informe_satisfaccion.DocumentoProveedor,
			informacion_informe_satisfaccion.CumplimientoTotal,
			informacion_informe_satisfaccion.TipoContrato,
			informacion_informe_satisfaccion.FechaInicio,
			strconv.Itoa(numero_contrato_suscrito),
			informacion_informe_satisfaccion.Cdp,
			informacion_informe_satisfaccion.VigenciaCdp,
			informacion_informe_satisfaccion.Rp,
			informacion_informe_satisfaccion.VigenciaRp,
			informacion_informe_satisfaccion.CargoSupervisor,
			tipo_factura,
			numero_cuenta_factura,
			total_contrato,
			periodo_inicio,
			periodo_fin,
			saldo_contrato,
			informacion_informe_satisfaccion.FechaFin,
			tipo_cuenta,
			numero_cuenta,
			nombre_banco,
			informacion_informe_satisfaccion.Supervisor,
			vigencia_contrato)
	}

	return informe_satisfaccion, nil
}
