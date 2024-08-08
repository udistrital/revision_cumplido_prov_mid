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
	"github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_estado_pago"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_soporte"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_contratacion"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_ordenador"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_supervisor"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/supervisor",
			beego.NSInclude(
				&controladores_supervisor.ContratosSupervisorController{}, &controladores_supervisor.SolicitudContratoController{}, &controladores_supervisor.InformeSeguimientoController{}, &controladores_supervisor.ListarTipoDocumentosCumplidoController{}, &controladores_supervisor.InformeSeguimientoController{},
			),
		),
		beego.NSNamespace("/solicitud-pago",
			beego.NSInclude(
				&controladores_soporte.SoportesCumplidoController{},
			),
		),
		beego.NSNamespace("/solicitud-pago",
			beego.NSInclude(
				&controladores_estado_pago.EstadoSoporteController{},
			),
		),
		beego.NSNamespace("/ping",
			beego.NSInclude(
				&controladores_ordenador.PingController{},
			),
		),
		beego.NSNamespace("/ordenador",
			beego.NSInclude(
				&controladores_ordenador.RevisionCumplidoOrdenadorController{},
			),
		),
		beego.NSNamespace("/contratacion",
			beego.NSInclude(
				&controladores_contratacion.RevisionCumplidoContratacionController{},
			),
		))

	beego.AddNamespace(ns)
}
