package models

import "time"

type DocumentosSoporte struct {
	Documento Documento
	Archivo   FileGestorDocumental
}

type DocumentosSoporteSimplificado struct {
	SoporteCumplidoId int
	Documento         DocumentoSimplificado
	Archivo           FileGestorDocumental
}

type DocumentoSimplificado struct {
	Id            int
	Nombre        string
	TipoDocumento string
	Descripcion   string
	Observaciones string
	FechaCreacion string
}

type Documento struct {
	Id                int
	Nombre            string
	Descripcion       string
	Enlace            string
	TipoDocumento     *TipoDocumento
	Metadatos         string
	Activo            bool
	FechaCreacion     string
	FechaModificacion string
}

type TipoDocumento struct {
	Id                   int
	Nombre               string
	Descripcion          string
	CodigoAbreviacion    string
	Activo               bool
	NumeroOrden          float64
	Tamano               float64
	Extension            string
	Workspace            string
	TipoDocumentoNuxeo   string
	FechaCreacion        string
	FechaModificacion    string
	DominioTipoDocumento *DominioTipoDocumento
}

type DominioTipoDocumento struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
	NumeroOrden       float64
	FechaCreacion     string
	FechaModificacion string
}

type FileGestorDocumental struct {
	// UIDUID          interface{} `json:"uid:uid"`
	// UIDMajorVersion int         `json:"uid:major_version"`
	// UIDMinorVersion int         `json:"uid:minor_version"`
	// ThumbThumbnail  struct {
	// 	Name            string      `json:"name"`
	// 	MimeType        string      `json:"mime-type"`
	// 	Encoding        interface{} `json:"encoding"`
	// 	DigestAlgorithm string      `json:"digestAlgorithm"`
	// 	Digest          string      `json:"digest"`
	// 	Length          string      `json:"length"`
	// 	Data            string      `json:"data"`
	// } `json:"thumb:thumbnail"`
	// FileContent struct {
	// 	Name            string      `json:"name"`
	// 	MimeType        string      `json:"mime-type"`
	// 	Encoding        interface{} `json:"encoding"`
	// 	DigestAlgorithm string      `json:"digestAlgorithm"`
	// 	Digest          string      `json:"digest"`
	// 	Length          string      `json:"length"`
	// 	Data            string      `json:"data"`
	// } `json:"file:content"`
	// CommonIconExpanded              interface{}   `json:"common:icon-expanded"`
	// CommonIcon                      string        `json:"common:icon"`
	// FilesFiles                      []interface{} `json:"files:files"`
	// DcDescription                   interface{}   `json:"dc:description"`
	// DcLanguage                      interface{}   `json:"dc:language"`
	// DcCoverage                      interface{}   `json:"dc:coverage"`
	// DcValid                         interface{}   `json:"dc:valid"`
	// DcCreator                       string        `json:"dc:creator"`
	// DcModified                      time.Time     `json:"dc:modified"`
	// DcLastContributor               string        `json:"dc:lastContributor"`
	// DcRights                        interface{}   `json:"dc:rights"`
	// DcExpired                       interface{}   `json:"dc:expired"`
	// DcFormat                        interface{}   `json:"dc:format"`
	// DcCreated                       time.Time     `json:"dc:created"`
	// DcTitle                         string        `json:"dc:title"`
	// DcIssued                        interface{}   `json:"dc:issued"`
	// DcNature                        interface{}   `json:"dc:nature"`
	// DcSubjects                      []interface{} `json:"dc:subjects"`
	// DcContributors                  []string      `json:"dc:contributors"`
	// DcSource                        interface{}   `json:"dc:source"`
	// DcPublisher                     interface{}   `json:"dc:publisher"`
	// RelatedtextRelatedtextresources []interface{} `json:"relatedtext:relatedtextresources"`
	// NxtagTags                       []interface{} `json:"nxtag:tags"`
	File string `json:"File"`
}

type ItemInformeTipoContrato struct {
	Id                int
	ItemInformeId     *ItemInforme
	TipoContratoId    int
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
}

type ItemInforme struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	NumeroOrden       float64
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
}

type DocumentoCumplido struct {
	IdTipoDocumento                int
	CodigoAbreviacionTipoDocumento string
	Nombre                         string
}

type DocumentosComprimido struct {
	Nombre string `json:"nombre"`
	File   string `json:"file"`
}
