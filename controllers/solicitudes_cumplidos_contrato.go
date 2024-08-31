package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
)

// Solicitud_contratoController operations for Solicitud_contrato
type SolicitudesCumplidosContratoController struct {
	beego.Controller
}

// URLMapping ...
func (c *SolicitudesCumplidosContratoController) URLMapping() {
	c.Mapping("ObtenerSolicitudesContrato", c.ObtenerSolicitudesContrato)
}

// ObtenerSolicitudesContrato ...
// @Title ObtenerSolicitudesContrato
// @Description get solicitudes de contrato
// @Param	numero_contrato		path 	string	true		"numero_contrato"
// @Param	vigencia			path 	string	true		"vigencia del contrato"
// @Success 200 {object} []models.CambioEstadoCumplido
// @Failure 403 :numero_contrato or vigencia is empty
// @router /solicitudes-contrato/:numero_contrato/:vigencia [get]
func (c *SolicitudesCumplidosContratoController) ObtenerSolicitudesContrato() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "ContratosSupervisorController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	numero_contrato := c.Ctx.Input.Param(":numero_contrato")
	vigencia := c.Ctx.Input.Param(":vigencia")

	data, err := services.ObtenerSolicitudesCumplidosContrato(numero_contrato, vigencia)
	if err == nil {
		c.Data["json"] = map[string]interface{}{"Succes": true, "Status:": 200, "Message": "Consulta completa", "Data": data}
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 404, "Message": err, "Data": []map[string]interface{}{}}
	}
	c.ServeJSON()
}
