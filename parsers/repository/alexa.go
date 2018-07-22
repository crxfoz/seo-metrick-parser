package repository

import (
	"fmt"
	"github.com/crxfoz/webclient"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
)

var (
	reGlobal = regexp.MustCompile(`<POPULARITY URL=".*?" TEXT="(.+?)" SOURCE=`)
	reLocal  = regexp.MustCompile(`COUNTRY CODE="(.+?)" NAME=".*?" RANK="(.+?)"`)
)

type AlexaResult struct {
	Global int    `json:"global"`
	GEO    string `json:"geo"`
	Local  int    `json:"local"`
}

func GetAlexa(client *webclient.Webclient, site string) (interface{}, error) {
	var result AlexaResult

	resp, body, err := client.Get("http://xml.alexa.com/data/").
		SetHeaders(baseHeaders).
		QueryParam("cli", "10").
		QueryParam("dat", "nsa").
		QueryParam("url", site).
		Do()

	if err != nil {
		return result, err
	}

	if resp.StatusCode != 200 {
		return result, fmt.Errorf("[alexa] unexpected StatusCode: %d", resp.StatusCode)
	}

	globalResult := reGlobal.FindStringSubmatch(body)
	if len(globalResult) == 2 {
		if v, err := strconv.Atoi(globalResult[1]); err == nil {
			result.Global = v
		} else {
			logrus.
				WithField("parser", "alexa").
				WithField("url", site).
				WithField("value", globalResult[1]).
				Error("cant convert string to int, alexa parser")
		}
	}

	localResult := reLocal.FindStringSubmatch(body)
	if len(localResult) == 3 {
		result.GEO = localResult[1]
		if v, err := strconv.Atoi(localResult[2]); err == nil {
			result.Local = v
		} else {
			logrus.
				WithField("parser", "alexa").
				WithField("url", site).
				WithField("value", localResult[2]).
				Error("Could not convert string to int")
		}
	}

	return &result, nil

}
