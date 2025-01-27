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
	NumeroContratoSuscrito int       `json:"NumeroContratoSuscrito"`
	VigenciaContrato       string    `json:"VigenciaContrato"`
	TipoPago               string    `json:"TipoPago"`
	PeriodoInicio          time.Time `json:"PeriodoInicio"`
	PeriodoFin             time.Time `json:"PeriodoFin"`
	TipoFactura            string    `json:"TipoFactura"`
	NumeroCuentaFactura    string    `json:"NumeroCuentaFactura"`
	ValorPagar             int       `json:"ValorPagar"`
	TipoCuenta             string    `json:"TipoCuenta"`
	NumeroCuenta           string    `json:"NumeroCuenta"`
	Banco                  string    `json:"Banco"`
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

type BodyHistoricoRequest struct {
	Anios          []int    `json:"Anios"`
	Meses          []int    `json:"Meses"`
	Vigencias      []int    `json:"Vigencias"`
	Proveedores    []int    `json:"Proveedores"`
	Estados        []string `json:"Estados"`
	Dependencias   []string `json:"Dependencias"`
	Contratos      []string `json:"Contratos"`
	TiposContratos []int    `json:"TiposContratos"`
}
