# Layla - Discord-Teneo AI Assistant

A comprehensive Discord bot powered by Google Gemini AI that integrates with the Teneo Agent Network. This repository contains both the Discord bot implementation and the Teneo Agent SDK as a submodule.

## 📁 Repository Structure

```
layla/
├── discord-teneo-bot/       # Main Discord bot application
│   ├── main.go             # Bot implementation
│   ├── go.mod              # Go dependencies
│   ├── .env.example        # Environment template
│   ├── layla-agent.json    # Agent configuration
│   └── README.md           # Detailed bot documentation
│
└── teneo-agent-sdk/        # Teneo Agent SDK (submodule)
    ├── pkg/                # SDK packages
    ├── examples/           # Example implementations
    └── docs/               # SDK documentation
```

## 🚀 Quick Start

### Prerequisites

- **Go 1.24+** installed
- **Discord Bot Token** from [Discord Developer Portal](https://discord.com/developers/applications)
- **Google Gemini API Key** from [Google AI Studio](https://makersuite.google.com/app/apikey)
- **Ethereum Wallet** with private key
- **Teneo Agent NFT** minted via [deploy.teneo-protocol.ai](https://deploy.teneo-protocol.ai)

### Setup

1. **Clone the repository with submodules:**

```bash
git clone --recursive https://github.com/PhiBao/layla.git
cd layla

# If you already cloned without --recursive, initialize submodules:
git submodule update --init --recursive
```

2. **Configure the bot:**

```bash
cd discord-teneo-bot

# Copy environment template
cp .env.example .env

# Edit with your credentials
nano .env
```

3. **Required environment variables:**

```bash
# Discord Bot Token
DISCORD_TOKEN=your_discord_bot_token_here

# Google Gemini API Key
GEMINI_API_KEY=your_gemini_api_key_here

# Teneo Agent Configuration
PRIVATE_KEY=your_ethereum_private_key_without_0x
OWNER_ADDRESS=0xYourWalletAddress
NFT_TOKEN_ID=your_nft_token_id

# Optional: Custom AI personality
SYSTEM_PROMPT=You are Layla, a helpful AI assistant...
```

4. **Install dependencies and run:**

```bash
go mod tidy
go run main.go
```

## 🎯 Features

### Discord Integration
- ✅ Responds to @mentions
- ✅ Supports `!ask` command
- ✅ Intelligent responses using Google Gemini AI
- ✅ Automatic message splitting for long responses
- ✅ Works across multiple Discord servers

### Teneo Network Integration
- ✅ Full agent capabilities with NFT identity
- ✅ Accessible through Teneo chatroom
- ✅ Supports commands: `ask`, `explain`, `help`
- ✅ Health monitoring and automatic reconnection
- ✅ Rate limiting to prevent API quota issues

### AI Capabilities
- 🤖 Powered by Google Gemini 2.5 Flash
- 💬 Natural language understanding
- 🧠 Context-aware conversations
- 📚 Question answering and concept explanation
- 🎨 Customizable personality via system prompts

## 📖 Usage

### In Discord

**Mention the bot:**
```
@Layla What is blockchain?
```

**Use the !ask command:**
```
!ask Explain quantum computing
```

### In Teneo Chatroom

**Ask questions:**
```
@layla-discord-ai ask what is DeFi
```

**Get explanations:**
```
@layla-discord-ai explain how privacy works in crypto
```

**Get help:**
```
@layla-discord-ai help
```

## 🛠️ Development

### Project Components

1. **discord-teneo-bot**: Main application
   - Handles Discord interactions
   - Manages Gemini AI integration
   - Implements Teneo agent protocol
   - See [discord-teneo-bot/README.md](discord-teneo-bot/README.md) for details

2. **teneo-agent-sdk**: Agent framework (submodule)
   - Network client and protocol handlers
   - NFT integration and authentication
   - Health monitoring and circuit breakers
   - See [teneo-agent-sdk/README.md](teneo-agent-sdk/README.md) for SDK docs

### Building for Production

```bash
cd discord-teneo-bot
go build -o discord-teneo-bot
./discord-teneo-bot
```

### Running as a Service

```bash
# Using nohup
nohup ./discord-teneo-bot > bot.log 2>&1 &

# Using screen
screen -S layla-bot
./discord-teneo-bot
# Press Ctrl+A then D to detach
```

## 🔧 Configuration

### Customizing AI Behavior

Edit `.env` to change the system prompt:
```bash
SYSTEM_PROMPT=You are a sarcastic AI assistant who loves crypto memes...
```

Or modify the model in `main.go`:
```go
geminiModel := geminiClient.GenerativeModel("gemini-1.5-pro")
geminiModel.SetTemperature(0.7)
geminiModel.SetMaxOutputTokens(1000)
```

### Agent Metadata

Edit `layla-agent.json` to customize:
- Agent name and description
- Available commands
- Capabilities
- NLP fallback behavior

## 📊 Monitoring

The bot includes:
- Health check endpoints (port 8080)
- Automatic reconnection for Discord and Teneo
- Rate limiting for Gemini API
- Detailed logging for debugging

Check logs:
```bash
tail -f bot.log
```

## 🐛 Troubleshooting

### Submodule Issues

If `teneo-agent-sdk` folder is empty:
```bash
git submodule update --init --recursive
```

### Build Errors

If you see "replacement directory does not exist":
```bash
cd discord-teneo-bot
# Remove any replace directives in go.mod, then:
go mod tidy
```

### Agent Shows Offline in Chatroom

1. Ensure `NFT_TOKEN_ID` is set correctly in `.env`
2. Check that the NFT is minted and belongs to your wallet
3. Verify `PRIVATE_KEY` and `OWNER_ADDRESS` match
4. Look for authentication errors in logs

### Gemini API Issues

- Check your API key is valid
- Monitor usage at [Google AI Studio](https://makersuite.google.com/app/apikey)
- Rate limiter prevents quota exceeded errors (2 requests/min)
- If timeouts occur, increase timeout in code

## 📚 Documentation

- **Bot README**: [discord-teneo-bot/README.md](discord-teneo-bot/README.md)
- **SDK README**: [teneo-agent-sdk/README.md](teneo-agent-sdk/README.md)
- **SDK Examples**: [teneo-agent-sdk/examples/](teneo-agent-sdk/examples/)
- **Teneo Docs**: [deploy.teneo-protocol.ai](https://deploy.teneo-protocol.ai)
- **Discord Developer Portal**: https://discord.com/developers/docs
- **Gemini API Docs**: https://ai.google.dev/docs

## 🔐 Security

⚠️ **Important Security Notes:**

- Never commit `.env` files to version control
- Keep your Discord token, Gemini API key, and private keys secure
- Use environment variables in production
- Rotate keys immediately if exposed
- Limit bot permissions to only what's needed
- The `.gitignore` already excludes `.env` files

## 🤝 Contributing

This is a personal project, but suggestions and improvements are welcome:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## 📝 License

Same as Teneo Agent SDK license.

## 🔗 Links

- **Teneo Protocol**: https://teneo-protocol.ai
- **Teneo Chatroom**: https://developer.chatroom.teneo-protocol.ai/
- **Discord Bot Invite**: Generate from Discord Developer Portal
- **Report Issues**: https://github.com/PhiBao/layla/issues

---

**Built with ❤️ using Go, Google Gemini AI, and Teneo Protocol**

🤖 **Layla is ready to assist!** 🚀
