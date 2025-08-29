# NewsBrief MVP API Specification v1.0

## 2-Week MVP Implementation - Practical API Documentation

**Audience:** Frontend developers building the MVP  
**Scope:** Essential endpoints for 2-week implementation  
**Architecture:** 3 microservices (Core API + Scraper + Summarizer)

---

## 🌐 MVP Service Architecture

| Service               | URL                     | Responsibility                                        |
| --------------------- | ----------------------- | ----------------------------------------------------- |
| **Core API** (Go Gin) | `http://localhost:8080` | Feed, Stories, Search, User Management                |
| **Auth Service** (Go) | `http://localhost:8083` | User authentication, JWT issuing, password management |
| **Scraper** (FastAPI) | `http://localhost:8001` | Content extraction from URLs                          |
| **Summarizer** (Go)   | `http://localhost:8002` | Dual summary generation (short + medium)              |
| **RabbitMQ**          | `amqp://localhost:5672` | Message broker for asynchronous tasks                 |

### **Architectural Changes:**

- **Dedicated Auth Service**: Authentication is now handled by a separate service, improving security and separation of concerns.
- **Asynchronous Processing**: RabbitMQ is introduced to decouple services. Scraping and summarization are now asynchronous jobs, making the system more resilient and scalable.
- **4 services + message broker** instead of 3 for a more robust, loosely coupled system.

---

## 🔐 Authentication Service (Go - Port 8083)

Handles user registration, login, and token management.

### **POST /v1/auth/register**

Create new account.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "SecureP@ssw0rd123!",
  "name": "John Doe",
  "lang": "am"
}
```

**Success Response (201):**

```json
{
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "name": "John Doe",
    "email_verified": false,
    "preferences": {
      "lang": "am",
      "topics": [],
      "data_saver": true,
      "brief_type": "short"
    }
  },
  "verification_token_id": "66c1234567890abcdef1234a"
}
```

### **POST /v1/auth/login**

User authentication.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "SecureP@ssw0rd123!"
}
```

