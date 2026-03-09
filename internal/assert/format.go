package assert

import (
	"fmt"
	"strings"
)

func format(msg []string, tpl string, args ...any) string {
	result := fmt.Sprintf(tpl, args...)

	if len(msg) == 0 {
		return result
	}

	return strings.Join(msg, " ") + ". " + result
}
