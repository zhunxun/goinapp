package ios

// Receipt type has the receipt property
type Receipt struct {
	// The app’s bundle identifier.
	// This corresponds to the value of CFBundleIdentifier in the Info.plist file.
	// Use this value to validate if the receipt was indeed generated for your app.
	BundleID string `json:"bundle_id,omitempty"`
	// The app’s version number.
	// This corresponds to the value of CFBundleVersion (in iOS) or CFBundleShortVersionString (in macOS) in the Info.plist
	ApplicationVersion string `json:"application_version,omitempty"`
	// The receipt for an in-app purchase.
	// In the JSON file, the value of this key is an array containing all in-app purchase receipts based on the in-app purchase
	// transactions present in the input base-64 receipt-data. For receipts containing auto-renewable subscriptions,
	// check the value of the latest_receipt_info key to get the status of the most recent renewal.
	//
	// The in-app purchase receipt for a consumable product is added to the receipt when the purchase is made. It is kept in the receipt
	// until your app finishes that transaction. After that point, it is removed from the receipt the next time the receipt is
	// updated - for example, when the user makes another purchase or if your app explicitly refreshes the receipt.
	// The in-app purchase receipt for a non-consumable product, auto-renewable subscription, non-renewing subscription,
	// or free subscription remains in the receipt indefinitely.
	InApp InApps `json:"in_app,omitempty"`
	// The version of the app that was originally purchased.
	// This corresponds to the value of CFBundleVersion (in iOS) or CFBundleShortVersionString (in macOS) in the Info.plist file
	// when the purchase was originally made. In the sandbox environment, the value of this field is always “1.0”
	OriginalApplicationVersion string `json:"original_application_version,omitempty"`
	// The date when the app receipt was created.
	// When validating a receipt, use this date to validate the receipt’s signature.
	ReceiptCreationDate    string `json:"receipt_creation_date,omitempty"`
	ReceiptCreationDateMS  int64  `json:"receipt_creation_date_ms,string,omitempty"`
	ReceiptCreationDatePST string `json:"receipt_creation_date_pst,omitempty"`
	// The date that the app receipt expires.
	// This key is present only for apps purchased through the Volume Purchase Program.
	// If this key is not present, the receipt does not expire.
	// When validating a receipt, compare this date to the current date to determine whether the receipt is expired.
	// Do not try to use this date to calculate any other information, such as the time remaining before expiration.
	ReceiptExpirationDate    string `json:"receipt_expiration_date,omitempty"`
	ReceiptExpirationDateMS  int64  `json:"receipt_expiration_date_ms,string,omitempty"`
	ReceiptExpirationDatePST string `json:"receipt_expiration_date_pst,omitempty"`
	// OriginalPurchaseDate type indicates the beginning of the subscription period
	OriginalPurchaseDate    string `json:"original_purchase_date,omitempty"`
	OriginalPurchaseDateMS  int64  `json:"original_purchase_date_ms,string,omitempty"`
	OriginalPurchaseDatePST string `json:"original_purchase_date_pst,omitempty"`
	// ReceiptRequestDate type indicates the date and time that the ValidationRequest was sent
	ReceiptRequestDate    string `json:"request_date,omitempty"`
	ReceiptRequestDateMS  int64  `json:"request_date_ms,string,omitempty"`
	ReceiptRequestDatePST string `json:"request_date_pst,omitempty"`
	// Undocumented field
	AdamID int `json:"adam_id,omitempty"`
	// Undocumented field
	AppItemID int `json:"app_item_id,omitempty"`
	// Undocumented field
	DownloadID int `json:"download_id,omitempty"`
	// Undocumented field
	VersionExternalIdentifier int `json:"version_external_identifier,omitempty"`
	// Undocumented field
	ReceiptType string `json:"receipt_type,omitempty"`
}
