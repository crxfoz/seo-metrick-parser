package worker

type Reporter map[string]map[string]interface{}

func (r Reporter) AddUrl(url string) map[string]interface{} {
	r[url] = make(map[string]interface{})
	return r[url]
}

func (r Reporter) AddForUrl(url map[string]interface{}, metrick string, data interface{}) {
	url[metrick] = data
}
