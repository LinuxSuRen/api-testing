package pkg

type TestCase struct {
	SuiteName string `json:"suiteName"`
	Name      string `json:"name"`
	API       string
	Method    string
	Body      string
	Header    string
	Query     string
	Form      string

	ExpectStatusCode int
	ExpectBody       string
	ExpectSchema     string
	ExpectHeader     string
	ExpectBodyFields string
	ExpectVerify     string
}

type TestSuite struct {
	Name string
	API  string
}
