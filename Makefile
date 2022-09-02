PROJECT_NAME = hashcat.launcher

CMD_DIR = ./cmd/

BIN_DIR = ./bin/

FRONTEND_DIR = ./frontend/

HASHCAT_VERSION = "6.2.6"

RESOURCES_DIR = ./resources/
RESOURCES_HASHCAT_DIR = ./resources/hashcat/

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
	npm --prefix $(FRONTEND_DIR)/$(PROJECT_NAME) install
	npm --prefix $(FRONTEND_DIR)/$(PROJECT_NAME) run build

build-linux:
	@echo "Building hashcat.launcher for Linux"
	@mkdir -p $(BIN_DIR_LINUX)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags $(LD_FLAGS_LINUX) -o $(BIN_DIR_LINUX)$(PROJECT_NAME) $(CMD_DIR)$(PROJECT_NAME)
	@7za -y x $(RESOURCES_HASHCAT_DIR)hashcat-$(HASHCAT_VERSION).7z -o$(BIN_DIR_LINUX)
	@mv $(BIN_DIR_LINUX)hashcat-*/ $(BIN_DIR_LINUX)hashcat/
	@7za -y a $(BIN_DIR)$(PROJECT_NAME)_$(GIT_TAG)_linux.7z $(BIN_DIR_LINUX)*

build-windows:
	@echo "Building hashcat.launcher for Windows"
	@mkdir -p $(BIN_DIR_WINDOWS)
	rsrc -arch amd64 -ico $(RESOURCES_DIR)Icon.ico -o $(CMD_DIR)$(PROJECT_NAME)/rsrc_windows_amd64.syso
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags $(LD_FLAGS_WINDOWS) -o $(BIN_DIR_WINDOWS)$(PROJECT_NAME).exe $(CMD_DIR)$(PROJECT_NAME)
	@7za -y x $(RESOURCES_HASHCAT_DIR)hashcat-$(HASHCAT_VERSION).7z -o$(BIN_DIR_WINDOWS)
	@mv $(BIN_DIR_WINDOWS)hashcat-*/ $(BIN_DIR_WINDOWS)hashcat/
	@7za -y a $(BIN_DIR)$(PROJECT_NAME)_$(GIT_TAG)_windows.7z $(BIN_DIR_WINDOWS)*

build-mac:
	@echo "Building hashcat.launcher for macOS"
	@mkdir -p $(BIN_DIR_MAC)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags $(LD_FLAGS_MAC) -o $(BIN_DIR_MAC)$(PROJECT_NAME) $(CMD_DIR)$(PROJECT_NAME)
	@7za -y x $(RESOURCES_HASHCAT_DIR)hashcat-$(HASHCAT_VERSION).7z -o$(BIN_DIR_MAC)
	@mv $(BIN_DIR_MAC)hashcat-*/ $(BIN_DIR_MAC)hashcat/
	@7za -y a $(BIN_DIR)$(PROJECT_NAME)_$(GIT_TAG)_mac.7z $(BIN_DIR_MAC)*

build: build-frontend build-linux build-windows build-mac
