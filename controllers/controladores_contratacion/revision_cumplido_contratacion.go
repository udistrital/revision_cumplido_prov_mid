package controladores_contratacion

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers/helpers_contratacion"
)

type RevisionCumplidoContratacionController struct {
	beego.Controller
}

func (c *RevisionCumplidoContratacionController) URLMapping() {
	c.Mapping("ObtenerPendientesContratacion", c.ObtenerPendientesContratacion)
}

//ObtenerPendientesRevisionOrdenador
//@Title ObtenerPendientesRevisionOrdenador
//@Description Metodo para que el personal de c  Obtenga los contratos para la aprobacion de pago
//Success 200 {object} models.Contrato
// @Failure 403 :document is empty
//@router /solicitudes-pago/ [get]
func (c *RevisionCumplidoContratacionController) ObtenerPendientesContratacion() {

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

	dependencias, err := helpers_contratacion.ObtenerCumplidosPendientesContratacion("PRC,Activo:true")

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
