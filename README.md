
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
- **External Merchant Platforms** – WooCommerce, Magento, etc.  
- **Elasticsearch** – Search and analytics engine  

---

## Core Use Cases

### UC-01: Onboard Merchant
Register a new merchant and configure how their product data will be collected (API, webhook, CSV, or scraping).  
**Outcome:** Merchant is ready to send or sync product data.

### UC-02: Ingest Product Data
Receive raw product data from merchant systems and store it unchanged for traceability.  
**Outcome:** Raw data is safely stored and queued for processing.

### UC-03: Normalize Product Data
Standardize product titles, descriptions, attributes, units, and currencies into a common structure.  
**Outcome:** Clean, consistent, comparable product data.

### UC-04: Automatically Map Categories
Map merchant-specific categories to marketplace categories using rules or ML.  
**Outcome:** Unified category structure across all merchants.

### UC-05: Detect Similar or Duplicate Products
Identify products that represent the same real-world item and link them without removing merchant ownership.  
**Outcome:** Related products are intelligently connected.

### UC-06: Create or Update Merchant Offers
Create merchant-specific offers containing price, stock, delivery, and warranty details.  
**Outcome:** Multiple offers can exist for the same product.

### UC-07: Index Data into Elasticsearch
Index products and offers for fast search and filtering.  
**Outcome:** High-performance, accurate search results.

---

## Buyer Use Cases

### UC-08: Search Products
Buyers search for products using keywords and filters.  
**Outcome:** Relevant products are returned quickly.

### UC-09: Compare Merchant Offers
Buyers compare prices, delivery times, warranties, and trust scores.  
**Outcome:** Buyers make informed purchasing decisions.

---

## Administrative & Control Use Cases

### UC-10: Review Product Matching Suggestions
Admins review uncertain duplicate matches without blocking the pipeline.  
**Outcome:** Matching accuracy improves over time.

### UC-11: Review Sensitive Products
Admins review flagged products (alcohol, medical, children).  
**Outcome:** Regulatory compliance is ensured.

### UC-12: Manage Category & Attribute Mappings
Admins adjust mappings to improve normalization quality.  
**Outcome:** Continuous data quality improvement.

### UC-13: Handle Processing Errors
Automatically retry failed jobs and store unresolved issues for review.  
**Outcome:** System remains reliable and fault-tolerant.

### UC-14: Track Data Source & History
View raw data and normalization steps for any product.  
**Outcome:** Full transparency and auditability.

### UC-15: Handle GDPR Requests
Support deletion and anonymization of data on request.  
**Outcome:** GDPR compliance and user trust.

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

## License

Educational / Academic Use
