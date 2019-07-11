package ios

import (
	"testing"
)

func TestEndpoint_Endpoint(t *testing.T) {
	type args struct {
		want string
		env  Env
	}

	tests := map[string]args{
		"Production": {want: prodURL, env: Production},
		"Sandbox":    {want: sandURL, env: Sandbox},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.env.Endpoint(); got != tt.want {
				t.Errorf("Env.Endpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppleEnv_String(t *testing.T) {
	type args struct {
		env AppleEnv
	}
	type test struct {
		args args
		want string
	}
	tests := map[string]test{
		"Production": {args{env: Production}, "Production"},
		"Sandbox":    {args{env: Sandbox}, "Sandbox"},
		"Custom":     {args{env: AppleEnv(3)}, "Custom"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tc.args.env.String(); got != tc.want {
				t.Errorf("AppleEnv.String() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestAppleEnv_UnmarshalJSON(t *testing.T) {
	type args struct {
		env []byte
	}
	type test struct {
		args args
		want AppleEnv
	}
	tests := map[string]test{
		"Production": {args{env: []byte(`"0"`)}, Production},
		"Sandbox":    {args{env: []byte(`"1"`)}, Sandbox},
	}

	var env AppleEnv

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if err := env.UnmarshalJSON(tc.args.env); err != nil {
				t.Errorf("%v", err)
			}
			if env != tc.want {
				t.Errorf("AppleEnv.UnmarshalJSON() shoult be %v but got %v", tc.want, env)
			}
		})
	}

	t.Run("Error", func(t *testing.T) {
		if err := env.UnmarshalJSON([]byte(`"3"`)); err == nil {
			t.Errorf("AppleEnv.UnmarshalJSON() should return error it this case")
		}
	})
}
