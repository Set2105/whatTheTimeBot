package wttbot

import (
	"fmt"
	"regexp"
)

var re = regexp.MustCompile(`"([^"]*)"|\S+`)

func SplitArgs(input string) []string {
	matches := re.FindAllStringSubmatch(input, -1)

	var result []string
	for _, match := range matches {
		if match[1] != "" {
			result = append(result, match[1])
		} else {
			result = append(result, match[0])
		}
	}

	return result
}

func parseArgs(text string, minLen int) ([]string, error) {
	args := SplitArgs(text)
	if len(args) > 0 {
		args = args[1:]
	}
	if len(args) < minLen {
		return nil, fmt.Errorf("not enough arguments, must be %d check /help for detailes", minLen)
	}
	return args, nil
}

func parseVariousArgs(text string, l ...int) ([]string, error) {
	args := SplitArgs(text)
	if len(args) > 0 {
		args = args[1:]
	}
	ok := false
	argsLen := len(args)
	for _, i := range l {
		if i == argsLen {
			ok = true
			break
		}
	}
	if !ok {
		return nil, fmt.Errorf("not enough arguments, must be in %v check /help for detailes", l)
	}
	return args, nil
}
