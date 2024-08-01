package controladores_ordenador

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers/helper_generar_documento"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers/helpers_ordenador"
	"io/ioutil"
	"net/http"
)

type RevisionCumplidoOrdenadorController struct {
	beego.Controller
}

// URLMapping asigna métodos a la estructura PingController
func (c *RevisionCumplidoOrdenadorController) URLMapping() {
	c.Mapping("ObtenerPendientesRevisionOrdenador", c.ObtenerPendientesRevisionOrdenador)
	c.Mapping("ListaCumplidosReversibles", c.ListaCumplidosReversibles)
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

	dependencias, err := helpers_ordenador.ObtenerSolicitudesCumplidos(documento_ordenador, "PRO,Activo:true")

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
// @Failure 403 :id_cumplido is empty
//@router /revertir-solicitud-pago/:id_cumplido [get]
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

	id_cumplido := c.GetString(":id_cumplido")
	print(id_cumplido)

	dependencias, err := helpers_ordenador.ListaCumplidosReversibles(id_cumplido)

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
//@Param body  body models.AutorizacionPago true "body para la autorizacion de pago"
//Success 200 {object}
// @Failure 403 :id_solicitud_pago is empty
//@router /certificado-aprobacion-pago/:id_solicitud_pago [get]
func (c *RevisionCumplidoOrdenadorController) ObtenerCertificado() {

	helper_generar_documento.GenerarPdf()
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

	id_solicitud_pago := c.GetString(":id_solicitud_pago")
	autorizacion, err := helpers_ordenador.GenerarAutorizacion(id_solicitud_pago)

	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = err
	} else if autorizacion == nil {
		println("2")
		c.Data["json"] = map[string]interface{}{"Succes": true, "Status:": 204, "Message": "No hay datos", "Data": autorizacion}
	} else {
		println("3")
		c.Data["json"] = map[string]interface{}{"Succes": true, "Status:": 200, "Message": "Consulta completa", "Data": autorizacion}
	}
	c.ServeJSON()
}

//generarDocumentoAutorizacion
//@Title GnerarAutorizaxionPago
//@Description Metodo
//Success 200 {object}
// @Failure 403
//@router /testpdf [get]
func (c *RevisionCumplidoOrdenadorController) GenerarPdf() {

	helper_generar_documento.GenerarPdf()

	// Leer el archivo PDF generado
	pdfFile, err := ioutil.ReadFile("tabla.pdf")
	if err != nil {
		beego.Error("Error al leer el archivo PDF:", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.ResponseWriter.Write([]byte("Error al generar el archivo PDF"))
		return
	}

	// Establecer encabezados para la respuesta HTTP
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/pdf")
	c.Ctx.ResponseWriter.Header().Set("Content-Disposition", "attachment; filename=tabla.pdf")

	// Escribir el contenido del PDF en la respuesta
	_, err = c.Ctx.ResponseWriter.Write(pdfFile)
	if err != nil {
		beego.Error("Error al escribir el archivo PDF en la respuesta:", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.ResponseWriter.Write([]byte("Error al enviar el archivo PDF"))
		return
	}

	// Opcional: Enviar una respuesta JSON si deseas informar sobre el éxito de la operación
	c.Data["json"] = map[string]string{"status": "success"}
	c.ServeJSON()
}
