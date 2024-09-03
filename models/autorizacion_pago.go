package models

type DocumentoAutorizacionPago struct {
	NombreArchivo        string `json:"NombreArchivo"`
	NombreResponsable    string `json:"NombreResponsable"` 
	CargoResponsable     string `json:"CargoResponsable"` 
	DescripcionDocumento string `json:"DescripcionDocumento"` 
	Archivo              string `json:"Archivo"`
}
