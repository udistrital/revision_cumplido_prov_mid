package helpers_soporte

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func EliminarSoporteCumplido(documento_id string) (response string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetDocumentosPagoMensual", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	var soportes_pagos_mensuales []models.SoportePago
	var respuesta_peticion map[string]interface{}

	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/soporte_cumplido/?limit=1&query=DocumentoId:"+documento_id+",Activo:true", &respuesta_peticion); (err == nil) && (response == 200) {
		LimpiezaRespuestaRefactor(respuesta_peticion, &soportes_pagos_mensuales)
		if (soportes_pagos_mensuales[0] == models.SoportePago{}) {
			return "No se encontró el soporte de pago o este ya se elimino con anterioridad", nil
		}
	} else {
		outputError = map[string]interface{}{"funcion": "/EliminarSoporteCumplido/soporte_cumplido", "err": err, "status": "502"}
		return "No se encontró el soporte de pago", outputError
	}
	var res map[string]interface{}
	delete_true := "Soporte pago eliminado correctamente"
	delect_false := "No se encontró el soporte de pago"

	if err := sendJson(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/soporte_cumplido/"+strconv.Itoa(soportes_pagos_mensuales[0].Id), "DELETE", &res, nil); err == nil {
		response = delete_true
		return response, nil
	} else {
		outputError = map[string]interface{}{"funcion": "/EliminarSoporteCumplido/soporte_cumplido", "err": err, "status": "502"}
		response = delect_false
		return response, outputError
	}

}
