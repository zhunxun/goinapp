// TODO package docs
package env

const (
	sandboxURL    = "https://sandbox.itunes.apple.com/verifyReceipt"
	productionURL = "https://buy.itunes.apple.com/verifyReceipt"
)

// Environment interface provide ability to chose one of the URL for validation in-app purchases.
type Environment interface {
	URL() string
}

// production type represent production environment URL option.
type production string

// URL implement Environment interface for production type.
func (p production) URL() string {
	return productionURL
}

// Production set validation endpoint to production URL by returning production type that implements Environment interface.
func Production() *production {
	return new(production)
}

// sandbox type represent sandbox environment URL option.
type sandbox string

// URL implement Environment interface for sandbox type.
func (s sandbox) URL() string {
	return sandboxURL
}

// Sandbox set validation endpoint to production URL by returning sandbox type that implements Environment interface.
func Sandbox() *sandbox {
	return new(sandbox)
}
