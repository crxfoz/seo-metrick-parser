package repository

import (
	"fmt"
	"github.com/crxfoz/webclient"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
)

var (
	reTyc = regexp.MustCompile(`tcy rang=".*?" value="(.+?)"`)
)

type TYCResult struct {
	Value int `json:"value"`
}

func GetTyc(client *webclient.Webclient, site string) (interface{}, error) {
	result := TYCResult{-1}

	resp, body, err := client.Get("http://bar-navig.yandex.ru/u").
		SetHeaders(baseHeaders).
		QueryParam("ver", "2").
		QueryParam("show", "1").
		QueryParam("url", site).Do()

	if err != nil {
		return result, err
	}

	if resp.StatusCode != 200 {
		return result, fmt.Errorf("[tyc] unexpected StatusCode: %d", resp.StatusCode)
	}

	tyc := reTyc.FindStringSubmatch(body)
	if len(tyc) == 2 {
		if v, err := strconv.Atoi(tyc[1]); err == nil {
			result.Value = v
			return result, nil
		} else {
			logrus.
				WithField("parser", "tyc").
				WithField("url", site).
				WithField("value", tyc[1]).
				Error("Could not convert string to int")
		}
	}

	return result, nil
}
