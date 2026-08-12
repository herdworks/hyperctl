package dev

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/herdworks/clif"
)

func MimetypeCommand() *clif.Command {
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "mimetype",
			Description: "Mimetype utilities.",
			Usage:       "hyperctl dev mimetype [options] <path>",
			Validator:   MimetypeCommandValidator,
			Handler:     MimetypeCommandHandler,
		},
		Options: clif.CommandOptions{},
	}
}

func MimetypeCommandValidator(cmd *clif.Command) (err error) {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("expected exactly 1 argument, got %v (%s)", len(cmd.Args), strings.Join(cmd.Args, ", "))
	}
	return nil
}

func MimetypeCommandHandler(cmd *clif.Command) (err error) {
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
