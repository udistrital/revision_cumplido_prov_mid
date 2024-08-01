package models

import "time"

type SoportePago struct {
	Id                  int
	DocumentoId         int
	CumplidoProveedorId *CumplidoProveedor
	FechaCreacion       time.Time
	FechaModificacion   time.Time
	Activo              bool
}

type Documento struct {
	Id                int            `json:"Id"`
	Nombre            string         `json:"Nombre"`
	Descripcion       string         `json:"Descripcion"`
	Enlace            string         `json:"Enlace"`
	TipoDocumento     *TipoDocumento `json:"TipoDocumento"`
	Metadatos         string         `json:"Metadatos"`
	Activo            bool           `json:"Activo"`
	FechaCreacion     time.Time      `json:"FechaCreacion"`
	FechaModificacion time.Time      `json:"FechaModificacion"`
}

// Estructura anidada para TipoDocumento
type TipoDocumento struct {
	Id                   int                   `json:"Id"`
	Nombre               string                `json:"Nombre"`
	Descripcion          string                `json:"Descripcion"`
	CodigoAbreviacion    string                `json:"CodigoAbreviacion"`
	Activo               bool                  `json:"Activo"`
	NumeroOrden          int                   `json:"NumeroOrden"`
	Tamano               int                   `json:"Tamano"`
	Extension            string                `json:"Extension"`
	Workspace            string                `json:"Workspace"`
	TipoDocumentoNuxeo   string                `json:"TipoDocumentoNuxeo"`
	FechaCreacion        time.Time             `json:"FechaCreacion"`
	FechaModificacion    time.Time             `json:"FechaModificacion"`
	DominioTipoDocumento *DominioTipoDocumento `json:"DominioTipoDocumento"`
}

// Estructura anidada para DominioTipoDocumento
type DominioTipoDocumento struct {
	Id                int       `json:"Id"`
	Nombre            string    `json:"Nombre"`
	Descripcion       string    `json:"Descripcion"`
	CodigoAbreviacion string    `json:"CodigoAbreviacion"`
	Activo            bool      `json:"Activo"`
	NumeroOrden       int       `json:"NumeroOrden"`
	FechaCreacion     time.Time `json:"FechaCreacion"`
	FechaModificacion time.Time `json:"FechaModificacion"`
}
