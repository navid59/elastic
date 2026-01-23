package normalization

// CategoryMap defines the rule-based mapping from merchant categories
// to marketplace categories.
// In a real system, this would be stored in a database and be manageable.
var categoryMap = map[string]string{
	"t-shirts": "Apparel > Shirts",
	"hoodies":  "Apparel > Hoodies & Sweatshirts",
	"albums":   "Music > Records & Vinyl",
	"posters":  "Home & Decor > Wall Art",
}

// LowConfidenceKeywords are keywords that suggest a mapping might be uncertain.
var lowConfidenceKeywords = []string{"accessory", "other", "general"}

// MapCategoryResult holds the result of a category mapping operation.
type MapCategoryResult struct {
	MarketplaceCategory string
	LowConfidence       bool
}

// MapCategory maps a merchant's category to the marketplace's canonical category.
func MapCategory(merchantCategory string) MapCategoryResult {
	// Simple rule-based mapping for now.
	// We'll just do a direct lookup. A real system would use fuzzy matching, etc.
	if marketplaceCat, ok := categoryMap[merchantCategory]; ok {
		return MapCategoryResult{
			MarketplaceCategory: marketplaceCat,
			LowConfidence:       false, // Direct match is high confidence
		}
	}

	// Basic low-confidence check
	for _, keyword := range lowConfidenceKeywords {
		if merchantCategory == keyword {
			return MapCategoryResult{
				MarketplaceCategory: "Uncategorized", // Default for low confidence
				LowConfidence:       true,
			}
		}
	}

	// If no mapping is found, we can assign a default and flag for review.
	// For now, we don't block it.
	return MapCategoryResult{
		MarketplaceCategory: "Uncategorized",
		LowConfidence:       true,
	}
}
