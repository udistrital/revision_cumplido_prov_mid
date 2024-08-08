package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_supervisor:ContratosSupervisorController"] = append(beego.GlobalControllerRouter["github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_supervisor:ContratosSupervisorController"],
        beego.ControllerComments{
            Method: "GetContratos",
            Router: "/:documento_supervisor",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
