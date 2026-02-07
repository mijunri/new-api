package claude_code

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/relay/channel"
	"github.com/QuantumNous/new-api/relay/channel/claude"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/types"

	"github.com/gin-gonic/gin"
)

type Adaptor struct {
}

func (a *Adaptor) ConvertGeminiRequest(*gin.Context, *relaycommon.RelayInfo, *dto.GeminiChatRequest) (any, error) {
	return nil, errors.New("not implemented")
}

func (a *Adaptor) ConvertClaudeRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.ClaudeRequest) (any, error) {
	return request, nil
}

func (a *Adaptor) ConvertAudioRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.AudioRequest) (io.Reader, error) {
	return nil, errors.New("not implemented")
}

func (a *Adaptor) ConvertImageRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (any, error) {
	return nil, errors.New("not implemented")
}

func (a *Adaptor) Init(info *relaycommon.RelayInfo) {
}

func (a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error) {
	baseURL := fmt.Sprintf("%s/v1/messages", info.ChannelBaseUrl)
	// Claude Code OAuth always uses beta=true
	baseURL = baseURL + "?beta=true"
	return baseURL, nil
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error {
	channel.SetupApiRequestHeader(info, c, req)

	// Use Bearer Token authentication for Claude Code (OAuth)
	// Authorization uses channel's API Key, not client's token
	req.Set("Authorization", "Bearer "+info.ApiKey)

	// Pass through Claude Code related headers from client
	// These headers are critical for Claude Code OAuth and vary by model/behavior
	claudeCodeHeaders := []string{
		"anthropic-version",
		"anthropic-beta",
		"x-app",
		"x-client-id",
		"x-session-id",
		"x-request-id",
		"x-trace-id",
		"x-device-id",
		"x-client-version",
		"x-platform",
		"x-os",
		"x-os-version",
	}

	for _, headerName := range claudeCodeHeaders {
		if headerValue := c.Request.Header.Get(headerName); headerValue != "" {
			req.Set(headerName, headerValue)
		}
	}

	// Set User-Agent: pass through if client provides one, otherwise use default
	userAgent := c.Request.Header.Get("User-Agent")
	if userAgent != "" {
		req.Set("User-Agent", userAgent)
	} else {
		req.Set("User-Agent", "claude-cli/2.1.6 (external, cli)")
	}

	// Set default values only if client didn't provide them
	if req.Get("anthropic-version") == "" {
		req.Set("anthropic-version", "2023-06-01")
	}

	if req.Get("x-app") == "" {
		req.Set("x-app", "cli")
	}

	// Ensure oauth beta flag is present in anthropic-beta header
	// This is REQUIRED for Claude Code OAuth authentication
	anthropicBeta := req.Get("anthropic-beta")
	oauthBeta := "oauth-2025-04-20"
	if anthropicBeta == "" {
		req.Set("anthropic-beta", oauthBeta)
	} else if !strings.Contains(anthropicBeta, oauthBeta) {
		req.Set("anthropic-beta", oauthBeta+","+anthropicBeta)
	}

	// Set Accept header if not already set
	if req.Get("Accept") == "" {
		req.Set("Accept", "application/json")
	}

	return nil
}

func (a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	return claude.RequestOpenAI2ClaudeMessage(c, *request)
}

func (a *Adaptor) ConvertRerankRequest(c *gin.Context, relayMode int, request dto.RerankRequest) (any, error) {
	return nil, nil
}

func (a *Adaptor) ConvertEmbeddingRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.EmbeddingRequest) (any, error) {
	return nil, errors.New("not implemented")
}

func (a *Adaptor) ConvertOpenAIResponsesRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.OpenAIResponsesRequest) (any, error) {
	return nil, errors.New("not implemented")
}

func (a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error) {
	return channel.DoApiRequest(a, c, info, requestBody)
}

func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError) {
	if info.IsStream {
		return claude.ClaudeStreamHandler(c, resp, info)
	} else {
		return claude.ClaudeHandler(c, resp, info)
	}
}

func (a *Adaptor) GetModelList() []string {
	return ModelList
}

func (a *Adaptor) GetChannelName() string {
	return ChannelName
}
