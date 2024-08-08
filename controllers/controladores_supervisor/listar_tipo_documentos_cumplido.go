package controladores_supervisor

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers/helpers_supervisor"
)

// Listar_tipo_documentos_cumplidoController operations for Listar_tipo_documentos_cumplido
type ListarTipoDocumentosCumplidoController struct {
	beego.Controller
}

// URLMapping ...
func (c *ListarTipoDocumentosCumplidoController) URLMapping() {
	c.Mapping("GetTiposDocumentosCumplido", c.GetTiposDocumentosCumplido)
}

// GetTiposDocumentosCumplido ...
// @Title GetTiposDocumentosCumplido
// @Description get tipos de documentos cumplido
// @Success 200 {object} []models.DocumentoCumplido
// @Failure 502 {object} map[string]interface{} "Error interno del servidor"
// @router /tipos-documentos-cumplido [get]
func (c *ListarTipoDocumentosCumplidoController) GetTiposDocumentosCumplido() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "GetTiposDocumentosCumplido" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	if data, err := helpers_supervisor.GetTiposDocumentosCumplido(); err == nil {
		if len(data) > 0 {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
		} else {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "No se encontro ningun tipo de documento", "Data": nil}
		}
		c.ServeJSON()
	} else {
		panic(err)
	}
}
