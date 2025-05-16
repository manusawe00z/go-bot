# go-bot

Discord Voice Bot with TTS (Text-to-Speech)

## Features

- Join Discord voice channels
- Convert text messages to speech using gTTS
- Play TTS audio in voice channels
- Support for multiple languages
- Message queue system
- Structured logging
- Health check system for container orchestration

## Quick Start with Docker Compose

1. Clone the repository:
```bash
git clone https://github.com/manusawe00z/go-bot
cd go-bot
```

2. Create and configure your environment file:
```bash
cp env.example .env
# Edit .env with your Discord token and preferences
```

3. Start the bot with Docker Compose:
```bash
docker-compose up -d
```

4. View logs:
```bash
docker-compose logs -f
```

5. To stop the bot:
```bash
docker-compose down
```

## Docker Compose Configuration

The project uses a single `docker-compose.yml` file with profiles for different environments

## Configuration Options

Main environment variables (set in .env file):

- `DISCORD_TOKEN` - Your Discord bot token (required)
- `COMMAND_PREFIX` - Prefix for commands (default: !)
- `TTS_LANGUAGE` - Default language for TTS (default: en)
- `QUEUE_SIZE` - Maximum size of message queue (default: 100)
- `AUDIO_QUALITY` - Audio quality (low/medium/high) (default: medium)
- `LOG_LEVEL` - Log level (0=debug to 4=fatal) (default: 1)

## Health Check System

The bot includes an HTTP server that provides a health check endpoint at:

```http
http://localhost:8080/health
```

This endpoint returns:

- HTTP 200 status when the bot is running properly
- JSON response with status, timestamp, and version information

The health check is configured in Docker Compose file to support container orchestration:

- Development profile: Every 15 seconds
- Production profile: Every 30 seconds

## Requirements (for manual installation)

- Go 1.22+
- Python 3.11+
- gtts Python package
- ffmpeg

## License

MIT
