package search

import (
	"log"
	"strings" // Added for string manipulation

	"data-normalizer/internal/offers"
	"data-normalizer/internal/products"
)

// ElasticsearchClient is a mock client for Elasticsearch operations.
type ElasticsearchClient struct {
	// In a real implementation, this would hold the actual Elasticsearch client.
	// For this mock, we'll need access to the normalized product store.
	// This is a hack for the mock; in real world, ES client wouldn't directly access app's stores.
	normalizedProductStore map[string]*products.NormalizedProduct
}

// NewElasticsearchClient creates a new mock Elasticsearch client.
func NewElasticsearchClient(productStore map[string]*products.NormalizedProduct) *ElasticsearchClient {
	return &ElasticsearchClient{
		normalizedProductStore: productStore,
	}
}

// IndexProduct simulates indexing a normalized product into Elasticsearch.
func (c *ElasticsearchClient) IndexProduct(product *products.NormalizedProduct) error {
	log.Printf("Simulating: Indexing product '%s' (ID: %s) into Elasticsearch 'products' index.", product.Title, product.ID)
	// In a real scenario, convert product to Elasticsearch document and index.
	return nil
}

// IndexOffer simulates indexing an offer into Elasticsearch.
func (c *ElasticsearchClient) IndexOffer(offer *offers.Offer) error {
	log.Printf("Simulating: Indexing offer '%s' (ProductID: %s) into Elasticsearch 'offers' index.", offer.ID, offer.ProductID)
	// In a real scenario, convert offer to Elasticsearch document and index.
	return nil
}

// SearchProducts simulates searching for products based on a keyword.
func (c *ElasticsearchClient) SearchProducts(keyword string) []*products.NormalizedProduct {
	var results []*products.NormalizedProduct
	searchLower := strings.ToLower(keyword)

	for _, prod := range c.normalizedProductStore {
		if strings.Contains(strings.ToLower(prod.Title), searchLower) ||
			strings.Contains(strings.ToLower(prod.SKU), searchLower) {
			results = append(results, prod)
		}
	}
	log.Printf("Simulating: Found %d products for keyword '%s'.", len(results), keyword)
	return results
}
