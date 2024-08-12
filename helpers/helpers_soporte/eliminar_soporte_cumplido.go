package helpers_soporte

import (
	"github.com/astaxie/beego"
)

func EliminarSoporteCumplido(id_soporte_pago string) (response string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetDocumentosPagoMensual", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	var res map[string]interface{}
	delete_true := "Soporte pago eliminado correctamente"
	delect_false := "No se encontró el soporte de pago"

	if err := sendJson(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/soporte_cumplido/"+id_soporte_pago, "DELETE", &res, nil); err == nil {
		response = delete_true
		return response, nil
	} else {
		outputError = map[string]interface{}{"funcion": "/EliminarSoporteCumplido/soporte_cumplido", "err": err, "status": "502"}
		response = delect_false
		return response, outputError
	}

}
