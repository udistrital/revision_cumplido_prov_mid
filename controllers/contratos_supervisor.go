package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
)

// ContratosSupervisorController operations for ContratosSupervisor
type ContratosSupervisorController struct {
	beego.Controller
}

// URLMapping ...
func (c *ContratosSupervisorController) URLMapping() {
	c.Mapping("ObtenerContratosSupervisor", c.ObtenerContratosSupervisor)

}

// ObtenerContratosSupervisor ...
// @Title GetContratosSupervisor
// @Description get GetContratosSupervisor by documento_supervisor
// @Param	documento_supervisor		path 	string	true		"documento_supervisor"
// @Success 200 {object} models.ContratoSupervisor
// @Failure 403 :documento_supervisor is empty
// @router /contratos-supervisor/:documento_supervisor [get]
func (c *ContratosSupervisorController) ObtenerContratosSupervisor() {

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

	documento_supervisor := c.Ctx.Input.Param(":documento_supervisor")

	if data, err := services.ObtenerContratosSupervisor(documento_supervisor); err == nil {
		if len(data.Contratos) > 0 {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
		} else {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "404", "Message": err, "Data": []map[string]interface{}{}}
		}
		c.ServeJSON()
	} else {
		panic(err)
	}

}
