package controladores_ordenador

import (
	"github.com/astaxie/beego"
)

type PingController struct {
	beego.Controller
}

// URLMapping asigna métodos a la estructura PingController
func (c *PingController) URLMapping() {
	c.Mapping("Ping", c.Ping)
}

// ping realiza una prueba de ping y devuelve una respuesta JSON
// @Title PingTest
// @Description Realiza una prueba de ping y devuelve el estado.
// @Success 200 {object} map[string]interface{}{"Success":true,"Status":"200","Message":"Successful","Data":{}}
// @Failure 500 {object} map[string]interface{}{"Success":false,"Status":"500","Message":"Internal Server Error"}
// @router /ping [get]
func (c *PingController) Ping() {
	beego.Info("Hola desde el método ping")

	var miCadena = "Hola, mundo"

	c.Ctx.Output.Header("Content-Type", "text/plain")
	c.Ctx.Output.Body([]byte(miCadena))
}
