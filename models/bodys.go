package models

import "time"

type BodySubirSoporte struct {
	IdTipoDocumento int    `json:"IdTipoDocumento"`
	Nombre          string `json:"nombre"`
	Metadatos       struct {
		NombreArchivo string `json:"nombre_archivo"`
		Tipo          string `json:"tipo"`
		Observaciones string `json:"observaciones"`
	} `json:"metadatos"`
	Descripcion string `json:"descripcion"`
	File        string `json:"file"`
}

type BodyCambioEstadoCumplido struct {
	EstadoCumplidoId     EstadoCumplido    `json:"EstadoCumplidoId"`
	CumplidoProveedorId  CumplidoProveedor `json:"CumplidoProveedorId"`
	DocumentoResponsable int               `json:"DocumentoResponsable"`
	CargoResponsable     string            `json:"CargoResponsable"`
	Activo               bool              `json:"Activo"`
	FechaCreacion        time.Time         `json:"FechaCreacion"`
	FechaModificacion    time.Time         `json:"FechaModificacion"`
}

type BodySoportePago struct {
	DocumentoId         int
	CumplidoProveedorId CumplidoProveedor
	FechaCreacion       time.Time
	FechaModificacion   time.Time
	Activo              bool
}

type BodyCumplidoSatisfaccion struct {
	NumeroContratoSuscrito int    `json:"NumeroContratoSuscrito"`
	VigenciaContrato       string `json:"VigenciaContrato"`
	TipoPagoId             string `json:"TipoPagoId"`
	PeiodoInicio           string `json:"PeriodoInicio"`
	PeriodoFin             string `json:"PeriodoFin"`
	TipoDocumentoCobroId   string `json:"TipoDocumentoCobroId"`
	NumeroCuentaFactura    string `json:"NumeroCuentaFactura"`
	ValorPagar             int    `json:"ValorPagar"`
	TipoCuenta             string `json:"TipoCuenta"`
	NumeroCuenta           string `json:"NumeroCuenta"`
	BancoId                int    `json:"BancoId"`
	CumplimientoContrato   string `json:"CumplimientoContrato"`
}

type BodyCumplidoRequest struct {
	CodigoAbreviacionEstadoCumplido string `json:"CodigoAbreviacionEstadoCumplido"`
	CumplidoProveedorID             int    `json:"CumplidoProveedorId"`
}

type BodySubirSoporteRequest struct {
	SolicitudPagoID int    `json:"SolicitudPagoID"`
	TipoDocumento   string `json:"TipoDocumento"`
	ItemID          int    `json:"ItemID"`
	Observaciones   string `json:"Observaciones"`
	NombreArchivo   string `json:"NombreArchivo"`
	Archivo         string `json:"Archivo"`
}

type AgregarComentarioSoporteRequest struct {
	SoporteId      string `json:"soporte_id"`
	CambioEstadoId string `json:"cambio_estado_id"`
	Comentario     string `json:"comentario"`
}
