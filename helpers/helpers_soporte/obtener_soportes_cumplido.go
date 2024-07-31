package helpers_soporte

import (
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func GetDocumentosPagoMensual(cumplido_proveedor_id string) (documentos []models.DocumentosSoporte, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetDocumentosPagoMensual", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	var soportes_pagos_mensuales []models.SoportePago
	var documentos_crud []models.Documento
	var fileGestor models.FileGestorDocumental
	var soporte models.DocumentosSoporte

	var respuesta_peticion map[string]interface{}

	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/soporte_pago/?limit=-1&query=CumplidoProveedorId.Id:"+cumplido_proveedor_id, &respuesta_peticion); (err == nil) && (response == 200) {
		LimpiezaRespuestaRefactor(respuesta_peticion, &soportes_pagos_mensuales)
		if len(soportes_pagos_mensuales) != 0 {
			var ids_documentos []string
			for _, soporte_pago_mensual := range soportes_pagos_mensuales {
				ids_documentos = append(ids_documentos, strconv.Itoa(soporte_pago_mensual.DocumentoId))
			}

			var ids_documentos_juntos = strings.Join(ids_documentos, "|")
			if response, err := getJsonTest(beego.AppConfig.String("UrlDocumentosCrud")+"/documento/?limit=-1&query=Activo:True,Id.in:"+ids_documentos_juntos, &documentos_crud); (err == nil) && (response == 200) {
				for _, documento_crud := range documentos_crud {
					soporte.Documento = documento_crud
					if response, err := getJsonTest(beego.AppConfig.String("UrlGestorDocumental")+"/document/"+documento_crud.Enlace, &fileGestor); (err == nil) && (response == 200) {
						soporte.Archivo = fileGestor
						documentos = append(documentos, soporte)
					} else {
						logs.Error(err)
						continue
					}
				}
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/GetDocumentosPagoMensual/documento", "err": err, "status": "502"}
				return nil, outputError
			}
		} else {
			return nil, outputError
		}

	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/GetDocumentosPagoMensual/soporte_pago_mensual", "err": err, "status": "502"}
		return nil, outputError
	}

	return
}
