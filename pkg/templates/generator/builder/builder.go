package builder

import (
	"bufio"
	"fmt"
	"os"

	"github.com/udmire/observability-operator/pkg/templates/generator"
	"github.com/udmire/observability-operator/pkg/templates/generator/model"
)

type appBuilder struct {
	app   *model.App
	gen   *generator.Generator
	stack *generator.Stack[generator.OperationHolder]

	scanner *bufio.Scanner
}

func NewBuilder(gen *generator.Generator) *appBuilder {
	app := &model.App{GenericModel: model.GenericModel{Base: model.Base{Labels: make(map[string]string)}}}
	stack := generator.NewStack[generator.OperationHolder]()
	stack.Push(app)
	return &appBuilder{
		app:     app,
		gen:     gen,
		stack:   stack,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (a *appBuilder) BuildApp() error {
	opt := a.stack.Peek()

	a.buildManifestArgs(opt)

	for {
		if nextOpt, ok := opt.(generator.NextOperation); ok {
			opts := nextOpt.NextOperation(a.scanner)
			if len(opts) == 0 {
				a.stack.Pop()
				return nil
			}
			fmt.Printf("下步操作:\n %s\n请选择: ", generator.BuildInlineOptionsWithModifier[generator.Operation](opts, generator.SigForOperation))
			a.scanner.Scan()
			input := a.scanner.Text()
			parsedIdx, err := generator.ParseOptions(input, opts)
			for err != nil {
				fmt.Print("输入错误，请重新选择: ")
				a.scanner.Scan()
				input = a.scanner.Text()
				parsedIdx, err = generator.ParseOptions(input, opts)
			}
			opt := generator.Operations[parsedIdx]
			var holder generator.OperationHolder
			switch opt {
			case generator.Create:
				holder = nextOpt.Create(a.scanner)
			case generator.Modify:
				holder = nextOpt.Modify(a.scanner)
			case generator.Remove:
				nextOpt.Remove(a.scanner)
			case generator.General:
				cmd := nextOpt.General(a.scanner)
				switch cmd {
				case generator.OP_Finish:
					a.stack.Pop()
					return nil
				case generator.OP_Generate:
					a.generateApp()
				case generator.OP_Cancel:
					continue
				case generator.OP_NOOP:
					continue
				}
			}

			if holder != nil {
				a.stack.Push(holder)
				err = a.BuildApp()
				if err != nil {
					return err
				}
			}
		} else {
			a.stack.Pop()
			return nil
		}
	}
}

func (a *appBuilder) buildManifestArgs(opt generator.OperationHolder) {
	if opt == nil {
		return
	}

	args := opt.Args()
	if len(args) <= 0 {
		return
	}
	fmt.Printf("%s参数样例: %s \n请输入: ", opt.Type(), opt.ArgsExample())
	a.scanner.Scan()
	userInput := a.scanner.Text()
	err := opt.ParseArgs(userInput)
	for err != nil {
		fmt.Printf("输入参数错误！%s参数样例:\n%s\n请重新输入: ", opt.Type(), opt.ArgsExample())
		a.scanner.Scan()
		userInput := a.scanner.Text()
		err = opt.ParseArgs(userInput)
	}
}

func (a *appBuilder) generateApp() {
Loop:
	fmt.Printf("请输入要生成文件目录: ")
	a.scanner.Scan()
	input := a.scanner.Text()
	info, err := os.Stat(input)

	for err != nil || !info.IsDir() {
		fmt.Print("输入错误，请重新输入: ")
		a.scanner.Scan()
		input = a.scanner.Text()
		info, err = os.Stat(input)
	}
	err = a.gen.Generate(a.app.Name, a.app.Version, a.app.Build(), input)
	if err != nil {
		fmt.Printf("应用生成失败！错误原因：%s", err)
		goto Loop
	}
}
