package client

import "strings"

const (
	CertificatePathCreate = ControlPlanePathGet + "/core-entities/certificates"
	CertificatePathGet    = CertificatePathCreate + "/%s"
)

type Certificate struct {
	Id                   string   `json:"id,omitempty"`
	Certificate          string   `json:"cert"`
	Key                  string   `json:"key"`
	AlternateCertificate string   `json:"cert_alt,omitempty"`
	AlternateKey         string   `json:"key_alt,omitempty"`
	Tags                 []string `json:"tags,omitempty"`
	ControlPlaneId       string   `json:"-"`
}

func (s *Certificate) CertificateEncodeId() string {
	return s.ControlPlaneId + IdSeparator + s.Id
}

func CertificateDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}

type CertificateCollection struct {
	Certificates []Certificate `json:"data"`
}
