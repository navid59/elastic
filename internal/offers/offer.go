package offers

// Offer represents a merchant's specific offer for a given product.
type Offer struct {
	ID          string  `json:"id"`
	ProductID   string  `json:"product_id"` // Links to the NormalizedProduct
	MerchantID  string  `json:"merchant_id"`
	Price       string  `json:"price"` // Using string to handle different formats initially
	Stock       int     `json:"stock"`
	Delivery    string  `json:"delivery_info"` // Placeholder for delivery details
	TrustScore  float64 `json:"trust_score"`
}
