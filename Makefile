all: html/main.wasm.br html/wasm_exec.js

clean:
	rm -f ./html/main.wasm ./html/main.wasm.br html/wasm_exec.js

html/wasm_exec.js: $(shell go env GOROOT)/misc/wasm/wasm_exec.js
	cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js ./html/wasm_exec.js

html/main.wasm.br: html/main.wasm
	rm -f ./html/main.wasm.br
	brotli -o ./html/main.wasm.br ./html/main.wasm

html/main.wasm: main.go Makefile
	CGO_ENABLED=0 GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o ./html/main.wasm ./main.go

serve:
	go run server.go