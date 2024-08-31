package controllers

import (
	"encoding/json"
	"log"

	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
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
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "CambioEstadoCumplidoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	// Estructura para recibir el cuerpo de la solicitud
	type BodyParams struct {
		CodigoAbreviacionEstadoCumplido string `json:"CodigoAbreviacionEstadoCumplido"`
		CumplidoProveedorID             int    `json:"CumplidoProveedorId"`
	}

	var v BodyParams

	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	// Llamada al helper para cambiar el estado de pago
	response, outputError := services.CambioEstadoCumplido(v.CodigoAbreviacionEstadoCumplido, v.CumplidoProveedorID)

	if outputError != nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "404", "Message": outputError, "Data": []map[string]interface{}{}}
	} else {
		c.Data["json"] = map[string]interface{}{
			"Success": true,
			"Status":  "200",
			"Message": "Successful",
			"Data":    response,
		}
		c.Ctx.Output.SetStatus(200)
	}

	c.ServeJSON()
}
