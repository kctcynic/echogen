all: shell 

build:
	go build -o echogen

shell:
	mkdir -p .cache/pkg
	docker run --rm -it --user "${UID}:${UID}" -v "${PWD}/.cache/pkg":/go/pkg -v "${PWD}":/go/src/echogen -w /go/src/echogen -e GO111MODULE=on golang:1.12
