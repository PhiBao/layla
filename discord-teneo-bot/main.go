package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/TeneoProtocolAI/teneo-agent-sdk/pkg/agent"
	"github.com/bwmarrin/discordgo"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

// DiscordTeneoAgent combines Discord bot with Teneo agent functionality
type DiscordTeneoAgent struct {
	geminiClient  *genai.Client
	geminiModel   *genai.GenerativeModel
	discord       *discordgo.Session
	botUserID     string
	systemPrompt  string
	rateLimiter   *RateLimiter
}

// RateLimiter implements a simple rate limiter
type RateLimiter struct {
	mu           sync.Mutex
	lastRequest  time.Time
	minInterval  time.Duration
}

func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	interval := time.Minute / time.Duration(requestsPerMinute)
	return &RateLimiter{
		minInterval: interval,
	}
}

func (rl *RateLimiter) Wait() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	elapsed := now.Sub(rl.lastRequest)
	
	if elapsed < rl.minInterval {
		waitTime := rl.minInterval - elapsed
		time.Sleep(waitTime)
	}
	
	rl.lastRequest = time.Now()
}

// ProcessTask implements the Teneo AgentHandler interface
// This handles tasks from the Teneo network
func (d *DiscordTeneoAgent) ProcessTask(ctx context.Context, task string) (string, error) {
	// Parse command from task
	parts := strings.Fields(task)
	if len(parts) == 0 {
		return "Please provide a command. Available commands: ask, explain, help", nil
	}
	
	command := strings.ToLower(parts[0])
	var question string
	if len(parts) > 1 {
		question = strings.Join(parts[1:], " ")
	}
	
	switch command {
	case "ask":
		if question == "" {
			return "Please provide a question. Usage: ask [your question]", nil
		}
		geminiCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		response, err := d.queryGemini(geminiCtx, question)
		if err != nil {
			return "", fmt.Errorf("failed to process question: %w", err)
		}
		return response, nil
		
	case "explain":
		if question == "" {
			return "Please provide a concept to explain. Usage: explain [concept]", nil
		}
		geminiCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		response, err := d.queryGemini(geminiCtx, question)
		if err != nil {
			return "", fmt.Errorf("failed to explain concept: %w", err)
		}
		return response, nil
		
	case "help":
		helpText := "🤖 Layla - AI Assistant\n\n" +
			"Available commands:\n" +
			"• ask [question] - Ask me anything\n" +
			"• explain [concept] - Get detailed explanation of a concept\n" +
			"• help - Show this help message\n\n" +
			"Example: ask what is the meaning of life"
		return helpText, nil
		
	default:
		// If no recognized command, treat entire task as a question
		geminiCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		response, err := d.queryGemini(geminiCtx, task)
		if err != nil {
			return "", fmt.Errorf("failed to process task: %w", err)
		}
		return response, nil
	}
}

// queryGemini sends a query to Google Gemini and returns the response
func (d *DiscordTeneoAgent) queryGemini(ctx context.Context, userMessage string) (string, error) {
	// Apply rate limiting
	d.rateLimiter.Wait()
	
	prompt := d.systemPrompt + "\n\nUser: " + userMessage
	
	resp, err := d.geminiModel.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}
	
	if len(resp.Candidates[0].Content.Parts) == 0 {
		// Check if response was blocked
		if resp.Candidates[0].FinishReason != 0 {
			return "I apologize, but I cannot provide a response to that query due to content safety filters.", nil
		}
		return "", fmt.Errorf("no response from Gemini")
	}

	// Extract text from response
	result := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	return result, nil
}

// handleDiscordMessage processes messages from Discord
func (d *DiscordTeneoAgent) handleDiscordMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot's own messages
	if m.Author.ID == d.botUserID {
		return
	}

	// Ignore other bots
	if m.Author.Bot {
		return
	}

	// Check if bot was mentioned or message starts with !ask
	mentioned := false
	for _, user := range m.Mentions {
		if user.ID == d.botUserID {
			mentioned = true
			break
		}
	}

	// Process if mentioned or command starts with !ask
	if !mentioned && !strings.HasPrefix(m.Content, "!ask") {
		return
	}

	// Clean the message (remove mention and command)
	content := m.Content
	content = strings.TrimPrefix(content, "!ask")
	content = strings.TrimSpace(content)
	
	// Remove bot mention from content
	for _, user := range m.Mentions {
		if user.ID == d.botUserID {
			content = strings.Replace(content, fmt.Sprintf("<@%s>", user.ID), "", -1)
			content = strings.Replace(content, fmt.Sprintf("<@!%s>", user.ID), "", -1)
		}
	}
	content = strings.TrimSpace(content)

	if content == "" {
		s.ChannelMessageSend(m.ChannelID, "👋 Hi! I'm Layla, your AI assistant! Ask me anything or use `!ask <your question>`")
		return
	}

	// Show typing indicator
	s.ChannelTyping(m.ChannelID)

	// Process the message through Gemini
	ctx := context.Background()
	response, err := d.queryGemini(ctx, content)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ Sorry, I encountered an error: %v", err))
		return
	}

	// Split response if it's too long (Discord has 2000 char limit)
	if len(response) > 2000 {
		chunks := splitMessage(response, 2000)
		for _, chunk := range chunks {
			s.ChannelMessageSend(m.ChannelID, chunk)
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, response)
	}
}

