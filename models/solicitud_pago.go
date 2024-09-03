package models

type AutorizacionPago struct {
	SolicitudPagoId string
	Item            string
	Observaciones   string
	NombreArchivo   string
	Archivo         string
}

type DatosAutorizacionPago struct {
	NombreOrdenador    string
	DocumentoOrdenador string
	Rubro              string
	NombreProveedor    string
	DocumentoProveedor string
	ValorPago          int
	DocumentosCargados []string
}
