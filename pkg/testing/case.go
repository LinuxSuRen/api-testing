package testing

// TestSuite represents a set of test cases
type TestSuite struct {
	Name  string     `yaml:"name"`
	API   string     `yaml:"api"`
	Items []TestCase `yaml:"items"`
}

// TestCase represents a test case
type TestCase struct {
	Name    string
	Group   string
	Prepare Prepare  `yaml:"prepare"`
	Request Request  `yaml:"request"`
	Expect  Response `yaml:"expect"`
	Clean   Clean    `yaml:"clean"`
}

// Prepare does the prepare work
type Prepare struct {
	Kubernetes []string `yaml:"kubernetes"`
}

// Request represents a HTTP request
type Request struct {
	API          string            `yaml:"api"`
	Method       string            `yaml:"method"`
	Query        map[string]string `yaml:"query"`
	Header       map[string]string `yaml:"header"`
	Form         map[string]string `yaml:"form"`
	Body         string            `yaml:"body"`
	BodyFromFile string            `yaml:"bodyFromFile"`
}

// Response is the expected response
type Response struct {
	StatusCode       int                    `yaml:"statusCode"`
	Body             string                 `yaml:"body"`
	Header           map[string]string      `yaml:"header"`
	BodyFieldsExpect map[string]interface{} `yaml:"bodyFieldsExpect"`
	Verify           []string               `yaml:"verify"`
}

// Clean represents the clean work after testing
type Clean struct {
	CleanPrepare bool `yaml:"cleanPrepare"`
}
