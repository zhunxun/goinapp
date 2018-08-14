package env

const (
	// sandboxURL is endpoint URL for sandbox environment.
	sandboxURL = "https://sandbox.itunes.apple.com/verifyReceipt"

	// productionURL is endpoint URL for production environment.
	productionURL = "https://buy.itunes.apple.com/verifyReceipt"
)

// Environment interface encapsulate different endpoint realization
type Environment interface {
	URL() string
}

// Production type represent Production environment URL option
type Production string

// URL implement Environment interface
func (p Production) URL() string {
	return productionURL
}

// Sandbox type represent Sandbox environment URL option
type Sandbox string

// URL implement Environment interface
func (s Sandbox) URL() string {
	return sandboxURL
}

// Endpoint type represent custom URL option for sending request to any provided endpoint
type Endpoint struct {
	url string
}

// URL implement Environment interface
func (e Endpoint) URL() string {
	return e.url
}
