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
	CargoResponsable     string            `json:"CargoReponsable"`
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

type BodyInformeSeguimiento struct {
	NumeroContratoSuscrito string `json:"numero_contrato_suscrito"`
	VigenciaContrato       string `json:"vigencia_contrato"`
	TipoPago               string `json:"tipo_pago"`
	PeiodoInicio           string `json:"periodo_inicio"`
	PeriodoFin             string `json:"periodo_fin"`
	TipoSoportePagar       string `json:"tipo_soporte_pagar"`
	NumeroCuentaFactura    string `json:"numero_cuenta_factura"`
	ValorPagar             string `json:"valor_pagar"`
	TipoCuenta             string `json:"tipo_cuenta"`
	NumeroCuenta           string `json:"numero_cuenta"`
	Banco                  string `json:"banco"`
}
