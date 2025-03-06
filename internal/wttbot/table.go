package wttbot

import "fmt"

func escapeTag(input, tag string) string {
	return fmt.Sprintf("<%s>%s</%s>", tag, input, tag)
}
