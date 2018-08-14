package env

const (
	// sandboxURL is endpoint URL for sandbox environment.
	sandboxURL = "https://sandbox.itunes.apple.com/verifyReceipt"

	// productionURL is endpoint URL for production environment.
	productionURL = "https://buy.itunes.apple.com/verifyReceipt"
)

// environment interface encapsulate different endpoint realization
type environment interface {
	URL() string
}

// productionEnv type represent productionEnv environment URL option
type productionEnv string

// URL implement environment interface for productionEnv type
func (p productionEnv) URL() string {
	return productionURL
}

func Production() *productionEnv {
	return new(productionEnv)
}

// sandboxEnv type represent sandboxEnv environment URL option
type sandboxEnv string

// URL implement environment interface
func (s sandboxEnv) URL() string {
	return sandboxURL
}

func Sandbox() *sandboxEnv {
	return new(sandboxEnv)
}

// endpointEnv type represent custom URL option for sending request to any provided endpoint
type endpointEnv struct {
	url string
}

// URL implement environment interface
func (e endpointEnv) URL() string {
	return e.url
}

func Endpoint(url string) *endpointEnv {
	return &endpointEnv{url: url}
}
