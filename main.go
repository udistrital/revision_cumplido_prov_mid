package main

import (
	_ "revision_cumplidos_proveedores_mid/routers"

	"github.com/astaxie/beego"
	apistatus "github.com/udistrital/utils_oas/apiStatusLib"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	apistatus.Init()
	beego.Run()
}
