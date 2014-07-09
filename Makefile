all:
	go install github.com/piffio/owm/owmapi
	go install github.com/piffio/owm/owmworker

test:
	go test github.com/piffio/owm/owmapi
	go test github.com/piffio/owm/owmworker
