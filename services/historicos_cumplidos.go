package services

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
	"strconv"
)

func ObtberHistoricoEstado(cumplido_proveedor_id string) (historicos []models.HistoricoCumplido, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var peticion_info_persona []models.InformacionPersonaNatural
	var respuesta_histroricos []models.CambioEstadoCumplido
	var peticion_historicos map[string]interface{}

	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/?query=CumplidoProveedorId.Id:"+cumplido_proveedor_id+"&sortby=FechaCreacion&order=desc&limit=-1", &peticion_historicos); err == nil && response == 200 {

		data := peticion_historicos["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) == 0 {

			return historicos, nil
		}
		helpers.LimpiezaRespuestaRefactor(peticion_historicos, &respuesta_histroricos)

	} else {
		outputError = fmt.Errorf("Error al consultar el historico del cumplido")
		return nil, outputError
	}

	urlRequet := "/informacion_persona_natural?fields=PrimerNombre,SegundoNombre,PrimerApellido,SegundoApellido&limit=0&query=Id:" + strconv.Itoa(respuesta_histroricos[0].DocumentoResponsable)

	for _, historico_estado := range respuesta_histroricos {
		fmt.Println(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores") + "/cambio_estado_cumplido/?query=CumplidoProveedorId.Id:" + cumplido_proveedor_id + "&sortby=FechaCreacion&order=desc")
		if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlcrudAgora")+urlRequet, &peticion_info_persona); response == 200 && err == nil {
			println(peticion_info_persona[0].PrimerApellido)
			historico := models.HistoricoCumplido{
				NombreResponsable: peticion_info_persona[0].PrimerNombre + " " + peticion_info_persona[0].SegundoNombre + " " + peticion_info_persona[0].PrimerApellido + " " + peticion_info_persona[0].SegundoApellido,
				Estado:            historico_estado.EstadoCumplidoId.Nombre,
				Fecha:             historico_estado.FechaCreacion,
				CargoResponsable:  historico_estado.CargoResponsable,
			}
			historicos = append(historicos, historico)
		} else {
			outputError = fmt.Errorf("Error al consultar la informacion de la persona")
			return nil, outputError
		}

	}
	return historicos, nil
}
