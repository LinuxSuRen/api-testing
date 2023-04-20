package testing

// TestSuite represents a set of test cases
type TestSuite struct {
	Name  string     `yaml:"name" json:"name"`
	API   string     `yaml:"api,omitempty" json:"api,omitempty"`
	Items []TestCase `yaml:"items" json:"items"`
}

// TestCase represents a test case
type TestCase struct {
	Name    string `yaml:"name" json:"name"`
	Group   string
	Prepare Prepare  `yaml:"prepare" json:"-"`
	Request Request  `yaml:"request" json:"request"`
	Expect  Response `yaml:"expect" json:"expect"`
	Clean   Clean    `yaml:"clean" json:"-"`
}

// Prepare does the prepare work
type Prepare struct {
	Kubernetes []string `yaml:"kubernetes"`
}

// Request represents a HTTP request
type Request struct {
	API          string            `yaml:"api" json:"api"`
	Method       string            `yaml:"method,omitempty" json:"method,omitempty" jsonschema:"enum=GET,enum=POST,enum=PUT,enum=DELETE"`
	Query        map[string]string `yaml:"query" json:"query,omitempty"`
	Header       map[string]string `yaml:"header" json:"header,omitempty"`
	Form         map[string]string `yaml:"form" json:"form,omitempty"`
	Body         string            `yaml:"body" json:"body,omitempty"`
	BodyFromFile string            `yaml:"bodyFromFile" json:"bodyFromFile,omitempty"`
}

// Response is the expected response
type Response struct {
	StatusCode       int                    `yaml:"statusCode" json:"statusCode,omitempty"`
	Body             string                 `yaml:"body" json:"body,omitempty"`
	Header           map[string]string      `yaml:"header" json:"header,omitempty"`
	BodyFieldsExpect map[string]interface{} `yaml:"bodyFieldsExpect" json:"bodyFieldsExpect,omitempty"`
	Verify           []string               `yaml:"verify" json:"verify,omitempty"`
	Schema           string                 `yaml:"schema" json:"schema,omitempty"`
}

// Clean represents the clean work after testing
type Clean struct {
	CleanPrepare bool `yaml:"cleanPrepare"`
}
