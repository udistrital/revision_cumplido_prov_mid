package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

type RevisionCumplidoOrdenadorController struct {
	beego.Controller
}

func (c *RevisionCumplidoOrdenadorController) URLMapping() {
	c.Mapping("ObtenerCumplidosPendientesRevisionOrdenador", c.ObtenerCumplidosPendientesRevisionOrdenador)
	c.Mapping("ListaCumplidosReversibles", c.ListaCumplidosReversibles)
	c.Mapping("GenerarAutorizacionGiro", c.GenerarAutorizacionGiro)
}

// ObtenerCumplidosPendientesRevisionOrdenador
// @Title ObtenerCumplidosPendientesRevisionOrdenador
// @Description Metodo para que el ordenador  Obtenga los contratos para la aprobacion de pago
// @Param documento_ordenador path string true  "Documento del ordenador"
// Success 200 {object} models.Contrato
// @Failure 403 :document is empty
// @router /solicitudes-pago/:documento_ordenador [get]
func (c *RevisionCumplidoOrdenadorController) ObtenerCumplidosPendientesRevisionOrdenador() {

	defer errorhandler.HandlePanic(&c.Controller)

	documento_ordenador := c.GetString(":documento_ordenador")

	data, err := services.ObtenerSolicitudesCumplidos(documento_ordenador)
	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, data)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// GetContratos
// @Title RevertirSolicitud
// @Description Metodo para que el ordenador  revierta el estado de un cumplido
// @Param id de pago path string true  "id_solictud_de_pago"
// Success 200 {object}
// @Failure 403 :id_cumplido is empty
// @router /revertir-solicitud-pago/:documento_ordenador [get]
func (c *RevisionCumplidoOrdenadorController) ListaCumplidosReversibles() {

	defer errorhandler.HandlePanic(&c.Controller)

	id_cumplido := c.GetString(":documento_ordenador")

	data, err := services.ListaCumplidosReversibles(id_cumplido)
	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, data)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// GenerarAutorizacionGiro ...
// @Title GenerarAutorizacionGiro
// @Description Generar la autorización de giro para un cumplido proveedor específico
// @Param	cumplido_proveedor_id	path	string	true	"ID del cumplido del proveedor"
// @Success 200 {object} models.DocumentoAutorizacionPago
// @Failure 404
// @router /autorizacion-giro/:cumplido_proveedor_id [get]
func (c *RevisionCumplidoOrdenadorController) GenerarAutorizacionGiro() {

	defer errorhandler.HandlePanic(&c.Controller)

	cumplido_proveedor_id := c.GetString(":cumplido_proveedor_id")
	data, err := services.GenerarAutorizacionGiro(cumplido_proveedor_id)
	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, data)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}
