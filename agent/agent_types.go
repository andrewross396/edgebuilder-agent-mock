package agent

type AgentRequest struct {
	Version       string           `json:"version"`
	Signature     string           `json:"signature"`
	RequestID     string           `json:"requestID"`
	RequesterID   string           `json:"requesterID"`
	ConcurrentOps bool             `json:"concurrentOps"`
	Ops           []AgentRequestOp `json:"ops"`
}

type AgentRequestOp struct {
	ID     string          `json:"id"`
	Name   string          `json:"name"`
	Params OperationParams `json:"params"`
}

type AgentResponseOp struct {
	ID            string             `json:"id"`
	Success       bool               `json:"success"`
	Result        map[string]string  `json:"result"`
	Errors        []AgentResponseErr `json:"errors"`
	StartTime     int                `json:"startTime"`
	ExecutionTime int                `json:"executionTime"`
}

type AgentResponseErr struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

type OperationParams = map[string]interface{}

type AgentResponse struct {
	Version     string            `json:"version"`
	Signature   string            `json:"signature"`
	RequestID   string            `json:"requestID"`
	ResponderID string            `json:"responderID"`
	Success     bool              `json:"success"`
	Ops         []AgentResponseOp `json:"ops"`
}
