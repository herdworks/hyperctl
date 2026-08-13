package cli

import (
	"fmt"
	"strings"

	"github.com/herdworks/clic"
)

func positionalArgumentsError(execution *clic.Execution, expected int) error {
	return fmt.Errorf(
		"%s expected %v argument(s), got %v (%s)",
		execution.Command.Meta.Name,
		expected,
		execution.Command.PositionalArguments(),
		strings.Join(execution.Command.State.Arguments, ", "),
	)
}
