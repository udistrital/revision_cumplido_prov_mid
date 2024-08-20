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
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func AgregarComentarioSoporte(soporte_id string, cambio_estado_id string, comentario string) (respuesta models.RespuestaComentarioSoporte, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/AgregarComentarioSoporte", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	var respuesta_peticion map[string]interface{}
	var comentario_soporte models.ComentarioSoporte
	var soporte_pago models.SoporteCumplido
	var cambio_estado_cumplido models.CambioEstadoCumplido

	if soporte_id == "" || cambio_estado_id == "" || comentario == "" {
		outputError = map[string]interface{}{"funcion": "/AgregarComentarioSoporte", "err": "Faltan datos en la solicitud", "status": "400"}
		return respuesta, outputError
	}

	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/soporte_cumplido/"+soporte_id, &respuesta_peticion); (err == nil) && (response == 200) {
		helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &soporte_pago)
		if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/"+cambio_estado_id, &respuesta_peticion); (err == nil) && (response == 200) {
			helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &cambio_estado_cumplido)
			comentario_soporte.Comentario = comentario
			comentario_soporte.SoporteCumplidoId = &soporte_pago
			comentario_soporte.CambioEstadoCumplidoId = &cambio_estado_cumplido
			if err := helpers.SendJson(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/comentario_soporte", "POST", &respuesta, comentario_soporte); err == nil {
				respuesta.SoportePagoId = soporte_pago.Id
				respuesta.CambioEstadoCumplidoId = cambio_estado_cumplido.Id
				respuesta.Comentario = comentario
				return respuesta, nil
			} else {
				outputError = map[string]interface{}{"funcion": "/AgregarComentarioSoporte/comentario_soporte", "err": err, "status": "502"}
				return respuesta, outputError
			}
		} else {
			outputError = map[string]interface{}{"funcion": "/AgregarComentarioSoporte/cambio_estado_cumplido", "err": err, "status": "502"}
			return respuesta, outputError
		}
	}
	return respuesta, outputError
}

func EliminarSoporteCumplido(documento_id string) (response string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetDocumentosPagoMensual", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	var soportes_pagos_mensuales []models.SoporteCumplido
	var respuesta_peticion map[string]interface{}

	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/soporte_cumplido/?limit=1&query=DocumentoId:"+documento_id+",Activo:true", &respuesta_peticion); (err == nil) && (response == 200) {
		helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &soportes_pagos_mensuales)
		if (soportes_pagos_mensuales[0] == models.SoporteCumplido{}) {
			return "No se encontró el soporte de pago o este ya se elimino con anterioridad", nil
		}
	} else {
		outputError = map[string]interface{}{"funcion": "/EliminarSoporteCumplido/soporte_cumplido", "err": err, "status": "502"}
		return "No se encontró el soporte de pago", outputError
	}
	var res map[string]interface{}
	delete_true := "Soporte pago eliminado correctamente"
	delect_false := "No se encontró el soporte de pago"

	if err := helpers.SendJson(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/soporte_cumplido/"+strconv.Itoa(soportes_pagos_mensuales[0].Id), "DELETE", &res, nil); err == nil {
		response = delete_true
		return response, nil
	} else {
		outputError = map[string]interface{}{"funcion": "/EliminarSoporteCumplido/soporte_cumplido", "err": err, "status": "502"}
		response = delect_false
		return response, outputError
	}

}

