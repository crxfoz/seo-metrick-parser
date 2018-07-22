Example project
---

#### Build API server:
```bash
# Clone the repo
git clone https://github.com/crxfoz/seo-metrick-parser.git
make run-server

# or manually

# Build binaries
make build-linux
# or build via docker container
make build-docker

# Create container
docker-compose build
# Start container
docker-compose up
```

#### Build worker:
If you would like to run api server and worker at the same host use: `make run`.
Otherwise you have to setup `worker.env` file or setup a network between containers.
```bash
git clone https://github.com/crxfoz/seo-metrick-parser.git

make build-linux
# or
make build-docker

cd ./deployment
docker build -t "worker" -f Dockerfile.worker .
docker run --env-file worker.env worker
```
Also don't forget to put yours ApiKeys in ./deployment/keys.env. 
```
# keys.env
MOZ_ACCESSID=member-de5df321
MOZ_SECRET_KEY=4532ea1ffd56e0200aabab8ad90ec21
```

### Add custom parser
* Go to parsers/repository.
* Create a new Go-file ` parsers/repository/awesome_metrick.go `
* Creaty a new function or method that match the prototype `func(*webclient.Webclient, string) (interface{}, error)`
```go
package repository
func GetAwesomeMetrick(client *webclient.Webclient, site string) (interface{}, error) {
	type AwesomeMetrickResult struct {
		Param1 string `json:"param_1"`
		Param2 string `json:"param_2"`
	}
	
	var result AwesomeMetrickResult
	
	return result, nil 
}
```
New type
```go
package repository
type AwesomeMetrickParser struct {
    ApiKey string
}

func (a *AwesomeMetrickParser) GetAwesomeMetrick(client *webclient.Webclient, site string) (interface{}, error) {
	
	// api request to get data using a.ApiKey
	
	type AwesomeMetrickResult struct {
		Param1 string `json:"param_1"`
		Param2 string `json:"param_2"`
	}
	
	var result AwesomeMetrickResult
	
	return result, nil 
}
```
* Register new parser. \
Edit `parsers/list.go`
```go
func init() {
	ParsersList = []Parser{
		// ...
		{
			Name:        "Awesome",
			Description: "Get awesome metrick",
			ParserFn:    repository.GetAwesomeMetrick,
			Timeout:     time.Second * 5,
			Status:      true,
		},
	}
}
```
```go
func init() {
	awesome := &AwesomeMetrickParser{ApiKey: os.Getenv("AWESOME_SERVICE_APIKEY")}
	
	ParsersList = []Parser{
		// ...
		{
			Name:        "Awesome",
			Description: "Get awesome metrick",
			ParserFn:    awesome.GetAwesomeMetrick,
			Timeout:     time.Second * 5,
			Status:      true,
		},
	}
}
```

#### Run tests
```bash
# Start test env
docker run --rm -p 5432:5432 -d --name seo_parser_postgres_test -e POSTGRES_DB=seo_parser_test postgres
docker run --rm -p 6379:6379 -d --name seo_parser_redis_test redis

# Import scheme
usql -f ./deployment/scheme.sql "pg://postgres@localhost:5432/seo_parser_test?sslmode=disable"

# Run tests
GOCACHE=off go test -v ./...

# Stop containers
docker stop seo_parser_postgres_test
docker stop seo_parser_redis_test

```