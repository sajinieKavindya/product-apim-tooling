package utils


// ------------------- Structs for YAML Config Files ----------------------------------

// For env_keys_all.yaml
// Not to be manually edited
type EnvKeysAll struct {
	Environments map[string]EnvKeys `yaml:"environments"`
}

// For env_endpoints_all.yaml
// To be manually edited by the user
type EnvEndpointsAll struct {
	Environments map[string]EnvEndpoints `yaml:"environments"`
}

type EnvKeys struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"` // to be encrypted (with the user's password) and stored
	Username string		`yaml:"username"`
}

type EnvEndpoints struct {
	APIManagerEndpoint   string `yaml:"api_manager_endpoint"`
	RegistrationEndpoint string `yaml:"registration_endpoint"`
	TokenEndpoint        string `yaml:"token_endpoint"`
}

// ---------------- End of Structs for YAML Config Files ---------------------------------

// variables
var envEndpointsAll EnvEndpointsAll
var envKeysAll EnvKeysAll

// Validates the configuration file
func (envEndpointsAll *EnvEndpointsAll) validate() {
	//
}



/**
Load the Environments Configuration file from the config.yaml file. If the file is not there
create a new config.yaml file and add default values
Validates the configuration, if it exists
*/
func LoadEnvConfig(envLocalConfig string) /* EnvEndpointsAll */ {
}

/*
// Returns a pointer to EnvEndpointsAll
func GetEnvEndpointsAll() *EnvEndpointsAll {
	if &envEndpointsAll == nil {
		HandleErrorAndExit("Env configuration is not available", nil)
	}
	return &envEndpointsAll
}

// Returns a pointer to EnvKeysAll
func GetEnvKeysAll() *EnvKeysAll {
	if &envKeysAll == nil {
		HandleErrorAndExit("EnvKeys configuration is not available", nil)
	}
	return &envKeysAll
}
*/


/*
env_keys_config.yaml (Programmatically edited)
===============
environments:
	dev:
		client_id: xxxxxxxxxx
		client_secret: xxxxxxxxxx
		username: xxxxxx

	staging:
		client_id: xxxxxxxxxx
		client_secret: xxxxxxxxxx
		username: xxxxxx
 */

/*
env_config.yaml (Manually edited)
===============
environments:
	dev:
		apim_endpoint: xxxxxxxxx
		registration_endpoint: xxxxxxxxxx
		token_endpoint: xxxxxxxxx

	staging:
		apim_endpoint: xxxxxxxxx
		registration_endpoint: xxxxxxxxxx
		token_endpoint: xxxxxxxxx
*/
