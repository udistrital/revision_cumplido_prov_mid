package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_ordenador:PingController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_ordenador:PingController"],
		beego.ControllerComments{
			Method:           "Ping",
			Router:           "/ping",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_ordenador:RevisionCumplidoOrdenadorController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_ordenador:RevisionCumplidoOrdenadorController"],
		beego.ControllerComments{
			Method:           "ObtenerCertificado",
			Router:           "/certificado-aprobacion-pago/:id_solicitud_pago",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_ordenador:RevisionCumplidoOrdenadorController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_ordenador:RevisionCumplidoOrdenadorController"],
		beego.ControllerComments{
			Method:           "ListaCumplidosReversibles",
			Router:           "/revertir-solicitud-pago/:id_cumplido",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_ordenador:RevisionCumplidoOrdenadorController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_ordenador:RevisionCumplidoOrdenadorController"],
		beego.ControllerComments{
			Method:           "ObtenerPendientesRevisionOrdenador",
			Router:           "/solicitudes-pago/:documento_ordenador",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_ordenador:RevisionCumplidoOrdenadorController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_ordenador:RevisionCumplidoOrdenadorController"],
		beego.ControllerComments{
			Method:           "GenerarPdf",
			Router:           "/testpdf",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
