package worker

type UrlConfig struct {
	url    string
	states map[string]bool
}

func NewUrlConfig(url string, states map[string]bool) UrlConfig {
	return UrlConfig{url: url, states: states}
}
