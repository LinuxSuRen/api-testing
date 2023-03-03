package testing

type TestSuite struct {
	Name  string     `yaml:"name"`
	Items []TestCase `yaml:"items"`
}

type TestCase struct {
	Name    string
	Group   string
	Prepare Prepare  `yaml:"prepare"`
	Request Request  `yaml:"request"`
	Expect  Response `yaml:"expect"`
	Clean   Clean    `yaml:"clean"`
}

type Prepare struct {
	Kubernetes []string `yaml:"kubernetes"`
}

type Request struct {
	API          string            `yaml:"api"`
	Method       string            `yaml:"method"`
	Query        map[string]string `yaml:"query"`
	Header       map[string]string `yaml:"header"`
	Body         string            `yaml:"body"`
	BodyFromFile string            `yaml:"bodyFromFile"`
}

type Response struct {
	StatusCode       int               `yaml:"statusCode"`
	Body             string            `yaml:"body"`
	Header           map[string]string `yaml:"header"`
	BodyFieldsExpect map[string]string `yaml:"bodyFieldsExpect"`
	Verify           []string          `yaml:"verify"`
}

type Clean struct {
	CleanPrepare bool `yaml:"cleanPrepare"`
}
