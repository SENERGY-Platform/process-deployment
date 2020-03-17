http api to deploy processes and there dependencies

## Tests
- all tests: `go test ./... -count=1`
- without test of kafka consumer and database-connection (with docker): `go test ./... -count=1 -short`
- the `-count=1` argument is used to prevent caching of the tests

## Deployment Model
to create uml for model:
- terminal: `goplantuml -show-aggregations lib/model | xclip -selection c`
- http://www.plantuml.com/plantuml/uml  -> `CTRL + V`

