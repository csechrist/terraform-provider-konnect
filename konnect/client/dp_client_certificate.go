package client

const (
	DataPlanePathCreate = ControlPlanePathGet + "/dp-client-certificates"
	DataPlanePathGet    = DataPlanePathCreate + "/%s"
)

type DataPlaneCertificate struct {
	Id        string `json:"id,omitempty"`
	CreatedAt int    `json:"created_at,omitempty"`
	UpdatedAt int    `json:"updated_at,omitempty"`
	Cert      string `json:"cert"`
}

type DataPlaneCertificateResponse struct {
	Item DataPlaneCertificate `json:"item"`
}

type DataPlaneCertificateCollection struct {
	DataPlaneCertificates []DataPlaneCertificate `json:"data"`
}
