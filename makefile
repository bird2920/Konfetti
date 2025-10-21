BINARY_NAME=konfetti
SRC=main.go
DIST=dist
VERSION ?= dev

.PHONY: all win linux osx linux-arm osx-arm clean version

all: preflight win linux osx
	@echo "🎉 All done! Your binaries are ready to party! 🎊"

preflight:
	@mkdir -p $(DIST)
	@echo "✨ Prepping the confetti cannons... ✨"

win:
	@echo "🎈 Building for Windows... Hope you like .exe files!"
	GOOS=windows GOARCH=amd64 go build -o $(DIST)/$(BINARY_NAME)-win.exe $(SRC)
	@echo "🥳 Windows build complete! Go break the registry!"

linux:
	@echo "🎈 Building for Linux... Penguins love confetti!"
	GOOS=linux GOARCH=amd64 go build -o $(DIST)/$(BINARY_NAME)-linux $(SRC)
	@echo "🐧 Linux build done! Tux is proud."

osx:
	@echo "🎈 Building for macOS... For all the fancy folks."
	GOOS=darwin GOARCH=amd64 go build -o $(DIST)/$(BINARY_NAME)-mac $(SRC)
	@echo "🍏 macOS build ready! Now go be creative."

linux-arm:
	@echo "🤖 Building for Linux ARM... Raspberry Pi confetti incoming!"
	GOOS=linux GOARCH=arm64 go build -o $(DIST)/$(BINARY_NAME)-linux-arm64 $(SRC)
	@echo "🍓 Linux ARM build done! Pi never tasted so festive."

osx-arm:
	@echo "🦄 Building for macOS ARM (M1/M2)... Apple silicon, meet confetti!"
	GOOS=darwin GOARCH=arm64 go build -o $(DIST)/$(BINARY_NAME)-mac-arm64 $(SRC)
	@echo "🚀 macOS ARM build ready! Unleash the unicorns."

version:
	@echo "🔢 Injecting version: $(VERSION) into the binary with confetti magic!"
	GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=$(VERSION)" -o $(DIST)/$(BINARY_NAME)-linux-versioned $(SRC)
	@echo "🪄 Versioned Linux build done! $(VERSION) never looked so good."

clean:
	@echo "🧹 Cleaning up all the confetti... (and binaries)"
	rm -rf $(DIST)
	@echo "😢 All clean! (But so empty...)"