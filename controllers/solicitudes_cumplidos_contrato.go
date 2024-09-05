package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
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
	defer errorhandler.HandlePanic(&c.Controller)

	numero_contrato := c.Ctx.Input.Param(":numero_contrato")
	vigencia := c.Ctx.Input.Param(":vigencia")

	data, err := services.ObtenerSolicitudesCumplidosContrato(numero_contrato, vigencia)
	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, data)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}
