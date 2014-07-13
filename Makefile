all:
	go install github.com/piffio/owm/owmapi
	go install github.com/piffio/owm/owmworker

test: all
	go test github.com/piffio/owm/owmapi