**Success Response (200):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "66c1234567890abcdef12349",
  "expires_in": 3600
}
```

### **POST /v1/auth/verify-email**

Email verification using unified token system.

**Request Body:**

```json
{
  "token": "66c1234567890abcdef1234a"
}
```

**Success Response (200):**

```json
{
  "message": "Email verified successfully",
  "email_verified": true
}
```

### **POST /v1/auth/forgot-password**

Request password reset.

**Request Body:**

```json
{
  "email": "user@example.com"
}
```

**Success Response (200):**

```json
{
  "message": "Reset instructions sent",
  "reset_token_id": "66c1234567890abcdef1234c"
}
```

### **PATCH /v1/me/password** 🔐

Change user password.

**Request Body:**

```json
{
  "current_password": "OldP@ssw0rd123!",
  "new_password": "NewSecureP@ssw0rd456!"
}
```

**Success Response (200):**

```json
{
  "message": "Password updated successfully",
  "new_access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "new_refresh_token": "66c1234567890abcdef1234d"
}
```

---

## 📱 Core API Endpoints (Go Gin - Port 8080)

### **GET /v1/feed**

Get paginated story feed with enhanced filtering.

**Query Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `lang` | `"am" \| "en"` | No | Summary language preference |
| `topic` | `string` | No | Topic filter |
| `source` | `string` | No | News outlet filter |
| `brief_type` | `"short" \| "medium"` | No | Summary type preference |
| `since` | `string` | No | ISO8601 timestamp |
| `limit` | `number` | No | Page size (1-50, default: 20) |
| `cursor` | `string` | No | Pagination cursor |

**Success Response (200):**

```json
{
  "items": [
    {
      "id": "507f1f77bcf86cd799439012",
      "title": "Ethiopia launches new agricultural initiative",
      "source": {
        "key": "addisstandard",
        "name": "Addis Standard",
        "logo_url": "https://cdn.newsbrief.et/sources/addisstandard.png"
      },
      "url": "https://addisstandard.com/news/ethiopia-agri-2025",
      "published_at": "2025-08-25T09:10:00Z",
      "summary_short": "Government announces $50M farming investment targeting 100,000 farmers.",
      "summary_bullets": [
        "Government announces $50M investment in rural farming",
        "Program targets 100,000 smallholder farmers nationwide",
        "Focus on drought-resistant crop varieties and irrigation"
      ],
      "summary_lang": "am",
      "topic_tags": ["agriculture", "economy"],
      "topic_image": "https://cdn.newsbrief.et/topics/agriculture.jpg",
      "processing_status": "completed",
      "reading_time": {
        "short": 1,
        "medium": 3
      }
    }
  ],
  "next_cursor": "eyJfaWQiOiI2NmMxMjM0NTY3ODkwYWJjZGVmMTIzNDUifQ==",
  "total_available": 156,
  "server_time": "2025-08-25T10:30:00Z"
}
```

---

### **GET /v1/story/:id**

Get detailed story by MongoDB ObjectId.

**Success Response (200):**

```json
{
  "id": "507f1f77bcf86cd799439012",
  "title": "Ethiopia launches new agricultural initiative",
  "source": {
    "key": "addisstandard",
    "name": "Addis Standard",
    "logo_url": "https://cdn.newsbrief.et/sources/addisstandard.png"
  },
  "url": "https://addisstandard.com/news/ethiopia-agri-2025",
  "published_at": "2025-08-25T09:10:00Z",
  "summary_short": "Government announces $50M farming investment.",
  "summary_bullets": [
    "Government announces $50M investment in rural farming",
    "Program targets 100,000 smallholder farmers nationwide",
    "Focus on drought-resistant crop varieties and irrigation"
  ],
  "summary_lang": "am",
  "topic_tags": ["agriculture", "economy"],
  "topic_image": "https://cdn.newsbrief.et/topics/agriculture.jpg",
  "reading_time": {
    "short": 1,
    "medium": 3
  },
  "word_count": 450,
  "scraped_at": "2025-08-25T09:05:00Z",
  "summarized_at": "2025-08-25T09:08:00Z"
}
```

---

### **GET /v1/search**

**MVP Implementation**: MongoDB text search (no semantic search).

**Query Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `q` | `string` | Yes | Search query (min 2 chars) |
| `topic` | `string` | No | Topic filter |
| `source` | `string` | No | Source filter |
| `limit` | `number` | No | Results per page (1-50) |

**Success Response (200):**

```json
{
  "query": "agriculture investment",
  "items": [
    {
      "id": "507f1f77bcf86cd799439012",
      "title": "Ethiopia launches new agricultural initiative",
      "source": {
        "key": "addisstandard",
        "name": "Addis Standard",
        "logo_url": "https://cdn.newsbrief.et/sources/addisstandard.png"
      },
      "published_at": "2025-08-25T09:10:00Z",
      "summary_short": "Government announces $50M farming investment.",
      "topic_tags": ["agriculture", "economy"],
      "text_score": 1.2,
      "matched_terms": ["agriculture", "investment"]
    }
  ],
  "total_matches": 12,
  "search_method": "mongodb_text_index"
}
```

---

### **GET /v1/topics**

Get topic categories with images and descriptions.

**Success Response (200):**

```json
{
  "topics": [
    {
      "key": "agriculture",
      "label": { "en": "Agriculture", "am": "ግብርና" },
      "description": {
        "en": "Farming, livestock, and agricultural development",
        "am": "የግብርና፣ የእንስሳት ሀብት እና የግብርና ልማት"
      },
      "image_url": "https://cdn.newsbrief.et/topics/agriculture.jpg",
      "story_count": 23
    },
    {
      "key": "economy",
      "label": { "en": "Economy", "am": "ኢኮኖሚ" },
      "description": {
        "en": "Business, finance, and economic news",
        "am": "የንግድ፣ የፋይናንስ እና የኢኮኖሚ ዜናዎች"
      },
      "image_url": "https://cdn.newsbrief.et/topics/economy.jpg",
      "story_count": 45
    }
  ],
  "total_topics": 6
}
```

---

### **GET /v1/sources**

Get available news outlets with subscription support.

**Query Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `subscribed_only` | `boolean` | No | Show only user's subscriptions |

**Success Response (200):**

```json
{
  "sources": [
    {
      "key": "addisstandard",
      "name": "Addis Standard",
      "description": "Independent news outlet covering Ethiopian politics",
      "url": "https://addisstandard.com",
      "logo_url": "https://cdn.newsbrief.et/sources/addisstandard.png",
      "languages": ["en", "am"],
      "topics": ["politics", "economy", "society"],
      "reliability_score": 0.92,
      "update_frequency": "hourly"
    }
  ],
  "total_sources": 12
}
```

---

## 🔑 Authentication Endpoints

### **POST /v1/auth/register**

Create new account.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "SecureP@ssw0rd123!",
  "name": "John Doe",
  "lang": "am"
}
```

