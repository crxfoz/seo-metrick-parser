package parsers

import (
	"github.com/crxfoz/webclient"
	"time"
)

type ParserFn func(*webclient.Webclient, string) (interface{}, error)

type Parser struct {
	ParserFn    ParserFn      `json:"-"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Timeout     time.Duration `json:"-"`
	Status      bool          `json:"status"`
}
