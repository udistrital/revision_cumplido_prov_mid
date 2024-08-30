package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
)

// CumplidoSatisfaccionController operations for CumplidoSatisfaccion
type CumplidoSatisfaccionController struct {
	beego.Controller
}

// URLMapping ...
func (c *CumplidoSatisfaccionController) URLMapping() {
	c.Mapping("ObtenerBalanceFinancieroContrato", c.ObtenerBalanceFinancieroContrato)
	c.Mapping("GenerarCumplidoSatisfaccion", c.GenerarCumplidoSatisfaccion)
}

// ObtenerBalanceFinancieroContrato ...
// @Title ObtenerBalanceFinancieroContrato
// @Description Obtener el balance financiero de un contrato
// @Param	numero_contrato_suscrito		path 	string	true		"Numero del contrato suscrito"
// @Param	vigencia_contrato			path 	string	true		"Vigencia del contrato"
// @Success 200 {object} models.BalanceContrato
// @Failure 502 Error procesando la solicitud
// @router /balance-financiero-contrato/:numero_contrato_suscrito/:vigencia_contrato [get]
func (c *CumplidoSatisfaccionController) ObtenerBalanceFinancieroContrato() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "ObtenerBalanceFinancieroContrato" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()
	numero_contrato_suscrito := c.Ctx.Input.Param(":numero_contrato_suscrito")
	vigencia_contrato := c.Ctx.Input.Param(":vigencia_contrato")

	if data, err := services.ObtenerBalanceFinancieroContrato(numero_contrato_suscrito, vigencia_contrato); err == nil {
		if (data == models.BalanceContrato{}) {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "No se encontro balance financiero", "Data": nil}
		} else {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
		}
		c.ServeJSON()
	} else {
		panic(err)
	}
}

// GenerarInformeSatisfaccion ...
// @Title GenerarInformeSatisfaccion
// @Description Genera un informe de seguimiento para un contrato suscrito
// @Param	body	body 	models.BodyInformeSeguimiento	true	"Parámetros necesarios para generar el informe de seguimiento"
// @Success 200 {object} models.InformeSeguimiento "Successful - Informe de seguimiento generado exitosamente"
// @Failure 404 "No se encontró el recurso solicitado"
// @Failure 502 "Error al intentar generar el informe de seguimiento"
// @router /cumplido-satisfaccion [post]
func (c *CumplidoSatisfaccionController) GenerarCumplidoSatisfaccion() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "ContratosSupervisorController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	var v models.BodyCumplidoSatisfaccion

	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	if data, err := services.CrearCumplidoSatisfaccion(v.NumeroContratoSuscrito, v.VigenciaContrato, v.CumplimientoContrato, v.TipoPagoId, v.PeiodoInicio, v.PeriodoFin, v.TipoDocumentoCobroId, v.NumeroCuentaFactura, v.ValorPagar, v.TipoCuenta, v.NumeroCuenta, v.BancoId); err == nil {
		if (data == models.CumplidoSatisfaccion{}) {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "No se pudo generar el informe de seguimiento", "Data": nil}
		} else {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
		}
		c.ServeJSON()
	} else {
		panic(err)
	}

}
