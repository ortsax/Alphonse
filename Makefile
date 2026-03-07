SHELL      := powershell.exe
.SHELLFLAGS := -NoProfile -Command

BINARY  := alphonse
VERSION ?= 0.0.1
COMMIT  := $(shell git rev-parse --short HEAD 2>$$null)
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT)

.PHONY: build run whatsmeow icon release tag help

## build: compile the bot binary
build:
	go build -ldflags "$(LDFLAGS)" -trimpath -o $(BINARY).exe .

## run: run the bot
run:
	go run .

## whatsmeow: pull latest whatsmeow, re-apply all patches, and verify the build
whatsmeow:
	go get go.mau.fi/whatsmeow@latest
	go mod tidy
	pwsh -NoProfile -File scripts/patch-whatsmeow.ps1
	go build ./...

## icon: regenerate Windows resource files (.syso) from winres/winres.json  (requires go-winres)
icon:
	go-winres make --product-version $(VERSION).0 --file-version $(VERSION).0

## release: build release archives for all platforms into dist/  (usage: make release VERSION=x.y.z)
release:
	pwsh -NoProfile -File scripts/release.ps1 -Version $(VERSION)

## tag: create and push an annotated git tag  (usage: make tag VERSION=x.y.z)
tag:
	git tag -a "v$(VERSION)" -m "Release v$(VERSION)"
	git push origin "v$(VERSION)"
	@Write-Host "Tagged and pushed v$(VERSION)"

## help: list available targets
help:
	@Select-String "^## " Makefile | ForEach-Object { $$_.Line -replace "^## ", "" }
