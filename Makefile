
setup: dev/workspace
	cd extension && yarn link qmux qrpc
	cd extension && yarn compile
	go install ./cmd/tractor

dev/workspace:
	mkdir -p dev
	cp -r data/workspace dev/workspace