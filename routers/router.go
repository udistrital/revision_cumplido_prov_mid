// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/controllers"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/supervisor",
			beego.NSInclude(
				&controllers.ContratosSupervisorController{}, &controllers.SolicitudesCumplidosContratoController{}, &controllers.InformeSeguimientoController{}, &controllers.ListarTipoDocumentosCumplidoController{}, &controllers.InformeSeguimientoController{},
			),
		),
		beego.NSNamespace("/solicitud-pago",
			beego.NSInclude(
				&controllers.SoportesCumplidoController{},
			),
		),
		beego.NSNamespace("/solicitud-pago",
			beego.NSInclude(
				&controllers.CambioEstadoCumplidoController{},
			),
		),
		beego.NSNamespace("/ordenador",
			beego.NSInclude(
				&controllers.RevisionCumplidoOrdenadorController{},
			),
		),
		beego.NSNamespace("/contratacion",
			beego.NSInclude(
				&controllers.RevisionCumplidoContratacionController{},
			),
		))

	beego.AddNamespace(ns)
}
