# go-bot

Discord Voice Bot with TTS (Text-to-Speech)

## Features
- Join Discord voice channels
- Convert text messages to speech using gTTS (Google Text-to-Speech)
- Play TTS audio in voice channels

## Requirements
- Go 1.22+
- Python 3.11+
- [gtts](https://pypi.org/project/gTTS/) Python package
- ffmpeg
- Discord Bot Token

## Getting Started

### 1. Clone the repository
```sh
git clone https://github.com/manusawe00z/go-bot
cd go-bot
```

### 2. Set up environment variables
Create a `.env` file or set the environment variable `DISCORD_TOKEN` with your Discord bot token.

### 3. Install Python dependencies
```sh
pip install -r requirements.txt
```

### 4. Build the Go application
```sh
go build -o app ./cmd/main.go
```

### 5. Run the bot (local)
```sh
DISCORD_TOKEN=your_token_here ./app
```

### 6. Run with Docker
Build and run the container:
```sh
docker build -t go-bot .
docker run -e DISCORD_TOKEN=your_token_here go-bot
```

## Deploy to Railway
1. Push your code to GitHub
2. Connect your repo to [Railway](https://railway.app)
3. Set the `DISCORD_TOKEN` variable in Railway's dashboard
4. Deploy (Railway will use the provided Dockerfile)

## File Structure
```
├── cmd/main.go           # Entry point
├── internal/bot/         # Bot logic
│   ├── bot.go
│   ├── voice.go
│   └── ...
├── internal/tts/tts.py   # TTS Python script
├── requirements.txt      # Python dependencies
├── Dockerfile            # Multi-stage build for Go + Python + ffmpeg
├── go.mod, go.sum        # Go dependencies
```

## Notes
- Ensure your bot has permission to join and speak in voice channels.
- ffmpeg is required for audio processing.
- gTTS requires internet access to generate speech.

## License
MIT
