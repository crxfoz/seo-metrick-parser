package parsers

import (
	"github.com/crxfoz/seo_metrick_parser/parsers/repository"
	"os"
	"time"
)

var ParsersList []Parser

func init() {
	moz := &repository.MozParser{
		ApiKeyAccessID:  os.Getenv("MOZ_ACCESSID"),
		ApiKeySecretKey: os.Getenv("MOZ_SECRET_KEY"),
	}

	ParsersList = []Parser{
		{
			Name:        "alexa",
			Description: "Get information from http://alexa.com. Return a global alexa rank, local alexa rank, GEO.",
			ParserFn:    repository.GetAlexa,
			Timeout:     time.Second * 5,
			Status:      true,
		},
		{
			Name:        "tyc",
			Description: "Get Yandex TYC",
			ParserFn:    repository.GetTyc,
			Timeout:     time.Second * 5,
			Status:      true,
		},
		{
			Name:        "moz",
			Description: "Get DA PA and backlinks from MozApi",
			ParserFn:    moz.GetMoz,
			Timeout:     time.Second * 12,
			Status:      true,
		},
	}
}
