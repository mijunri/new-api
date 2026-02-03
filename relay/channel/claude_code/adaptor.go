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
	RequestMode int
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
	if strings.HasPrefix(info.UpstreamModelName, "claude-2") || strings.HasPrefix(info.UpstreamModelName, "claude-instant") {
		a.RequestMode = claude.RequestModeCompletion
	} else {
		a.RequestMode = claude.RequestModeMessage
	}
}

func (a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error) {
	baseURL := ""
	if a.RequestMode == claude.RequestModeMessage {
		baseURL = fmt.Sprintf("%s/v1/messages", info.ChannelBaseUrl)
	} else {
		baseURL = fmt.Sprintf("%s/v1/complete", info.ChannelBaseUrl)
	}
	// Claude Code OAuth always uses beta=true
	baseURL = baseURL + "?beta=true"
	return baseURL, nil
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error {
	channel.SetupApiRequestHeader(info, c, req)

	// Use Bearer Token authentication for Claude Code (OAuth)
	req.Set("Authorization", "Bearer "+info.ApiKey)

	// Set anthropic-version
	anthropicVersion := c.Request.Header.Get("anthropic-version")
	if anthropicVersion == "" {
		anthropicVersion = "2023-06-01"
	}
	req.Set("anthropic-version", anthropicVersion)

	// Set anthropic-beta header with oauth flag and interleaved-thinking
	anthropicBeta := c.Request.Header.Get("anthropic-beta")
	requiredBetas := []string{"oauth-2025-04-20", "interleaved-thinking-2025-05-14"}
	if anthropicBeta != "" {
		// Add required betas if not already included
		for _, beta := range requiredBetas {
			if !strings.Contains(anthropicBeta, beta) {
				anthropicBeta = beta + "," + anthropicBeta
			}
		}
	} else {
		anthropicBeta = strings.Join(requiredBetas, ",")
	}
	req.Set("anthropic-beta", anthropicBeta)

	// Set x-app header - CRITICAL for Claude Code OAuth
	xApp := c.Request.Header.Get("x-app")
	if xApp == "" {
		xApp = "cli"
	}
	req.Set("x-app", xApp)

	// Set User-Agent to mimic Claude CLI client
	userAgent := c.Request.Header.Get("User-Agent")
	if userAgent == "" || !strings.Contains(userAgent, "claude") {
		req.Set("User-Agent", "claude-cli/2.1.6 (external, cli)")
	}

	// Set Accept header
	if req.Get("Accept") == "" {
		req.Set("Accept", "application/json")
	}

	// Apply Claude-specific headers from model settings
	claude.CommonClaudeHeadersOperation(c, req, info)

	return nil
}

func (a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	if a.RequestMode == claude.RequestModeCompletion {
		return claude.RequestOpenAI2ClaudeComplete(*request), nil
	} else {
		return claude.RequestOpenAI2ClaudeMessage(c, *request)
	}
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
		return claude.ClaudeStreamHandler(c, resp, info, a.RequestMode)
	} else {
		return claude.ClaudeHandler(c, resp, info, a.RequestMode)
	}
}

func (a *Adaptor) GetModelList() []string {
	return ModelList
}

func (a *Adaptor) GetChannelName() string {
	return ChannelName
}
