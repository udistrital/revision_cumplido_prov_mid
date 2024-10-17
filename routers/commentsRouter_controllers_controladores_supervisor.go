package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_supervisor:ContratosSupervisorController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_supervisor:ContratosSupervisorController"],
		beego.ControllerComments{
			Method:           "GetContratosSupervisor",
			Router:           "/contratos-supervisor/:documento_supervisor",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_supervisor:InformeSeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_supervisor:InformeSeguimientoController"],
		beego.ControllerComments{
			Method:           "GetBalanceFinancieroContrato",
			Router:           "/balance-financiero-contrato/:numero_contrato_suscrito/:vigencia_contrato",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_supervisor:InformeSeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_supervisor:InformeSeguimientoController"],
		beego.ControllerComments{
			Method:           "GenerateInformeSeguimiento",
			Router:           "/informe-seguimiento",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_supervisor:ListarTipoDocumentosCumplidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_supervisor:ListarTipoDocumentosCumplidoController"],
		beego.ControllerComments{
			Method:           "GetTiposDocumentosCumplido",
			Router:           "/tipos-documentos-cumplido",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_supervisor:SolicitudContratoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_supervisor:SolicitudContratoController"],
		beego.ControllerComments{
			Method:           "GetSolicitudesContrato",
			Router:           "/solicitudes-contrato/:numero_contrato/:vigencia",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
