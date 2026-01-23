# Mock Data Handling Locations

This document shows exactly where mock data is generated, stored, and processed in the application.

---

## 📍 Overview: Where Mock Data Lives

The application uses **3 types of mock implementations**:

1. **WooCommerce Client** - Returns dummy product data
2. **In-Memory Stores** - Replaces database (data lost on restart)
3. **Elasticsearch Client** - Simulates search (uses in-memory string matching)

---

## 1. WooCommerce Mock Data Generation

### Location: `internal/integrations/woocommerce/client.go`

**Line 47**: This is where the dummy product data is created:

```go
func (c *Client) GetProducts() ([]byte, error) {
    url := fmt.Sprintf("%s/wp-json/wc/v3/products", c.baseURL)
    
    // ❌ NO ACTUAL HTTP REQUEST HERE
    // ❌ NO AUTHENTICATION
    // ❌ NO REAL API CALL
    
    fmt.Printf("Making request to: %s\n", url)
    
    // ✅ THIS IS THE MOCK DATA:
    dummyResponse := []byte(`[{"id": 1, "name": "Sample Product"}]`)
    return dummyResponse, nil
}
```

**What it does:**
- Always returns the same hardcoded JSON: `[{"id": 1, "name": "Sample Product"}]`
- Does NOT make any HTTP request
- Does NOT authenticate with WooCommerce
- Does NOT fetch real products

**When it's called:**
- Called from `cmd/api/main.go` line 84 in `merchantProductSyncHandler()`

---

## 2. In-Memory Data Stores (No Database)

### Location: `cmd/api/main.go`

**Lines 20-25**: These are the in-memory stores that replace a database:

```go
var (
    merchantStore          []merchants.Merchant                    // Line 21
    rawProductStore        map[string][]RawProduct                  // Line 22
    normalizedProductStore map[string]*products.NormalizedProduct   // Line 23
    offerStore             map[string]*offers.Offer                 // Line 24
    esClient               *search.ElasticsearchClient              // Line 25
)
```

**Lines 36-40**: Initialization (happens on app start):

```go
func init() {
    rawProductStore = make(map[string][]RawProduct)              // Line 37
    normalizedProductStore = make(map[string]*products.NormalizedProduct) // Line 38
    offerStore = make(map[string]*offers.Offer)                  // Line 39
    esClient = search.NewElasticsearchClient(normalizedProductStore) // Line 40
}
```

### How Data Flows Through These Stores:

#### A. Merchant Storage (Line 56)
```go
merchantStore = append(merchantStore, merchant)
```
- Stores merchants in a slice
- ❌ Lost when server restarts

#### B. Raw Product Storage (Line 112)
```go
rawProductStore[ingestID] = rawProductsToStore
```
- Stores raw WooCommerce data (from mock client)
- Key: `ingestID` (UUID)
- Value: Array of `RawProduct` structs

#### C. Normalized Product Storage (Line 139)
```go
normalizedProductStore[normalizedProd.ID] = normalizedProd
```
- Stores normalized products after processing
- Key: Product ID (e.g., `"prod-{sku}"`)
- Value: `NormalizedProduct` struct

#### D. Offer Storage (Line 178)
```go
offerStore[offer.ID] = offer
```
- Stores merchant offers
- Key: Offer ID (UUID)
- Value: `Offer` struct

---

## 3. Elasticsearch Mock Implementation

### Location: `internal/search/elasticsearch.go`

**Lines 11-16**: Mock client structure:

```go
// ElasticsearchClient is a mock client for Elasticsearch operations.
type ElasticsearchClient struct {
    // In a real implementation, this would hold the actual Elasticsearch client.
    // For this mock, we'll need access to the normalized product store.
    normalizedProductStore map[string]*products.NormalizedProduct
}
```

### Mock Operations:

#### A. IndexProduct() - Line 27-31
```go
func (c *ElasticsearchClient) IndexProduct(product *products.NormalizedProduct) error {
    log.Printf("Simulating: Indexing product '%s' (ID: %s) into Elasticsearch 'products' index.", 
               product.Title, product.ID)
    // ❌ NO ACTUAL INDEXING - Just logs!
    return nil
}
```
- **Does NOT** connect to Elasticsearch
- **Does NOT** create indices
- **Only logs** a message

#### B. IndexOffer() - Line 34-38
```go
func (c *ElasticsearchClient) IndexOffer(offer *offers.Offer) error {
    log.Printf("Simulating: Indexing offer '%s' (ProductID: %s) into Elasticsearch 'offers' index.", 
               offer.ID, offer.ID)
    // ❌ NO ACTUAL INDEXING - Just logs!
    return nil
}
```
- Same as above - only logs

