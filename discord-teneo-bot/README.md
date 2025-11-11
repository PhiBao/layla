# Discord-Teneo AI Assistant Bot

A Discord bot powered by OpenAI that connects to the Teneo Agent Network. This bot responds to mentions and commands in your Discord server while also being accessible through the Teneo network.

## Features

✅ **Dual Connectivity**: Works in Discord AND on Teneo network  
✅ **AI-Powered**: Uses OpenAI GPT-4 for intelligent responses  
✅ **Flexible Interaction**: Responds to @mentions or `!ask` commands  
✅ **Message Splitting**: Automatically handles long responses  
✅ **Teneo Integration**: Full agent capabilities with NFT identity  
✅ **Production Ready**: Health monitoring, reconnection, error handling

## Prerequisites

1. **Go 1.24+** installed
2. **Discord Bot Token** - [Create a Discord Application](https://discord.com/developers/applications)
3. **OpenAI API Key** - [Get from OpenAI](https://platform.openai.com/api-keys)
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

# OpenAI API Key
OPENAI_API_KEY=sk-proj-xxxxxxxxxxxxx

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

Your bot is also accessible as a Teneo agent! Other agents on the network can send tasks to it, and it will process them through the same OpenAI logic.

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

Edit `main.go` line 47:
```go
Model: openai.GPT4TurboPreview,  // Change to GPT4o, GPT35Turbo, etc.
```

### Adjust Response Length

Edit `main.go` line 49:
```go
MaxTokens: 1000,  // Increase for longer responses
```

### Modify System Prompt

Either edit `.env`:
```bash
SYSTEM_PROMPT=You are a sarcastic AI assistant who loves memes...
```

Or edit `main.go` directly (around line 186).

### Change Command Prefix

Edit `main.go` line 119 to change from `!ask` to something else:
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

### OpenAI API errors

- Verify API key is correct
- Check you have credits/quota remaining
- Ensure API key has proper permissions

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
- `github.com/sashabaranov/go-openai` - OpenAI API client
- `github.com/joho/godotenv` - Environment variable loader
- `github.com/TeneoProtocolAI/teneo-agent-sdk` - Teneo agent framework

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

- **Teneo SDK Docs**: `/home/kiter/layla/teneo-agent-sdk/README.md`
- **Discord.js Docs**: https://discord.com/developers/docs
- **OpenAI Docs**: https://platform.openai.com/docs

## License

Same as Teneo Agent SDK license.

---

**Enjoy your Discord-Teneo AI Assistant! 🤖🚀**
