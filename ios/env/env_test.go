package env

import "testing"

func TestEnvironment_URL(t *testing.T) {
	pwant := productionURL
	swant := sandboxURL

	t.Run("production URL", func(t *testing.T) {
		if got := Production().URL(); got != pwant {
			t.Errorf("production.URL() = %v, want %v", got, pwant)
		}
	})

	t.Run("sandbox URL", func(t *testing.T) {
		if got := Sandbox().URL(); got != swant {
			t.Errorf("sandbox.URL() = %v, want %v", got, swant)
		}
	})
}
