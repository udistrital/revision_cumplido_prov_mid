package models

type AutorizacionPago struct {
	SolicitudPagoId string
	Item            string
	Observaciones   string
	NombreArchivo   string
	Archivo         string
}

type DocuementoAutorizacionPago struct {
	NombreOrdenador    string
	DocumentoOrdenador string
	NombreProveedor    string
	DocumentoProveedor string
	ValorPago          string
}
