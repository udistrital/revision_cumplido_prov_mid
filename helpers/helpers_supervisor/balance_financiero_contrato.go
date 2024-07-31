package helpers_supervisor

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func GetBalanceFinancieroContrato(numero_contrato_suscrito string, vigencia_contrato string) (balance_contrato models.BalanceContrato, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetBalanceContrato", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	informacion_contrato, err := GetInformacionContratoContratista(numero_contrato_suscrito, vigencia_contrato)
	if err == nil {
		informacion_contratista, err := ContratosContratista(informacion_contrato.InformacionContratista.Documento.Numero)
		if err == nil {
			contrato, err := GetContrato(numero_contrato_suscrito, vigencia_contrato)
			if err == nil {
				valor_girado, err := GetValorGiradoPorCdp(informacion_contratista[0].NumeroCdp, informacion_contratista[0].VigenciaCdp, contrato.Contrato.UnidadEjecutora)
				if err == nil {
					total_contrato, err := strconv.ParseFloat(informacion_contrato.InformacionContratista.ValorContrato, 64)
					if err == nil {
						saldo_contrato := int(total_contrato) - valor_girado
						balance_contrato.TotalContrato = informacion_contrato.InformacionContratista.ValorContrato
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

func GetValorGiradoPorCdp(cdp string, vigencia_cdp string, unidad_ejecucion string) (valor_girado int, err error) {
	var temp_giros_tercero map[string]interface{}
	var giros_tercero models.GirosTercero
	valor_girado = 0
	fmt.Println(beego.AppConfig.String("UrlFinancieraJBPM") + "/" + "giros_tercero/" + cdp + "/" + vigencia_cdp + "/" + unidad_ejecucion)
	if response, err := getJsonWSO2Test(beego.AppConfig.String("UrlFinancieraJBPM")+"/"+"giros_tercero/"+cdp+"/"+vigencia_cdp+"/"+unidad_ejecucion, &temp_giros_tercero); (err == nil) && (response == 200) {
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
