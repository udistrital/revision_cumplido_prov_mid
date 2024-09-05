package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// CumplidoSatisfaccionController operations for CumplidoSatisfaccion
type CumplidoSatisfaccionController struct {
	beego.Controller
}

// URLMapping ...
func (c *CumplidoSatisfaccionController) URLMapping() {
	c.Mapping("ObtenerBalanceFinancieroContrato", c.ObtenerBalanceFinancieroContrato)
	c.Mapping("GenerarCumplidoSatisfaccion", c.GenerarCumplidoSatisfaccion)
}

// ObtenerBalanceFinancieroContrato ...
// @Title ObtenerBalanceFinancieroContrato
// @Description Obtener el balance financiero de un contrato
// @Param	numero_contrato_suscrito		path 	string	true		"Numero del contrato suscrito"
// @Param	vigencia_contrato			path 	string	true		"Vigencia del contrato"
// @Success 200 {object} models.BalanceContrato
// @Failure 404 Error procesando la solicitud
// @router /balance-financiero-contrato/:numero_contrato_suscrito/:vigencia_contrato [get]
func (c *CumplidoSatisfaccionController) ObtenerBalanceFinancieroContrato() {

	defer errorhandler.HandlePanic(&c.Controller)

	numero_contrato_suscrito := c.Ctx.Input.Param(":numero_contrato_suscrito")
	vigencia_contrato := c.Ctx.Input.Param(":vigencia_contrato")

	data, err := services.ObtenerBalanceFinancieroContrato(numero_contrato_suscrito, vigencia_contrato)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, data)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()

}

// GenerarInformeSatisfaccion ...
// @Title GenerarInformeSatisfaccion
// @Description Genera un informe de seguimiento para un contrato suscrito
// @Param	body	body 	models.BodyInformeSeguimiento	true	"Parámetros necesarios para generar el informe de seguimiento"
// @Success 200 {object} models.InformeSeguimiento "Successful - Informe de seguimiento generado exitosamente"
// @Failure 404 "No se encontró el recurso solicitado"
// @Failure 404 "Error al intentar generar el informe de seguimiento"
// @router /cumplido-satisfaccion [post]
func (c *CumplidoSatisfaccionController) GenerarCumplidoSatisfaccion() {
	defer errorhandler.HandlePanic(&c.Controller)

	var v models.BodyCumplidoSatisfaccion

	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	data, err := services.CrearCumplidoSatisfaccion(v.NumeroContratoSuscrito, v.VigenciaContrato, v.TipoPago, v.PeriodoInicio, v.PeriodoFin, v.TipoFactura, v.NumeroCuentaFactura, v.ValorPagar, v.TipoCuenta, v.NumeroCuenta, v.Banco)
	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, data)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()

}
