package config

type Config struct {
	Gateway   *Gateway   `json:"gateway"`   // gateway settings
	Directory *Directory `json:"directory"` // target directory settings
	Auth0     *Auth0     `json:"auth0"`     // source auth0 settings
	Loader    *Loader    `json:"loader"`    // ds-load* settings
}

type Gateway struct {
	Port int    `json:"port"` // gateway listen port
	Path string `json:"path"` // gateway URL path
}

type Directory struct {
	Host     string `json:"host"`      // directory gRPC host address
	APIKey   string `json:"api_key"`   // directory read-write API key
	TenantID string `json:"tenant_id"` // directory tenant ID
	Insecure bool   `json:"insecure"`  // skip TLS validation
}

type Auth0 struct {
	Domain       string `json:"domain"`        // Auth0 domain
	ClientID     string `json:"client_id"`     // Auth0 client ID
	ClientSecret string `json:"client_secret"` // Auth0 client secret
	Connection   string `json:"connection"`    // Auth0 connection name
	UserPID      string `json:"user_pid"`      // Auth0 user ID of user to export
	UserEmail    string `json:"user_email"`    // Auth0 email name of user to export
	InclRoles    bool   `json:"incl_roles"`    // Auth0 include roles in export
}

type Loader struct {
	BinPath  string `json:"bin_path"` // ds-load absolute path to binaries (ds-load & ds-load-auth0)
	Template string `json:"template"` // ds-load absolute path to template file
}
