package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"data-normalizer/internal/integrations/woocommerce" // Corrected import
	"data-normalizer/internal/merchants"
	"data-normalizer/internal/normalization"
	"data-normalizer/internal/offers"
	"data-normalizer/internal/products"
	"data-normalizer/internal/search"

	"github.com/google/uuid"
)

var (
	merchantStore          []merchants.Merchant
	rawProductStore        map[string][]RawProduct                // Maps ingestId to raw products
	normalizedProductStore map[string]*products.NormalizedProduct // Maps product ID to normalized product
	offerStore             map[string]*offers.Offer               // Maps offer ID to offer
	esClient               *search.ElasticsearchClient            // Added Elasticsearch client
)

// RawProduct represents the raw data received from a merchant.
type RawProduct struct {
	MerchantID        string          `json:"merchant_id"`
	ExternalProductID string          `json:"external_product_id"`
	RawPayload        json.RawMessage `json:"raw_payload"`
	ReceivedAt        time.Time       `json:"received_at"`
}

func init() {
	rawProductStore = make(map[string][]RawProduct)
	normalizedProductStore = make(map[string]*products.NormalizedProduct)
	offerStore = make(map[string]*offers.Offer)
	esClient = search.NewElasticsearchClient(normalizedProductStore) // Pass normalizedProductStore
}

func createMerchantHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var merchant merchants.Merchant
	err := json.NewDecoder(r.Body).Decode(&merchant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	merchantStore = append(merchantStore, merchant)
	log.Printf("Merchant created: %+v", merchant)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(merchant)
}

func merchantProductSyncHandler(w http.ResponseWriter, r *http.Request, merchantID string) {
	log.Printf("Starting product sync for merchant: %s", merchantID)

	var merchant merchants.Merchant
	found := false
	for _, m := range merchantStore {
		if m.ID == merchantID {
			merchant = m
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Merchant not found", http.StatusNotFound)
		return
	}

	baseURL := "https://example.com"
	client := woocommerce.NewClient(baseURL, merchant.Config.WooCommerce.ConsumerKey, merchant.Config.WooCommerce.ConsumerSecret)

	rawPayload, err := client.GetProducts()
	if err != nil {
		http.Error(w, "Failed to fetch products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ingestID := uuid.New().String()
	var productsData []map[string]interface{}
	if err := json.Unmarshal(rawPayload, &productsData); err != nil {
		http.Error(w, "Failed to parse products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var rawProductsToStore []RawProduct
	for _, p := range productsData {
		externalID := ""
		if id, ok := p["id"].(float64); ok {
			externalID = string(int(id))
		}

		productBytes, _ := json.Marshal(p)
		rawProductsToStore = append(rawProductsToStore, RawProduct{
			MerchantID:        merchantID,
			ExternalProductID: externalID,
			RawPayload:        productBytes,
			ReceivedAt:        time.Now(),
		})
	}
	rawProductStore[ingestID] = rawProductsToStore
	log.Printf("Stored %d raw products under ingest ID: %s", len(rawProductsToStore), ingestID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Sync initiated.",
		"ingest_id": ingestID,
	})
}

func normalizeBatchHandler(w http.ResponseWriter, r *http.Request, ingestID string) {
	log.Printf("Starting normalization for ingest batch: %s", ingestID)

	rawProducts, ok := rawProductStore[ingestID]
	if !ok {
		http.Error(w, "Ingest ID not found", http.StatusNotFound)
		return
	}

	normalizedCount := 0
	offerCount := 0
	for _, rawProduct := range rawProducts {
		normalizedProd, err := normalization.Normalize(rawProduct.RawPayload)
		if err != nil {
			log.Printf("Failed to normalize product %s: %v", rawProduct.ExternalProductID, err)
			continue
		}
		normalizedProductStore[normalizedProd.ID] = normalizedProd
		normalizedCount++
		// Index the normalized product
		if err := esClient.IndexProduct(normalizedProd); err != nil {
			log.Printf("Failed to index product %s: %v", normalizedProd.ID, err)
		}

		// --- Step 9: Offer Creation ---
		var wooProduct map[string]interface{}
		json.Unmarshal(rawProduct.RawPayload, &wooProduct)

		price := ""
		if p, ok := wooProduct["price"].(string); ok {
			price = p
		}

		stock := 0
		if s, ok := wooProduct["stock_quantity"].(float64); ok {
			stock = int(s)
		}

		// Find merchant trust score
		var trustScore float64
		for _, m := range merchantStore {
			if m.ID == rawProduct.MerchantID {
				trustScore = m.TrustScore
				break
			}
		}

		offer := &offers.Offer{
			ID:         uuid.New().String(),
			ProductID:  normalizedProd.ID,
			MerchantID: rawProduct.MerchantID,
			Price:      price,
			Stock:      stock,
			Delivery:   "standard", // Placeholder
			TrustScore: trustScore,
		}
		offerStore[offer.ID] = offer
		offerCount++
		// Index the offer
		if err := esClient.IndexOffer(offer); err != nil {
			log.Printf("Failed to index offer %s: %v", offer.ID, err)
		}
	}

	log.Printf("Successfully normalized %d products, created %d offers, and indexed both from ingest batch %s", normalizedCount, offerCount, ingestID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":             "Normalization completed.",
		"ingest_id":           ingestID,
		"normalized_products": normalizedCount,
		"created_offers":      offerCount,
	})
}

func searchProductsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	keyword := r.URL.Query().Get("keyword")
	log.Printf("Searching for products with keyword: '%s'", keyword)

	results := esClient.SearchProducts(keyword)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) >= 1 && parts[0] == "merchants" {
		if len(parts) == 4 && parts[2] == "sync" && parts[3] == "products" {
			if r.Method == http.MethodPost {
				merchantProductSyncHandler(w, r, parts[1])
				return
			}
		} else if len(parts) == 1 {
			if r.Method == http.MethodPost {
				createMerchantHandler(w, r)
				return
			}
		}
	}

	if len(parts) == 2 && parts[0] == "normalize" {
		if r.Method == http.MethodPost {
			normalizeBatchHandler(w, r, parts[1])
			return
		}
	}

	if len(parts) == 2 && parts[0] == "search" && parts[1] == "products" {
		if r.Method == http.MethodGet {
			searchProductsHandler(w, r)
			return
		}
	}

	http.NotFound(w, r)
}

func main() {
	log.Println("Application starting...")
	http.HandleFunc("/", apiHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}