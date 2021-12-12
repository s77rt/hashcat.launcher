PROJECT_NAME = hashcat.launcher

CMD_DIR = ./cmd/

BIN_DIR = ./bin/

FRONTEND_DIR = ./frontend/

BIN_DIR_MAC = ./bin/mac/
BIN_DIR_LINUX = ./bin/linux/
BIN_DIR_WINDOWS = ./bin/windows/

GIT_TAG = "$(shell git describe --tags | cut -c 2-)"

LD_FLAGS_LINUX = "-X 'github.com/s77rt/$(PROJECT_NAME).Version=$(GIT_TAG)'"
LD_FLAGS_WINDOWS = "-X 'github.com/s77rt/$(PROJECT_NAME).Version=$(GIT_TAG)' -H windowsgui"
LD_FLAGS_MAC = "-X 'github.com/s77rt/$(PROJECT_NAME).Version=$(GIT_TAG)'"

all: clean dep build

clean:
	@echo -n "Cleaning: "
	@rm -rf $(BIN_DIR)
	@echo "[OK]"

dep:
	@echo -n "Downloading Dependencies: "
	@go mod tidy
	@go install github.com/akavel/rsrc@latest
	@echo "[OK]"

build-frontend:
	@echo "Building frontend"
	npm --prefix $(FRONTEND_DIR)/$(PROJECT_NAME) run build

build-linux:
	@echo "Building hashcat.launcher for Linux"
	@mkdir -p $(BIN_DIR_LINUX)
	GOOS=linux GOARCH=amd64 go build -ldflags $(LD_FLAGS_LINUX) -o $(BIN_DIR_LINUX)$(PROJECT_NAME) $(CMD_DIR)$(PROJECT_NAME)
	@zip -j $(BIN_DIR)$(PROJECT_NAME)_$(GIT_TAG)_linux.zip $(BIN_DIR_LINUX)$(PROJECT_NAME)
	@echo "Building for Linux [OK]"

build-windows:
	@echo "Building hashcat.launcher for windows"
	@mkdir -p $(BIN_DIR_WINDOWS)
	rsrc -arch amd64 -ico Icon.ico -o $(CMD_DIR)$(PROJECT_NAME)/rsrc_windows_amd64.syso
	GOOS=windows GOARCH=amd64 go build -ldflags $(LD_FLAGS_WINDOWS) -o $(BIN_DIR_WINDOWS)$(PROJECT_NAME).exe $(CMD_DIR)$(PROJECT_NAME)
	@zip -j $(BIN_DIR)$(PROJECT_NAME)_$(GIT_TAG)_windows.zip $(BIN_DIR_WINDOWS)$(PROJECT_NAME).exe
	@echo "Building for windows [OK]"

build-mac:
	@echo "Building hashcat.launcher for macOS"
	@mkdir -p $(BIN_DIR_MAC)
	GOOS=darwin GOARCH=amd64 go build -ldflags $(LD_FLAGS_MAC) -o $(BIN_DIR_MAC)$(PROJECT_NAME) $(CMD_DIR)$(PROJECT_NAME)
	@zip -j $(BIN_DIR)$(PROJECT_NAME)_$(GIT_TAG)_mac.zip $(BIN_DIR_MAC)$(PROJECT_NAME)
	@echo "Building for macOS [OK]"

#build: build-linux build-windows build-mac
build: build-frontend build-linux build-windows