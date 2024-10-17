package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_soporte:SoportesCumplidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_soporte:SoportesCumplidoController"],
		beego.ControllerComments{
			Method:           "AgregarComentarioSoporte",
			Router:           "/comentario-soporte",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_soporte:SoportesCumplidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_soporte:SoportesCumplidoController"],
		beego.ControllerComments{
			Method:           "SubirSoporte",
			Router:           "/soportes",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_soporte:SoportesCumplidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_soporte:SoportesCumplidoController"],
		beego.ControllerComments{
			Method:           "GetSoportesComprimido",
			Router:           "/soportes-comprimido/:id_cumplido_proveedor",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_soporte:SoportesCumplidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_soporte:SoportesCumplidoController"],
		beego.ControllerComments{
			Method:           "GetDocumentosPagoMensual",
			Router:           "/soportes/:cumplido_proveedor_id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_soporte:SoportesCumplidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_soporte:SoportesCumplidoController"],
		beego.ControllerComments{
			Method:           "EliminarSoporteCumplido",
			Router:           "/soportes/:soporte_pago_id",
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
