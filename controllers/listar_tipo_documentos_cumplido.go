package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// ListarTipoDocumentosCumplidoController operations for ListarTipoDocumentosCumplidoController
type ListarTipoDocumentosCumplidoController struct {
	beego.Controller
}

// URLMapping ...
func (c *ListarTipoDocumentosCumplidoController) URLMapping() {
	c.Mapping("ObtenerTiposDocumentosCumplido", c.ObtenerTiposDocumentosCumplido)
}

// ObtenerTiposDocumentosCumplido ...
// @Title ObtenerTiposDocumentosCumplido
// @Description get tipos de documentos cumplido
// @Success 200 {object} []models.DocumentoCumplido
// @Failure 404 {object} map[string]interface{} "Error interno del servidor"
// @router /tipos-documentos-cumplido [get]
func (c *ListarTipoDocumentosCumplidoController) ObtenerTiposDocumentosCumplido() {
	defer errorhandler.HandlePanic(&c.Controller)

	data, err := services.ObtenerTiposDocumentosCumplido()
	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, data)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}
