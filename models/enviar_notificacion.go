package models

type NotificacionEmail struct {
	Source              string            `json:"Source"`
	Template            string            `json:"Template"`
	Destinations        []DestinationItem `json:"Destinations"`
	DefaultTemplateData TemplateData      `json:"DefaultTemplateData"`
}

type DestinationItem struct {
	Destination             Destination  `json:"Destination"`
	ReplacementTemplateData TemplateData `json:"ReplacementTemplateData"`
	Attachments             []string     `json:"Attachments"`
}

type Destination struct {
	BccAddresses []string `json:"BccAddresses"`
	CcAddresses  []string `json:"CcAddresses"`
	ToAddresses  []string `json:"ToAddresses"`
}

type TemplateData struct {
	EstadoCumplido                  string `json:"estado_cumplido"`
	NombreEmisorCumplido            string `json:"nombre_emisor_cumplido"`
	NombreProveedor                 string `json:"nombre_proveedor"`
	RevisionCumplidosProveedoresUrl string `json:"revision_cumplidos_proveedores_url"`
}
