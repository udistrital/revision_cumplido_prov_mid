package helpers_supervisor

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func GetDependenciasSupervisor(documento string) (dependenciasList []models.Dependencia, errorOutput interface{}) {
	defer func() {
		if err := recover(); err != nil {
			errorOutput = map[string]interface{}{
				"Success": true,
				"Status":  502,
				"Message": "Error al consultar las dependencias: " + documento,
				"Error":   err,
			}
		}
	}()
	var respuesta map[string]interface{}
	response, err := getJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/dependencias_supervisor/"+documento, &respuesta)
	if err == nil && response == 200 {
		if respuesta != nil {
			if dependenciasMap, ok := respuesta["dependencias"].(map[string]interface{}); ok {

				for _, depList := range dependenciasMap {

					if list, ok := depList.([]interface{}); ok {

						for _, dep := range list {

							depMap := dep.(map[string]interface{})
							dependencia := models.Dependencia{

								Codigo: depMap["codigo"].(string),
								Nombre: depMap["nombre"].(string),
							}
							dependenciasList = append(dependenciasList, dependencia)
						}

					}
				}
			}
		}
	} else {
		return nil, errorOutput
	}
	if dependenciasList != nil {
		return dependenciasList, nil
	}

	return nil, nil
}
