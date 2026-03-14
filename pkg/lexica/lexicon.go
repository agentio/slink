package lexica

// https://atproto.com/specs/lexicon#lexicon-files
type Lexicon struct {
	Lexicon     int             `json:"lexicon"`
	Id          string          `json:"id"`
	Description string          `json:"description"`
	Defs        map[string]*Def `json:"defs"`
}

type Def struct {
	Type        string `json:"type"`
	Description string `json:"description"`

	// object
	Required   []string            `json:"required,omitempty"`
	Nullable   []string            `json:"nullable,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`

	// query
	Parameters *Parameters `json:"parameters,omitempty"`
	Output     *Output     `json:"output,omitempty"`

	// procedure
	Input *Input `json:"input,omitempty"`

	// array
	Items *Items `json:"items,omitempty"`

	// permission-set
	Permissions []*Permission `json:"permissions,omitempty"`

	// record
	Record *Schema `json:"record,omitempty"`

	// subscription
	Message *Message `json:"message,omitempty"`
}

type Parameters struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
}

type Output struct {
	Encoding string `json:"encoding"`
	Schema   Schema `json:"schema"`
}

type Input struct {
	Encoding string `json:"encoding"`
	Schema   Schema `json:"schema"`
}

type Message struct {
	Description string `json:"description"`
	Schema      Schema `json:"schema"`
}

type Schema struct {
	Type       string              `json:"type"`
	Ref        string              `json:"ref,omitempty"`
	Refs       []string            `json:"refs,omitempty"`
	Required   []string            `json:"required,omitempty"`
	Nullable   []string            `json:"nullable,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
}

type Property struct {
	Type       string              `json:"type,omitempty"`
	Ref        string              `json:"ref,omitempty"`
	Refs       []string            `json:"refs,omitempty"`
	Items      *Items              `json:"items,omitempty"`
	Minimum    int64               `json:"minimum,omitempty"`
	Maximum    int64               `json:"maximum,omitempty"`
	Default    any                 `json:"default,omitempty"`
	Required   []string            `json:"required,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
}

type Items struct {
	Type string   `json:"type,omitempty"`
	Ref  string   `json:"ref,omitempty"`
	Refs []string `json:"refs,omitempty"`
}

type Permission struct {
	Type           string   `json:"type,omitempty"`
	Resource       string   `json:"resource,omitempty"`
	Actions        []string `json:"action,omitempty"`
	Collections    []string `json:"collection,omitempty"`
	InheritAud     bool     `json:"inheritAud,omitempty"`
	LexiconMethods []string `json:"lxm,omitempty"`
}
