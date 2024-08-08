package controladores_soporte

import (
	"encoding/json"
	"log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers/helpers_soporte"
)

// SoportesCumplidoController operations for SoportesCumplido
type SoportesCumplidoController struct {
	beego.Controller
}

// URLMapping ...
func (c *SoportesCumplidoController) URLMapping() {
	c.Mapping("SubirSoporte", c.SubirSoporte)
	c.Mapping("GetDocumentosPagoMensual", c.GetDocumentosPagoMensual)
	c.Mapping("EliminarSoporteCumplido", c.EliminarSoporteCumplido)
	c.Mapping("AgregarComentarioSoporte", c.AgregarComentarioSoporte)
	c.Mapping("GetSoportesComprimido", c.GetSoportesComprimido)

}

// Post ...
// @Title SubirSoporte
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
func (c *SoportesCumplidoController) SubirSoporte() {
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
	var soporteReq struct {
		SolicitudPagoID int    `json:"SolicitudPagoID"`
		TipoDocumento   string `json:"TipoDocumento"`
		ItemID          int    `json:"ItemID"`
		Observaciones   string `json:"Observaciones"`
		NombreArchivo   string `json:"NombreArchivo"`
		Archivo         string `json:"Archivo"`
	}

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
	soporte, err := helpers_soporte.SubirSoporte(soporteReq.SolicitudPagoID, soporteReq.TipoDocumento, soporteReq.ItemID, soporteReq.Observaciones, soporteReq.NombreArchivo, soporteReq.Archivo)
	if err != nil {
		panic(err)
	}

	// Respuesta exitosa
	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": soporte}
	c.ServeJSON()
}

// GetDocumentosPagoMensual ...
// @Title GetDocumentosPagoMensual
// @Description Obtener los documentos de soporte de pago mensual
// @Param	cumplido_proveedor_id		path 	string	true		"ID del cumplido proveedor"
// @Success 200 {object} []models.DocumentosSoporte
// @Failure 404 No se encontraron documentos de soporte
// @router /soportes/:cumplido_proveedor_id [get]
func (c *SoportesCumplidoController) GetDocumentosPagoMensual() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "GetDocumentosPagoMensual" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	cumplido_proveedor_id := c.Ctx.Input.Param(":cumplido_proveedor_id")

	if data, err := helpers_soporte.GetDocumentosPagoMensual(cumplido_proveedor_id); err == nil {
		if data == nil {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "No se encontraron documentos de soporte", "Data": nil}
		} else {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
		}
		c.ServeJSON()
	} else {
		panic(err)
	}
}

// EliminarSoporteCumplido ...
// @Title EliminarSoporteCumplido
// @Description Eliminar un soporte de pago cumplido
// @Param	soporte_pago_id		path 	string	true		"ID del soporte de pago a eliminar"
// @Success 200 {object} map[string]interface{} "Soporte de pago eliminado exitosamente"
// @Failure 404 "No se encontró el soporte de pago"
// @Failure 502 "Error al intentar eliminar el soporte de pago"
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

	if data, err := helpers_soporte.EliminarSoporteCumplido(soporte_pago_id); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Data": data}
		c.ServeJSON()
	} else {
		c.Data["json"] = err
		c.ServeJSON()
	}
}

// AgregarComentarioSoporte ...
// @Title AgregarComentarioSoporte
// @Description Agregar un comentario a un soporte de pago
// @Param	soporte_id			path 	string	true		"ID del soporte de pago"
// @Param	cambio_estado_id	path 	string	true		"ID del cambio de estado"
// @Param	comentario			body 	string	true		"Comentario a agregar"
// @Success 200 {object} models.RespuestaComentarioSoporte "Comentario agregado exitosamente"
// @Failure 404 "No se encontró el soporte de pago o cambio de estado"
// @Failure 502 "Error al intentar agregar el comentario"
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
				c.Abort("502")
			}
		}
	}()

	type BodyParams struct {
		SoporteId      string `json:"soporte_id"`
		CambioEstadoId string `json:"cambio_estado_id"`
		Comentario     string `json:"comentario"`
	}

	var v BodyParams

	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	if v.SoporteId == "" || v.CambioEstadoId == "" || v.Comentario == "" {
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Parámetros incompletos"}
		c.ServeJSON()
		return
	}

	if res, err := helpers_soporte.AgregarComentarioSoporte(v.SoporteId, v.CambioEstadoId, v.Comentario); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Comentario agregado exitosamente", "Data": res}
		c.ServeJSON()
	} else {
		c.Data["json"] = err
		c.ServeJSON()
	}
}

// GetSoportesComprimido ...
// @Title GetSoportesComprimido
// @Description Obtener los documentos de soporte de pago mensual y comprimirlos en un archivo ZIP
// @Param	id_cumplido_proveedor	path 	string	true		"ID del cumplido proveedor"
// @Success 200 {object} models.DocumentosComprimido "Documentos comprimidos en formato base64"
// @Failure 404 "No se encontraron documentos de soporte"
// @Failure 502 "Error al intentar obtener o comprimir los documentos"
// @router /soportes-comprimido/:id_cumplido_proveedor [get]
func (c *SoportesCumplidoController) GetSoportesComprimido() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "GetSoportesComprimido" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	id_cumplido_proveedor := c.Ctx.Input.Param(":id_cumplido_proveedor")

	data, err := helpers_soporte.SoportesComprimido(id_cumplido_proveedor)
	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Documentos comprimidos exitosamente", "Data": data}
		c.ServeJSON()
	} else {
		c.Data["json"] = err
		c.ServeJSON()
	}
}
