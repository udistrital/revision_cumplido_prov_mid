package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// SoportesCumplidoController operations for SoportesCumplido
type SoportesCumplidoController struct {
	beego.Controller
}

// URLMapping ...
func (c *SoportesCumplidoController) URLMapping() {
	c.Mapping("SubirSoporteCumplido", c.SubirSoporteCumplido)
	c.Mapping("ObtenerDocumentosPagoMensual", c.ObtenerDocumentosPagoMensual)
	c.Mapping("ObtenerComprimidoSoportes", c.ObtenerComprimidoSoportes)

}

// Post ...
// @Title SubirSoporteCumplido
// @Description Subir un soporte de pago
// @Param	solicitud_pago_id	body 	string	true		"ID de la solicitud de pago"
// @Param	tipo_documento		body 	string	true		"Tipo de documento (debe ser 'application/pdf')"
// @Param	item_id				body 	string	true		"ID del tipo de documento"
// @Param	observaciones		body 	string	false		"Observaciones del documento"
// @Param	nombre_archivo		body 	string	true		"Nombre del archivo"
// @Param	archivo				body 	string	true		"Archivo en base64"
// @Success 200 {object} models.SoportePago
// @Failure 403 body is empty
// @router /soportes [post]
func (c *SoportesCumplidoController) SubirSoporteCumplido() {
	defer errorhandler.HandlePanic(&c.Controller)

	// Estructura para recibir el cuerpo de la solicitud
	var soporteReq models.BodySubirSoporteRequest

	// Parsear el cuerpo de la solicitud
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &soporteReq); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{
			"Success": false,
			"Status":  "400",
			"Message": "Bad Request: " + err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Llamada al helper para subir el soporte
	data, err := services.SubirSoporteCumplido(soporteReq.SolicitudPagoID, soporteReq.TipoDocumento, soporteReq.ItemID, soporteReq.Observaciones, soporteReq.NombreArchivo, soporteReq.Archivo)
	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, data)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// ObtenerDocumentosPagoMensual ...
// @Title ObtenerDocumentosPagoMensual
// @Description Obtener los documentos de soporte de pago mensual
// @Param	cumplido_proveedor_id		path 	string	true		"ID del cumplido proveedor"
// @Success 200 {object} []models.DocumentosSoporte
// @Failure 404 No se encontraron documentos de soporte
// @router /soportes/:cumplido_proveedor_id [get]
func (c *SoportesCumplidoController) ObtenerDocumentosPagoMensual() {
	defer errorhandler.HandlePanic(&c.Controller)

	cumplido_proveedor_id := c.Ctx.Input.Param(":cumplido_proveedor_id")

	data, err := services.ObtenerSoportesCumplido(cumplido_proveedor_id)
	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, data)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()

}

// ObtenerComprimidoSoportes ...
// @Title ObtenerComprimidoSoportes
// @Description Obtener los documentos de soporte de pago mensual y comprimirlos en un archivo ZIP
// @Param	id_cumplido_proveedor	path 	string	true		"ID del cumplido proveedor"
// @Success 200 {object} models.DocumentosComprimido "Documentos comprimidos en formato base64"
// @Failure 404 "No se encontraron documentos de soporte"
// @Failure 404 "Error al intentar obtener o comprimir los documentos"
// @router /soportes-comprimido/:id_cumplido_proveedor [get]
func (c *SoportesCumplidoController) ObtenerComprimidoSoportes() {
	defer errorhandler.HandlePanic(&c.Controller)

	id_cumplido_proveedor := c.Ctx.Input.Param(":id_cumplido_proveedor")

	data, err := services.ObtenerComprimidoSoportes(id_cumplido_proveedor)
	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, data)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}
