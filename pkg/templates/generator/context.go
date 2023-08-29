package generator

import (
	"bufio"

	"github.com/udmire/observability-operator/pkg/apps/manifest"
)

const (
	StatusAppBuilding       int = iota // Start to build an app, should given the appname and version
	StatusComponentBuilding            // Start to build a component in the app.
	StatusFileBuilding                 // Start to build a file in app layer or component layer
	StatusFileBuilded                  // File builded, to start a new session or return parent
	StatusComponentBuilded             // Component builded, to start a new component session
	StatusAppBuilded                   // App builded, to start the content generating.
	StatusAppGenerateFailed            // App content writing, incase failure.
	StatusAppGenerated                 // Generated, next to finished or modify
)

// Build AppManifests from model.
type AppBuilder interface {
	Build() *manifest.AppManifests
}

type NextOperation interface {
	NextOperation(scanner *bufio.Scanner) []Operation
	// 列出可选类型，并选择
	Create(scanner *bufio.Scanner) OperationHolder
	Modify(scanner *bufio.Scanner) OperationHolder
	Remove(scanner *bufio.Scanner)
	General(scanner *bufio.Scanner) GeneralCommand
}

type OperationHolder interface {
	Type() string

	Args() []string
	ArgsExample() string
	ParseArgs(input string) error
	String() string
}

type Operation int

const (
	Create Operation = iota
	Modify
	Remove
	General
)

var Operations = []Operation{Create, Modify, Remove, General}

type GeneralCommand string

const (
	OP_NOOP     GeneralCommand = "noop"
	OP_Finish   GeneralCommand = "finished"
	OP_Cancel   GeneralCommand = "cancel"
	OP_Generate GeneralCommand = "generate"
)
