package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"

	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtenerBalanceFinancieroContrato(numero_contrato_suscrito string, vigencia_contrato string) (balance_contrato models.BalanceContrato, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ObtenerBalanceFinancieroContrato", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	informacion_contrato, err := helpers.ObtenerInformacionContratoProveedor(numero_contrato_suscrito, vigencia_contrato)
	if err == nil {
		contrato_general, err := helpers.ObtenerContratoGeneralProveedor(numero_contrato_suscrito, vigencia_contrato)
		if err == nil {
			valor_girado, err := ObtenerValorGiradoPorCdp(informacion_contrato[0].NumeroCdp, informacion_contrato[0].VigenciaCdp, strconv.Itoa(contrato_general.UnidadEjecutora))
			if err == nil {
				total_contrato := contrato_general.ValorContrato
				if err == nil {
					saldo_contrato := int(total_contrato) - valor_girado
					balance_contrato.TotalContrato = strconv.FormatFloat(total_contrato, 'f', 0, 64)
					balance_contrato.Saldo = strconv.Itoa(saldo_contrato)
					return balance_contrato, nil
				} else {
					outputError = map[string]interface{}{"funcion": "/ObtenerBalanceFinancieroContrato", "err": err, "status": "502"}
					return balance_contrato, outputError
				}

			} else {
				outputError = map[string]interface{}{"funcion": "/ObtenerBalanceFinancieroContrato", "err": err, "status": "502"}
				return balance_contrato, outputError
			}
		}
	}

	return balance_contrato, nil
}

func ObtenerValorGiradoPorCdp(cdp string, vigencia_cdp string, unidad_ejecucion string) (valor_girado int, err error) {
	var temp_giros_tercero map[string]interface{}
	var giros_tercero models.GirosTercero
	valor_girado = 0
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlFinancieraJBPM")+"/giros_tercero/"+cdp+"/"+vigencia_cdp+"/"+unidad_ejecucion, &temp_giros_tercero); (err == nil) && (response == 200) {
		json_giros_tercero, error_json := json.Marshal(temp_giros_tercero)
		if error_json == nil {
			if err := json.Unmarshal(json_giros_tercero, &giros_tercero); err == nil {
				for _, giro := range giros_tercero.Giros.Tercero {
					total_girado, err := strconv.Atoi(giro.ValorBrutoGirado)
					if err == nil {
						valor_girado = valor_girado + total_girado
					}
				}
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

func ObtenerInformacionCumplidoSatisfaccion(numero_contrato_suscrito string, vigencia string) (informacion_informe models.InformacionInformeSatisfaccion, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ObtenerInformacionInformeSatisfaccion", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var informacion_proveedor []models.InformacionProveedor
	informacion_contrato, err := helpers.ObtenerInformacionContratoProveedor(numero_contrato_suscrito, vigencia)
	fmt.Println("informacion_contrato", informacion_contrato)
	if err == nil {
		contrato_general, err := helpers.ObtenerContratoGeneralProveedor(numero_contrato_suscrito, vigencia)
		if err == nil {
			if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/informacion_proveedor/?query=Id:"+strconv.Itoa(contrato_general.Contratista), &informacion_proveedor); (err == nil) && (response == 200) {
				if len(informacion_proveedor) > 0 {
					vigencia_contrato, _ := strconv.Atoi(vigencia)
					acta_inicio, err := helpers.ObtenerActaInicio(numero_contrato_suscrito, vigencia_contrato)
					if err == nil {
						balance_general, err := ObtenerBalanceFinancieroContrato(numero_contrato_suscrito, vigencia)
						if err == nil {
							total_contrato, _ := strconv.Atoi(balance_general.TotalContrato)
							saldo, _ := strconv.Atoi(balance_general.Saldo)
							informacion_informe.Dependencia = informacion_contrato[0].NombreDependencia
							informacion_informe.NombreProveedor = informacion_contrato[0].NombreProveedor
							informacion_informe.DocumentoProveedor = informacion_proveedor[0].NumDocumento
							informacion_informe.TipoContrato = informacion_contrato[0].TipoContrato
							informacion_informe.FechaInicio = acta_inicio.FechaInicio
							informacion_informe.NumeroContratoSuscrito = numero_contrato_suscrito
							informacion_informe.Cdp = informacion_contrato[0].NumeroCdp
							informacion_informe.VigenciaCdp = informacion_contrato[0].VigenciaCdp
							informacion_informe.Rp = informacion_contrato[0].NumeroRp
							informacion_informe.VigenciaRp = informacion_contrato[0].VigenciaRp
							informacion_informe.CargoSupervisor = contrato_general.Supervisor.Cargo
							informacion_informe.ValorTotalContrato = total_contrato
							informacion_informe.SaldoContrato = saldo
							informacion_informe.FechaFin = acta_inicio.FechaFin
							informacion_informe.Supervisor = contrato_general.Supervisor.Nombre
							informacion_informe.DocumentoSupervisor = strconv.Itoa(contrato_general.Supervisor.Documento)
						} else {
							outputError = map[string]interface{}{"funcion": "/ObtenerInformacionInformeSatisfaccion", "message": "No fue posible obtener el balance financiero del contrato", "err": err, "status": "502"}
							return informacion_informe, outputError
						}

					} else {
						outputError = map[string]interface{}{"funcion": "/ObtenerInformacionInformeSatisfaccion", "message": "No fue posible obtener el acta de inicio", "err": err, "status": "502"}
						return informacion_informe, outputError
					}
				} else {
					outputError = map[string]interface{}{"funcion": "/ObtenerInformacionInformeSatisfaccion", "message": "No fue posible obtener la informacion del proveedor", "status": "404"}
					return informacion_informe, outputError
				}
			} else {
				outputError = map[string]interface{}{"funcion": "/ObtenerInformacionInformeSatisfaccion", "message": "No fue posible obtener la informacion del proveedor", "err": err, "status": "502"}
				return informacion_informe, outputError
			}
		} else {
			outputError = map[string]interface{}{"funcion": "/ObtenerInformacionInformeSatisfaccion", "message": "No fue posible obtener el contrato general", "err": err, "status": "502"}
			return informacion_informe, outputError
		}

	} else {
		outputError = map[string]interface{}{"funcion": "/ObtenerInformacionInformeSatisfaccion", "message": "No fue posible obtener la informacion del contrato", "err": err, "status": "502"}
		return informacion_informe, outputError
	}

	return informacion_informe, nil
}

func ObtenerBanco(banco_id int) (banco models.Banco, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
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

func CrearCumplidoSatisfaccion(numero_contrato_suscrito int, vigencia_contrato string, cumplimiento_contrato string, tipo_pago string, periodo_inicio string, periodo_fin string, tipo_factura string, numero_cuenta_factura string, valor_pagar int, tipo_cuenta string, numero_cuenta string, banco_id int) (informe_satisfaccion models.CumplidoSatisfaccion, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/CrearInformeSatisfaccion", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	nombre_banco := ""

	informacion_informe_satisfaccion, err := ObtenerInformacionCumplidoSatisfaccion(strconv.Itoa(numero_contrato_suscrito), vigencia_contrato)
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

		informe_satisfaccion, _ = helpers.CrearPdfCumplidoSatisfaccion(informacion_informe_satisfaccion.Dependencia,
			informacion_informe_satisfaccion.NombreProveedor,
			informacion_informe_satisfaccion.DocumentoProveedor,
			cumplimiento_contrato,
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
			vigencia_contrato,
			informacion_informe_satisfaccion.DocumentoSupervisor)
	}

	return informe_satisfaccion, nil
}
