package fake

import (
	"io"
	"strings"
)

type StructWriterTo interface {
	StructWriteTo(io.Writer) (int, error)
}

func DescribeStruct(s StructWriterTo, n string) string {
	var builder strings.Builder
	builder.WriteString(n)
	builder.WriteByte('(')
	_, _ = s.StructWriteTo(&builder)
	builder.WriteByte(')')

	return builder.String()
}
