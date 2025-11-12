# Discord-Teneo AI Assistant Bot

A Discord bot powered by Google Gemini that connects to the Teneo Agent Network. This bot responds to mentions and commands in your Discord server while also being accessible through the Teneo network.

## Features

✅ **Dual Connectivity**: Works in Discord AND on Teneo network  
✅ **AI-Powered**: Uses Google Gemini 2.5 Flash for intelligent responses  
✅ **Flexible Interaction**: Responds to @mentions or `!ask` commands  
✅ **Message Splitting**: Automatically handles long responses  
✅ **Teneo Integration**: Full agent capabilities with NFT identity  
✅ **Production Ready**: Health monitoring, reconnection, error handling

## Prerequisites

1. **Go 1.24+** installed
2. **Discord Bot Token** - [Create a Discord Application](https://discord.com/developers/applications)
3. **Google Gemini API Key** - [Get from Google AI Studio](https://aistudio.google.com/app/apikey)
4. **Ethereum Wallet** - Private key for Teneo network authentication
5. **Teneo Agent NFT** - Mint via [deploy.teneo-protocol.ai](https://deploy.teneo-protocol.ai)

## Discord Bot Setup

### 1. Create Discord Application

1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Click "New Application" and give it a name
3. Go to "Bot" section and click "Add Bot"
4. Under "Token", click "Reset Token" and copy it (you'll need this for `.env`)
5. Enable these Privileged Gateway Intents:
   - ✅ MESSAGE CONTENT INTENT
   - ✅ SERVER MEMBERS INTENT (optional)

### 2. Invite Bot to Your Server

1. Go to "OAuth2" → "URL Generator"
2. Select scopes:
   - ✅ `bot`
   - ✅ `applications.commands`
3. Select bot permissions:
   - ✅ Read Messages/View Channels
   - ✅ Send Messages
   - ✅ Read Message History
   - ✅ Add Reactions (optional)
4. Copy the generated URL and open it in browser
5. Select your server and authorize

## Installation

### 1. Clone and Setup

```bash
cd /home/kiter/layla/discord-teneo-bot

# Copy environment template
cp .env.example .env

# Edit .env with your credentials
nano .env
```

### 2. Configure Environment Variables

Edit `.env` file with your actual values:

```bash
# Discord Bot Token (from Discord Developer Portal)
DISCORD_TOKEN=your_discord_bot_token_here

# Google Gemini API Key (from Google AI Studio)
GEMINI_API_KEY=your_gemini_api_key_here

# Teneo Agent Configuration
PRIVATE_KEY=your_ethereum_private_key_without_0x
OWNER_ADDRESS=0xYourWalletAddress
NFT_TOKEN_ID=your_nft_token_id

# Optional: Custom AI personality
SYSTEM_PROMPT=You are a helpful AI assistant in a Discord server. Be friendly, concise, and helpful.
```

### 3. Install Dependencies

```bash
go mod tidy
```

## Running the Bot

### Development Mode

```bash
go run main.go
```

### Production Mode

```bash
# Build the binary
go build -o discord-teneo-bot

# Run it
./discord-teneo-bot
```

### Run as Background Service

```bash
# Using nohup
nohup ./discord-teneo-bot > bot.log 2>&1 &

# Or using screen
screen -S discord-bot
./discord-teneo-bot
# Press Ctrl+A then D to detach
```

## Usage

### In Discord

The bot responds to two interaction methods:

#### 1. Mention the Bot
```
@YourBot What is the meaning of life?
```

#### 2. Use the !ask Command
```
!ask Explain quantum computing
```

#### Examples
```
@YourBot Tell me a joke
!ask What's the weather like in Tokyo?
@YourBot Write a Python function to reverse a string
!ask Explain Docker in simple terms
```

### On Teneo Network

Your bot is also accessible as a Teneo agent! Other agents on the network can send tasks to it using commands:
- `ask [question]` - Ask any question
- `explain [concept]` - Get detailed explanations
- `help` - Show available commands

## Bot Output

When running, you'll see:

```
🤖 Discord-Teneo Bot 'YourBotName' is now running!
📱 Discord: Connected and listening for messages
💬 Usage: Mention the bot or use !ask <question>
🌐 Starting Teneo agent...
✅ Teneo agent: Connected to network
Press CTRL+C to stop the bot
```

When processing messages:
```
💬 Discord Message from username: What is AI?
✅ Response sent to Discord
📥 Teneo Task Received: Explain blockchain
📤 Teneo Response: Blockchain is a distributed ledger...
```

## Customization

### Change AI Model

Edit `main.go` around line 289:
```go
geminiModel := geminiClient.GenerativeModel("gemini-2.5-flash")
// Other options: "gemini-1.5-pro", "gemini-2.0-flash-exp"
```

### Adjust Response Length

Edit `main.go` around line 291:
```go
geminiModel.SetMaxOutputTokens(1000)  // Increase for longer responses
```

### Modify System Prompt

Either edit `.env`:
```bash
SYSTEM_PROMPT=You are a helpful AI assistant who specializes in blockchain...
```

Or edit `main.go` directly around line 340.

### Change Command Prefix

Edit `main.go` to change from `!ask` to something else:
```go
if !mentioned && !strings.HasPrefix(m.Content, "!mycommand") {
```

## Troubleshooting

### Bot doesn't respond in Discord

1. Check MESSAGE CONTENT INTENT is enabled in Discord Developer Portal
2. Verify bot has permissions in the channel (can read and send messages)
3. Check logs for errors: `go run main.go`

### "Required environment variable not set"

Ensure `.env` file exists and all required variables are filled:
```bash
cat .env
```

### Gemini API errors

- Verify API key is correct from [Google AI Studio](https://aistudio.google.com/app/apikey)
- Check your usage quota at [Google AI Studio](https://aistudio.google.com/)
- Free tier has rate limits: 10 requests per minute for gemini-2.5-flash
- For production use, consider upgrading to a paid plan

### Teneo connection issues

- Verify `PRIVATE_KEY`, `OWNER_ADDRESS`, and `NFT_TOKEN_ID` are correct
- Ensure your NFT is minted (check [deploy.teneo-protocol.ai](https://deploy.teneo-protocol.ai))
- Check network connectivity

### Bot crashes or disconnects

The bot has automatic reconnection for Discord and Teneo. If it still crashes:
1. Check logs for specific error messages
2. Ensure all dependencies are installed: `go mod tidy`
3. Verify Go version: `go version` (needs 1.24+)

## Project Structure

```
discord-teneo-bot/
├── main.go           # Main bot implementation
├── go.mod           # Go module dependencies
├── .env             # Environment variables (create from .env.example)
├── .env.example     # Environment template
└── README.md        # This file
```

## Dependencies

- `github.com/bwmarrin/discordgo` - Discord API wrapper
- `github.com/google/generative-ai-go` - Google Gemini API client
- `github.com/joho/godotenv` - Environment variable loader
- `github.com/TeneoProtocolAI/teneo-agent-sdk` - Teneo agent framework

## About the Teneo Agent SDK

This project depends on the `teneo-agent-sdk`, which has been added as a git submodule. The SDK is located at `/teneo-agent-sdk` in this repository.

### Cloning this repository

When cloning this repository, make sure to initialize the submodules:

```bash
git clone https://github.com/PhiBao/layla.git
cd layla
git submodule update --init --recursive
```

### Working with the SDK

The SDK is included via git submodule, which keeps it synchronized with the upstream repository. To update the SDK to the latest version:

```bash
cd teneo-agent-sdk
git pull origin main
cd ..
git add teneo-agent-sdk
git commit -m "Update teneo-agent-sdk submodule"
```

### Local Development

For local development, the `go.mod` file uses a replace directive to point to the local SDK:

```go
// inside discord-teneo-bot/go.mod
replace github.com/TeneoProtocolAI/teneo-agent-sdk => ../teneo-agent-sdk
```

This allows you to make changes to both projects simultaneously. For production builds or CI, ensure the submodule is initialized.

## Security Notes

⚠️ **Never commit your `.env` file to version control!**

Add to `.gitignore`:
```bash
echo ".env" >> .gitignore
```

- Keep your Discord token, OpenAI key, and private key secure
- Use environment variables in production
- Rotate keys if they're exposed
- Limit bot permissions to only what's needed

## Advanced Features

### Rate Limiting

The Teneo agent has built-in rate limiting. Configure in `.env`:
```bash
RATE_LIMIT_PER_MINUTE=60  # Limit tasks from Teneo network
```

### Health Monitoring

The Teneo agent includes health endpoints (check SDK docs for details).

### Multiple Servers

The same bot can work in multiple Discord servers simultaneously!

## Support

- **Teneo SDK Docs**: [teneo-agent-sdk/README.md](../teneo-agent-sdk/README.md)
- **Discord Developer Docs**: https://discord.com/developers/docs
- **Google Gemini Docs**: https://ai.google.dev/gemini-api/docs

## License

Same as Teneo Agent SDK license.

---

**Enjoy your Discord-Teneo AI Assistant! 🤖🚀**
