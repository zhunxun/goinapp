package ios

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Validator type represent http client for validation in-app purchases.
type Validator struct {
	client   *http.Client
	password string
}

// NewValidator return a new instance of Validator type.
func NewValidator(opts ...ValidatorOption) *Validator {
	validator := &Validator{
		password: "",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(validator)
	}

	return validator
}

// ValidatorOption represents optional function, which could be passed to NewValidator() func to change the
// default properties of returned Validator type.
type ValidatorOption func(*Validator)

// WithHTTPClient represents the optional function, which returns ValidatorOption function type.
// Receives the http.Client, which will be set to Validator client field.
func WithHTTPClient(c *http.Client) func(*Validator) {
	return func(v *Validator) {
		v.client = c
	}
}

// WithPassword represents the optional function, which returns ValidatorOption function type.
// Receives the string, which will be set to Validator password field.
func WithPassword(password string) func(*Validator) {
	return func(v *Validator) {
		v.password = password
	}
}

// Validate sends http POST with JSON body, which is represented by ValidationRequest struct to AppStore backend
// and parse the response with JSON body to ValidationResponse struct.
//
// The receipt must be a valid base64 encoded string from your StoreKit.
//
// The env must implement the Env interface.
// You can use AppleEnv type, which is represented by two constants: Production and Sandbox.
//
// You also can implement Env interface to send receipt to your custom endpoint. In that
// case the custom endpoint should take care about in-app purchases validation and returning the valid response.
func (v *Validator) Validate(ctx context.Context, receipt string, env Env) (*ValidationResponse, error) {
	payload := ValidationRequest{
		ReceiptData: receipt,
		Password:    v.password,
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(payload); err != nil {
		return nil, fmt.Errorf("body payload encoding error: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, env.Endpoint(), &body)
	if err != nil {
		return nil, fmt.Errorf("http request creation error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	res, err := v.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("http request failure: %v", err)
	}
	defer res.Body.Close()

	var response ValidationResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (v *Validator) ValidateAuto(ctx context.Context, receipt string) (*ValidationResponse, error) {
	resp, err := v.Validate(ctx, receipt, Production)
	if err != nil {
		return nil, fmt.Errorf("validation with auto env failed: %v", err)
	}
	if !resp.IsValid() && resp.StatusError() == ErrProductionOnSandbox {
		retryResp, retryErr := v.Validate(ctx, receipt, Sandbox)
		if retryErr != nil {
			return nil, fmt.Errorf("validation with auto env failed: %v", err)
		}
		return retryResp, nil
	}
	return resp, nil
}

// ValidationRequest type has the request properties.
// Submit this struct as JSON payload of an HTTP POST request to AppStore backend.
// In the test environment, use https://sandbox.itunes.apple.com/verifyReceipt as the url.
// In production, use https://buy.itunes.apple.com/verifyReceipt as the url.
// See Apple docs:
// https://developer.apple.com/library/archive/releasenotes/General/ValidateAppStoreReceipt/Chapters/ValidateRemotely.html
type ValidationRequest struct {
	// Base64 encoded json receipt.
	ReceiptData string `json:"receipt-data"`
	// Only used for receipts that contain auto-renewable subscriptions.
	Password string `json:"password,omitempty"`
	// Only used for iOS7 style app receipts that contain auto-renewable or non-renewing subscriptions.
	// If value is true, ValidationResponse includes only the latest renewal transaction for any subscriptions.
	ExcludeOldTransactions bool `json:"exclude-old-transactions,omitempty"`
}

// ValidationResponse type has the ValidationResponse properties.
// See Apple docs:
// https://developer.apple.com/library/archive/releasenotes/General/ValidateAppStoreReceipt/Chapters/ValidateRemotely.html
type ValidationResponse struct {
	// Status represent validation status codes from Apple validation endpoint.
	//
	// Either 0 if the receipt is valid, or one of the error codes. Look at StatusError() method below.
	// For iOS 6 style transaction receipts, the status code reflects the status of the specific transaction’s receipt.
	//
	// iOS 7 style app receipts, the status code is reflects the status of the app receipt as a whole.
	// For example, if you send a valid app receipt that contains an expired subscription, the response is 0
	// because the receipt as a whole is valid.
	Status int `json:"status"`
	// Environment represent receipt environment against which the validation has been performed.
	Environment AppleEnv `json:"environment,omitempty"`
	// A JSON representation of the receipt that was sent for verification
	Receipt Receipt `json:"receipt,omitempty"`
	// Base64 encoded string. Only returned for receipts containing auto-renewable subscriptions
	LatestReceipt string `json:"latest_receipt,omitempty"`
	// Only returned for receipts containing auto-renewable subscriptions
	LatestReceiptInfo InApps `json:"latest_receipt_info,omitempty"`
	// Only returned for iOS 6 style transaction receipts, for an auto-renewable subscription.
	// The JSON representation of the receipt for the expired subscription.
	LatestExpiredReceiptInfo InApps `json:"latest_expired_receipt_info,omitempty"`
	// A pending renewal may refer to a renewal that is scheduled in the future,
	// or a renewal that failed in the past for some reason.
	PendingRenewalInfo PendingRenewalInfos `json:"pending_renewal_info,omitempty"`
	// Retry validation for this receipt. Only applicable to status codes 21100-21199
	IsRetryable bool `json:"is-retryable,string,omitempty"`
}

// PendingRenewalInfos
type PendingRenewalInfos []PendingRenewalInfo

// PendingRenewalInfo represents a pending renewal, which may refer to a renewal that is scheduled in the future,
// or a renewal that failed in the past for some reason.
type PendingRenewalInfo struct {
	ProductID                      string `json:"product_id"`
	SubscriptionExpirationIntent   string `json:"expiration_intent"`
	SubscriptionAutoRenewProductID string `json:"auto_renew_product_id"`
	SubscriptionRetryFlag          string `json:"is_in_billing_retry_period"`
	SubscriptionAutoRenewStatus    string `json:"auto_renew_status"`
	SubscriptionPriceConsentStatus string `json:"price_consent_status"`
}

var (
	ErrMalformedJSON        = errors.New("the App Store could not read the JSON object you provided")
	ErrMalformedReceiptData = errors.New("data in the receipt-data property was malformed or missing")
	ErrNotAuthenticated     = errors.New("receipt could not be authenticated")
	ErrUnauthorizedReceipt  = errors.New("receipt could not be authorized")
	ErrIncorrectSecret      = errors.New("wrong secret")
	ErrServerNotAvailable   = errors.New("receipt server is not currently available")
	ErrSubscriptionExpired  = errors.New("receipt is valid but the subscription has expired")
	ErrSandboxOnProduction  = errors.New("wrong environment: sandbox receipt was sent to production environment")
	ErrProductionOnSandbox  = errors.New("wrong environment: production receipt was sent to sandbox environment")
	ErrInternalDataAccess   = errors.New("internal data access error")
	ErrUnknown              = errors.New("an unknown error occurred")
)

// IsRenewable returns true if receipt containing auto-renewable subscriptions.
func (r *ValidationResponse) IsRenewable() bool {
	if r.LatestReceipt != "" && len(r.LatestReceiptInfo) > 0 {
		return true
	}
	return false
}

// IsValid returns true if validation was successful.
func (r *ValidationResponse) IsValid() bool {
	switch r.Status {
	case 0:
		return true
	default:
		return false
	}
}

// StatusError returns error based on Status property of ValidationResponse.
func (r *ValidationResponse) StatusError() error {
	errs := map[int]error{
		0:     nil,
		21000: ErrMalformedJSON,
		21002: ErrMalformedReceiptData,
		21003: ErrNotAuthenticated,
		21004: ErrIncorrectSecret,
		21005: ErrServerNotAvailable,
		21006: ErrSubscriptionExpired,
		21007: ErrSandboxOnProduction,
		21008: ErrProductionOnSandbox,
		21010: ErrUnauthorizedReceipt,
	}

	err, ok := errs[r.Status]
	if ok {
		return err
	}

	if r.Status >= 21100 && r.Status <= 21199 {
		return ErrInternalDataAccess
	}
	return ErrUnknown
}
