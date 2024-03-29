package client

import "strings"

const (
	DataPlanePathCreate = ControlPlanePathGet + "/dp-client-certificates"
	DataPlanePathGet    = DataPlanePathCreate + "/%s"
)

type DataPlaneCertificate struct {
	Id             string `json:"id,omitempty"`
	CreatedAt      int    `json:"created_at,omitempty"`
	UpdatedAt      int    `json:"updated_at,omitempty"`
	Cert           string `json:"cert"`
	ControlPlaneId string `json:"-"`
}

type DataPlaneCertificateResponse struct {
	Item DataPlaneCertificate `json:"item"`
}

func (s *DataPlaneCertificate) DataPlaneCertificateEncodeId() string {
	return s.ControlPlaneId + IdSeparator + s.Id
}

func DataPlaneCertificateDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}

type DataPlaneCertificateCollection struct {
	DataPlaneCertificates []DataPlaneCertificate `json:"data"`
}