func ObtenerDocumentosPagoMensual(cumplido_proveedor_id string) (documentos []models.DocumentosSoporteCorto, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetDocumentosPagoMensual", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	var soportes_pagos_mensuales []models.SoporteCumplido
	var documentos_crud []models.Documento
	var fileGestor models.FileGestorDocumental
	var soporte models.DocumentosSoporteCorto
	var documento_individual models.DocumentoCorto

	var respuesta_peticion map[string]interface{}

	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/soporte_cumplido/?limit=-1&query=CumplidoProveedorId.Id:"+cumplido_proveedor_id+",Activo:true", &respuesta_peticion); (err == nil) && (response == 200) {
		helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &soportes_pagos_mensuales)
		if len(soportes_pagos_mensuales) != 0 {
			var ids_documentos []string
			for _, soporte_pago_mensual := range soportes_pagos_mensuales {
				ids_documentos = append(ids_documentos, strconv.Itoa(soporte_pago_mensual.DocumentoId))
			}

			var ids_documentos_juntos = strings.Join(ids_documentos, "|")
			//fmt.Println("URL documentos juntos: ", beego.AppConfig.String("UrlDocumentosCrud")+"/documento/?limit=-1&query=Activo:True,Id.in:"+ids_documentos_juntos)
			if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlDocumentosCrud")+"/documento/?limit=-1&query=Activo:True,Id.in:"+ids_documentos_juntos, &documentos_crud); (err == nil) && (response == 200) {
				if len(documentos_crud) > 0 {
					for _, documento_crud := range documentos_crud {
						var observaciones map[string]interface{}
						documento_individual.Id = documento_crud.Id
						documento_individual.Nombre = documento_crud.Nombre
						documento_individual.TipoDocumento = documento_crud.TipoDocumento.Nombre
						documento_individual.Descripcion = documento_crud.Descripcion
						documento_individual.FechaCreacion = documento_crud.FechaCreacion
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
					outputError = map[string]interface{}{"funcion": "/GetDocumentosPagoMensual/documento", "err": err, "status": "502", "Message": "No se encontraron documentos asociados al cumplido proveedor"}
					return nil, outputError
				}

			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/GetDocumentosPagoMensual/documento", "err": err, "status": "502"}
				return nil, outputError
			}
		} else {
			return nil, outputError
		}

	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/GetDocumentosPagoMensual/soporte_pago_mensual", "err": err, "status": "502"}
		return nil, outputError
	}

	return
}

func ObtenerComprimidoSoportes(id_cumplido_proveedor string) (documentos_comprimido models.DocumentosComprimido, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"Succes":  false,
				"Status":  502,
				"Message": "Error al descargar el .zip",
				"Funcion": "SoportesComprimido",
			}
			panic(outputError)
		}
	}()

	//Realizar la solicitud opara obtener los documentos asociados al pago

	var respuesta_peticion map[string]interface{}
	var cumplidos_proveedor []models.CumplidoProveedor
	documentos, error := ObtenerDocumentosPagoMensual(id_cumplido_proveedor)

	if error != nil {
		outputError = map[string]interface{}{
			"Succes":  false,
			"Status":  502,
			"Message": "Error al obtener los documentos del pago",
			"Funcion": "GetDocumentosPagoMensual",
			"Error":   error,
		}
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
			outputError = map[string]interface{}{
				"Succes":  false,
				"Status":  502,
				"Message": "Error al decodificar el archivo base64",
				"Error":   error,
			}
			return documentos_comprimido, outputError
		}

		// Crear una entrada en el ZIP para cada archivo PDF con su nombre específico y un índice único
		fileName := fmt.Sprintf("%s_%d.pdf", filepath.Base(documento.Documento.Nombre), i)
		zipEntry, err := zipWriter.Create(fileName)
		if err != nil {
			outputError = map[string]interface{}{
				"Success": false,
				"Status":  502,
				"Message": "Error al crear la entrada en el archivo ZIP",
				"Error":   err.Error(),
			}
			return documentos_comprimido, outputError
		}

		// Escribir los datos del PDF en la entrada del ZIP
		_, err = zipEntry.Write(pdfData)
		if err != nil {
			outputError := map[string]interface{}{
				"Success": false,
				"Status":  502,
				"Message": "Error al escribir el contenido del PDF en el archivo ZIP",
				"Error":   err.Error(),
			}
			return documentos_comprimido, outputError
		}
	}

	// Cerrar el writer del ZIP
	err := zipWriter.Close()
	if err != nil {
		outputError := map[string]interface{}{
			"Success": false,
			"Status":  502,
			"Message": "Error al cerrar el archivo ZIP",
			"Error":   err.Error(),
		}
		return documentos_comprimido, outputError
	}

	documentos_comprimido.File = base64.StdEncoding.EncodeToString(buf.Bytes())

	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cumplido_proveedor/?query=Id:"+id_cumplido_proveedor, &respuesta_peticion); (err == nil) && (response == 200) {
		if respuesta_peticion != nil {
			helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &cumplidos_proveedor)
		}
	}

	var informacion_contrato_contratista models.InformacionContratoContratista
	for _, cumplido := range cumplidos_proveedor {
		informacion_contrato_contratista, error = helpers.ObtenerInformacionContratoProveedor(cumplido.NumeroContrato, strconv.Itoa(cumplido.VigenciaContrato))

		if error == nil {
			fmt.Println("Nombre: ", informacion_contrato_contratista)
			documentos_comprimido.Nombre = informacion_contrato_contratista.InformacionContratista.NombreCompleto + "_" + cumplido.NumeroContrato + "_" + informacion_contrato_contratista.InformacionContratista.Documento.Numero + "_" + strconv.Itoa(int(cumplido.FechaCreacion.Month())) + "_" + strconv.Itoa(cumplido.FechaCreacion.Year())
		} else {
			outputError := map[string]interface{}{
				"Success": false,
				"Status":  502,
				"Message": "Error al Buscar los datos del contratista",
			}
			return documentos_comprimido, outputError
		}
	}

	return documentos_comprimido, nil
}

