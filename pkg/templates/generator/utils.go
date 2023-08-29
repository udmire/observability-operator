package generator

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const NamePatternString = "[a-zA-Z][a-z0-9A-Z-]{0,62}[a-z0-9A-Z]"

var namePattern = regexp.MustCompile(NamePatternString)

func ParseOptions[T any](input string, options []T) (int, error) {
	parsedIdx, err := strconv.Atoi(input)
	if err != nil {
		return -1, err
	}

	if parsedIdx >= 0 && parsedIdx <= len(options)-1 {
		return parsedIdx, nil
	}
	return -1, fmt.Errorf("invalid input")
}

func BuildOptionsWithModifier[T any](ojbs []T, fun func(T) string) string {
	var line []string
	for idx, typ := range ojbs {
		line = append(line, fmt.Sprintf("%d %s", idx, fun(typ)))
	}
	return strings.Join(line, "\n\t")
}
func BuildOptions[T any](ojbs []T) string {
	return BuildInlineOptionsWithModifier[T](ojbs, func(t T) string {
		return fmt.Sprintf("%v", t)
	})
}

func BuildInlineOptionsWithModifier[T any](ojbs []T, fun func(T) string) string {
	var line []string
	for idx, typ := range ojbs {
		line = append(line, fmt.Sprintf("%d %s", idx, fun(typ)))
	}
	return strings.Join(line, "\t")
}

func BuildInlineOptions[T any](ojbs []T) string {
	return BuildInlineOptionsWithModifier[T](ojbs, func(t T) string {
		return fmt.Sprintf("%v", t)
	})
}

func SigForOperation(op Operation) string {
	switch op {
	case Create:
		return "Create"
	case Remove:
		return "Delete"
	case Modify:
		return "Modify"
	case General:
		return "Operation"
	}
	return ""
}

// 列出可选类型，并选择
func ChooseOperation(scanner *bufio.Scanner) Operation {
	fmt.Printf("下步操作:\n %s\n请选择: ", BuildInlineOptionsWithModifier[Operation](Operations, SigForOperation))
	scanner.Scan()
	input := scanner.Text()
	parsedIdx, err := ParseOptions(input, Operations)
	for err != nil {
		fmt.Print("输入错误，请重新选择: ")
		scanner.Scan()
		input = scanner.Text()
		parsedIdx, err = ParseOptions(input, Operations)
	}
	return Operations[parsedIdx]
}

func ChooseType[T any](scanner *bufio.Scanner, opts []T) T {
	fmt.Printf("可选类型列表:\n%s\n请选择: ", BuildInlineOptions[T](opts))
	scanner.Scan()
	input := scanner.Text()
	parsedIdx, err := ParseOptions(input, opts)
	for err != nil {
		fmt.Print("输入错误，请重新选择: ")
		scanner.Scan()
		input = scanner.Text()
		parsedIdx, err = ParseOptions(input, opts)
	}
	return opts[parsedIdx]
}

func ChooseIndex(scanner *bufio.Scanner, typ string, holders []OperationHolder) int {
	fmt.Printf("请选择要操作的%s:\n %s", typ, BuildOptions(holders))
	input := scanner.Text()
	parsedIdx, err := ParseOptions(input, holders)
	for err != nil {
		fmt.Print("输入错误，请重新选择: ")
		scanner.Scan()
		input = scanner.Text()
		parsedIdx, err = ParseOptions(input, holders)
	}
	return parsedIdx
}

func IsValidName(val string) bool {
	return namePattern.MatchString(val)
}
