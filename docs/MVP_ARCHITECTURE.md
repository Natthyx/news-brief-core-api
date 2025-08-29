# NewsBrief MVP Architecture (2 Weeks)

## Updated Simplified Service Layout

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Core API      │    │   Scraper        │    │  Summarizer     │
│   (Go Gin)      │◄──►│   (FastAPI)      │◄──►│  (Go HTTP)      │
│   • Auth        │    │   • curl_cffi    │    │  • Gemini API   │
│   • Stories     │    │   • BeautifulSoup│    │  • Short/Medium │
│   • Feed mgmt   │    │   • Content Hash │    │    Summaries    │
│   • User Subs   │    │   • Basic Vector │    │  • Simple HTTP  │
│   • Topics      │    │     Storage      │    └─────────────────┘
│   • Sources     │    └──────────────────┘            ▲
│   • CRON JOBS   │            ▲                        │
└─────────────────┘            │                        │
         │                     │                        │
         ▼                     │                        │
┌─────────────────────────────────────────────────────────────────┐
│                     MongoDB 7.0                                │
│  • users (with embedded preferences & subscriptions)           │
│  • tokens (unified: refresh_token|email_verify|password_reset) │
│  • stories (with short/medium summaries & topic images)       │
│  • sources (news outlets with logos & reliability scores)     │
│  • topics (with images & bilingual descriptions)              │
│  • daily_briefs (AI-curated morning/evening collections)     │
│  • chat_queries (for future chatbot - simple storage)        │
│  • jobs (background task tracking)                            │
└─────────────────────────────────────────────────────────────────┘
```

## MVP Feature Scope (Essential for 2 Weeks)

### ✅ **Week 1 Priority Features**

1. **Enhanced User Management**

   - Registration/Login with embedded preferences
   - News outlet subscriptions (max 5 sources for MVP)
   - Topic selection with images
   - Password change functionality

2. **Core Content System**

   - Short & medium summary generation
   - Topic images and source logos
   - Basic story processing pipeline
   - MongoDB text search

3. **Feed System**
   - Personalized feed based on subscriptions
   - Topic filtering with visual elements
   - Source filtering
   - Brief type selection (short/medium)

### ✅ **Week 2 Priority Features**

4. **News Source Management**

   - Browse available sources with logos
   - Subscribe/unsubscribe to outlets
   - Topic preferences per source
   - Basic source reliability scoring

5. **Enhanced Content Processing**

   - Dual summary generation (short + medium)
   - Content hashing for deduplication
   - Processing status tracking
   - Reading time estimation

6. **Basic Analytics**
   - User reading stats
   - Subscription analytics
   - Simple recommendations

### 🚫 **Deferred for Post-MVP (Month 2-3)**

- **Advanced Chatbot**: Web search + vector database integration
- **Full Vector Database**: Semantic search and content caching
- **Advanced Analytics**: Detailed user behavior analysis
- **Premium Features**: Advanced subscriptions and billing
- **Real-time Notifications**: Push notification system
- **Data Export**: GDPR compliance features
- **Advanced Search**: Full-text search with highlighting

## MVP Simplifications

### 1. Background Jobs → Enhanced Cron Jobs in Core API

```go
// Enhanced background processing for MVP
func main() {
    // Start web server
    go startWebServer()

    // Start background tasks
    go startFeedIngestion()     // Every 15 minutes
    go startSummaryGeneration() // Process pending stories
    go startBriefCuration()     // Daily at 6AM, 6PM
    go startAnalyticsUpdate()   // Update user stats hourly

    select {} // Keep alive
}

func startFeedIngestion() {
    ticker := time.NewTicker(15 * time.Minute)
    for range ticker.C {
        // Enhanced ingestion with subscriptions
        ingestSubscribedFeeds()
    }
}

func ingestSubscribedFeeds() {
    // Get all active source subscriptions
    sources := getActiveSubscriptions()

    for _, source := range sources {
        // Fetch RSS → Scrape → Generate short/medium summaries
        articles := fetchRSSFeed(source.FeedURL)
        for _, article := range articles {
            if !articleExists(article.URL) {
                content := scrapeArticle(article.URL)
                summaries := generateBothSummaries(content)
                storeStoryWithSummaries(article, content, summaries)
            }
        }
    }
}
```

### 2. Enhanced MongoDB Schema for MVP

```javascript
// Enhanced collections for MVP features
db.users.insertOne({
  _id: ObjectId("507f1f77bcf86cd799439011"),
  email: "user@example.com",
  password_hash: "bcrypt_hash",
  preferences: {
    lang: "am",
    topics: ["economy", "agriculture", "politics"],
    subscribed_sources: ["addisstandard", "ethiopianherald"],
    brief_type: "short", // "short" | "medium"
    data_saver: true,
  },
  subscription: {
    plan: "free",
    source_limit: 5,
    features: ["basic_summaries"],
  },
  stats: {
    stories_read: 0,
    last_active: new Date(),
  },
});

