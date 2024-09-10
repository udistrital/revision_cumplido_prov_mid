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

// URLMapping ...
func (c *HistoricoCumplidosController) URLMapping() {
	c.Mapping("ObtenerHistoricoCumplidos", c.ObtenerHistoricoCumplidos)
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
// @router / [post]
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
