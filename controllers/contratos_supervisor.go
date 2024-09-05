package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
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

	defer errorhandler.HandlePanic(&c.Controller)

	documento_supervisor := c.Ctx.Input.Param(":documento_supervisor")

	data, err := services.ObtenerContratosSupervisor(documento_supervisor)
	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, data)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}
