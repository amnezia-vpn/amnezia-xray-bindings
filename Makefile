LIB_OUT = libxray.a
HDR_OUT = libxray.h

all:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o $(LIB_OUT) -buildmode=c-archive

PHONY: clean
clean:
	rm -f $(LIB_OUT) $(HDR_OUT)