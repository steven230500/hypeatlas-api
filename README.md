# HypeAtlas API – Esports Analytics & Live Streaming Intelligence Platform

HypeAtlas is a comprehensive API platform that provides real-time analytics for League of Legends and Valorant esports, combining live streaming data with Riot Games API integration. The platform offers intelligent insights into meta-game analysis, event management, and live viewership tracking.

## 🚀 Features

### 🎮 **Signal Module** - Meta-Game Intelligence

- **Champion Data**: Complete champion database with stats, abilities, and images
- **Patch Analysis**: Track changes between game versions
- **Meta-Game Reports**: AI-powered meta analysis with champion rotations
- **League Rankings**: Challenger league statistics and player data
- **Data Dragon Integration**: Full access to Riot's static game data
- **Pagination Support**: Efficient data retrieval with limit/offset

### 🔥 **Relay Module** - Live Streaming & Events

- **Event Management**: Complete CRUD operations for esports events
- **Live Co-Streams**: Track streamers broadcasting esports events
- **HypeMap**: Real-time viewership rankings across platforms
- **Event Filtering**: Filter by game, league, and status (upcoming/live/past)
- **Viewership Analytics**: Aggregate viewer counts and trending events

### 📊 **Advanced Features**

- **Pagination**: Consistent pagination across all list endpoints
- **Real-time Data**: Live updates from Riot APIs and streaming platforms
- **Rate Limiting**: Built-in rate limiter with automatic retry
- **Swagger Documentation**: Auto-generated API docs at `/docs`
- **CI/CD Pipeline**: Automated testing and deployment

---

## 🌐 API Endpoints

### **Signal Module** - Game Data & Meta-Game

#### Champion & Game Data

- `GET /v1/signal/champions?version={version}` - List all champions (paginated)
- `GET /v1/signal/regions` - List available regions (paginated)
- `GET /v1/signal/patches` - List game patches
- `GET /v1/signal/changes` - Patch change history
- `GET /v1/signal/comps` - Champion composition analysis
- `GET /v1/signal/leagues` - League information

#### Riot API Integration

- `GET /v1/signal/riot/versions` - Available game versions (paginated, 480+ versions)
- `GET /v1/signal/riot/champions/{version}/list` - Complete champion list
- `GET /v1/signal/riot/champions/{version}/{championID}` - Champion details
- `GET /v1/signal/riot/items/{version}` - Items data
- `GET /v1/signal/riot/runes/{version}` - Runes data
- `GET /v1/signal/riot/summoner-spells/{version}` - Summoner spells
- `GET /v1/signal/riot/patch-notes/{fromVersion}/{toVersion}` - Compare patches

#### Meta-Game Analysis

- `GET /v1/signal/riot/metagame/rotation/{platform}` - Champion rotation analysis
- `GET /v1/signal/riot/metagame/league/{platform}/{queue}` - League rankings analysis
- `GET /v1/signal/riot/metagame/report/{platform}` - Comprehensive meta report

#### Image URLs

- `GET /v1/signal/riot/images/champions/{version}/{championID}` - Champion images
- `GET /v1/signal/riot/images/items/{version}/{itemID}` - Item images
- `GET /v1/signal/riot/images/spells/{version}/{spellName}` - Spell images
- `GET /v1/signal/riot/images/runes/{runeIcon}` - Rune images

### **Relay Module** - Events & Live Streaming

#### Event Management

- `GET /v1/relay/events/` - List events (filter by game, league, status)
- `GET /v1/relay/events/{slug}` - Get event details
- `POST /v1/relay/events/` - Create event
- `PUT /v1/relay/events/{slug}` - Update event
- `DELETE /v1/relay/events/{slug}` - Delete event

#### Live Streaming

- `GET /v1/relay/costreams?event_id={slug}&lang={lang}` - Co-streams for an event
- `GET /v1/hypemap/live?game={game}&lang={lang}` - Live viewership rankings
- `GET /v1/hypemap/summary?game={game}&lang={lang}` - Event summary with totals

---

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
# Edit .env with your Riot API key
```

3. **Start the platform**

```bash
docker-compose up --build
```

The API will be available at `http://localhost:8080`

---

## 📚 API Documentation

- **Swagger UI**: https://api.hypeatlas.app/docs
- **OpenAPI Spec**: https://api.hypeatlas.app/openapi.json
- **Health Check**: https://api.hypeatlas.app/healthz

---

## 💡 Example Requests

### Pagination Examples

```bash
# Get first 20 game versions
curl "https://api.hypeatlas.app/v1/signal/riot/versions?limit=20&offset=0"

# Get champions with pagination
curl "https://api.hypeatlas.app/v1/signal/champions?version=14.14.1&limit=10"
```

### Event Management

```bash
# List all events
curl "https://api.hypeatlas.app/v1/relay/events/"

# Filter live events
curl "https://api.hypeatlas.app/v1/relay/events/?status=live&game=val"

# Get specific event
curl "https://api.hypeatlas.app/v1/relay/events/vct-emea-final"
```

### Live Streaming Data

```bash
# Live viewership rankings
curl "https://api.hypeatlas.app/v1/hypemap/live?game=val&limit=10"

# Event summary
curl "https://api.hypeatlas.app/v1/hypemap/summary?game=lol"
```

### Champion Data

```bash
# Get all champions for a version
curl "https://api.hypeatlas.app/v1/signal/champions?version=14.14.1"

# Get champion images
curl "https://api.hypeatlas.app/v1/signal/riot/images/champions/14.14.1/Ahri"
```

---

## 🏗 Architecture

### Modules

- **Signal**: Game data, meta-game analysis, Riot API integration
- **Relay**: Event management, live streaming, viewership tracking
- **Shared**: Common utilities (pagination, HTTP helpers, database)

### Technologies

- **Go 1.22+**: High-performance backend
- **PostgreSQL**: Robust data storage with GORM
- **Chi Router**: Lightweight HTTP routing
- **Swagger**: Auto-generated API documentation
- **Docker**: Production-ready containerization
- **GitHub Actions**: Automated CI/CD pipeline

### Clean Architecture

- **Domain Layer**: Business logic and entities
- **Infrastructure Layer**: HTTP handlers, repositories
- **Ports & Adapters**: Hexagonal architecture pattern

---

## 🚀 Deployment

### Production

```bash
docker-compose -f docker-compose.prod.yml up -d
```

### CI/CD Pipeline

- **go-ci.yml**: Build and test on push/PR
- **docker-publish.yml**: Build and push Docker images
- **deploy.yml**: Auto-deploy to production on main branch

---

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
```

---

## 📊 API Response Format

### Paginated Response

```json
{
  "items": [...],
  "pagination": {
    "limit": 20,
    "offset": 0,
    "total": 480,
    "hasMore": true
  }
}
```

### Event Response

```json
{
  "uuid": "...",
  "slug": "vct-emea-final",
  "title": "VCT EMEA Final",
  "game": "val",
  "league": "VCT EMEA",
  "starts_at": "2025-11-23T10:00:00Z",
  "ends_at": "2025-11-23T18:00:00Z"
}
```

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

---

## 📄 License

MIT License - see LICENSE file for details.

---

## 🙏 Acknowledgments

- **Riot Games**: For comprehensive League of Legends APIs
- **Data Dragon**: For static game data and image assets
- **Open Source Community**: For the amazing Go ecosystem

---

**HypeAtlas** - Transforming esports data into actionable intelligence.
