package rpc

type exampleService2 struct{}

type ExampleService2 interface {
	ExampleCall1(ExampleMessage1) ReturnType
	ExampleCall2(ExampleMessage2) ReturnType
}
