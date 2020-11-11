package rpc

type SomeEnum int32

const (
	SomeEnum_FIRST  SomeEnum = 0
	SomeEnum_SECOND SomeEnum = 1
)

type NestedEnum int32

const (
	NestedEnum_FIRST NestedEnum = 0
)

// ExampleMessage1 - Example Leading Comment for ExampleMessage1
type ExampleMessage1 struct {
	MyString         string              `json:"my_string"`
	MySliceOfStrings []string            `json:"MySliceOfStrings"`
	MyMap            map[string]SomeEnum `json:"MyMap"`
	EnumValue        SomeEnum            `json:"EnumValue"`
	NestedEnumValue  NestedEnum          `json:"NestedEnumValue"`
}

// ExampleMessage2 - Example Leading Comment for ExampleMessage2
type ExampleMessage2 struct {
	// MyInt just some int
	MyInt  int32            `json:"MyInt"`
	Nested []*ExampleNested `json:"nested"`
}

// ExampleNested - Example nested comment
type ExampleNested struct {
	Data   []byte        `json:"data"`
	WhyNot *DoubleNested `json:"whyNot"`
}

type DoubleNested struct {
	WhyNot string `json:"whyNot"`
}

// ReturnType some return type
type ReturnType struct {
}
