package env

import "testing"

func TestEnvironment_URL(t *testing.T) {
	pwant := productionURL
	swant := sandboxURL
	ewant := "http://google.com"

	t.Run("Production URL", func(t *testing.T) {
		if got := new(Production).URL(); got != pwant {
			t.Errorf("Production.URL() = %v, want %v", got, pwant)
		}
	})

	t.Run("Sandbox URL", func(t *testing.T) {
		if got := new(Sandbox).URL(); got != swant {
			t.Errorf("Sandbox.URL() = %v, want %v", got, swant)
		}
	})

	t.Run("Endpoint URL", func(t *testing.T) {
		end := &Endpoint{url: "http://google.com"}
		if got := end.URL(); got != ewant {
			t.Errorf("Endpoint.URL() = %v, want %v", got, ewant)
		}
	})
}
