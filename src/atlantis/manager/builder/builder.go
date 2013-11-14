package builder

import (
	. "atlantis/common"
	"io"
)

var DefaultBuilder Builder

type Builder interface {
	Build(*Task, string, string, string) (io.ReadCloser, error)
}
