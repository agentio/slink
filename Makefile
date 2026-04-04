all:	submodules bootstrap slink

bootstrap:	
	go run ./cmd/slink-generate lint
	go run ./cmd/slink-generate xrpc
	go run ./cmd/slink-generate call
	go run ./cmd/slink-generate check

slink:
	go install -tags jwx_es256k ./cmd/slink

manifest:
	go run ./cmd/slink-generate lint
	go run ./cmd/slink-generate xrpc -m sample-manifest.json
	go run ./cmd/slink-generate call -m sample-manifest.json
	go run ./cmd/slink-generate check -m sample-manifest.json
	go install ./cmd/slink

submodules:
	git submodule update --init --recursive