#### C. SearchProducts() - Line 41-53
```go
func (c *ElasticsearchClient) SearchProducts(keyword string) []*products.NormalizedProduct {
    var results []*products.NormalizedProduct
    searchLower := strings.ToLower(keyword)

    // ❌ NOT USING ELASTICSEARCH!
    // ✅ Using in-memory string matching instead:
    for _, prod := range c.normalizedProductStore {
        if strings.Contains(strings.ToLower(prod.Title), searchLower) ||
           strings.Contains(strings.ToLower(prod.SKU), searchLower) {
            results = append(results, prod)
        }
    }
    log.Printf("Simulating: Found %d products for keyword '%s'.", len(results), keyword)
    return results
}
```
- **Does NOT** query Elasticsearch
- **Does** search the in-memory `normalizedProductStore` map
- Uses simple string matching (`strings.Contains`)
- No advanced search features (fuzzy, filters, etc.)

---

## 4. Complete Data Flow with Mock Components

```
┌─────────────────────────────────────────────────────────────┐
│ 1. POST /merchants/{id}/sync/products                      │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ merchantProductSyncHandler()                                │
│ - Calls: woocommerce.Client.GetProducts()                   │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 🔴 MOCK: woocommerce/client.go:47                           │
│ Returns: [{"id": 1, "name": "Sample Product"}]             │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ Stores in: rawProductStore[ingestID]                        │
│ 🔴 IN-MEMORY (no database)                                  │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. POST /normalize/{ingestId}                              │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ normalizeBatchHandler()                                     │
│ - Normalizes each raw product                               │
│ - Creates NormalizedProduct                                 │
│ - Stores in: normalizedProductStore[productID]              │
│   🔴 IN-MEMORY (no database)                                │
│ - Creates Offer                                             │
│ - Stores in: offerStore[offerID]                            │
│   🔴 IN-MEMORY (no database)                                │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 🔴 MOCK: esClient.IndexProduct()                            │
│ - Only logs, doesn't index to Elasticsearch                 │
│ 🔴 MOCK: esClient.IndexOffer()                              │
│ - Only logs, doesn't index to Elasticsearch                 │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. GET /search/products?keyword=...                        │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 🔴 MOCK: esClient.SearchProducts()                          │
│ - Searches normalizedProductStore (in-memory map)           │
│ - Uses strings.Contains() - simple string matching          │
│ - Does NOT query Elasticsearch                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 5. Summary: All Mock Data Locations

| Component | File | Line(s) | What It Does |
|-----------|------|---------|--------------|
| **WooCommerce Mock** | `internal/integrations/woocommerce/client.go` | 47 | Returns `[{"id": 1, "name": "Sample Product"}]` |
| **Merchant Store** | `cmd/api/main.go` | 21, 56 | In-memory slice (no DB) |
| **Raw Product Store** | `cmd/api/main.go` | 22, 37, 112 | In-memory map (no DB) |
| **Normalized Product Store** | `cmd/api/main.go` | 23, 38, 139 | In-memory map (no DB) |
| **Offer Store** | `cmd/api/main.go` | 24, 39, 178 | In-memory map (no DB) |
| **ES Index Mock** | `internal/search/elasticsearch.go` | 27-31 | Only logs, doesn't index |
| **ES Offer Index Mock** | `internal/search/elasticsearch.go` | 34-38 | Only logs, doesn't index |
| **ES Search Mock** | `internal/search/elasticsearch.go` | 41-53 | In-memory string matching |

---

## 6. What Happens When Server Restarts?

**ALL DATA IS LOST** because everything is stored in memory:

1. ✅ `merchantStore` - **EMPTY** (all merchants gone)
2. ✅ `rawProductStore` - **EMPTY** (all raw products gone)
3. ✅ `normalizedProductStore` - **EMPTY** (all normalized products gone)
4. ✅ `offerStore` - **EMPTY** (all offers gone)

This is why you need a **real database** for production!

---

## 7. How to Verify Mock Data

### Test the WooCommerce Mock:
```bash
# This will always return the same dummy product
curl -X POST http://localhost:8080/merchants/test-merchant/sync/products
```

### Check what's stored in memory:
- Add logging to see what's in the stores
- Or add a debug endpoint to inspect stores
- Currently, you can only see data through search (which uses the mock ES)

### Test the Elasticsearch Mock:
```bash
# This searches in-memory, not Elasticsearch
curl "http://localhost:8080/search/products?keyword=sample"
```

---

## 8. To Replace Mocks with Real Implementations

### Replace WooCommerce Mock:
- **File**: `internal/integrations/woocommerce/client.go`
- **Change**: Implement real HTTP request with Basic Auth
- **Use**: `http.Client` with `req.SetBasicAuth()`

### Replace In-Memory Stores:
- **File**: `cmd/api/main.go`
- **Change**: Add database (PostgreSQL recommended)
- **Use**: Repository pattern with SQL queries

### Replace Elasticsearch Mock:
- **File**: `internal/search/elasticsearch.go`
- **Change**: Use official Elasticsearch Go client
- **Use**: `github.com/elastic/go-elasticsearch/v8`

---

*Last Updated: 2024*
