// Package ios contains client for validating in-app purchase receipt via AppStore backend.
// It also contain a bunch of usable methods for work with in-app purchase data.
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

const (
	// SandboxURL is the endpoint for sandbox environment.
	SandboxURL = "https://sandbox.itunes.apple.com/verifyReceipt"
	// ProductionURL is the endpoint for production environment.
	ProductionURL = "https://buy.itunes.apple.com/verifyReceipt"
)

// Validator type represent
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
// environment must be a string value of: "Production", "Sandbox" or you can pass any valid URL,
// to send request to your proxy for example.
func (v *Validator) Validate(ctx context.Context, receipt string, environment string) (*ValidationResponse, error) {
	payload := ValidationRequest{
		ReceiptData: receipt,
		Password:    v.password,
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(payload); err != nil {
		return nil, fmt.Errorf("body payload encoding error: %v", err)
	}

	var url string
	switch environment {
	case "Production":
		url = ProductionURL
	case "Sandbox":
		url = SandboxURL
	default:
		url = environment
	}

	req, err := http.NewRequest("POST", url, &body)
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

// CheckStatus method return an error is Status field of ValidationResponse not equal to 0.
func (c *ValidationResponse) CheckStatus() error {
	var message string
	switch c.Status {
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
		if c.Status >= 21100 && c.Status <= 21199 {
			message = "Internal data access error."
		} else {
			message = "An unknown error occurred."
		}
	}
	return errors.New(message)
}

// ValidationRequest type has the request properties.
// Submit this struct as JSON payload of an HTTP POST request to AppStore backend.
// In the test environment, use https://sandbox.itunes.apple.com/verifyReceipt as the URL.
// In production, use https://buy.itunes.apple.com/verifyReceipt as the URL.
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
	LatestReceiptInfo []InApp `json:"latest_receipt_info,omitempty"`
	// Only returned for iOS 6 style transaction receipts, for an auto-renewable subscription.
	// The JSON representation of the receipt for the expired subscription.
	LatestExpiredReceiptInfo []InApp `json:"latest_expired_receipt_info,omitempty"`
	// Retry validation for this receipt. Only applicable to status codes 21100-21199
	IsRetryable string `json:"is-retryable,omitempty"`
}

//
func (c ValidationResponse) TrialOriginalTransactionId() string {
	var trialOTID string
	if c.LatestReceipt != "" && len(c.LatestReceiptInfo) > 0 {
		for _, receipt := range c.LatestReceiptInfo {
			if receipt.IsTrialPeriod == "false" {
				continue
			}
			if receipt.IsTrialPeriod == "true" {
				trialOTID = receipt.OriginalTransactionID
				break
			}
		}
	}
	return trialOTID
}

func (c ValidationResponse) PaidOriginalTransactionId() string {
	var paidOTID string
	if c.LatestReceipt != "" && len(c.LatestReceiptInfo) > 0 {
		for _, receipt := range c.LatestReceiptInfo {
			if receipt.IsTrialPeriod == "true" {
				continue
			}
			if receipt.IsTrialPeriod == "false" {
				paidOTID = receipt.OriginalTransactionID
				break
			}
		}
	}
	return paidOTID
}

// Receipt type has the receipt property
type Receipt struct {
	AdamID    int `json:"adam_id"`
	AppItemID int `json:"app_item_id"`
	// This corresponds to the value of CFBundleVersion (in iOS) or CFBundleShortVersionString (in macOS) in the Info.plist
	ApplicationVersion string `json:"application_version"`
	// This corresponds to the value of CFBundleIdentifier in the Info.plist file
	// Use this value to validate if the receipt was indeed generated for your app
	BundleID   string `json:"bundle_id"`
	DownloadID int    `json:"download_id"`
	// In the JSON file, the value of this key is an array containing all in-app purchase receipts based on the in-app purchase transactions present in the input base-64 receipt-data
	// For receipts containing auto-renewable subscriptions, check the value of the latest_receipt_info key to get the status of the most recent renewal
	InApp []InApp `json:"in_app"`
	//This corresponds to the value of CFBundleVersion (in iOS) or CFBundleShortVersionString (in macOS) in the Info.plist file when the purchase was originally made
	//In the sandbox environment, the value of this field is always “1.0”
	OriginalApplicationVersion string `json:"original_application_version"`
	OriginalPurchaseDate
	// When validating a receipt, use this date to validate the receipt’s signature
	ReceiptCreationDate
	ReceiptExpirationDate
	ReceiptType string `json:"receipt_type"`
	ReceiptRequestDate
	VersionExternalIdentifier int `json:"version_external_identifier"`
}

type InApp struct {
	// The number of items purchased.
	Quantity string `json:"quantity"`
	// The product identifier of the item that was purchased.
	ProductID string `json:"product_id"`
	// The transaction identifier of the item that was purchased.
	// For a transaction that restores a previous transaction, this value is different from the transaction identifier of the original purchase transaction.
	// In an auto-renewable subscription receipt, a new value for the transaction identifier is generated every time the subscription automatically renews or is restored on a new device.
	TransactionID string `json:"transaction_id"`
	// For a transaction that restores a previous transaction, the transaction identifier of the original transaction.
	// Otherwise, identical to the transaction identifier.
	// This value is the same for all receipts that have been generated for a specific subscription.
	// This value is useful for relating together multiple iOS 6 style transaction receipts for the same individual customer’s subscription.
	OriginalTransactionID string `json:"original_transaction_id"`
	// The date and time that the item was purchased.
	// For a transaction that restores a previous transaction, the purchase date is the same as the original purchase date. Use Original Purchase Date to get the date of the original transaction.
	// In an auto-renewable subscription receipt, the purchase date is the date when the subscription was either purchased or renewed (with or without a lapse).
	// For an automatic renewal that occurs on the expiration date of the current period, the purchase date is the start date of the next period, which is identical to the end date of the current period.
	PurchaseDate
	// For a transaction that restores a previous transaction, the date of the original transaction.
	// In an auto-renewable subscription receipt, this indicates the beginning of the subscription period, even if the subscription has been renewed.
	OriginalPurchaseDate
	// The expiration date for the subscription, expressed as the number of milliseconds since January 1, 1970, 00:00:00 GMT.
	// This key is only present for auto-renewable subscription receipts. Use this value to identify the date when the subscription will renew or expire, to determine if a customer should have access to content or service.
	// After validating the latest receipt, if the subscription expiration date for the latest renewal transaction is a past date, it is safe to assume that the subscription has expired.
	ExpiresDate
	// For an expired subscription, the reason for the subscription expiration.
	// “1” - Customer canceled their subscription.
	// “2” - Billing error; for example customer’s payment information was no longer valid.
	// “3” - Customer did not agree to a recent price increase.
	// “4” - Product was not available for purchase at the time of renewal.
	// “5” - Unknown error.
	// This key is only present for a receipt containing an expired auto-renewable subscription.
	// You can use this value to decide whether to display appropriate messaging in your app for customers to resubscribe.
	ExpirationIntent string `json:"expiration_intent"`
	// For an expired subscription, whether or not Apple is still attempting to automatically renew the subscription.
	// “1” - App Store is still attempting to renew the subscription.
	// “0” - App Store has stopped attempting to renew the subscription.
	// This key is only present for auto-renewable subscription receipts.
	// If the customer’s subscription failed to renew because the App Store was unable to complete the transaction,
	// this value will reflect whether or not the App Store is still trying to renew the subscription.
	IsInBillingRetryPeriod string `json:"is_in_billing_retry_period"`
	// For a subscription, whether or not it is in the free trial period.
	// This key is only present for auto-renewable subscription receipts.
	// The value for this key is "true" if the customer’s subscription is currently in the free trial period, or "false" if not.
	// Note: If a previous subscription period in the receipt has the value “true” for either the is_trial_period or the is_in_intro_offer_period key,
	// the user is not eligible for a free trial or introductory price within that subscription group.
	IsTrialPeriod string `json:"is_trial_period"`
	// For an auto-renewable subscription, whether or not it is in the introductory price period.
	// This key is only present for auto-renewable subscription receipts.
	// The value for this key is "true" if the customer’s subscription is currently in an introductory price period, or "false" if not.
	// Note: If a previous subscription period in the receipt has the value “true” for either the is_trial_period or the is_in_intro_offer_period key,
	// the user is not eligible for a free trial or introductory price within that subscription group.
	IsInIntroOfferPeriod string `json:"is_in_intro_offer_period"`
	// For a transaction that was canceled by Apple customer support, the time and date of the cancellation.
	// For an auto-renewable subscription plan that was upgraded, the time and date of the upgrade transaction.
	// Treat a canceled receipt the same as if no purchase had ever been made.
	// Note: A canceled in-app purchase remains in the receipt indefinitely.
	// Only applicable if the refund was for a non-consumable product, an auto-renewable subscription,
	// a non-renewing subscription, or for a free subscription.
	CancellationDate
	// For a transaction that was canceled, the reason for cancellation.
	// “1” - Customer canceled their transaction due to an actual or perceived issue within your app.
	// “0” - Transaction was canceled for another reason, for example, if the customer made the purchase accidentally.
	// Use this value along with the cancellation date to identify possible issues in your app that may lead customers to contact Apple customer support.
	CancellationReason string `json:"cancellation_reason,omitempty"`
	// A string that the App Store uses to uniquely identify the application that created the transaction.
	// If your server supports multiple applications, you can use this value to differentiate between them.
	// Apps are assigned an identifier only in the production environment, so this key is not present for receipts created in the test environment.
	// This field is not present for Mac apps.
	// See also Bundle Identifier.
	AppItemId string `json:"app_item_id"`
	// An arbitrary number that uniquely identifies a revision of your application.
	// This key is not present for receipts created in the test environment.
	// Use this value to identify the version of the app that the customer bought.
	VersionExternalIdentifier string `json:"version_external_identifier"`
	// The primary key for identifying subscription purchases.
	// This value is a unique ID that identifies purchase events across devices, including subscription renewal purchase events.
	WebOrderLineItemID string `json:"web_order_line_item_id,omitempty"`
	// The current renewal status for the auto-renewable subscription.
	// “1” - Subscription will renew at the end of the current subscription period.
	// “0” - Customer has turned off automatic renewal for their subscription.
	// This key is only present for auto-renewable subscription receipts, for active or expired subscriptions.
	// The value for this key should not be interpreted as the customer’s subscription status.
	// You can use this value to display an alternative subscription product in your app, for example, a lower level subscription plan that the customer can downgrade to from their current plan.
	AutoRenewStatus string `json:"auto_renew_status"`
	// The current renewal preference for the auto-renewable subscription.
	// This key is only present for auto-renewable subscription receipts. The value for this key corresponds to the productIdentifier property of the product that the customer’s subscription renews.
	// You can use this value to present an alternative service level to the customer before the current subscription period ends.
	AutoRenewProductId string `json:"auto_renew_product_id"`
	// The current price consent status for a subscription price increase.
	// “1” - Customer has agreed to the price increase. Subscription will renew at the higher price.
	// “0” - Customer has not taken action regarding the increased price. Subscription expires if the customer takes no action before the renewal date.
	// This key is only present for auto-renewable subscription receipts if the subscription price was increased without keeping the existing price for active subscribers.
	// You can use this value to track customer adoption of the new price and take appropriate action.
	PriceConsentStatus string `json:"price_consent_status"`
}

// ReceiptCreationDate type indicates the date when the app receipt was created.
type ReceiptCreationDate struct {
	CreationDate    string `json:"receipt_creation_date,omitempty"`
	CreationDateMS  string `json:"receipt_creation_date_ms,omitempty"`
	CreationDatePST string `json:"receipt_creation_date_pst,omitempty"`
}

// ReceiptExpirationDate type indicates the date that the app receipt expires
type ReceiptExpirationDate struct {
	ExpirationDate    string `json:"receipt_expiration_date,omitempty"`
	ExpirationDateMS  string `json:"receipt_expiration_date_ms,omitempty"`
	ExpirationDatePST string `json:"receipt_expiration_date_pst,omitempty"`
}

// ReceiptRequestDate type indicates the date and time that the request was sent
type ReceiptRequestDate struct {
	RequestDate    string `json:"request_date,omitempty"`
	RequestDateMS  string `json:"request_date_ms,omitempty"`
	RequestDatePST string `json:"request_date_pst,omitempty"`
}

// OriginalPurchaseDate type indicates the beginning of the subscription period
type OriginalPurchaseDate struct {
	OriginalPurchaseDate    string `json:"original_purchase_date,omitempty"`
	OriginalPurchaseDateMS  string `json:"original_purchase_date_ms,omitempty"`
	OriginalPurchaseDatePST string `json:"original_purchase_date_pst,omitempty"`
}

// PurchaseDate type indicates the date and time that the item was purchased
type PurchaseDate struct {
	PurchaseDate    string `json:"purchase_date,omitempty"`
	PurchaseDateMS  string `json:"purchase_date_ms,omitempty"`
	PurchaseDatePST string `json:"purchase_date_pst,omitempty"`
}

// The ExpiresDate type indicates the expiration date for the subscription
type ExpiresDate struct {
	ExpiresDate             string `json:"expires_date,omitempty"`
	ExpiresDateMS           string `json:"expires_date_ms,omitempty"`
	ExpiresDatePST          string `json:"expires_date_pst,omitempty"`
	ExpiresDateFormatted    string `json:"expires_date_formatted,omitempty"`
	ExpiresDateFormattedPST string `json:"expires_date_formatted_pst,omitempty"`
}

// A pending renewal may refer to a renewal that is scheduled in the future or a renewal that failed in the past for some reason.
type PendingRenewalInfo struct {
	ProductID                      string `json:"product_id"`
	SubscriptionExpirationIntent   string `json:"expiration_intent"`
	SubscriptionAutoRenewProductID string `json:"auto_renew_product_id"`
	SubscriptionRetryFlag          string `json:"is_in_billing_retry_period"`
	SubscriptionAutoRenewStatus    string `json:"auto_renew_status"`
	SubscriptionPriceConsentStatus string `json:"price_consent_status"`
}

// The CancellationDate type indicates the time and date of the cancellation by Apple customer support
type CancellationDate struct {
	CancellationDate    string `json:"cancellation_date,omitempty"`
	CancellationDateMS  string `json:"cancellation_date_ms,omitempty"`
	CancellationDatePST string `json:"cancellation_date_pst,omitempty"`
}