db.stories.insertOne({
  _id: ObjectId("507f1f77bcf86cd799439012"),
  title: "Ethiopia launches agricultural initiative",
  url: "https://addisstandard.com/news/...",
  source: {
    key: "addisstandard",
    name: "Addis Standard",
    logo_url: "https://cdn.newsbrief.et/sources/addisstandard.png",
  },
  content_hash: "sha256:abc123...",
  summary_short: "Government announces $50M farming investment.",
  summary_bullets: [
    "Government announces $50M investment in rural farming",
    "Program targets 100,000 smallholder farmers nationwide",
    "Focus on drought-resistant crop varieties and irrigation",
  ],
  topic_tags: ["agriculture", "economy"],
  topic_image: "https://cdn.newsbrief.et/topics/agriculture.jpg",
  processing_status: "completed",
  reading_time: { short: 1, medium: 3 },
  published_at: new Date(),
  processed_at: new Date(),
});

db.sources.insertOne({
  _id: ObjectId("507f1f77bcf86cd799439013"),
  key: "addisstandard",
  name: "Addis Standard",
  description: "Independent news outlet covering Ethiopian politics",
  url: "https://addisstandard.com",
  logo_url: "https://cdn.newsbrief.et/sources/addisstandard.png",
  rss_feeds: ["https://addisstandard.com/feed/"],
  languages: ["en", "am"],
  topics: ["politics", "economy", "society"],
  reliability_score: 0.92,
  update_frequency: "hourly",
  active: true,
});

db.topics.insertOne({
  _id: ObjectId("507f1f77bcf86cd799439014"),
  key: "agriculture",
  label: { en: "Agriculture", am: "ግብርና" },
  description: {
    en: "Farming, livestock, and agricultural development",
    am: "የግብርና፣ የእንስሳት ሀብት እና የግብርና ልማት",
  },
  image_url: "https://cdn.newsbrief.et/topics/agriculture.jpg",
  story_count: 0,
});
```

### 3. Simplified Summarizer for MVP

- **Dual Summary Generation**: Both short (1 sentence) and medium (3-5 bullets)
- **Single HTTP endpoint**: `POST /summarize`
- **Simple retry logic**: 3 attempts with exponential backoff
- **No batch processing**: Process stories one by one for MVP

```go
// Summarizer service - enhanced for MVP
type SummaryRequest struct {
    Text         string   `json:"text"`
    Title        string   `json:"title"`
    Source       string   `json:"source"`
    TargetLang   string   `json:"target_lang"`
    SummaryTypes []string `json:"summary_types"` // ["short", "medium"]
    MaxBullets   int      `json:"max_bullets"`
}

type SummaryResponse struct {
    SummaryShort   string   `json:"summary_short"`
    SummaryBullets []string `json:"summary_bullets"`
    SummaryLang    string   `json:"summary_lang"`
    ConfidenceScore float64 `json:"confidence_score"`
    ReadingTime    struct {
        Short  int `json:"short"`
        Medium int `json:"medium"`
    } `json:"reading_time"`
    ProcessingTimeMs int `json:"processing_time_ms"`
}
```

### 4. Enhanced Development Timeline

**Week 1 (Days 1-7):**

**Days 1-2: Enhanced User System**

- Core API with enhanced auth endpoints
- MongoDB setup with new schema (users, sources, topics, stories)
- User subscription management (subscribe/unsubscribe to sources)
- Topic selection with image URLs
- Password change functionality

**Days 3-4: Content Processing Pipeline**

- Enhanced story model with dual summaries (short + medium)
- Content hashing for deduplication
- Processing status tracking
- Topic image integration
- Source logo integration

**Days 5-7: Feed System Enhancement**

- Subscription-based feed filtering
- Topic-based filtering with images
- Source-based filtering
- Brief type selection (short vs medium)
- Basic reading time calculation

**Week 2 (Days 8-14):**

**Days 8-9: Source Management**

- Sources API with logos and metadata
- Source discovery and search
- Subscription limits (free: 5 sources)
- Reliability scoring display
- RSS feed management

**Days 10-11: Summarizer Enhancement**

- Dual summary generation (short + medium)
- Reading time estimation
- Enhanced Gemini prompts for consistency
- Basic content quality scoring

**Days 12-14: Integration & Analytics**

- User analytics (reading stats, preferences)
- Basic recommendations (new sources/topics)
- Performance optimization
- End-to-end testing
- Basic monitoring and health checks

## Enhanced MVP Deployment

```yaml
# docker-compose.yml for enhanced MVP
version: "3.8"
services:
  mongodb:
    image: mongo:7.0
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: password
      MONGO_INITDB_DATABASE: newsbrief
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
      - ./init-scripts:/docker-entrypoint-initdb.d

  core-api:
    build: ./core-api
    ports:
      - "8080:8080"
    environment:
      MONGODB_URI: mongodb://admin:password@mongodb:27017/newsbrief?authSource=admin
      GEMINI_API_KEY: ${GEMINI_API_KEY}
      CDN_BASE_URL: https://cdn.newsbrief.et
      JWT_SECRET: ${JWT_SECRET}
    depends_on:
      - mongodb
    volumes:
      - ./static:/app/static # For topic images and source logos

  scraper:
    build: ./scraper
    ports:
      - "8001:8001"
    environment:
      STORAGE_TYPE: local # Simple file storage for MVP
      CONTENT_HASH_ENABLED: true

  summarizer:
    build: ./summarizer
    ports:
      - "8002:8002"
    environment:
      GEMINI_API_KEY: ${GEMINI_API_KEY}
      SUMMARY_TYPES: "short,medium"
      MAX_RETRIES: 3

  # Simple CDN service for images (MVP only)
  cdn:
    image: nginx:alpine
    ports:
      - "8090:80"
    volumes:
      - ./static:/usr/share/nginx/html
    command: ["nginx", "-g", "daemon off;"]

