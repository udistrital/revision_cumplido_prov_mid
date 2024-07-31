package helpers_soporte

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers/helpers_supervisor"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func SoportesComprimido(id_cumplido_proveedor string) (documentos_comprimido models.DocumentosComprimido, outputError map[string]interface{}) {
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
	documentos, error := GetDocumentosPagoMensual(id_cumplido_proveedor)

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
		fileName := fmt.Sprintf("%s_%d.pdf", filepath.Base(documento.Documento.TipoDocumento.Nombre), i)
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

	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cumplido_proveedor/?query=Id:"+id_cumplido_proveedor, &respuesta_peticion); (err == nil) && (response == 200) {
		if respuesta_peticion != nil {
			LimpiezaRespuestaRefactor(respuesta_peticion, &cumplidos_proveedor)
		}
	}

	var informacion_contrato_contratista models.InformacionContratoContratista
	for _, cumplido := range cumplidos_proveedor {
		informacion_contrato_contratista, error = helpers_supervisor.GetInformacionContratoContratista(cumplido.NumeroContrato, strconv.Itoa(cumplido.VigenciaContrato))

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
