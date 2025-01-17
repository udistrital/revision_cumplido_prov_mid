package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// EstadoSoporteController operations for EstadoSoporte
type CambioEstadoCumplidoController struct {
	beego.Controller
}

// URLMapping ...
func (c *CambioEstadoCumplidoController) URLMapping() {
	c.Mapping("CambioEstadoCumplido", c.CambioEstadoCumplido)

}

// @Title CambioEstadoCumplido
// @Description Cambia el estado de pago del proveedor
// @Param estado_cumplido_id query string true "ID del estado cumplido"
// @Param cumplido_proveedor_id query string true "ID del cumplido proveedor"
// @Param documento_responsable query string true "Número del documento responsable"
// @Param cargo_responsable query string true "Cargo del responsable"
// @Success 200 {object} models.CambioEstadoCumplidoResponse
// @Failure 404 {object} map[string]interface{}
// @router /cambio-estado [post]
func (c *CambioEstadoCumplidoController) CambioEstadoCumplido() {
	defer errorhandler.HandlePanic(&c.Controller)

	// Estructura para recibir el cuerpo de la solicitud

	var v models.BodyCumplidoRequest

	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	// Llamada al helper para cambiar el estado de pago
	response, err := services.CambioEstadoCumplido(v.CodigoAbreviacionEstadoCumplido, v.CumplidoProveedorID)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, response)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}
