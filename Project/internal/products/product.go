package products

// NormalizedProduct is the canonical, standardized representation of a product
// in the marketplace, stripped of any merchant-specific data.
type NormalizedProduct struct {
	ID          string            `json:"id"`
	SKU         string            `json:"sku"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Categories  []string          `json:"categories"`
	Brand       string            `json:"brand"`
	Attributes  map[string]string `json:"attributes"`
	// Images, dimensions, and other universal fields would go here.
}
