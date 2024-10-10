package services

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtenerSoportesCumplido(cumplido_proveedor_id string) (documentos []models.DocumentosSoporteSimplificado, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var soportes_pagos_mensuales []models.SoporteCumplido
	var documentos_crud []models.Documento
	var fileGestor models.FileGestorDocumental
	var soporte models.DocumentosSoporteSimplificado
	var documento_individual models.DocumentoSimplificado
	var a = map[int]int{}

	var respuesta_peticion map[string]interface{}

	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/soporte_cumplido/?limit=-1&query=CumplidoProveedorId.Id:"+cumplido_proveedor_id+",Activo:true", &respuesta_peticion); (err == nil) && (response == 200) {
		data := respuesta_peticion["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) == 0 {
			outputError = fmt.Errorf("No se encontraron soportes de pago")
			return nil, outputError
		}
		helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &soportes_pagos_mensuales)
		if len(soportes_pagos_mensuales) != 0 {
			var ids_documentos []string
			for _, soporte_pago_mensual := range soportes_pagos_mensuales {
				ids_documentos = append(ids_documentos, strconv.Itoa(soporte_pago_mensual.DocumentoId))
				a[soporte_pago_mensual.DocumentoId] = soporte_pago_mensual.Id
			}

			var ids_documentos_juntos = strings.Join(ids_documentos, "|")
			//fmt.Println("URL documentos juntos: ", beego.AppConfig.String("UrlDocumentosCrud")+"/documento/?limit=-1&query=Activo:True,Id.in:"+ids_documentos_juntos)
			if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlDocumentosCrud")+"/documento/?limit=-1&query=Activo:True,Id.in:"+ids_documentos_juntos, &documentos_crud); (err == nil) && (response == 200) {
				if len(documentos_crud) > 0 {
					for _, documento_crud := range documentos_crud {
						soporte.SoporteCumplidoId = a[documento_crud.Id]
						var observaciones map[string]interface{}
						documento_individual.Id = documento_crud.Id
						documento_individual.Nombre = documento_crud.Nombre
						documento_individual.TipoDocumento = documento_crud.TipoDocumento.Nombre
						documento_individual.Descripcion = documento_crud.Descripcion
						documento_individual.FechaCreacion = documento_crud.FechaCreacion
						documento_individual.CodigoAbreviacionTipoDocumento = documento_crud.TipoDocumento.CodigoAbreviacion
						if err := json.Unmarshal([]byte(documento_crud.Metadatos), &observaciones); err == nil {
							documento_individual.Observaciones = observaciones["observaciones"].(string)
						}
						soporte.Documento = documento_individual
						if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlGestorDocumental")+"/document/"+documento_crud.Enlace, &fileGestor); (err == nil) && (response == 200) {
							soporte.Archivo = fileGestor
							documentos = append(documentos, soporte)
						} else {
							logs.Error(err)
							continue
						}
					}
				} else {
					logs.Error(err)
					outputError = fmt.Errorf("No se encontraron documentos asociados al cumplido proveedor")
					return nil, outputError
				}

			} else {
				logs.Error(err)
				outputError = fmt.Errorf("Error al obtener los documentos del pago")
				return nil, outputError
			}
		} else {
			return nil, outputError
		}

	} else {
		logs.Error(err)
		outputError = fmt.Errorf("Error al ontener el soporte cumplido proveedor")
		return nil, outputError
	}

	return
}

func ObtenerComentariosSoporte(soporte_id int) (comentarios []models.ComentarioSoporte) {
	var respuesta_peticion map[string]interface{}
	var comentarios_soporte []models.ComentarioSoporte
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/comentario_soporte/?limit=-1&query=SoporteCumplidoId.Id:"+strconv.Itoa(soporte_id)+",Activo:true&sortby=FechaCreacion&order=desc", &respuesta_peticion); (err == nil) && (response == 200) {
		data := respuesta_peticion["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) == 0 {
			return comentarios
		}
		helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &comentarios_soporte)
		return comentarios_soporte
	} else {
		return comentarios
	}
}

