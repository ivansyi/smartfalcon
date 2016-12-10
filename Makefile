.PHONY: default agent clean all
export GOPATH:=$(shell pwd)

BUILDTAGS=debug
default: all

all: fmt falcon-agent

deps:
	go get github.com/toolkits/nux
	go get github.com/toolkits/net
	go get github.com/toolkits/time
	go get github.com/toolkits/core

fmt:
	go fmt falcon-agent/...

falcon-agent:
	go install -tags '$(BUILDTAGS)' falcon-agent
	cp src/falcon-agent/control bin/agent-ctl
	cp src/falcon-agent/cfg.example.json bin/agent.cfg

bin/go-bindata:
	GOOS="" GOARCH="" go get github.com/jteeuwen/go-bindata/go-bindata

release-agent: BUILDTAGS=release
release-agent: falcon-agent

release-all: fmt release-agent

clean: clean-agent

clean-agent:
	go clean -i -r falcon-agent/...
	rm -f bin/agent-ctl
	rm -f bin/agent.cfg
