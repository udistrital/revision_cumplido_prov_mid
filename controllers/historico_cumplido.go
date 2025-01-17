package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

type HistoricoCumplidosController struct {
	beego.Controller
}

func (c *HistoricoCumplidosController) URLMapping() {
	c.Mapping("ObtenerHistorico", c.ObtenerHistorico)
	c.Mapping("ObtenerHistoricoCumplidos", c.ObtenerHistoricoCumplidos)

}

// ObtenerHistoricoDeCumplido ...
// @Title ObtenerHistoricoDeCumplido
// @Description Generar el historico de los estado por lo que ha pasado un cumplido proveedor
// @Param	cumplido_proveedor_id	path	string	true	"ID del cumplido del proveedor"
// @Success 200 {object} models.HistoricoCumplido
// @Failure 404
// @router /historico_cumplido/:cumplido_proveedor_id [get]
func (c *HistoricoCumplidosController) ObtenerHistorico() {

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

// @Title ObtenerHistoricoCumplidos
// @Description Obtiene el historico de los cumplidos dependiendo de varios filtros
// @Param Anios query int true "años de los cumplidos a consultar"
// @Param Meses query int true "meses de los cumplidos a consultar"
// @Param Vigencias query string true "vigencias de los cumplidos a consultar"
// @Param Proveedores query string true "proveedores de los cumplidos a consultar"
// @Param Estados query string true "estados de los cumplidos a consultar"
// @Param Dependencias query string true "dependencias de los cumplidos a consultar"
// @Param Contratos query string true "contratos de los cumplidos a consultar"
// @Success 200 {object} []models.CumplidosFiltrados
// @Failure 404 {object} map[string]interface{}
// @router /filtro-cumplidos [post]
func (c *HistoricoCumplidosController) ObtenerHistoricoCumplidos() {
	defer errorhandler.HandlePanic(&c.Controller)

	var v models.BodyHistoricoRequest

	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	data, err := services.ObtenerHistoricoCumplidosFiltro(v.Anios, v.Meses, v.Vigencias, v.Proveedores, v.Estados, v.Dependencias, v.Contratos, v.TiposContratos)
	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, data)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}
