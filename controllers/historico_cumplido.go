package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
	"github.com/udistrital/utils_oas/requestresponse"
)

type HistoricoCumplidoController struct {
	beego.Controller
}

func (c *HistoricoCumplidoController) URLMapping() {
	c.Mapping("ObtenerHistorico", c.ObtenerHistorico)

}

// ObtenerHistoricoDeCumplido ...
// @Title ObtenerHistoricoDeCumplido
// @Description Generar el historico de los estado por lo que ha pasado un cumplido proveedor
// @Param	cumplido_proveedor_id	path	string	true	"ID del cumplido del proveedor"
// @Success 200 {object} models.HistoricoCumplido
// @Failure 404
// @router /historico_cumplido/:cumplido_proveedor_id [get]
func (c *HistoricoCumplidoController) ObtenerHistorico() {

	response, err := services.ObtberHistoricoEstado(c.Ctx.Input.Param(":cumplido_proveedor_id"))

	if err == nil {
		if len(response) > 0 {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = requestresponse.APIResponseDTO(true, 200, response)
		} else {

			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = requestresponse.APIResponseDTO(true, 200, []models.HistoricoCumplido{}, "No se encontraron registros")
		}
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}
