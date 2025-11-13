package libs

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

var output = zerolog.ConsoleWriter{
	Out:        os.Stdout,
	TimeFormat: "",
	NoColor:    false,
	FormatLevel: func(i any) string {
		return ""
	},
	FormatTimestamp: func(i any) string {
		return ""
	},
	FormatMessage: func(i any) string {
		return fmt.Sprintf("%s", i)
	},
}

var Logger = zerolog.New(output).With().Logger()
