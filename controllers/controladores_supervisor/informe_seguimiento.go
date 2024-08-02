package controladores_supervisor

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers/helpers_supervisor"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

// Informe_seguimientoController operations for Informe_seguimiento
type InformeSeguimientoController struct {
	beego.Controller
}

// URLMapping ...
func (c *InformeSeguimientoController) URLMapping() {
	c.Mapping("GetBalanceFinancieroContrato", c.GetBalanceFinancieroContrato)
	c.Mapping("GenerateInformeSeguimiento", c.GenerateInformeSeguimiento)
}

// GetBalanceFinancieroContrato ...
// @Title GetBalanceFinancieroContrato
// @Description Obtener el balance financiero de un contrato
// @Param	numero_contrato_suscrito		path 	string	true		"Numero del contrato suscrito"
// @Param	vigencia_contrato			path 	string	true		"Vigencia del contrato"
// @Success 200 {object} models.BalanceContrato
// @Failure 502 Error procesando la solicitud
// @router /balance-financiero-contrato/:numero_contrato_suscrito/:vigencia_contrato [get]
func (c *InformeSeguimientoController) GetBalanceFinancieroContrato() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "GetBalanceFinancieroContrato" + "/" + (localError["funcion"]).(string))
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

	if data, err := helpers_supervisor.GetBalanceFinancieroContrato(numero_contrato_suscrito, vigencia_contrato); err == nil {
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

// GenerateInformeSeguimiento ...
// @Title GenerateInformeSeguimiento
// @Description Genera un informe de seguimiento para un contrato suscrito
// @Param	body	body 	models.BodyInformeSeguimiento	true	"Parámetros necesarios para generar el informe de seguimiento"
// @Success 200 {object} models.InformeSeguimiento "Successful - Informe de seguimiento generado exitosamente"
// @Failure 404 "No se encontró el recurso solicitado"
// @Failure 502 "Error al intentar generar el informe de seguimiento"
// @router /informe-seguimiento [post]
func (c *InformeSeguimientoController) GenerateInformeSeguimiento() {
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

	var v models.BodyInformeSeguimiento

	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	if data, err := helpers_supervisor.CreateInformeSeguimiento(v.NumeroContratoSuscrito, v.VigenciaContrato, v.TipoPago, v.PeiodoInicio, v.PeriodoFin, v.TipoFactura, v.NumeroCuentaFactura, v.ValorPagar, v.TipoCuenta, v.NumeroCuenta, v.Banco); err == nil {
		if (data == models.InformeSeguimiento{}) {
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
