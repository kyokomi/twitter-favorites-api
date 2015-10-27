init:
	mkdir -p vendor

install:
	GOPATH=$(shell pwd)/vendor goapp get github.com/kyokomi/expcache
	GOPATH=$(shell pwd)/vendor goapp get github.com/unrolled/render

serve:
	GOPATH=$(shell pwd)/vendor goapp serve src/app.yaml
deploy:
	GOPATH=$(shell pwd)/vendor goapp deploy src/app.yaml

