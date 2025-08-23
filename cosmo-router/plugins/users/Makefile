
.PHONY: build test generate install-wgc

make: build

install-wgc:
	@which wgc > /dev/null 2>&1 || npm install -g wgc@latest

test: install-wgc
	wgc router plugin test .

generate: install-wgc
	wgc router plugin generate .

publish: generate
	wgc router plugin publish .

build: install-wgc
	wgc router plugin build . --debug
