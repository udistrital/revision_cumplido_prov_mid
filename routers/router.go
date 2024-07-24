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
	"github.com/udistrital/revision_cumplidos_proveedores_mid/controllers/controladores_supervisor"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/dependencias_supervisor",
			beego.NSInclude(
				&controladores_supervisor.ContratosSupervisorController{},
			),
		))
	beego.AddNamespace(ns)
}
