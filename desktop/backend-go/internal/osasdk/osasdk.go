// Package osasdk is a local stub that replaces the private github.com/Miosa-osa/sdk-go
// module. It defines all types and interfaces consumed by the BusinessOS backend so
// that the project compiles without network access to the original SDK repository.
//
// NewCloudClient and NewLocalClient return errors; all cloud/OSA functionality is
// unavailable in standalone mode. The rest of the codebase gracefully degrades
// because it already checks for nil clients and error returns.
package osasdk

import (
	"context"
	"fmt"
	"time"
)

// ---------------------------------------------------------------------------
// Client interface
// ---------------------------------------------------------------------------

// Client is the interface implemented by both the local and cloud SDK clients.
type Client interface {
	Health(ctx context.Context) (*HealthResponse, error)
	GenerateApp(ctx context.Context, req AppGenerationRequest) (*AppGenerationResponse, error)
	GetAppStatus(ctx context.Context, appID string) (*AppStatusResponse, error)
	Orchestrate(ctx context.Context, req OrchestrateRequest) (*OrchestrateResponse, error)
	GetWorkspaces(ctx context.Context) (*WorkspacesResponse, error)
	Stream(ctx context.Context, sessionID string) (<-chan Event, error)
	LaunchSwarm(ctx context.Context, req SwarmRequest) (*SwarmResponse, error)
	ListSwarms(ctx context.Context) ([]SwarmStatus, error)
	GetSwarm(ctx context.Context, swarmID string) (*SwarmStatus, error)
	CancelSwarm(ctx context.Context, swarmID string) error
	DispatchInstruction(ctx context.Context, agentID string, instruction Instruction) error
	ListTools(ctx context.Context) ([]ToolDefinition, error)
	ExecuteTool(ctx context.Context, toolName string, params map[string]interface{}) (*ToolResult, error)
	GenerateAppFromTemplate(ctx context.Context, req GenerateFromTemplateRequest) (*AppGenerationResponse, error)
	Close() error
}

// ---------------------------------------------------------------------------
// Configuration types
// ---------------------------------------------------------------------------

// LocalConfig configures a local (sidecar) OSA client.
type LocalConfig struct {
	BaseURL      string
	SharedSecret string
	Timeout      time.Duration
	Resilience   *ResilienceConfig
}

// CloudConfig configures a MIOSA Cloud client.
type CloudConfig struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
}

// ResilienceConfig controls the SDK's built-in resilience behaviour.
type ResilienceConfig struct {
	Enabled bool
}

// ---------------------------------------------------------------------------
// Constructors (stubs)
// ---------------------------------------------------------------------------

// NewLocalClient returns an error because the MIOSA SDK is not available in
// standalone mode.
func NewLocalClient(cfg LocalConfig) (Client, error) {
	return nil, fmt.Errorf("osasdk: local client not available in standalone mode (MIOSA SDK stub)")
}

// NewCloudClient returns an error because cloud sync is not available in
// standalone mode.
func NewCloudClient(cfg CloudConfig) (Client, error) {
	return nil, fmt.Errorf("osasdk: cloud sync not available in standalone mode (MIOSA SDK stub)")
}

// ---------------------------------------------------------------------------
// Request / Response types
// ---------------------------------------------------------------------------

// OrchestrateRequest is the input to the Orchestrate RPC.
type OrchestrateRequest struct {
	Input       string
	UserID      string
	WorkspaceID string
	Phase       string
	SessionID   string
	Context     map[string]interface{}
}

// OrchestrateResponse is the output of the Orchestrate RPC.
type OrchestrateResponse struct {
	Success    bool
	Output     string
	AgentsUsed []string
	ExecutionMS int64
	Metadata   map[string]interface{}
	NextStep   string
	SessionID  string
}

// AppGenerationRequest is the input for app generation.
type AppGenerationRequest struct {
	UserID      string
	WorkspaceID string
	Name        string
	Description string
	Type        string
	Parameters  map[string]interface{}
}

// AppGenerationResponse is the output of app generation.
type AppGenerationResponse struct {
	AppID       string
	Status      string
	WorkspaceID string
	Message     string
	Data        map[string]interface{}
	CreatedAt   string
}

// AppStatusResponse is the output of GetAppStatus.
type AppStatusResponse struct {
	AppID       string
	Status      string
	Progress    float64
	CurrentStep string
	Output      string
	Error       string
	Metadata    map[string]interface{}
	UpdatedAt   string
}

// GenerateFromTemplateRequest is the input for template-based app generation.
type GenerateFromTemplateRequest struct {
	TemplateName string
	Variables    map[string]interface{}
	UserID       string
	WorkspaceID  string
}

// HealthResponse is the output of the Health RPC.
type HealthResponse struct {
	Status   string
	Version  string
	Provider string
}

// WorkspacesResponse is the output of GetWorkspaces.
type WorkspacesResponse struct {
	Workspaces []WorkspaceInfo
	Total      int
}

// WorkspaceInfo represents a single workspace returned by the SDK.
type WorkspaceInfo struct {
	ID          string
	Name        string
	Description string
	OwnerID     string
	CreatedAt   string
	UpdatedAt   string
}

// ---------------------------------------------------------------------------
// Swarm types
// ---------------------------------------------------------------------------

// SwarmRequest is the input for launching a swarm.
type SwarmRequest struct {
	Pattern   string
	Task      string
	Config    map[string]interface{}
	MaxAgents int
	SessionID string
}

// SwarmResponse is the output of LaunchSwarm.
type SwarmResponse struct {
	SwarmID   string
	Status    string
	SessionID string
}

// SwarmStatus represents the current state of a swarm.
type SwarmStatus struct {
	SwarmID   string
	Status    string
	Pattern   string
	Agents    []string
	CreatedAt string
	UpdatedAt string
}

// ---------------------------------------------------------------------------
// Instruction / Event types
// ---------------------------------------------------------------------------

// Instruction is a CloudEvents-style instruction sent to a fleet agent.
type Instruction struct {
	SpecVersion string
	Type        string
	Source      string
	ID          string
	Data        map[string]interface{}
}

// Event is an OSA real-time streaming event.
type Event struct {
	Type string
	Data map[string]interface{}
}

// Event type constants.
const (
	EventThinking       = "thinking"
	EventResponse       = "response"
	EventSkillStarted   = "skill_started"
	EventSkillCompleted = "skill_completed"
	EventSkillFailed    = "skill_failed"
	EventError          = "error"
	EventConnected      = "connected"
	EventSignal         = "signal"
)

// ---------------------------------------------------------------------------
// Tool types
// ---------------------------------------------------------------------------

// ToolDefinition describes an available tool.
type ToolDefinition struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
}

// ToolResult is the output of ExecuteTool.
type ToolResult struct {
	Output   string
	Metadata map[string]interface{}
}

// ---------------------------------------------------------------------------
// Error types
// ---------------------------------------------------------------------------

// APIError represents an error response from the OSA API.
type APIError struct {
	StatusCode int
	ErrorCode  string
	Details    string
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Details != "" {
		return e.Details
	}
	if e.ErrorCode != "" {
		return e.ErrorCode
	}
	return fmt.Sprintf("API error (status %d)", e.StatusCode)
}
