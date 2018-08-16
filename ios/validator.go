// Package ios contains client for validating in-app purchase receipt via AppStore backend.
// It also contain a bunch of usable methods for work with in-app purchase data.
package ios

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/heartwilltell/goinapp/ios/env"
	"net/http"
	"time"
)

// Validator type represent http client for validation in-app purchases
type Validator struct {
	client   *http.Client
	password string
}

// NewValidator return instance of Validator struct that contain http.Client with timeout of 10 seconds.
func NewValidator(password string) *Validator {
	return &Validator{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		password: password,
	}
}

// NewValidatorWithClient return instance of Validator struct that contain specified http.Client.
func NewValidatorWithClient(client *http.Client, password string) *Validator {
	return &Validator{
		client:   client,
		password: password,
	}
}

// Validate send http POST request with receipt to AppStore backend and parse the response.
// receipt must be a base64 encoded string from your StoreKit.
// environment must be a string value of: "productionEnv", "sandboxEnv" or you can pass any valid url,
// to send request to your proxy for example.
func (v *Validator) Validate(ctx context.Context, receipt string, environment env.Environment) (*ValidationResponse, error) {
	payload := ValidationRequest{
		ReceiptData: receipt,
		Password:    v.password,
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(payload); err != nil {
		return nil, fmt.Errorf("body payload encoding error: %v", err)
	}

	req, err := http.NewRequest("POST", environment.URL(), &body)
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

// ValidationRequest type has the request properties.
// Submit this struct as JSON payload of an HTTP POST request to AppStore backend.
// In the test environment, use https://sandbox.itunes.apple.com/verifyReceipt as the url.
// In production, use https://buy.itunes.apple.com/verifyReceipt as the url.
// See Apple docs:
// https://developer.apple.com/library/archive/releasenotes/General/ValidateAppStoreReceipt/Chapters/ValidateRemotely.html
type ValidationRequest struct {
	// Base64 encoded json receipt
	ReceiptData string `json:"receipt-data"`
	// Only used for receipts that contain auto-renewable subscriptions
	Password string `json:"password,omitempty"`
	// Only used for iOS7 style app receipts that contain auto-renewable or non-renewing subscriptions.
	// If value is true, response includes only the latest renewal transaction for any subscriptions.
	ExcludeOldTransactions bool `json:"exclude-old-transactions,omitempty"`
}

// ValidationResponse type has the response properties.
// See Apple docs:
// https://developer.apple.com/library/archive/releasenotes/General/ValidateAppStoreReceipt/Chapters/ValidateRemotely.html
type ValidationResponse struct {
	// Represent apple validation status codes
	Status int `json:"status"`
	// Represent receipt environment
	Environment string `json:"environment"`
	// A JSON representation of the receipt that was sent for verification
	Receipt Receipt `json:"receipt"`
	// Only returned for receipts containing auto-renewable subscriptions
	// Base64 encoded string
	LatestReceipt string `json:"latest_receipt,omitempty"`
	// Only returned for receipts containing auto-renewable subscriptions
	LatestReceiptInfo InApps `json:"latest_receipt_info,omitempty"`
	// Only returned for iOS 6 style transaction receipts, for an auto-renewable subscription.
	// The JSON representation of the receipt for the expired subscription.
	LatestExpiredReceiptInfo InApps `json:"latest_expired_receipt_info,omitempty"`
	// A pending renewal may refer to a renewal that is scheduled in the future or a renewal that failed in the past for some reason.
	PendingRenewalInfo PendingRenewalInfos `json:"pending_renewal_info,omitempty"`
	// Retry validation for this receipt. Only applicable to status codes 21100-21199
	IsRetryable bool `json:"is-retryable,string,omitempty"`
}

type PendingRenewalInfos []PendingRenewalInfo

// A pending renewal may refer to a renewal that is scheduled in the future or a renewal that failed in the past for some reason.
type PendingRenewalInfo struct {
	ProductID                      string `json:"product_id"`
	SubscriptionExpirationIntent   string `json:"expiration_intent"`
	SubscriptionAutoRenewProductID string `json:"auto_renew_product_id"`
	SubscriptionRetryFlag          string `json:"is_in_billing_retry_period"`
	SubscriptionAutoRenewStatus    string `json:"auto_renew_status"`
	SubscriptionPriceConsentStatus string `json:"price_consent_status"`
}

// ValidationStatus method return an error is Status field of ValidationResponse not equal to 0.
func (v *ValidationResponse) ValidationStatus() error {
	var message string
	switch v.Status {
	case 0:
		return nil
	case 21000:
		message = "The App Store could not read the JSON object you provided."
	case 21002:
		message = "The data in the receipt-data property was malformed or missing."
	case 21003:
		message = "The receipt could not be authenticated."
	case 21004:
		message = "The shared secret you provided does not match the shared secret on file for your account."
	case 21005:
		message = "The receipt server is not currently available."
	case 21007:
		message = "This receipt is from the test environment, but it was sent to the production environment for verification. Send it to the test environment instead."
	case 21008:
		message = "This receipt is from the production environment, but it was sent to the test environment for verification. Send it to the production environment instead."
	case 21010:
		message = "This receipt could not be authorized. Treat this the same as if a purchase was never made."
	default:
		if v.Status >= 21100 && v.Status <= 21199 {
			message = "Internal data access error."
		} else {
			message = "An unknown error occurred."
		}
	}
	return errors.New(message)
}

// Renewable return true if receipt containing auto-renewable subscriptions
func (v *ValidationResponse) Renewable() bool {
	if v.LatestReceipt != "" && len(v.LatestReceiptInfo) > 0 {
		return true
	}
	return false
}
