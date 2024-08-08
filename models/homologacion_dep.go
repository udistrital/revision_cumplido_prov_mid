package models

type HomologacionDepOikos struct {
	Dependencias struct {
		Dependencia []struct {
			IDMaster string      `json:"id_master"`
			IDGedep  interface{} `json:"id_gedep"`
			IDArgo   string      `json:"id_argo"`
			IDAcad   interface{} `json:"id_acad"`
		} `json:"dependencia"`
	} `json:"dependencias"`
}
