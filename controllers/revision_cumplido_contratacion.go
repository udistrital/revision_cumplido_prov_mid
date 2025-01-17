package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
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

	defer errorhandler.HandlePanic(&c.Controller)

	data, err := services.ObtenerCumplidosPendientesContratacion()
	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, data)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}