**Success Response (201):**

```json
{
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "name": "John Doe",
    "email_verified": false,
    "preferences": {
      "lang": "am",
      "topics": [],
      "data_saver": true,
      "brief_type": "short"
    }
  },
  "verification_token_id": "66c1234567890abcdef1234a"
}
```

### **POST /v1/auth/login**

User authentication.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "SecureP@ssw0rd123!"
}
```

**Success Response (200):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "66c1234567890abcdef12349",
  "expires_in": 3600,
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "preferences": {
      "lang": "am",
      "topics": ["economy", "agriculture"],
      "subscribed_sources": ["addisstandard"],
      "brief_type": "short"
    }
  }
}
```

### **POST /v1/auth/verify-email**

Email verification using unified token system.

**Request Body:**

```json
{
  "token": "66c1234567890abcdef1234a"
}
```

**Success Response (200):**

```json
{
  "message": "Email verified successfully",
  "email_verified": true
}
```

### **POST /v1/auth/forgot-password**

Request password reset.

**Request Body:**

```json
{
  "email": "user@example.com"
}
```

**Success Response (200):**

```json
{
  "message": "Reset instructions sent",
  "reset_token_id": "66c1234567890abcdef1234c"
}
```

---

## 👤 User Management Endpoints

### **GET /v1/me** 🔐

Get user profile and preferences.

**Success Response (200):**

```json
{
  "id": "507f1f77bcf86cd799439011",
  "email": "user@example.com",
  "name": "John Doe",
  "email_verified": true,
  "preferences": {
    "lang": "am",
    "topics": ["economy", "agriculture", "politics"],
    "subscribed_sources": ["addisstandard", "ethiopianherald"],
    "brief_type": "short",
    "data_saver": true
  },
  "subscription": {
    "plan": "free",
    "source_limit": 5,
    "current_subscriptions": 2
  },
  "stats": {
    "stories_read": 147,
    "last_active": "2025-08-25T09:15:00Z"
  }
}
```

### **PATCH /v1/me/preferences** 🔐

Update user preferences.

**Request Body:**

```json
{
  "lang": "en",
  "topics": ["economy", "agriculture", "technology"],
  "brief_type": "medium",
  "data_saver": false
}
```

**Success Response (200):**

```json
{
  "preferences": {
    "lang": "en",
    "topics": ["economy", "agriculture", "technology"],
    "brief_type": "medium",
    "data_saver": false
  },
  "updated_at": "2025-08-25T10:30:00Z"
}
```

### **GET /v1/me/subscriptions** 🔐

Get user's source subscriptions.

**Success Response (200):**

```json
{
  "subscriptions": [
    {
      "source_key": "addisstandard",
      "source_name": "Addis Standard",
      "subscribed_at": "2025-08-01T10:00:00Z",
      "topics": ["politics", "economy"]
    }
  ],
  "total_subscriptions": 2,
  "subscription_limit": 5
}
```

### **POST /v1/me/subscriptions** 🔐

Subscribe to a news outlet.

**Request Body:**

```json
{
  "source_key": "ethiopianherald",
  "topics": ["all"]
}
```

**Success Response (201):**

```json
{
  "subscription": {
    "source_key": "ethiopianherald",
    "source_name": "Ethiopian Herald",
    "topics": ["all"],
    "subscribed_at": "2025-08-25T10:30:00Z"
  },
  "total_subscriptions": 3
}
```

