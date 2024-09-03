package models

type OrdenadorContratoProveedor struct {
	Contratos struct {
		Ordenador []struct {
			RolOrdenador    string `json:"rol_ordenador"`
			Documento       string `json:"documento"`
			NombreOrdenador string `json:"nombre_ordenador"`
		} `json:"ordenador"`
	} `json:"contratos"`
}
