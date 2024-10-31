build:
	GOARCH=wasm GOOS=js go build -o output/app.wasm

serve:
	npx http-server ./output -o

run: build serve