package controladores_supervisor

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers/helpers_supervisor"
)

// Solicitud_contratoController operations for Solicitud_contrato
type SolicitudContratoController struct {
	beego.Controller
}

// URLMapping ...
func (c *SolicitudContratoController) URLMapping() {
	c.Mapping("GetSolicitudesContrato", c.GetSolicitudesContrato)
}

// GetSolicitudesContrato ...
// @Title GetSolicitudesContrato
// @Description get solicitudes de contrato
// @Param	numero_contrato		path 	string	true		"numero_contrato"
// @Param	vigencia			path 	string	true		"vigencia del contrato"
// @Success 200 {object} []models.CambioEstadoCumplido
// @Failure 403 :numero_contrato or vigencia is empty
// @router /solicitudes-contrato/:numero_contrato/:vigencia [get]
func (c *SolicitudContratoController) GetSolicitudesContrato() {
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

	if data, err := helpers_supervisor.GetSolicitudesCumplidosProveedor(numero_contrato, vigencia); err == nil {
		if len(data) > 0 {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
		} else {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "No se encontro ninguna solicitud de contrato", "Data": nil}
		}
		c.ServeJSON()
	} else {
		panic(err)
	}
}
