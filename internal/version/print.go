package version

import (
	"fmt"
	"io"
)

func Run(w io.Writer) int {
	_, _ = fmt.Fprintln(w, Version())

	return 0
}
