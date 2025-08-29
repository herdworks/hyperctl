package cli

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func MimetypeCommand() *Command {
	return &Command{
		Config: CommandConfig{
			Name:        "mimetype",
			Description: "Get file mimetype information.",
			Usage:       "hyperctl mimetype [options] <path>",
			Validator:   MimetypeCommandValidator,
			Handler:     MimetypeCommandHandler,
		},
		Options: CommandOptions{},
	}
}

func MimetypeCommandValidator(cmd *Command) (err error) {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("expected exactly 1 argument, got %v (%s)", len(cmd.Args), strings.Join(cmd.Args, ", "))
	}
	return nil
}

func MimetypeCommandHandler(cmd *Command) (err error) {
	var uri string = cmd.Args[0]
	var file *os.File
	var buffer []byte = make([]byte, 512)
	if file, err = os.Open(uri); err == nil {
		defer file.Close()
	} else {
		return err
	}
	if _, err = file.Read(buffer); err != nil {
		return err
	}
	mimetype := http.DetectContentType(buffer)
	fmt.Printf("file \"%s\" has \"%s\" mimetype\n", uri, mimetype)
	return nil
}
