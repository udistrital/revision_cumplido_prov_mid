package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
)

type RevisionCumplidoContratacionController struct {
	beego.Controller
}

func (c *RevisionCumplidoContratacionController) URLMapping() {
	c.Mapping("ObtenerCumplidosPendientesContratacion", c.ObtenerCumplidosPendientesContratacion)
}

// ObtenerPendientesRevisionOrdenador
// @Title ObtenerPendientesRevisionOrdenador
// @Description Metodo para que el personal de c  Obtenga los contratos para la aprobacion de pago
// Success 200 {object} models.Contrato
// @Failure 403 :document is empty
// @router /solicitudes-pago [get]
func (c *RevisionCumplidoContratacionController) ObtenerCumplidosPendientesContratacion() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = beego.AppConfig.String("appname") + "/" + "RevisionCumplidoOrdenadorController" + "/" + (localError["funcion"]).(string)
			c.Data["data"] = localError["err"]
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500")
			}
		}
	}()

	dependencias, err := services.ObtenerCumplidosPendientesContratacion()

	if err == nil {
		c.Data["json"] = map[string]interface{}{"Succes": true, "Status:": 204, "Message": "Cumplidos pendientes obtenidos satisfactoriamente", "Data": dependencias}
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 404, "Message": err, "Data": []map[string]interface{}{}}
	}
	c.ServeJSON()
}
