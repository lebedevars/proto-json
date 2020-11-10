package rpc

type exampleService struct{}

type ExampleService interface {
	ExampleCall1(ExampleMessage1) ReturnType
	ExampleCall2(ExampleMessage2) ReturnType
}

type SomeEnum int32

const (
	FIRST  SomeEnum = 0
	SECOND SomeEnum = 1
)

type MyMapEntry struct {
	Key   string   `json:"key"`
	Value SomeEnum `json:"value"`
}

type ExampleMessage1 struct {
	MyString         string              `json:"my_string"`
	MySliceOfStrings []string            `json:"MySliceOfStrings"`
	MyMap            map[string]SomeEnum `json:"MyMap"`
	EnumValue        SomeEnum            `json:"EnumValue"`
}

type ExampleNested struct {
	Data []byte `json:"data"`
}

type ExampleMessage2 struct {
	MyInt  int32            `json:"MyInt"`
	Nested []*ExampleNested `json:"nested"`
}

type ReturnType struct {
}
