package ios

import (
	"math/rand"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestNewValidator(t *testing.T) {
	type args struct {
		opts []ValidatorOption
	}
	type test struct {
		args args
		want *Validator
	}

	tests := map[string]test{
		"Default": {
			args{[]ValidatorOption{}},
			&Validator{client: &http.Client{Timeout: 10 * time.Second}, password: ""},
		},
		"WithHTTPClient": {
			args{[]ValidatorOption{WithHTTPClient(&http.Client{Timeout: 20 * time.Second})}},
			&Validator{client: &http.Client{Timeout: 20 * time.Second}, password: ""},
		},
		"WithPassword": {
			args{[]ValidatorOption{WithPassword("pass")}},
			&Validator{client: &http.Client{Timeout: 10 * time.Second}, password: "pass"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := NewValidator(tc.args.opts...); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("NewValidator() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestValidationResponse_StatusError(t *testing.T) {
	type args struct {
		status int
	}
	type test struct {
		args args
		want error
	}

	tests := map[string]test{
		"OK":                        {args{status: 0}, nil},
		"ErrMalformedJSON":          {args{status: 21000}, ErrMalformedJSON},
		"ErrMalformedReceiptData":   {args{status: 21002}, ErrMalformedReceiptData},
		"ErrNotAuthenticated":       {args{status: 21003}, ErrNotAuthenticated},
		"ErrIncorrectSecret":        {args{status: 21004}, ErrIncorrectSecret},
		"ErrServerNotAvailable":     {args{status: 21005}, ErrServerNotAvailable},
		"ErrSubscriptionExpired":    {args{status: 21006}, ErrSubscriptionExpired},
		"ErrSandboxOnProduction":    {args{status: 21007}, ErrSandboxOnProduction},
		"ErrProductionOnSandbox":    {args{status: 21008}, ErrProductionOnSandbox},
		"ErrUnauthorizedReceipt":    {args{status: 21010}, ErrUnauthorizedReceipt},
		"ErrInternalDataAccess":     {args{status: 21100}, ErrInternalDataAccess},
		"ErrInternalDataAccessRand": {args{status: randStatus(21100, 21199)}, ErrInternalDataAccess},
		"ErrInternalDataAccessMid":  {args{status: 21133}, ErrInternalDataAccess},
		"ErrInternalDataAccessEdge": {args{status: 21199}, ErrInternalDataAccess},
		"ErrUnknown":                {args{status: 21200}, ErrUnknown},
		"ErrUnknownRand":            {args{status: randStatus(21200, 50000)}, ErrUnknown},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := &ValidationResponse{
				Status: tt.args.status,
			}
			if err := r.StatusError(); err != tt.want {
				t.Errorf("ValidationResponse.StatusError() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func randStatus(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
