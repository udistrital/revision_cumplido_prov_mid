package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
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
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "ObtenerTiposDocumentosCumplido" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	if data, err := services.ObtenerTiposDocumentosCumplido(); err == nil {
		if len(data) > 0 {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
		} else {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "404", "Message": err, "Data": []map[string]interface{}{}}
		}
		c.ServeJSON()
	} else {
		panic(err)
	}
}