func SubirSoporteCumplido(solicitud_pago_id int, tipo_documento string, item_id int, observaciones string, nombre_archivo string, archivo string) (soporte_pago models.SoporteCumplido, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/SubirSoporteCumplido", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	// Verificar que se envien todos los datos y que el archivo sea un PDF
	if tipo_documento != "application/pdf" || archivo == "" || item_id == 0 || solicitud_pago_id == 0 {
		outputError = map[string]interface{}{"funcion": "/SubirSoporteCumplido", "status": "502", "mensaje": "El archivo debe ser un PDF y no debe estar vacío"}
		return soporte_pago, outputError
	}

	// Convertir archivo Base64 a binario
	decodedFile, err := base64.StdEncoding.DecodeString(archivo)
	if err != nil {
		outputError = map[string]interface{}{"funcion": "/SubirSoporteCumplido", "err": err, "status": "502", "mensaje": "Error al decodificar el archivo Base64"}
		return soporte_pago, outputError
	}

	// Verificar tamaño del archivo (máximo 1MB)
	if len(decodedFile) > 1000000 {
		outputError = map[string]interface{}{"funcion": "/SubirSoporteCumplido", "err": err, "status": "502", "mensaje": "El archivo no debe superar 1MB"}
		return soporte_pago, outputError
	}

	var respuesta_peticion map[string]interface{}
	var cumplido_proveedor []models.CumplidoProveedor
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cumplido_proveedor/?query=Id:"+strconv.Itoa(solicitud_pago_id), &respuesta_peticion); err == nil && response == 200 {
		helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &cumplido_proveedor)
		fmt.Println("Cumplido proveedor: ", cumplido_proveedor)
	} else {
		outputError = map[string]interface{}{"funcion": "/SubirSoporteCumplido", "status": "502", "mensaje": "Error al consultar el cumplido proveedor"}
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

	fmt.Println("Data: ", data)

	var respuesta map[string]interface{}

	// Realizar la solicitud

	fmt.Println("URL Subir documento gestor documental: ", beego.AppConfig.String("UrlGestorDocumental")+"/document/upload")
	if err := helpers.SendJson(beego.AppConfig.String("UrlGestorDocumental")+"/document/upload", "POST", &respuesta, data); err == nil {
		id := respuesta["res"].(map[string]interface{})["Id"].(float64)
		soporte := models.BodySoportePago{
			DocumentoId:         int(id),
			CumplidoProveedorId: cumplido_proveedor[0],
			FechaCreacion:       time.Now(),
			FechaModificacion:   time.Now(),
			Activo:              true,
		}
		fmt.Println("Soporte: ", soporte)
		var res map[string]interface{}
		if err == nil {
			if err := helpers.SendJson(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/soporte_cumplido", "POST", &res, soporte); err == nil {
				helpers.LimpiezaRespuestaRefactor(res, &soporte_pago)
				return soporte_pago, nil
			} else {
				outputError = map[string]interface{}{"funcion": "/SubirSoporte", "status": "502", "error": err, "mensaje": "Error al subir el soporte"}
				return soporte_pago, outputError
			}
		} else {
			outputError = map[string]interface{}{"funcion": "/SubirSoporte", "status": "502", "error": err, "mensaje": "Error al convertir soporte a JSON"}
			return soporte_pago, outputError
		}

	} else {
		outputError = map[string]interface{}{"funcion": "/SubirSoporte", "status": "502", "error": err, "mensaje": "Error al subir el soporte"}
		return soporte_pago, outputError
	}

	return soporte_pago, nil
}