### **DELETE /v1/me/subscriptions/:source_key** 🔐

Unsubscribe from news outlet.

**Success Response (200):**

```json
{
  "message": "Successfully unsubscribed",
  "source_key": "addisstandard",
  "remaining_subscriptions": 2
}
```

### **PATCH /v1/me/password** 🔐

Change user password.

**Request Body:**

```json
{
  "current_password": "OldP@ssw0rd123!",
  "new_password": "NewSecureP@ssw0rd456!"
}
```

**Success Response (200):**

```json
{
  "message": "Password updated successfully",
  "new_access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "new_refresh_token": "66c1234567890abcdef1234d"
}
```

---

## 📊 Basic Analytics (MVP)

### **GET /v1/me/analytics** 🔐

Simple reading statistics.

**Success Response (200):**

```json
{
  "period": "month",
  "reading_stats": {
    "articles_read": 89,
    "time_spent_minutes": 267,
    "completion_rate": 0.76
  },
  "topic_breakdown": [
    { "topic": "economy", "articles": 34, "percentage": 38.2 },
    { "topic": "agriculture", "articles": 28, "percentage": 31.5 }
  ],
  "source_breakdown": [
    { "source": "addisstandard", "articles": 45, "percentage": 50.6 }
  ]
}
```

---

## 🔧 Asynchronous Operations (RabbitMQ)

The Core API uses RabbitMQ to delegate long-running tasks to worker services. This ensures the API remains responsive.

### **Scraping Queue**

- **Queue**: `scrape_tasks`
- **Trigger**: A new URL is discovered.
- **Action**: Core API publishes a message to the `scrape_tasks` queue.
- **Consumer**: The Scraper service listens for messages, scrapes the content, and stores it.

**Sample Message:**

```json
{
  "url": "https://addisstandard.com/news/ethiopia-agri-2025",
  "source_key": "addisstandard"
}
```

### **Summarization Queue**

- **Queue**: `summarize_tasks`
- **Trigger**: The Scraper service successfully scrapes new content.
- **Action**: Scraper service publishes a message to the `summarize_tasks` queue.
- **Consumer**: The Summarizer service listens for messages, generates summaries, and updates the story record.

**Sample Message:**

```json
{
  "story_id": "507f1f77bcf86cd799439012",
  "text": "The Ethiopian government today announced...",
  "title": "Ethiopia launches new agricultural initiative",
  "target_lang": "am"
}
```

---

## 🐳 MVP Deployment

```yaml
# docker-compose.yml
version: "3.8"
services:
  mongodb:
    image: mongo:7.0
    ports:
      - "27017:27017"

  rabbitmq:
    image: rabbitmq:3-management
    ports:
      - "5672:5672"
      - "15672:15672"

  core-api:
    build: ./core-api
    ports:
      - "8080:8080"
    depends_on:
      - mongodb
      - rabbitmq

  auth-service:
    build: ./auth-service
    ports:
      - "8083:8083"
    depends_on:
      - mongodb

  scraper:
    build: ./scraper
    ports:
      - "8001:8001"
    depends_on:
      - rabbitmq

  summarizer:
    build: ./summarizer
    ports:
      - "8002:8002"
    depends_on:
      - rabbitmq
    environment:
      GEMINI_API_KEY: ${GEMINI_API_KEY}
```

---

## 📝 Error Responses (Same as Full API)

All endpoints use consistent error format:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Field validation failed",
    "details": { "field": "email", "reason": "invalid_format" },
    "request_id": "req_550e8400-e29b-41d4-a716-446655440000",
    "timestamp": "2025-08-25T10:30:00Z"
  }
}
```

**Common Status Codes:**

- `200` - Success
- `400` - Validation Error
- `401` - Authentication Required
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict (duplicate data)
- `429` - Rate Limited
- `500` - Internal Error

---

This MVP API specification provides **everything needed for the 2-week implementation** while clearly indicating what features will come later. Frontend developers can build against these endpoints knowing they match the actual 3-service architecture.
