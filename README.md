# HypeAtlas API – League of Legends Meta-Game Analysis Platform

HypeAtlas is a comprehensive API platform that provides real-time meta-game analysis for League of Legends, combining live streaming data with Riot Games API integration. The platform offers intelligent insights into champion rotations, league rankings, and strategic recommendations powered by data analysis.

## 🚀 Features

### 🔥 Meta-Game Analysis System
- **Champion Rotation Analysis**: Detailed analysis of weekly free champion rotations with tier classification
- **League Rankings Intelligence**: Statistical analysis of Challenger league data including win rates and LP distribution
- **Strategic Recommendations**: AI-powered suggestions based on free champion availability
- **Real-time Data Integration**: Live data from Riot Games APIs with automatic synchronization

### 🎮 League of Legends Integration
- **Champion-V3 API**: Weekly free champion rotation data
- **League-V4 API**: Challenger league statistics and rankings
- **Data Dragon**: Static game data including champion information
- **Rate Limiting**: Built-in rate limiter (18 req/s, 95 req/min) with automatic retry

### 📊 Advanced Analytics
- **Impact Scoring**: Algorithm to calculate meta impact of champion rotations
- **Probability Calculations**: Meta shift probability based on free champion tiers
- **Role Analysis**: Detection of role concentrations in free champion pools
- **Performance Metrics**: Pick rates, win rates, and ban rate analysis

## 🌐 API Endpoints

### Meta-Game Analysis
- `GET /v1/signal/riot/metagame/rotation/{platform}` - Analyze weekly champion rotation
- `GET /v1/signal/riot/metagame/league/{platform}/{queue}` - Analyze league rankings
- `GET /v1/signal/riot/metagame/report/{platform}` - Generate comprehensive meta report

### Data Synchronization
- `POST /v1/signal/riot/sync/patches` - Synchronize patch data from Riot
- `GET /v1/signal/riot/patches/{version}` - Get detailed patch information

### Live Streaming Data
- `GET /v1/hypemap/live` - Live co-streaming rankings
- `GET /v1/hypemap/summary` - Event summary with aggregated data
- `GET /v1/relay/costreams` - Co-streaming data by event

### Game Data
- `GET /v1/signal/changes` - Patch change history
- `GET /v1/signal/comps` - Champion composition analysis
- `GET /v1/signal/leagues` - League information
- `GET /v1/signal/patches` - Available patches

## 🛠 Installation & Setup

### Prerequisites
- Docker & Docker Compose
- PostgreSQL (via Docker)
- Riot Games API Key

### Quick Start

1. **Clone the repository**
```bash
git clone <repository-url>
cd hypeatlas-api
```

2. **Configure environment**
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. **Set up Riot Games API Key**
```bash
# Add to your .env file
RIOT_API_KEY=RGAPI-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

4. **Start the platform**
```bash
docker-compose up --build
```

The API will be available at `http://localhost:8080`

## 📚 Documentation

### API Documentation
- **Swagger UI**: http://localhost:8080/docs
- **OpenAPI Spec**: http://localhost:8080/openapi.yaml
- **Health Check**: http://localhost:8080/healthz

### Example Requests

#### Champion Rotation Analysis
```bash
curl -s "http://localhost:8080/v1/signal/riot/metagame/rotation/na1" | jq .
```

#### League Rankings Analysis
```bash
curl -s "http://localhost:8080/v1/signal/riot/metagame/league/na1/RANKED_SOLO_5x5" | jq .
```

#### Live Streaming Data
```bash
curl -s "http://localhost:8080/v1/hypemap/live?game=lol&limit=10" | jq .
```

## 🏗 Architecture

### Clean Architecture
- **Domain Layer**: Business logic and entities
- **Infrastructure Layer**: External dependencies (HTTP, Database)
- **Presentation Layer**: HTTP handlers and routing

### Key Components
- **MetaGameService**: Core analysis engine
- **RiotClient**: API client with rate limiting
- **Repository Pattern**: Data access abstraction
- **Docker Containerization**: Production-ready deployment

### Technologies
- **Go 1.21+**: High-performance backend
- **PostgreSQL**: Robust data storage
- **Chi Router**: Lightweight HTTP routing
- **GORM**: ORM with migrations
- **Swagger**: API documentation
- **Docker**: Containerization

## 🔧 Configuration

### Environment Variables
```env
# Database
STORAGE=postgres
POSTGRES_URL=postgres://user:password@db:5432/hypeatlas_dev

# Riot Games API
RIOT_API_KEY=RGAPI-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx

# Server
PORT=8080
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com

# Worker
WORKER_INTERVAL_SEC=30
```

### Riot Games API Key
1. Visit [Riot Developer Portal](https://developer.riotgames.com/)
2. Create a new application
3. Copy the API key to your `.env` file
4. Ensure proper rate limits for your use case

## 🚀 Deployment

### Production Setup
```bash
# Production compose file
docker-compose -f docker-compose.prod.yml up -d

# With custom environment
docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d
```

### Health Monitoring
```bash
# Health check
curl http://localhost:8080/healthz

# Detailed health
curl http://localhost:8080/v1/signal/riot/_health
```

## 📈 Performance & Scaling

### Rate Limiting
- **Riot API**: 18 requests/second, 95 requests/minute
- **Automatic Retry**: Failed requests are retried with backoff
- **Circuit Breaker**: Protection against API outages

### Caching Strategy
- **Static Data**: Champion information cached locally
- **API Responses**: Intelligent caching for frequently requested data
- **Database Indexing**: Optimized queries for real-time performance

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- **Riot Games**: For providing comprehensive League of Legends APIs
- **Data Dragon**: For static game data and assets
- **Open Source Community**: For the amazing Go ecosystem

---

**HypeAtlas** - Transforming League of Legends data into actionable gaming intelligence.
