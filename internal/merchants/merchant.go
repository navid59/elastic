package merchants

// Merchant represents a merchant in the marketplace.
type Merchant struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Platform    string `json:"platform"`
	Config      Config `json:"config"`
	TrustScore  float64 `json:"trust_score"`
	Status      string `json:"status"`
}

// Config holds the integration-specific settings for a merchant.
type Config struct {
	WooCommerce WooCommerceConfig `json:"woocommerce"`
}

// WooCommerceConfig holds the credentials for the WooCommerce API.
type WooCommerceConfig struct {
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
}
