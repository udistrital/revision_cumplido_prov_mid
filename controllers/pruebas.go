package controllers

import (
	"log"

	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
)

// PruebasController operations for Pruebas
type PruebasController struct {
	beego.Controller
}

// URLMapping ...
func (c *PruebasController) URLMapping() {
	c.Mapping("ObtenerSupervisorContrato", c.ObtenerSupervisorContrato)
	c.Mapping("ObtenerOrdenadorContrato", c.ObtenerOrdenadorContrato)
}

// @Title ObtenerSupervisorContrato
// @Description Obtiene la información del supervisor de un contrato
// @Param numero_contrato_suscrito query string true "Número del contrato suscrito"
// @Param vigencia query string true "Vigencia del contrato"
// @Success 200 {object} models.SupervisorContratoProveedor
// @Failure 502 {object} map[string]interface{}
// @router /obtener-supervisor-contrato/:numero_contrato_suscrito/:vigencia [get]
func (c *PruebasController) ObtenerSupervisorContrato() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "PruebasController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("502")
			}
		}
	}()

	// Estructura para recibir los parámetros de la solicitud
	numeroContratoSuscrito := c.Ctx.Input.Param(":numero_contrato_suscrito")
	vigencia := c.Ctx.Input.Param(":vigencia")

	// Llamada a la función para obtener el supervisor del contrato
	supervisorContrato, outputError := services.ObtenerSupervisorContrato(numeroContratoSuscrito, vigencia)

	if outputError != nil {
		c.Data["json"] = outputError
		c.Ctx.Output.SetStatus(200)
	} else {
		c.Data["json"] = map[string]interface{}{
			"Success": true,
			"Status":  "200",
			"Message": "Successful",
			"Data":    supervisorContrato,
		}
		c.Ctx.Output.SetStatus(200)
	}

	c.ServeJSON()
}

// @Title ObtenerOrdenadorContrato
// @Description Obtiene la información del ordenador de un contrato
// @Param numero_contrato_suscrito query string true "Número del contrato suscrito"
// @Param vigencia query string true "Vigencia del contrato"
// @Success 200 {object} models.OrdenadorContratoProveedor
// @Failure 502 {object} map[string]interface{}
// @router /obtener-ordenador-contrato/:numero_contrato_suscrito/:vigencia [get]
func (c *PruebasController) ObtenerOrdenadorContrato() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "PruebasController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("502")
			}
		}
	}()

	// Estructura para recibir los parámetros de la solicitud
	numeroContratoSuscrito := c.Ctx.Input.Param(":numero_contrato_suscrito")
	vigencia := c.Ctx.Input.Param(":vigencia")

	// Llamada a la función para obtener el ordenador del contrato
	ordenadorContrato, outputError := services.ObtenerOrdenadorContrato(numeroContratoSuscrito, vigencia)

	if outputError != nil {
		c.Data["json"] = outputError
		c.Ctx.Output.SetStatus(200)
	} else {
		c.Data["json"] = map[string]interface{}{
			"Success": true,
			"Status":  "200",
			"Message": "Successful",
			"Data":    ordenadorContrato,
		}
		c.Ctx.Output.SetStatus(200)
	}

	c.ServeJSON()
}
