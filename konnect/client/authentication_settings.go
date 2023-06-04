package client

const (
	AuthenticationSettingsPath = "authentication-settings"
)

type AuthenticationSettings struct {
	Id                    string `json:"-"`
	BasicAuthEnabled      bool   `json:"basic_auth_enabled"`
	OIDCAuthEnabled       bool   `json:"oidc_auth_enabled"`
	IDPMappingEnabled     bool   `json:"idp_mapping_enabled"`
	KonnectMappingEnabled bool   `json:"konnect_mapping_enabled"`
}