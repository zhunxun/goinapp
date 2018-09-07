package ios

// Notification represent Apple status update notification ValidationRequest.
type Notification struct {
	// Specifies whether the notification is for a sandbox or a production environment: Sandbox and PROD
	Environment string `json:"environment"`
	// Represent the type of Apple status update notification.
	// INITIAL_BUY - Initial purchase of the subscription.
	// CANCEL - Subscription was canceled by Apple customer support.
	// RENEWAL - Automatic renewal was successful for an expired subscription.
	// INTERACTIVE_RENEWAL - Customer renewed a subscription interactively after it lapsed.
	// DID_CHANGE_RENEWAL_PREF - Customer changed the plan that takes affect at the next subscription renewal. Current active plan is not affected.
	NotificationType string `json:"notification_type"`
	// This value is the same as the shared secret you POST when validating receipts.
	Password string `json:"password"`
	// This value is the same as the Original Transaction Identifier in the receipt.
	// You can use this value to relate multiple iOS 6-style transaction receipts for an individual customer’s subscription.
	OriginalTransactionId string `json:"original_transaction_id"`
	// The time and date that a transaction was cancelled by Apple customer support. Posted only if the notification_type is CANCEL.
	CancellationDate    string `json:"cancellation_date,omitempty"`
	CancellationDateMS  string `json:"cancellation_date_ms,omitempty"`
	CancellationDatePST string `json:"cancellation_date_pst,omitempty"`
	// The primary key for identifying a subscription purchase. Posted only if the notification_type is CANCEL.
	WebOrderLineItemID string `json:"web_order_line_item_id"`
	// Posted if the notification_type is RENEWAL or INTERACTIVE_RENEWAL, and only if the renewal is successful.
	// Posted also if the notification_type is INITIAL_BUY. Not posted for notification_type CANCEL.
	LatestReceipt string `json:"latest_receipt"`
	// The JSON representation of the receipt for the most recent renewal. Posted only if renewal is successful.
	// Not posted for notification_type CANCEL.
	LatestReceiptInfo Receipt `json:"latest_receipt_info"`
	// The base-64 encoded transaction receipt for the most recent renewal transaction. Posted only if the subscription expired.
	LatestExpiredReceipt string `json:"latest_expired_receipt"`
	// The JSON representation of the receipt for the most recent renewal transaction.
	// Posted only if the notification_type is RENEWAL or CANCEL or if renewal failed and subscription expired.
	LatestExpiredReceiptInfo Receipt `json:"latest_expired_receipt_info"`
	// A Boolean value indicated by strings “true” or “false”. This is the same as the auto renew status in the receipt.
	AutoRenewStatus bool `json:"auto_renew_status"`
	// A Boolean value indicated by strings “true” or “false”. This is the same as the auto renew status in the receipt.
	AutoRenewAdamId string `json:"auto_renew_adam_id"`
	// This is the same as the Subscription Auto Renew Preference in the receipt.
	AutoRenewProductId string `json:"auto_renew_product_id"`
	// This is the same as the Subscription Expiration Intent in the receipt.
	// Posted only if notification_type is RENEWAL or INTERACTIVE_RENEWAL.
	ExpirationIntent string `json:"expiration_intent"`
}
