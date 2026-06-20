package output

type Status string

const (
	StatusOK      Status = "ok"
	StatusUpdate  Status = "update"
	StatusMissing Status = "missing"
	StatusSkipped Status = "skipped"
	StatusError   Status = "error"
)

type EventType string

const (
	EventTypeEvent   EventType = "event"
	EventTypeStatus  EventType = "status"
	EventTypeSummary EventType = "summary"
)

type EventName string

const (
	EventFetch   EventName = "fetch"
	EventClone   EventName = "clone"
	EventCompare EventName = "compare"
	EventInstall EventName = "install"
	EventRemove  EventName = "remove"
	EventRepair  EventName = "repair"
	EventPlan    EventName = "plan"
)

type SkillStatus struct {
	Name         string `json:"name"`
	Status       Status `json:"status"`
	SourceURL    string `json:"sourceUrl,omitempty"`
	RemoteHash   string `json:"remoteHash,omitempty"`
	InstalledDir string `json:"-"`
	CompareDir   string `json:"-"`
}

type PlannedAction struct {
	Action string `json:"action"`
	Name   string `json:"name,omitempty"`
	Target string `json:"target,omitempty"`
	Path   string `json:"path,omitempty"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Target  string `json:"target,omitempty"`
}

type Summary struct {
	OK         bool            `json:"ok"`
	Command    string          `json:"command"`
	Entrypoint string          `json:"entrypoint"`
	DryRun     bool            `json:"dryRun"`
	Statuses   []SkillStatus   `json:"statuses"`
	Actions    []PlannedAction `json:"actions"`
	Errors     []ErrorDetail   `json:"errors"`
}

type Event struct {
	Type       EventType `json:"type"`
	Event      EventName `json:"event,omitempty"`
	Name       string    `json:"name,omitempty"`
	Status     Status    `json:"status,omitempty"`
	SourceURL  string    `json:"sourceUrl,omitempty"`
	RemoteHash string    `json:"remoteHash,omitempty"`
	OK         bool      `json:"ok,omitempty"`
	Changed    int       `json:"changed,omitempty"`
	Errors     int       `json:"errors,omitempty"`
}
