.PHONY: build
build: build_linux build_windows

.PHONY: build_windows
build_windows:
	wails build -platform windows/amd64 -webview2 embed

.PHONY: build_linux
build_linux:
	wails build -platform linux/amd64
