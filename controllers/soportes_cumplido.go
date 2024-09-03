package controllers

import (
	"encoding/json"
	"log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/services"
)

// SoportesCumplidoController operations for SoportesCumplido
type SoportesCumplidoController struct {
	beego.Controller
}

// URLMapping ...
func (c *SoportesCumplidoController) URLMapping() {
	c.Mapping("SubirSoporteCumplido", c.SubirSoporteCumplido)
	c.Mapping("ObtenerDocumentosPagoMensual", c.ObtenerDocumentosPagoMensual)
	c.Mapping("EliminarSoporteCumplido", c.EliminarSoporteCumplido)
	c.Mapping("AgregarComentarioSoporte", c.AgregarComentarioSoporte)
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
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "SoportesCumplidoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

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
	soporte, err := services.SubirSoporteCumplido(soporteReq.SolicitudPagoID, soporteReq.TipoDocumento, soporteReq.ItemID, soporteReq.Observaciones, soporteReq.NombreArchivo, soporteReq.Archivo)
	if err == nil {
		c.Data["json"] = map[string]interface{}{"Succes": true, "Status:": 200, "Message": "Consulta completa", "Data": soporte}
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 404, "Message": err, "Data": []map[string]interface{}{}}
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
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "ObtenerDocumentosPagoMensual" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	cumplido_proveedor_id := c.Ctx.Input.Param(":cumplido_proveedor_id")

	data, err := services.ObtenerDocumentosPagoMensual(cumplido_proveedor_id)
	if err == nil {
		c.Data["json"] = map[string]interface{}{"Succes": true, "Status:": 200, "Message": "Consulta completa", "Data": data}
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 404, "Message": err, "Data": []map[string]interface{}{}}
	}
	c.ServeJSON()

}

// EliminarSoporteCumplido ...
// @Title EliminarSoporteCumplido
// @Description Eliminar un soporte de pago cumplido
// @Param	soporte_pago_id		path 	string	true		"ID del soporte de pago a eliminar"
// @Success 200 {object} map[string]interface{} "Soporte de pago eliminado exitosamente"
// @Failure 404 "No se encontró el soporte de pago"
// @Failure 404 "Error al intentar eliminar el soporte de pago"
// @router /soportes/:soporte_pago_id [delete]
func (c *SoportesCumplidoController) EliminarSoporteCumplido() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "EliminarSoporteCumplido" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	soporte_pago_id := c.Ctx.Input.Param(":soporte_pago_id")

	data, err := services.EliminarSoporteCumplido(soporte_pago_id)
	if err == nil {
		c.Data["json"] = map[string]interface{}{"Succes": true, "Status:": 200, "Message": "Consulta completa", "Data": data}
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 404, "Message": err, "Data": []map[string]interface{}{}}
	}
	c.ServeJSON()
}

// AgregarComentarioSoporte ...
// @Title AgregarComentarioSoporte
// @Description Agregar un comentario a un soporte de pago
// @Param	soporte_id			path 	string	true		"ID del soporte de pago"
// @Param	cambio_estado_id	path 	string	true		"ID del cambio de estado"
// @Param	comentario			body 	string	true		"Comentario a agregar"
// @Success 200 {object} models.RespuestaComentarioSoporte "Comentario agregado exitosamente"
// @Failure 404 "No se encontró el soporte de pago o cambio de estado"
// @Failure 404 "Error al intentar agregar el comentario"
// @router /comentario-soporte [post]
func (c *SoportesCumplidoController) AgregarComentarioSoporte() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "SoportesCumplidoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	var v models.AgregarComentarioSoporteRequest

	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	if v.SoporteId == "" || v.CambioEstadoId == "" || v.Comentario == "" {
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Parámetros incompletos"}
		c.ServeJSON()
		return
	}

	data, err := services.AgregarComentarioSoporte(v.SoporteId, v.CambioEstadoId, v.Comentario)
	if err == nil {
		c.Data["json"] = map[string]interface{}{"Succes": true, "Status:": 200, "Message": "Consulta completa", "Data": data}
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 404, "Message": err, "Data": []map[string]interface{}{}}
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
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "ObtenerComprimidoSoportes" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	id_cumplido_proveedor := c.Ctx.Input.Param(":id_cumplido_proveedor")

	data, err := services.ObtenerComprimidoSoportes(id_cumplido_proveedor)
	if err == nil {
		c.Data["json"] = map[string]interface{}{"Succes": true, "Status:": 200, "Message": "Consulta completa", "Data": data}
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 404, "Message": err, "Data": []map[string]interface{}{}}
	}
}
