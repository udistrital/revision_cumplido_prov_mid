package controladores_ordenador

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers/helpers_ordenador"
)

type RevisionCumplidoOrdenadorController struct {
	beego.Controller
}

// URLMapping asigna métodos a la estructura PingController
func (c *RevisionCumplidoOrdenadorController) URLMapping() {
	c.Mapping("ObtenerPendientesRevisionOrdenador", c.ObtenerPendientesRevisionOrdenador)
	c.Mapping("RevertirSolicitud", c.RevertirSolicitud)
}

//ObtenerPendientesRevisionOrdenador
//@Title ObtenerPendientesRevisionOrdenador
//@Description Metodo para que el ordenador  Obtenga los contratos para la aprobacion de pago
//@Param documento_ordenador path string true  "Documento del ordenador"
//Success 200 {object} models.Contrato
// @Failure 403 :document is empty
//@router /solicitudes-pago/:documento_ordenador [get]
func (c *RevisionCumplidoOrdenadorController) ObtenerPendientesRevisionOrdenador() {

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

	dependencias, err := helpers_ordenador.ObternerContratos(documento_ordenador, "PRO")

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

//GetContratos
//@Title RevertirSolicitud
//@Description Metodo para que el ordenador  revierta el estado de un cumplido
//@Param id de pago path string true  "id_solictud_de_pago"
//Success 200 {object}
// @Failure 403 :id_solicitud_pago is empty
//@router /revertir-solicitud-pago/:id_solicitud_pago [post]
func (c *RevisionCumplidoOrdenadorController) RevertirSolicitud() {

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

	id_pago := c.GetString(":id_solicitud_pago")
	print(id_pago)

	dependencias, err := helpers_ordenador.RevertirAprobadoSupervisor(id_pago)

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

//
//@Title ObtenerCertificado firmado de aprobacion pago.
//@Description Metodo encargado de retornar el certificado firmado de aprobacion pago.
//@Param id de pago path string true  "id_solictud_de_pago"
//Success 200 {object}
// @Failure 403 :id_solicitud_pago is empty
//@router /certificado-aprobacion-pago/:id_solicitud_pago [post]
func (c *RevisionCumplidoOrdenadorController) ObtenerCertificado() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = beego.AppConfig.String("appname") + "/" + "RevisionCumplidoOrdenadorController" + "/" + (localError["funcion"]).(string)
			c.Data["data"] = localError["err"]
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	id_solicitud_pago := c.Ctx.Input.Param(":id_solicitud_pago")

	cumplido, err := helpers_ordenador.ObtenerInfoContratoPorId(id_solicitud_pago)

	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = err
	} else if cumplido == nil {
		println("2")
		c.Data["json"] = map[string]interface{}{"Succes": true, "Status:": 204, "Message": "No hay datos", "Data": cumplido}
	} else {
		println("3")
		c.Data["json"] = map[string]interface{}{"Succes": true, "Status:": 200, "Message": "Consulta completa", "Data": cumplido}
	}
	c.ServeJSON()
}
