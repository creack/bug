all: build

DIST ?= "dist"

clean:
	rm -rf ${DIST}

build:
	mkdir -p ${DIST}
	go build -o ${DIST}/bugcli .

install: build
	mv ${DIST}/bugcli ${GOPATH}/bin/
