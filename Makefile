LIPO := /usr/bin/x86_64-apple-darwin-lipo
GIT_COMMIT_HASH := $(shell git rev-list -1 HEAD)
LDFLAGS_COMMON := -s -w
all: clean \
	ugrsical-windows-amd64 ugrsical-linux-amd64 ugrsical-linux-arm64 ugrsical-darwin-amd64 ugrsical-darwin-arm64 \
	ugrsicalsrv-windows-amd64 ugrsicalsrv-linux-amd64 ugrsicalsrv-linux-arm64 ugrsicalsrv-darwin-amd64 ugrsicalsrv-darwin-arm64 \
	merge-macos-binary

ugrsical-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS_COMMON) -X ugrs-ical/internal/ugrsical.version=$(GIT_COMMIT_HASH)" -o build/ugrsical-windows-amd64.exe ugrs-ical/cmd/ugrsical

ugrsical-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS_COMMON) -X ugrs-ical/internal/ugrsical.version=$(GIT_COMMIT_HASH)" -o build/ugrsical-linux-amd64 ugrs-ical/cmd/ugrsical

ugrsical-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS_COMMON) -X ugrs-ical/internal/ugrsical.version=$(GIT_COMMIT_HASH)" -o build/ugrsical-linux-arm64 ugrs-ical/cmd/ugrsical

ugrsical-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS_COMMON) -X ugrs-ical/internal/ugrsical.version=$(GIT_COMMIT_HASH)" -o build/ugrsical-darwin-amd64 ugrs-ical/cmd/ugrsical

ugrsical-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS_COMMON) -X ugrs-ical/internal/ugrsical.version=$(GIT_COMMIT_HASH)" -o build/ugrsical-darwin-arm64 ugrs-ical/cmd/ugrsical

ugrsicalsrv-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS_COMMON)" -o build/ugrsicalsrv-windows-amd64.exe ugrs-ical/cmd/ugrsicalsrv

ugrsicalsrv-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS_COMMON)" -o build/ugrsicalsrv-linux-amd64 ugrs-ical/cmd/ugrsicalsrv

ugrsicalsrv-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS_COMMON)" -o build/ugrsicalsrv-linux-arm64 ugrs-ical/cmd/ugrsicalsrv

ugrsicalsrv-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS_COMMON)" -o build/ugrsicalsrv-darwin-amd64 ugrs-ical/cmd/ugrsicalsrv

ugrsicalsrv-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS_COMMON)" -o build/ugrsicalsrv-darwin-arm64 ugrs-ical/cmd/ugrsicalsrv

merge-macos-binary:
	$(LIPO) -create build/ugrsical-darwin-amd64 build/ugrsical-darwin-arm64 -o build/ugrsical-darwin-universal
	$(LIPO) -create build/ugrsicalsrv-darwin-amd64 build/ugrsicalsrv-darwin-arm64 -o build/ugrsicalsrv-darwin-universal

clean:
	-rm -f build/*

.PHONY: all clean \
	ugrsical-windows-amd64 ugrsical-linux-amd64 ugrsical-linux-arm64 ugrsical-darwin-amd64 ugrsical-darwin-arm64 \
	ugrsicalsrv-windows-amd64 ugrsicalsrv-linux-amd64 ugrsicalsrv-linux-arm64 ugrsicalsrv-darwin-amd64 ugrsicalsrv-darwin-arm64 \
	merge-macos-binary