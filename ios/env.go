package ios

import (
	"encoding/json"
	"fmt"
)

const (
	prodURL = "https://sandbox.itunes.apple.com/verifyReceipt"
	sandURL = "https://buy.itunes.apple.com/verifyReceipt"
)

// Env interface provide ability to choose an environment for validation in-app purchases.
// This package comes with Apple environment endpoint by default.
// Implementing this interface will give possibility to send request to custom endpoint.
// In that case custom endpoint should take care about in-app purchases validation
// and returning the valid response.
type Env interface {
	// Endpoint returns endpoint URL of concrete environment.
	Endpoint() string
}

// AppleEnv represents enumeration of Apple environments for validation in-app purchases.
type AppleEnv int

const (
	// Production represents endpoint of Apple in-app purchases validation server production environment.
	Production AppleEnv = iota
	// Sandbox represents endpoint of Apple in-app purchases validation server sandbox environment.
	Sandbox
)

func (e AppleEnv) Endpoint() string {
	envs := map[AppleEnv]string{
		Production: prodURL,
		Sandbox:    sandURL,
	}
	return envs[e]
}

func (e AppleEnv) String() string {
	envs := map[AppleEnv]string{
		Production: "Production",
		Sandbox:    "Sandbox",
	}
	env, ok := envs[e]
	if !ok {
		return "Custom"
	}
	return env
}

func (e *AppleEnv) UnmarshalJSON(b []byte) error {
	var env string
	if err := json.Unmarshal(b, &env); err != nil {
		return err
	}

	switch env {
	case "0":
		*e = Production
		return nil
	case "1":
		*e = Sandbox
		return nil
	default:
		return fmt.Errorf("can't unmarshall json value: %s to AppleEnv type", env)
	}
}

func (e AppleEnv) MarshalJSON() ([]byte, error) {
	return []byte(`"` + e.String() + `"`), nil
}
