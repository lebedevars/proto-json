package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	generator "protoc-gen-protojson"

	"google.golang.org/protobuf/compiler/protogen"
)

var (
	version = "v1.0.0"
	flags   flag.FlagSet
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Fprintf(os.Stderr, "%s %s\n", filepath.Base(os.Args[0]), version)
		os.Exit(0)
	}

	o := protogen.Options{
		ParamFunc:         flags.Set,
		ImportRewriteFunc: nil,
	}
	o.Run(func(gen *protogen.Plugin) error {
		modelGen := generator.NewJsonRpcModelGenerator(gen)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			modelGen.Generate(f)
		}

		return nil
	})
}
