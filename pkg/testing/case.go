package testing

// TestSuite represents a set of test cases
type TestSuite struct {
	Name  string     `yaml:"name,omitempty" json:"name"`
	API   string     `yaml:"api,omitempty" json:"api,omitempty"`
	Items []TestCase `yaml:"items" json:"items"`
}

// TestCase represents a test case
type TestCase struct {
	Name    string   `yaml:"name,omitempty" json:"name"`
	Group   string   `yaml:"group,omitempty" json:"group"`
	Before  Job      `yaml:"before,omitempty" json:"before"`
	After   Job      `yaml:"after,omitempty" json:"after"`
	Request Request  `yaml:"request" json:"request"`
	Expect  Response `yaml:"expect,omitempty" json:"expect"`
	Clean   Clean    `yaml:"clean,omitempty" json:"-"`
}

// InScope returns true if the test case is in scope with the given items.
// Returns true if the items is empty.
func (c *TestCase) InScope(items []string) bool {
	if len(items) == 0 {
		return true
	}
	for _, item := range items {
		if item == c.Name {
			return true
		}
	}
	return false
}

// Job contains a list of jobs
type Job struct {
	Items []string `yaml:"items"`
}

// Request represents a HTTP request
type Request struct {
	API          string            `yaml:"api" json:"api"`
	Method       string            `yaml:"method,omitempty" json:"method,omitempty" jsonschema:"enum=GET,enum=POST,enum=PUT,enum=DELETE"`
	Query        map[string]string `yaml:"query,omitempty" json:"query,omitempty"`
	Header       map[string]string `yaml:"header,omitempty" json:"header,omitempty"`
	Form         map[string]string `yaml:"form,omitempty" json:"form,omitempty"`
	Body         string            `yaml:"body,omitempty" json:"body,omitempty"`
	BodyFromFile string            `yaml:"bodyFromFile,omitempty" json:"bodyFromFile,omitempty"`
}

// Response is the expected response
type Response struct {
	StatusCode       int                    `yaml:"statusCode,omitempty" json:"statusCode,omitempty"`
	Body             string                 `yaml:"body,omitempty" json:"body,omitempty"`
	Header           map[string]string      `yaml:"header,omitempty" json:"header,omitempty"`
	BodyFieldsExpect map[string]interface{} `yaml:"bodyFieldsExpect,omitempty" json:"bodyFieldsExpect,omitempty"`
	Verify           []string               `yaml:"verify,omitempty" json:"verify,omitempty"`
	Schema           string                 `yaml:"schema,omitempty" json:"schema,omitempty"`
}

// Clean represents the clean work after testing
type Clean struct {
	CleanPrepare bool `yaml:"cleanPrepare"`
}
