package models

type GirosTercero struct {
	Giros struct {
		Tercero []Tercero `json:"tercero"`
	} `json:"giros"`
}

type Tercero struct {
	TipoId                   string `json:"tipo_id"`
	Identificacion           string `json:"identificacion"`
	VigenciaPago             string `json:"vigencia_pago"`
	UnidadEjecutora          string `json:"unidad_ejecutora"`
	VigenciaPresupuestal     string `json:"vigencia_presupuestal"`
	Disponibilidad           string `json:"disponibilidad"`
	Registro                 string `json:"registro"`
	OrdenesPagoGiradas       string `json:"ordenes_pago_giradas"`
	ValorBrutoGirado         string `json:"valor_bruto_girado"`
	ValorBrutoGiradoNumerico int
}
