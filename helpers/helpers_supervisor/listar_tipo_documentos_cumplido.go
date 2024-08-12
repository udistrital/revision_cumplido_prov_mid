package helpers_supervisor

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func GetTiposDocumentosCumplido() (tipos_documento []models.DocumentoCumplido, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetTiposDocumentosCumplido", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	//var respuesta_peticion map[string]interface{}
	var tipo_documento []models.TipoDocumento

	//fmt.Println("UrlcrudAgora: ", beego.AppConfig.String("UrlcrudCore")+"/tipo_documento/?query=DominioTipoDocumento.Id:12")
	if response, err := getJsonTest(beego.AppConfig.String("UrlcrudCore")+"/tipo_documento/?query=DominioTipoDocumento.Id:12&limit=0", &tipo_documento); err == nil && response == 200 {
		//LimpiezaRespuestaRefactor(respuesta_peticion, &tipo_documento)
		for _, tipo := range tipo_documento {
			fmt.Println("Tipo: ", tipo)
			var documento models.DocumentoCumplido
			documento.IdTipoDocumento = tipo.Id
			documento.Nombre = tipo.Nombre
			if tipo.NumeroOrden == 0 {
				documento.Tipo = "Obligatorio"
			} else if tipo.NumeroOrden == 1 {
				documento.Tipo = "Autogenerado-Obligatorio"
			} else if tipo.NumeroOrden == 2 {
				documento.Tipo = "Aprobacion Pago"
			} else {
				documento.Tipo = "Acta liquidacion y evaluacion"
			}
			tipos_documento = append(tipos_documento, documento)
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/GetTiposDocumentosCumplido", "err": "No se pudo obtener los tipos de documentos", "status": "502"}
		return nil, outputError
	}
	return tipos_documento, nil
}
