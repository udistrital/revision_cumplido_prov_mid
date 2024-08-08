package helpers_soporte

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func GetDocumentosPagoMensual(cumplido_proveedor_id string) (documentos []models.DocumentosSoporteCorto, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetDocumentosPagoMensual", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	var soportes_pagos_mensuales []models.SoportePago
	var documentos_crud []models.Documento
	var fileGestor models.FileGestorDocumental
	var soporte models.DocumentosSoporteCorto
	var documento_individual models.DocumentoCorto

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
				if len(documentos_crud) > 0 {
					for _, documento_crud := range documentos_crud {
						var observaciones map[string]interface{}
						documento_individual.Id = documento_crud.Id
						documento_individual.Nombre = documento_crud.Nombre
						documento_individual.TipoDocumento = documento_crud.TipoDocumento.Nombre
						documento_individual.Descripcion = documento_crud.Descripcion
						documento_individual.FechaCreacion = documento_crud.FechaCreacion
						if err := json.Unmarshal([]byte(documento_crud.Metadatos), &observaciones); err == nil {
							documento_individual.Observaciones = observaciones["observaciones"].(string)
						}
						soporte.Documento = documento_individual
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
					outputError = map[string]interface{}{"funcion": "/GetDocumentosPagoMensual/documento", "err": err, "status": "502", "Message": "No se encontraron documentos asociados al cumplido proveedor"}
					return nil, outputError
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
