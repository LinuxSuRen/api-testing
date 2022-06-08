package testing

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
	API    string            `yaml:"api"`
	Method string            `yaml:"method"`
	Query  map[string]string `yaml:"query"`
	Header map[string]string `yaml:"header"`
	Body   string            `yaml:"body"`
}

type Response struct {
	StatusCode int               `yaml:"statusCode"`
	Body       string            `yaml:"body"`
	Header     map[string]string `yaml:"header"`
}

type Clean struct {
	CleanPrepare bool `yaml:"cleanPrepare"`
}
