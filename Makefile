all:	submodules bootstrap slink

bootstrap:	
	go run ./cmd/slink-generate lint -i lexicons-bluesky
	go run ./cmd/slink-generate xrpc -i lexicons-bluesky
	go run ./cmd/slink-generate call -i lexicons-bluesky
	go run ./cmd/slink-generate check -i lexicons-bluesky

slink:
	go install -tags jwx_es256k ./cmd/slink

manifest:
	go run ./cmd/slink-generate lint -i lexicons-bluesky
	go run ./cmd/slink-generate xrpc -m sample-manifest.json -i lexicons-bluesky
	go run ./cmd/slink-generate call -m sample-manifest.json -i lexicons-bluesky
	go run ./cmd/slink-generate check -m sample-manifest.json -i lexicons-bluesky
	go install ./cmd/slink

submodules:
	git submodule update --init --recursive
