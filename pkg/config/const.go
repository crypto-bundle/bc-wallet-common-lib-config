package config

const (
	BaseConfigPrefix = "APP"
)

// The struct tags for using in your config structures. See the
// library examples for details.
const (
	tagSecret    = "secret"
	tagEnvconfig = "envconfig"
	tagVaultKey  = "vault_key"
	tagRequired  = "required"
	tagIgnored   = "ignored"
	tagDefault   = "default"
)
