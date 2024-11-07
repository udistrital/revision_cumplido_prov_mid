package models

type InformacionPersonaNatural struct {
	PrimerApellido              string
	SegundoApellido             string
	PrimerNombre                string
	SegundoNombre               string
	IdCiudadExpedicionDocumento int
	TipoDocumento               struct {
		Abreviatura    string
		ValorParametro string
	}
}
