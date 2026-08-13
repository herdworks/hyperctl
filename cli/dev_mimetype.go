package cli

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/herdworks/clic"
)

func DevMimetypeCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name: "mimetype", Description: "Mimetype utilities.",
			Usage: []string{
				"hyperctl dev mimetype [options] <path>",
			},
		},
		Spec: clic.CommandSpec{
			Key: "mimetype", Help: clic.DefaultHelp,
			Validator: mimetypeValidator, Handler: mimetypeHandler,
		},
	}
}

func mimetypeValidator(execution *clic.Execution) error {
	arguments := execution.Command.State.Arguments
	if len(arguments) != 1 {
		return fmt.Errorf("expected exactly 1 argument, got %v (%s)", len(arguments), strings.Join(arguments, ", "))
	}
	return nil
}

func mimetypeHandler(execution *clic.Execution) error {
	path := execution.Command.State.Arguments[0]
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	buffer := make([]byte, 512)
	if _, err = file.Read(buffer); err != nil {
		return err
	}
	_, err = fmt.Fprintf(execution.Stdout, "file %q has %q mimetype\n", path, http.DetectContentType(buffer))
	return err
}
