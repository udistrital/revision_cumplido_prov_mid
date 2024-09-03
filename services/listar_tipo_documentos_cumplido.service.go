package services

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtenerTiposDocumentosCumplido() (tipos_documento []models.DocumentoCumplido, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetTiposDocumentosCumplido", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	var tipo_documento []models.TipoDocumento
	//fmt.Println("UrlcrudAgora: ", beego.AppConfig.String("UrlDocumentosCrud")+"/tipo_documento/?query=DominioTipoDocumento.Id:12")
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlDocumentosCrud")+"/tipo_documento/?query=DominioTipoDocumento.Id:12&limit=0", &tipo_documento); err == nil && response == 200 {
		if len(tipo_documento) == 0 {
			outputError = map[string]interface{}{"funcion": "/GetTiposDocumentosCumplido", "err": "No se encontraron tipos de documentos", "status": "404"}
			return nil, outputError
		}
		for _, tipo := range tipo_documento {
			var documento models.DocumentoCumplido
			documento.IdTipoDocumento = tipo.Id
			documento.Nombre = tipo.Nombre
			tipos_documento = append(tipos_documento, documento)
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/GetTiposDocumentosCumplido", "err": "No se pudo obtener los tipos de documentos", "status": "404"}
		return nil, outputError
	}
	return tipos_documento, nil
}
