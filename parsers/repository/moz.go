package repository

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crxfoz/webclient"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type MozParser struct {
	ApiKeyAccessID  string
	ApiKeySecretKey string
}

type MozResult struct {
	DA        int `json:"da"`
	PA        int `json:"pa"`
	Backlinks int `json:"backlinks"`
}

type mozResultLocal struct {
	DA        float64 `json:"pda"`
	PA        float64 `json:"upa"`
	Backlinks int     `json:"uid"`
}

func (m *MozParser) GetMoz(client *webclient.Webclient, site string) (interface{}, error) {
	// https://github.com/seomoz/SEOmozAPISamples/blob/master/python/mozscape.py
	expires := time.Now().Unix() + 3000
	signStr := fmt.Sprintf("%s\n%d", m.ApiKeyAccessID, expires)

	hm := hmac.New(sha1.New, []byte(m.ApiKeySecretKey))
	hm.Write([]byte(signStr))

	signature := base64.StdEncoding.EncodeToString(hm.Sum(nil))

	urlToQuery := fmt.Sprintf("http://lsapi.seomoz.com/linkscape/url-metrics/%s", url.QueryEscape(site))

	resp, body, err := client.Get(urlToQuery).
		QueryParam("AccessID", m.ApiKeyAccessID).
		QueryParam("Expires",
			fmt.Sprintf("%d", expires),
		).
		QueryParam("Signature", signature).
		QueryParam("Cols", "103079217152").
		SetHeaders(baseHeaders).
		SetHeader("Cache-Control", "max-age=0").Do()

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 || strings.Contains(body, "Permission denied") {
		return nil, errors.New("api error")
	}

	var result mozResultLocal

	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return nil, errors.New("could not unmarshal")
	}

	da, err := roundFloat(result.DA)
	if err != nil {
		return nil, errors.New("could not convert float to int")
	}

	pa, err := roundFloat(result.PA)
	if err != nil {
		return nil, errors.New("could not convert float to int")
	}

	return &MozResult{
		DA:        da,
		PA:        pa,
		Backlinks: result.Backlinks,
	}, nil
}

func roundFloat(in float64) (int, error) {
	s := fmt.Sprintf("%.0f", in)
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return i, nil
}
