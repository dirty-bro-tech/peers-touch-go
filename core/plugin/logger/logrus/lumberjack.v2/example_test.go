package lumberjack_test

import (
	"log"

	"github.com/dirty-bro-tech/peers-touch-go/core/plugin/logger/logrus/lumberjack.v2"
)

// To use lumberjack with the standard library's log package, just pass it into
// the SetOutput function when your application starts.
func Example() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "/var/log/myapp/foo.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   // days
		Compress:   true, // disabled by default
	})
}