volumes:
  mongodb_data:
```

## What This MVP Includes vs Removes

### ✅ **MVP Features (2 Weeks)**

**Essential User Features:**

- Enhanced user registration/login with embedded preferences
- News outlet subscription management (limit: 5 sources for free tier)
- Topic selection with visual images
- Password management
- Personalized feed with dual summary types (short/medium)
- Source and topic filtering with visual elements
- Basic user analytics and reading stats

**Content Processing:**

- Dual summary generation (short + medium)
- Content deduplication via hashing
- Topic images and source logos
- Processing status tracking
- Reading time estimation
- Basic content quality scoring

**Technical Features:**

- Unified token system (refresh, verify, reset)
- MongoDB with enhanced schema
- Clean REST API with comprehensive error handling
- Basic background job processing
- Simple CDN for static assets

### 🚫 **Deferred Features (Post-MVP)**

**Advanced Features:**

- AI-powered chatbot with web search
- Vector database for semantic search
- Advanced analytics and ML recommendations
- Real-time push notifications
- GDPR data export functionality
- Premium subscription billing
- Advanced search with highlighting
- Horizontal scaling and load balancing

**Complex Infrastructure:**

- Redis for job queues and caching
- Sophisticated retry mechanisms
- Real-time job monitoring
- Advanced security features
- Full-text search indexing
- Message queues (RabbitMQ/Kafka)

## Post-MVP Migration Path (Month 2-3)

When you need advanced features:

1. **Add Vector Database**: Implement semantic search and content caching
2. **Add Chatbot**: Integrate web search + AI query processing
3. **Add Redis**: Job queues and advanced caching
4. **Add Real-time Features**: WebSocket connections and push notifications
5. **Add Advanced Analytics**: ML-based recommendations and behavior tracking
6. **Add Premium Features**: Billing, advanced subscriptions, and feature gates

## Estimated Development Complexity

**MVP (2 Weeks):**

- **Complexity**: Medium
- **Features**: 80% of core user value
- **Technical Debt**: Minimal
- **Scalability**: Handles 1K-10K users

**Full Version (2 Months):**

- **Complexity**: High
- **Features**: 100% of planned features
- **Technical Debt**: Production-ready
- **Scalability**: Handles 100K+ users

## Key MVP Success Metrics

**Technical Metrics:**

- [ ] All core API endpoints functional
- [ ] Dual summary generation working
- [ ] User subscription system operational
- [ ] MongoDB schema optimized for performance
- [ ] Background job processing stable

**User Experience Metrics:**

- [ ] User can register and manage subscriptions
- [ ] Personalized feed shows relevant content
- [ ] Topic and source filtering works smoothly
- [ ] Short/medium summaries provide value
- [ ] Basic analytics show user engagement

**Performance Targets (MVP):**

- API response time: < 500ms for 95% requests
- Summary generation: < 5s per article
- Feed loading: < 1s for 20 articles
- User registration: < 2s end-to-end
- Database queries: All use proper indexes

This enhanced MVP architecture provides a solid foundation with the new features while maintaining the 2-week timeline through strategic feature prioritization and technical simplifications.
