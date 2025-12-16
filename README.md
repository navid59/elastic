
# Data Normalizer & Integrator Microservice
**Elasticsearch-powered product ingestion and normalization for a multi-merchant marketplace**

---

## Overview

This project implements a **Data Normalizer + Integrator** microservice for an e-commerce marketplace that does not own products itself.
The marketplace works with small and medium merchants using different platforms such as **WooCommerce, Magento, PrestaShop, and OpenCart**.

The service collects product data from merchants, normalizes and standardizes it, handles missing and inconsistent data, detects similar or duplicate products, and indexes the cleaned data into **Elasticsearch**.
Each merchant keeps full ownership of their products through **offers**, while buyers can easily search, compare, and choose the best option.

---

## Key Goals

- Integrate product data from heterogeneous merchant platforms
- Normalize formats, attributes, units, and currencies
- Automatically map merchant categories to marketplace categories
- Detect and link similar or duplicate products without removing merchant data
- Preserve merchant rights through offer-based modeling
- Provide fast and accurate search using Elasticsearch
- Support scalability, monitoring, and GDPR compliance

---

## System Actors

- **Merchant** – Provides product data and manages offers
- **Buyer** – Searches and compares products
- **Marketplace Admin** – Manages merchants, mappings, and reviews
- **Data Reviewer** – Reviews uncertain matches and sensitive products
- **External Merchant Platforms** – WooCommerce, Magento, PrestaShop, OpenCart
- **Elasticsearch** – Search and analytics engine

---

## Core Use Cases

### UC-01: Onboard Merchant
**Actor:** Marketplace Admin

**Description:**  
Register a new merchant and configure how their product data will be integrated into the system. This ensures secure and reliable data exchange.

**Steps:**
- Admin creates a merchant account
- Selects integration method (API, webhook, CSV, scraping)
- Stores credentials or connection details
- Assigns initial trust score
- Activates merchant

**Outcome:**  
Merchant is ready to send or be synced for product data.

---

### UC-02: Ingest Product Data
**Actors:** Merchant Platform, System

**Description:**  
Receive raw product data from merchants either by push or pull mechanisms and store it unchanged for traceability.

**Steps:**
- Receive or fetch product data
- Validate basic structure
- Store raw payload
- Queue data for processing

**Outcome:**  
Raw merchant data is safely stored and ready for normalization.

---

### UC-03: Normalize Product Data
**Actor:** System

**Description:**  
Clean and standardize incoming product data so it follows a unified internal format.

**Steps:**
- Normalize text fields
- Convert units and currencies
- Standardize attributes
- Detect missing fields
- Add normalization flags

**Outcome:**  
Consistent and comparable product data across merchants.

---

### UC-04: Map Categories Automatically
**Actor:** System

**Description:**  
Automatically map merchant categories to marketplace categories using rules or machine learning.

**Steps:**
- Analyze merchant category
- Apply mapping logic
- Assign marketplace category
- Flag low-confidence results

**Outcome:**  
Unified product categorization.

---

### UC-05: Detect Similar or Duplicate Products
**Actor:** System

**Description:**  
Identify products that represent the same real-world item while preserving merchant ownership.

**Steps:**
- Generate match candidates
- Calculate similarity scores
- Link related products
- Flag uncertain cases

**Outcome:**  
Related products are linked without removing merchant data.

---

### UC-06: Create or Update Merchant Offers
**Actor:** System

**Description:**  
Create offers that represent merchant-specific selling conditions.

**Steps:**
- Create or update offer
- Assign price, stock, delivery
- Attach trust score
- Link to product

**Outcome:**  
Multiple merchant offers per product.

---

### UC-07: Index Data into Elasticsearch
**Actor:** System

**Description:**  
Index normalized products and offers for fast search and filtering.

**Steps:**
- Index products
- Index offers
- Update search documents

**Outcome:**  
Fast and accurate search results.

---

## Buyer Use Cases

### UC-08: Search Products
**Actor:** Buyer

**Description:**  
Allow buyers to search for products using keywords and filters.

**Steps:**
- Enter query
- Execute search
- Display results

**Outcome:**  
Relevant products returned quickly.

---

### UC-09: Compare Merchant Offers
**Actor:** Buyer

**Description:**  
Enable buyers to compare different offers for the same product.

**Steps:**
- View offers
- Sort and filter
- Choose merchant

**Outcome:**  
Transparent and informed purchasing decisions.

---

## Administrative & Control Use Cases

### UC-10: Review Product Matching Suggestions
**Actors:** Admin / Data Reviewer

**Description:**  
Review system-suggested product matches when confidence is low.

**Steps:**
- View suggestions
- Approve or reject
- Save decision

**Outcome:**  
Improved matching accuracy.

---

### UC-11: Review Sensitive Products
**Actor:** Admin

**Description:**  
Review products that require special approval such as alcohol or medical items.

**Steps:**
- Review flagged products
- Approve or restrict

**Outcome:**  
Regulatory compliance without blocking ingestion.

---

### UC-12: Manage Category and Attribute Mappings
**Actor:** Admin

**Description:**  
Improve automatic mappings through manual adjustments.

**Steps:**
- Edit mappings
- Apply changes
- Reprocess if needed

**Outcome:**  
Higher data quality.

---

### UC-13: Handle Processing Errors
**Actor:** System

**Description:**  
Automatically handle ingestion and processing errors.

**Steps:**
- Detect error
- Retry processing
- Store failures
- Notify admin

**Outcome:**  
System reliability.

---

### UC-14: Track Data Source and History
**Actor:** Admin

**Description:**  
Trace products back to their original source.

**Steps:**
- View raw data
- View normalization history

**Outcome:**  
Full transparency and auditability.

---

### UC-15: Handle GDPR Requests
**Actor:** Admin

**Description:**  
Support GDPR-related data requests.

**Steps:**
- Receive request
- Delete or anonymize data
- Log actions

**Outcome:**  
GDPR compliance.

---

## High-Level Architecture

Merchant Platforms  
→ Data Ingestion Layer  
→ Message Queue  
→ Data Normalizer & Integrator  
→ PostgreSQL (Raw + Normalized Data)  
→ Elasticsearch (Products & Offers)  
→ Marketplace Search & Admin Tools

---

## Technology Stack (Suggested)

- Backend: Python (FastAPI) or Node.js (NestJS)
- Queue: RabbitMQ or Kafka
- Database: PostgreSQL
- Search: Elasticsearch
- Admin UI: React
- Deployment: Docker / Docker Compose

---

## Project Context

This project is developed as part of a **special topic course** and focuses on real-world challenges in marketplace data integration, normalization, and search.

---

## License

Educational / Academic Use
