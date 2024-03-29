package client

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
}

type CertificateCollection struct {
	Certificates []Certificate `json:"data"`
}
