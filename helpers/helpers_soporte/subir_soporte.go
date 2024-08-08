package helpers_soporte

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func SubirSoporte(solicitud_pago_id int, tipo_documento string, item_id int, observaciones string, nombre_archivo string, archivo string) (soporte_pago models.SoportePago, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/SubirSoporte", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	// Verificar que se envien todos los datos y que el archivo sea un PDF
	if tipo_documento != "application/pdf" || archivo == "" || item_id == 0 || solicitud_pago_id == 0 {
		outputError = map[string]interface{}{"funcion": "/SubirSoporte", "status": "502", "mensaje": "El archivo debe ser un PDF y no debe estar vacío"}
		return soporte_pago, outputError
	}

	// Convertir archivo Base64 a binario
	decodedFile, err := base64.StdEncoding.DecodeString(archivo)
	if err != nil {
		outputError = map[string]interface{}{"funcion": "/SubirSoporte", "err": err, "status": "502", "mensaje": "Error al decodificar el archivo Base64"}
		return soporte_pago, outputError
	}

	// Verificar tamaño del archivo (máximo 1MB)
	if len(decodedFile) > 1000000 {
		outputError = map[string]interface{}{"funcion": "/SubirSoporte", "err": err, "status": "502", "mensaje": "El archivo no debe superar 1MB"}
		return soporte_pago, outputError
	}

	var respuesta_peticion map[string]interface{}
	var cumplido_proveedor []models.CumplidoProveedor
	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cumplido_proveedor/?query=Id:"+strconv.Itoa(solicitud_pago_id), &respuesta_peticion); err == nil && response == 200 {
		LimpiezaRespuestaRefactor(respuesta_peticion, &cumplido_proveedor)
		fmt.Println("Cumplido Proveedor: ", &cumplido_proveedor[0])
	} else {
		outputError = map[string]interface{}{"funcion": "/SubirSoporte", "status": "502", "mensaje": "Error al consultar el cumplido proveedor"}
		return soporte_pago, outputError
	}

	//var respuesta map[string]interface{}
	var tipo []models.TipoDocumento
	fmt.Println("Item ID: ", item_id)
	fmt.Println("URL: ", beego.AppConfig.String("UrlDocumentosCrud")+"/tipo_documento/?query=Id:"+strconv.Itoa(item_id))
	if response, err := getJsonTest(beego.AppConfig.String("UrlDocumentosCrud")+"/tipo_documento/?query=Id:"+strconv.Itoa(item_id), &tipo); err == nil && response == 200 {
		//LimpiezaRespuestaRefactor(respuesta, &tipo)
		fmt.Println("Tipo Documento: ", tipo)
	} else {
		outputError = map[string]interface{}{"funcion": "/SubirSoporte", "status": "502", "error": err, "mensaje": "Error al consultar el tipo de documento"}
		return soporte_pago, outputError
	}

	data := []models.BodySubirSoporte{
		{
			IdTipoDocumento: tipo[0].Id,
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
			Descripcion: tipo[0].Descripcion,
			File:        archivo,
		},
	}

	var respuesta map[string]interface{}

	// Realizar la solicitud

	if err := sendJson(beego.AppConfig.String("UrlGestorDocumental")+"/document/upload", "POST", &respuesta, data); err == nil {
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
			if err := sendJson(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/soporte_pago", "POST", &res, soporte); err == nil {
				LimpiezaRespuestaRefactor(res, &soporte_pago)
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
