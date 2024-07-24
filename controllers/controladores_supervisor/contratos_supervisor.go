package controladores_supervisor

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers/helpers_supervisor"
)

// ContratosSupervisorController operations for ContratosSupervisor
type ContratosSupervisorController struct {
	beego.Controller
}

// URLMapping ...
func (c *ContratosSupervisorController) URLMapping() {
	c.Mapping("GetContratos", c.GetContratos)
}

// GetContratos ...
// @Title GetContratos
// @Description get DependenciasSupervisor by documento_supervisor
// @Param	documento_supervisor		path 	string	true		"documento_supervisor"
// @Success 200 {object} models.ContratosSupervisor
// @Failure 403 :id is empty
// @router /:documento_supervisor [get]
func (c *ContratosSupervisorController) GetContratos() {

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

	if data, err := helpers_supervisor.GetDependenciasSupervisor(documento_supervisor); err == nil {
		if data != nil {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
		} else {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "No se encontro ninguna dependencia para el supervisor", "Data": nil}
		}
		c.ServeJSON()
	} else {
		panic(err)
	}

}
