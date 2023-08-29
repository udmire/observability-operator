package main

import (
	"github.com/udmire/observability-operator/pkg/templates/generator"
	"github.com/udmire/observability-operator/pkg/templates/generator/builder"
	util_log "github.com/udmire/observability-operator/pkg/utils/log"
)

func main() {
	gen := generator.NewGenerator(util_log.Logger)
	builder := builder.NewBuilder(gen)
	err := builder.BuildApp()
	for err != nil {
		err = builder.BuildApp()
	}
}
