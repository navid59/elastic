package normalization

import (
	"encoding/json"
	"fmt"
	"log"

	"data-normalizer/internal/products"
)

// Normalize takes a raw product payload and transforms it into a NormalizedProduct.
// This is a simplified version for now.
func Normalize(rawPayload json.RawMessage) (*products.NormalizedProduct, error) {
	var wooProduct map[string]interface{}
	if err := json.Unmarshal(rawPayload, &wooProduct); err != nil {
		return nil, fmt.Errorf("could not parse raw product: %w", err)
	}

	// Basic normalization: Extract title and SKU.
	// A real implementation would handle type assertions gracefully, clean text,
	// map categories, etc.
	
	title := ""
	if t, ok := wooProduct["name"].(string); ok {
		title = t
	}

	sku := ""
	if s, ok := wooProduct["sku"].(string); ok {
		sku = s
	}

	// ... (previous normalization logic for title and sku)

	// Category normalization
	mappedCategories := []string{}
	if cats, ok := wooProduct["categories"].([]interface{}); ok {
		for _, cat := range cats {
			if catMap, ok := cat.(map[string]interface{}); ok {
				if catName, ok := catMap["name"].(string); ok {
					// Map the category
					result := MapCategory(catName)
					// For now, just add the mapped category.
					// A more complex system might handle multiple categories
					// or low confidence flags differently.
					mappedCategories = append(mappedCategories, result.MarketplaceCategory)
				}
			}
		}
	}
	if len(mappedCategories) == 0 {
		mappedCategories = append(mappedCategories, "Uncategorized")
	}


	normalized := &products.NormalizedProduct{
		// A real ID would be generated based on matching logic.
		// For now, we use a placeholder.
		ID: fmt.Sprintf("prod-%s", sku),
		SKU: sku,
		Title: title,
		Categories: mappedCategories,
		// Other fields like Description, etc., would be normalized here.
	}
	
	log.Printf("Normalized product: %+v", normalized)
	return normalized, nil
}
