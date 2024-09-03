package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
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

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = beego.AppConfig.String("appname") + "/" + "RevisionCumplidoOrdenadorController" + "/" + (localError["funcion"]).(string)
			c.Data["data"] = localError["err"]
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500")
			}
		}
	}()

	documento_ordenador := c.GetString(":documento_ordenador")

	solicitudes, err := services.ObtenerSolicitudesCumplidos(documento_ordenador)

	if err == nil {
		c.Data["json"] = map[string]interface{}{"Succes": true, "Status:": 200, "Message": "Consulta completa", "Data": solicitudes}
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 404, "Message": err, "Data": []map[string]interface{}{}}
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

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = beego.AppConfig.String("appname") + "/" + "RevertirSolicitud" + "/" + (localError["funcion"]).(string)
			c.Data["data"] = localError["err"]
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500")
			}
		}
	}()

	id_cumplido := c.GetString(":documento_ordenador")

	cumplidos_reversibles, err := services.ListaCumplidosReversibles(id_cumplido)

	if err == nil {
		c.Data["json"] = map[string]interface{}{"Succes": true, "Status:": 200, "Message": "Consulta completa", "Data": cumplidos_reversibles}
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 404, "Message": err, "Data": []map[string]interface{}{}}
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

	cumplido_proveedor_id := c.GetString(":cumplido_proveedor_id")
	autorizacion, err := services.GenerarAutorizacionGiro(cumplido_proveedor_id)
	if err == nil {
		c.Data["json"] = map[string]interface{}{"Succes": true, "Status:": 200, "Message": "Consulta completa", "Data": autorizacion}
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 404, "Message": err, "Data": []map[string]interface{}{}}
	}
	c.ServeJSON()
}
