

MAINNAME=gin_demo

all: build

build: buildversion
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${MAINNAME}
	test -d ./bin || mkdir -p ./bin
	\cp -f ${MAINNAME} ./bin
	rm -f ${MAINNAME}

buildversion:
	@echo package main > buildversion.go
	@echo const BuildVersion = \"`git rev-parse HEAD 2>/dev/null | cut -c 1-8`\" >> buildversion.go
	@echo const BuildBranch = \"`git symbolic-ref --short -q HEAD`\" >> buildversion.go
	@echo const BuildDate = \"`date +'%F %T %z'`\" >> buildversion.go

version:
	@echo `git rev-parse HEAD 2>/dev/null | cut -c 1-8`

clean:
	rm -f ./bin/${MAINNAME}
	rm -f buildversion.go

