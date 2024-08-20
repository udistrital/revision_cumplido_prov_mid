package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
)

type RevisionCumplidoContratacionController struct {
	beego.Controller
}

func (c *RevisionCumplidoContratacionController) URLMapping() {
	c.Mapping("ObtenerCumplidosPendientesContratacion", c.ObtenerCumplidosPendientesContratacion)
}

// ObtenerPendientesRevisionOrdenador
// @Title ObtenerPendientesRevisionOrdenador
// @Description Metodo para que el personal de c  Obtenga los contratos para la aprobacion de pago
// Success 200 {object} models.Contrato
// @Failure 403 :document is empty
// @router /solicitudes-pago/ [get]
func (c *RevisionCumplidoContratacionController) ObtenerCumplidosPendientesContratacion() {

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

	dependencias, err := services.ObtenerCumplidosPendientesContratacion("PRC,Activo:true")

	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = err
	} else if dependencias == nil {
		println("2")
		c.Data["json"] = map[string]interface{}{"Succes": true, "Status:": 204, "Message": "No hay datos", "Data": dependencias}
	} else {
		println("3")
		c.Data["json"] = map[string]interface{}{"Succes": true, "Status:": 200, "Message": "Consulta completa", "Data": dependencias}
	}
	c.ServeJSON()
}

// generarDocumentoAutorizacion
// @Title GnerarAutorizaxionPago
// @Description Metodo
// Success 200 {object}
// @Failure 403
// @router /certificado-aprobacion-pago/:id_solicitud_pago [get]
func (c *RevisionCumplidoContratacionController) GenerarPdf() {

	id_solicitud_pago := c.GetString(":id_solicitud_pago")
	autorizacion, err := services.GenerarAutorizacionPago(id_solicitud_pago)

	if err != nil {
		beego.Error("Error al leer el archivo PDF:", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.ResponseWriter.Write([]byte("Error al generar el archivo PDF"))
		return
	}

	file := helpers.GenerarPdfAutorizacionPago(autorizacion)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = err
	} else if file == "" {
		c.Data["json"] = map[string]interface{}{"Succes": true, "Status:": 204, "Message": "No hay datos", "Data": ""}
	} else {
		c.Data["json"] = map[string]interface{}{"Succes": true, "Status:": 200, "Message": "Consulta completa", "Data": file}
	}
	c.ServeJSON()
}
