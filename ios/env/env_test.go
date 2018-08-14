package env

import "testing"

func TestEnvironment_URL(t *testing.T) {
	pwant := productionURL
	swant := sandboxURL
	ewant := "http://google.com"

	t.Run("productionEnv URL", func(t *testing.T) {
		if got := new(productionEnv).URL(); got != pwant {
			t.Errorf("productionEnv.URL() = %v, want %v", got, pwant)
		}
	})

	t.Run("sandboxEnv URL", func(t *testing.T) {
		if got := new(sandboxEnv).URL(); got != swant {
			t.Errorf("sandboxEnv.URL() = %v, want %v", got, swant)
		}
	})

	t.Run("endpointEnv URL", func(t *testing.T) {
		end := &endpointEnv{url: "http://google.com"}
		if got := end.URL(); got != ewant {
			t.Errorf("endpointEnv.URL() = %v, want %v", got, ewant)
		}
	})
}
