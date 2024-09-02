package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:CambioEstadoCumplidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:CambioEstadoCumplidoController"],
        beego.ControllerComments{
            Method: "CambioEstadoCumplido",
            Router: "/cambio-estado",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:ContratosSupervisorController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:ContratosSupervisorController"],
        beego.ControllerComments{
            Method: "ObtenerContratosSupervisor",
            Router: "/contratos-supervisor/:documento_supervisor",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:CumplidoSatisfaccionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:CumplidoSatisfaccionController"],
        beego.ControllerComments{
            Method: "ObtenerBalanceFinancieroContrato",
            Router: "/balance-financiero-contrato/:numero_contrato_suscrito/:vigencia_contrato",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:CumplidoSatisfaccionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:CumplidoSatisfaccionController"],
        beego.ControllerComments{
            Method: "GenerarCumplidoSatisfaccion",
            Router: "/cumplido-satisfaccion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:ListarTipoDocumentosCumplidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:ListarTipoDocumentosCumplidoController"],
        beego.ControllerComments{
            Method: "ObtenerTiposDocumentosCumplido",
            Router: "/tipos-documentos-cumplido",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:RevisionCumplidoContratacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:RevisionCumplidoContratacionController"],
        beego.ControllerComments{
            Method: "ObtenerCumplidosPendientesContratacion",
            Router: "/solicitudes-pago/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:RevisionCumplidoOrdenadorController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:RevisionCumplidoOrdenadorController"],
        beego.ControllerComments{
            Method: "GenerarPdfAutorizacionPago",
            Router: "/certificado-aprobacion-pago/:id_solicitud_pago",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:RevisionCumplidoOrdenadorController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:RevisionCumplidoOrdenadorController"],
        beego.ControllerComments{
            Method: "ListaCumplidosReversibles",
            Router: "/revertir-solicitud-pago/:documento_ordenador",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:RevisionCumplidoOrdenadorController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:RevisionCumplidoOrdenadorController"],
        beego.ControllerComments{
            Method: "ObtenerCumplidosPendientesRevisionOrdenador",
            Router: "/solicitudes-pago/:documento_ordenador",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:SolicitudesCumplidosContratoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:SolicitudesCumplidosContratoController"],
        beego.ControllerComments{
            Method: "ObtenerSolicitudesContrato",
            Router: "/solicitudes-contrato/:numero_contrato/:vigencia",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:SoportesCumplidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:SoportesCumplidoController"],
        beego.ControllerComments{
            Method: "AgregarComentarioSoporte",
            Router: "/comentario-soporte",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:SoportesCumplidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:SoportesCumplidoController"],
        beego.ControllerComments{
            Method: "SubirSoporteCumplido",
            Router: "/soportes",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:SoportesCumplidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:SoportesCumplidoController"],
        beego.ControllerComments{
            Method: "ObtenerComprimidoSoportes",
            Router: "/soportes-comprimido/:id_cumplido_proveedor",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:SoportesCumplidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:SoportesCumplidoController"],
        beego.ControllerComments{
            Method: "ObtenerDocumentosPagoMensual",
            Router: "/soportes/:cumplido_proveedor_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:SoportesCumplidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers:SoportesCumplidoController"],
        beego.ControllerComments{
            Method: "EliminarSoporteCumplido",
            Router: "/soportes/:soporte_pago_id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
