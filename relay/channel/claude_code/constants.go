package claude_code

// ModelList contains the models supported by Claude Code (OAuth) channel
// Claude Code (OAuth) only supports Claude 4.5+ models
var ModelList = []string{
	// Claude 4.5 series
	"claude-sonnet-4-5-20250929",
	"claude-opus-4-5-20251101",
	"claude-haiku-4-5-20251001",
}

const ChannelName = "claude_code"