func ObtenerComprimidoSoportes(id_cumplido_proveedor string) (documentos_comprimido models.DocumentosComprimido, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	//Realizar la solicitud opara obtener los documentos asociados al pago

	var respuesta_peticion map[string]interface{}
	var cumplidos_proveedor []models.CumplidoProveedor
	documentos, error := ObtenerSoportesCumplido(id_cumplido_proveedor)

	if error != nil {
		outputError = fmt.Errorf("Error al obtener los documentos del pago")
		return documentos_comprimido, outputError
	} else if len(documentos) == 0 {
		return documentos_comprimido, nil
	}

	//Crear un archivo ZIP

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Decodificar los archivos en base64 y agregarlos al .zip
	for i, documento := range documentos {
		pdfData, error := base64.StdEncoding.DecodeString(documento.Archivo.File)
		if error != nil {
			outputError = fmt.Errorf("Error al decodificar el archivo base64")
			return documentos_comprimido, outputError
		}

		// Crear una entrada en el ZIP para cada archivo PDF con su nombre específico y un índice único
		fileName := fmt.Sprintf("%s_%d.pdf", filepath.Base(documento.Documento.Nombre), i)
		zipEntry, err := zipWriter.Create(fileName)
		if err != nil {
			outputError = fmt.Errorf("Error al crear la entrada en el archivo ZIP")
			return documentos_comprimido, outputError
		}

		// Escribir los datos del PDF en la entrada del ZIP
		_, err = zipEntry.Write(pdfData)
		if err != nil {
			outputError := fmt.Errorf("Error al escribir el contenido del PDF en el archivo ZIP")
			return documentos_comprimido, outputError
		}
	}

	// Cerrar el writer del ZIP
	err := zipWriter.Close()
	if err != nil {
		outputError := fmt.Errorf("Error al cerrar el archivo ZIP")
		return documentos_comprimido, outputError
	}

	documentos_comprimido.File = base64.StdEncoding.EncodeToString(buf.Bytes())

	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cumplido_proveedor/?query=Id:"+id_cumplido_proveedor+"&limit=-1", &respuesta_peticion); (err == nil) && (response == 200) {
		data := respuesta_peticion["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) > 0 {
			helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &cumplidos_proveedor)
		}
	}

	for _, cumplido := range cumplidos_proveedor {
		informacion_contrato_proveedor, error := helpers.ObtenerInformacionContratoProveedor(cumplido.NumeroContrato, strconv.Itoa(cumplido.VigenciaContrato))

		if error == nil {
			documentos_comprimido.Nombre = informacion_contrato_proveedor[0].NombreProveedor + "_" + cumplido.NumeroContrato + "_" + "_" + strconv.Itoa(int(cumplido.FechaCreacion.Month())) + "_" + strconv.Itoa(cumplido.FechaCreacion.Year())
		} else {
			logs.Error(err)
			outputError = fmt.Errorf("Error al Buscar los datos del proveedor")
			continue
		}
	}

	return documentos_comprimido, nil
}

func SubirSoporteCumplido(solicitud_pago_id int, tipo_documento string, item_id int, observaciones string, nombre_archivo string, archivo string) (soporte_pago models.SoporteCumplido, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	// Verificar que se envien todos los datos y que el archivo sea un PDF
	if tipo_documento != "application/pdf" || archivo == "" || item_id == 0 || solicitud_pago_id == 0 {
		outputError = fmt.Errorf("El archivo debe ser un PDF y no debe estar vacío")
		return soporte_pago, outputError
	}

	// Convertir archivo Base64 a binario
	decodedFile, err := base64.StdEncoding.DecodeString(archivo)
	if err != nil {
		outputError = fmt.Errorf("Error al decodificar el archivo Base64")
		return soporte_pago, outputError
	}

	// Verificar tamaño del archivo (máximo 5MB)
	if len(decodedFile) > 5000000 {
		outputError = fmt.Errorf("El archivo no debe superar 5MB")
		return soporte_pago, outputError
	}

	var respuesta_peticion map[string]interface{}
	var cumplido_proveedor []models.CumplidoProveedor
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cumplido_proveedor/?query=Id:"+strconv.Itoa(solicitud_pago_id)+"&limit=-1", &respuesta_peticion); err == nil && response == 200 {
		data := respuesta_peticion["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) == 0 {
			outputError = fmt.Errorf("El cumplido proveedor no existe")
			return soporte_pago, outputError
		}
		helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &cumplido_proveedor)
	} else {
		outputError = fmt.Errorf("Error al consultar el cumplido proveedor")
		return soporte_pago, outputError
	}

	data := []models.BodySubirSoporte{
		{
			IdTipoDocumento: item_id,
			Nombre:          nombre_archivo,
			Metadatos: struct {
				NombreArchivo string `json:"nombre_archivo"`
				Tipo          string `json:"tipo"`
				Observaciones string `json:"observaciones"`
			}{
				NombreArchivo: nombre_archivo,
				Tipo:          "Archivo",
				Observaciones: observaciones,
			},
			Descripcion: observaciones,
			File:        archivo,
		},
	}

	var respuesta map[string]interface{}

	// Realizar la solicitud

	//fmt.Println("URL Subir documento gestor documental: ", beego.AppConfig.String("UrlGestorDocumental")+"/document/upload")
	if err := helpers.SendJson(beego.AppConfig.String("UrlGestorDocumental")+"/document/upload", "POST", &respuesta, data); err == nil {
		id := respuesta["res"].(map[string]interface{})["Id"].(float64)
		soporte := models.BodySoportePago{
			DocumentoId:         int(id),
			CumplidoProveedorId: cumplido_proveedor[0],
		}
		var res map[string]interface{}
		if err := helpers.SendJson(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/soporte_cumplido", "POST", &res, soporte); err == nil {
			helpers.LimpiezaRespuestaRefactor(res, &soporte_pago)
			return soporte_pago, nil
		} else {
			outputError = fmt.Errorf("Error al subir el soporte")
			return soporte_pago, outputError
		}

	} else {
		outputError = fmt.Errorf("Error al subir el soporte")
		return soporte_pago, outputError
	}

}