// splitMessage splits a long message into chunks
func splitMessage(message string, maxLen int) []string {
	if len(message) <= maxLen {
		return []string{message}
	}

	var chunks []string
	for len(message) > maxLen {
		// Try to split at a newline or space
		splitAt := maxLen
		for i := maxLen - 1; i > maxLen-200 && i > 0; i-- {
			if message[i] == '\n' || message[i] == ' ' {
				splitAt = i
				break
			}
		}
		chunks = append(chunks, message[:splitAt])
		message = message[splitAt:]
	}
	if len(message) > 0 {
		chunks = append(chunks, message)
	}
	return chunks
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	// Validate required environment variables
	requiredEnvVars := []string{
		"DISCORD_TOKEN",
		"GEMINI_API_KEY",
		"PRIVATE_KEY",
	}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("❌ Required environment variable %s is not set", envVar)
		}
	}
	
	// Warn if Teneo vars are missing but continue
	if os.Getenv("OWNER_ADDRESS") == "" {
		log.Println("⚠️  OWNER_ADDRESS not set - will be derived from private key")
	}
	if os.Getenv("NFT_TOKEN_ID") == "" {
		log.Println("⚠️  NFT_TOKEN_ID not set - Teneo agent features may be limited")
	}

	// Initialize Gemini client
	ctx := context.Background()
	geminiClient, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatalf("❌ Error creating Gemini client: %v", err)
	}
	defer geminiClient.Close()

	// Create Gemini model - gemini-2.5-flash works fine
	geminiModel := geminiClient.GenerativeModel("gemini-2.5-flash")
	geminiModel.SetTemperature(0.7)
	geminiModel.SetMaxOutputTokens(1000)
	// Set safety settings to allow more content
	geminiModel.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
	}

	// Create Discord session
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatalf("❌ Error creating Discord session: %v", err)
	}

	// Get bot user info
	user, err := discord.User("@me")
	if err != nil {
		log.Fatalf("❌ Error getting bot user info: %v", err)
	}

	// Create agent handler
	systemPrompt := os.Getenv("SYSTEM_PROMPT")
	if systemPrompt == "" {
		systemPrompt = "You are Layla, a helpful AI assistant in a Discord server. Be friendly, concise, and helpful. Answer questions accurately and engage naturally with users."
	}

	// Create rate limiter for Gemini API (8 requests per minute - well within the 10 RPM limit)
	rateLimiter := NewRateLimiter(8)

	agentHandler := &DiscordTeneoAgent{
		geminiClient: geminiClient,
		geminiModel:  geminiModel,
		discord:      discord,
		botUserID:    user.ID,
		systemPrompt: systemPrompt,
		rateLimiter:  rateLimiter,
	}

	// Configure Teneo Agent using default config
	agentConfig := agent.DefaultConfig()
	agentConfig.Name = "Layla"
	agentConfig.Description = "AI-powered Discord assistant connected to Teneo network"
	agentConfig.Version = "1.0.0"
	agentConfig.Capabilities = []string{"discord_chat", "question_answering", "concept_explanation", "conversation", "teneo_network"}
	agentConfig.PrivateKey = os.Getenv("PRIVATE_KEY")

	// Parse NFT Token ID as uint64
	var tokenID uint64
	if nftID := os.Getenv("NFT_TOKEN_ID"); nftID != "" {
		fmt.Sscanf(nftID, "%d", &tokenID)
		agentConfig.NFTTokenID = nftID  // Also set NFTTokenID string in Config
	}

	// Create Enhanced Teneo agent (needed for chatroom connectivity)
	enhancedConfig := &agent.EnhancedAgentConfig{
		Config:       agentConfig,
		AgentHandler: agentHandler,
		Mint:         false,
		TokenID:      tokenID,  // Use TokenID in EnhancedAgentConfig, not NFTTokenID in Config
	}

	teneoAgent, err := agent.NewEnhancedAgent(enhancedConfig)
	if err != nil {
		log.Fatalf("❌ Error creating Teneo agent: %v", err)
	}

	// Register Discord message handler
	discord.AddHandler(agentHandler.handleDiscordMessage)

	// Set intents
	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	// Open Discord connection
	err = discord.Open()
	if err != nil {
		log.Fatalf("❌ Error opening Discord connection: %v", err)
	}
	defer discord.Close()

	log.Printf("🤖 Discord-Teneo Bot '%s' is now running!", user.Username)
	log.Println("📱 Discord: Connected and listening for messages")
	log.Println("💬 Usage: Mention the bot or use !ask <question>")
	log.Println("🌐 Starting Teneo agent...")

	// Start Teneo agent in background (Run blocks until stopped)
	go func() {
		if err := teneoAgent.Run(); err != nil {
			log.Printf("❌ Teneo agent error: %v", err)
		}
	}()

	// Give the agent time to fully connect and register
	time.Sleep(2 * time.Second)
	
	log.Println("✅ Teneo agent should now be visible in chatroom")
	log.Println("⏳ Check: https://developer.chatroom.teneo-protocol.ai/")
	log.Println("Press CTRL+C to stop the bot")

	// Wait for interrupt signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("🛑 Shutting down gracefully...")
}
