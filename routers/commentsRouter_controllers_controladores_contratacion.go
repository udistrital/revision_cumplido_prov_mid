package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_contratacion:RevisionCumplidoContratacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_contratacion:RevisionCumplidoContratacionController"],
		beego.ControllerComments{
			Method:           "GenerarPdf",
			Router:           "/certificado-aprobacion-pago/:id_solicitud_pago",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_contratacion:RevisionCumplidoContratacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_contratacion:RevisionCumplidoContratacionController"],
		beego.ControllerComments{
			Method:           "ObtenerPendientesContratacion",
			Router:           "/solicitudes-pago/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
