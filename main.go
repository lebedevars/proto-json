package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	req := &pluginpb.CodeGeneratorRequest{}
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Panic(err)
	}

	if err := proto.Unmarshal(input, req); err != nil {
		log.Panic(err)
	}

	opts := protogen.Options{}
	plugin, err := opts.New(req)
	if err != nil {
		log.Panic(err)
	}

	generator := &jsonRpcGenerator{plugin: plugin}
	err = generator.generate()
	if err != nil {
		log.Panic(err)
	}

	stdout := plugin.Response()
	out, err := proto.Marshal(stdout)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(os.Stdout, string(out))
}
