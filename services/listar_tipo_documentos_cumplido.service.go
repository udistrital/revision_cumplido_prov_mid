package services

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtenerTiposDocumentosCumplido() (tipos_documento []models.DocumentoCumplido, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var tipo_documento []models.TipoDocumento
	//fmt.Println("UrlcrudAgora: ", beego.AppConfig.String("UrlDocumentosCrud")+"/tipo_documento/?query=DominioTipoDocumento.Id:12")
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlDocumentosCrud")+"/tipo_documento/?query=DominioTipoDocumento.CodigoAbreviacion:CUMP_PROV&limit=0", &tipo_documento); err == nil && response == 200 {
		if len(tipo_documento) == 0 {
			outputError = fmt.Errorf("No se encontraron tipos de documentos")
			return nil, outputError
		}
		for _, tipo := range tipo_documento {
			var documento models.DocumentoCumplido
			documento.IdTipoDocumento = tipo.Id
			documento.CodigoAbreviacionTipoDocumento = tipo.CodigoAbreviacion
			documento.Nombre = tipo.Nombre
			tipos_documento = append(tipos_documento, documento)
		}
	} else {
		logs.Error(err)
		outputError = fmt.Errorf("No se pudo obtener los tipos de documentos")
		return nil, outputError
	}
	return tipos_documento, nil
}
