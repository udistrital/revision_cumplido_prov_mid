package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"

	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtenerBalanceFinancieroContrato(numero_contrato_suscrito string, vigencia_contrato string) (balance_contrato models.BalanceContrato, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
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
				saldo_contrato := int(total_contrato) - valor_girado
				balance_contrato.TotalContrato = strconv.FormatFloat(total_contrato, 'f', 0, 64)
				balance_contrato.Saldo = strconv.Itoa(saldo_contrato)
				return balance_contrato, nil

			} else {
				outputError = fmt.Errorf("Error al obtener el valor girado por el CDP")
				return balance_contrato, outputError
			}
		} else {
			outputError = fmt.Errorf("Error al obtener el contrato general")
			return balance_contrato, outputError
		}
	} else {
		outputError = fmt.Errorf("Error al obtener la información del contrato proveedor")
		return balance_contrato, outputError
	}
}

func ObtenerValorGiradoPorCdp(cdp string, vigencia_cdp string, unidad_ejecucion string) (valor_girado int, err error) {
	var temp_giros_tercero map[string]interface{}
	var giros_tercero models.GirosTercero
	valor_girado = 0
	fmt.Println(beego.AppConfig.String("UrlFinancieraJBPM") + "/giros_tercero/" + cdp + "/" + vigencia_cdp + "/" + unidad_ejecucion)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlFinancieraJBPM")+"/giros_tercero/"+cdp+"/"+vigencia_cdp+"/"+unidad_ejecucion, &temp_giros_tercero); (err == nil) && (response == 200) {
		if temp_giros_tercero == nil {
			err = errors.New("error en la consulta de giros_tercero")
			return valor_girado, err
		}
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
				err = errors.New("error Unmarshal giros_tercero")
				return valor_girado, err
			}

		} else {
			err = errors.New("error Marshal giros_tercero")
			return valor_girado, err
		}

	} else {
		return valor_girado, err
	}
}

func ObtenerInformacionCumplidoSatisfaccion(numero_contrato_suscrito string, vigencia string) (informacion_informe models.InformacionInformeSatisfaccion, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var informacion_proveedor []models.InformacionProveedor
	informacion_contrato, err := helpers.ObtenerInformacionContratoProveedor(numero_contrato_suscrito, vigencia)
	//fmt.Println("informacion_contrato", informacion_contrato)
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
							informacion_informe.VigenciaCdp = informacion_contrato[0].CDPFechaExpedicion
							informacion_informe.Rp = informacion_contrato[0].NumeroRp
							informacion_informe.VigenciaRp = informacion_contrato[0].RPFechaRegistro
							informacion_informe.CargoSupervisor = contrato_general.Supervisor.Cargo
							informacion_informe.ValorTotalContrato = total_contrato
							informacion_informe.SaldoContrato = saldo
							informacion_informe.FechaFin = acta_inicio.FechaFin
							informacion_informe.Supervisor = contrato_general.Supervisor.Nombre
							informacion_informe.DocumentoSupervisor = strconv.Itoa(contrato_general.Supervisor.Documento)
						} else {
							outputError = fmt.Errorf("No fue posible obtener el balance financiero del contrato")
							return informacion_informe, outputError
						}

					} else {
						outputError = fmt.Errorf("No fue posible obtener el acta de inicio")
						return informacion_informe, outputError
					}
				} else {
					outputError = fmt.Errorf("No fue posible obtener la informacion del proveedor")
					return informacion_informe, outputError
				}
			} else {
				outputError = fmt.Errorf("No fue posible obtener la informacion del proveedor")
				return informacion_informe, outputError
			}
		} else {
			outputError = fmt.Errorf("No fue posible obtener el contrato general")
			return informacion_informe, outputError
		}

	} else {
		outputError = fmt.Errorf("No fue posible obtener la informacion del contrato")
		return informacion_informe, outputError
	}

	return informacion_informe, nil
}

func CrearCumplidoSatisfaccion(numero_contrato_suscrito int, vigencia_contrato string, tipo_pago string, periodo_inicio time.Time, periodo_fin time.Time, tipo_factura string, numero_cuenta_factura string, valor_pagar int, tipo_cuenta string, numero_cuenta string, banco string) (cumplido_satisfaccion models.CumplidoSatisfaccion, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	informacion_informe_satisfaccion, err := ObtenerInformacionCumplidoSatisfaccion(strconv.Itoa(numero_contrato_suscrito), vigencia_contrato)
	if err == nil {
		archivo_cumplido_satisfaccion := helpers.CrearPdfCumplidoSatisfaccion(informacion_informe_satisfaccion.Dependencia,
			informacion_informe_satisfaccion.NombreProveedor,
			informacion_informe_satisfaccion.DocumentoProveedor,
			tipo_pago,
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
			valor_pagar,
			informacion_informe_satisfaccion.ValorTotalContrato,
			periodo_inicio,
			periodo_fin,
			informacion_informe_satisfaccion.SaldoContrato,
			informacion_informe_satisfaccion.FechaFin,
			tipo_cuenta,
			numero_cuenta,
			banco,
			informacion_informe_satisfaccion.Supervisor,
			vigencia_contrato,
			informacion_informe_satisfaccion.DocumentoSupervisor)

		nombre := "CumplidoSatisfaccion_" + strings.Join(strings.Fields(informacion_informe_satisfaccion.NombreProveedor), "") + "_" + strconv.Itoa(numero_contrato_suscrito) + "_" + vigencia_contrato

		cumplido_satisfaccion.NombreArchivo = nombre
		cumplido_satisfaccion.Archivo = archivo_cumplido_satisfaccion
		cumplido_satisfaccion.NombreResponsable = informacion_informe_satisfaccion.Supervisor
		cumplido_satisfaccion.CargoResponsable = informacion_informe_satisfaccion.CargoSupervisor
		cumplido_satisfaccion.DescripcionDocumento = "Cumplido satisfaccion del contrato suscrito " + strconv.Itoa(numero_contrato_suscrito) + " de " + vigencia_contrato + " con actividades compredidas entre " + helpers.FormatearFecha(periodo_inicio) + " al " + helpers.FormatearFecha(periodo_fin)

		return cumplido_satisfaccion, nil
	} else {
		outputError = fmt.Errorf("No fue posible obtener la informacion del contrato")
		return cumplido_satisfaccion, outputError
	}

}
