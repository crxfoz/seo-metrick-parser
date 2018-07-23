package worker

import (
	"github.com/crxfoz/seo_metrick_parser/parsers"
)

type ParserService struct {
	Parser parsers.Parser `json:"parser"`
	Worker *Worker        `json:"-"`
}
