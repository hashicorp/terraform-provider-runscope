package runscope

import (
  "fmt"
	"log"
  "strings"

	runscope "github.com/ewilde/go-runscope"
)

// Config contains runscope provider settings
type config struct {
	AccessToken string
	APIURL      string
}

func (c *config) client() (*runscope.Client, error) {
	client := runscope.NewClient(c.APIURL, c.AccessToken)
  runscope.RegisterLogHandlers(levelLogHandler("DEBUG"), levelLogHandler("INFO"), levelLogHandler("ERROR"))

	log.Printf("[INFO] runscope client configured for server %s", c.APIURL)

	return client, nil
}

func levelLogHandler(errorLevel string) func(level int, format string, args ...interface{}) {
	return func(level int, format string, args ...interface{}) {
		log.Printf("[%s] %s %s\n", errorLevel, strings.Repeat("\t", level-1), fmt.Sprintf(format, args...))
	}
}
