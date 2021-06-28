PROJECT_NAME = hashcat.launcher
PACKAGE_NAME = hashcatlauncher

CMD_DIR = ./cmd/

BIN_DIR = ./bin/

GIT_TAG = "$(shell git describe --tags --abbrev=0 | cut -c2-)"

all: clean dep compile

clean:
	@echo -n "Cleaning: "
	@rm -rf $(BIN_DIR)
	@rm -f $(CMD_DIR)$(PROJECT_NAME)/$(PROJECT_NAME)
	@rm -f $(CMD_DIR)$(PROJECT_NAME)/$(PROJECT_NAME).exe
	@rm -f $(CMD_DIR)$(PROJECT_NAME)/*.syso
	@rm -f $(PROJECT_NAME).tar.gz
	@echo "[OK]"

dep:
	@echo -n "Downloading Dependencies: "
	@go get -d ./...
	@echo "[OK]"

compile:
	@echo "Compiling: "
	@mkdir -p $(BIN_DIR)

	# Linux (64bit)
	CC=gcc fyne package -appBuild 1 -appID s77rt.hashcat.launcher -appVersion $(GIT_TAG) -icon $(CMD_DIR)$(PROJECT_NAME)/../../Icon.png -os linux -sourceDir $(CMD_DIR)$(PROJECT_NAME)
	@zip -j $(BIN_DIR)$(PROJECT_NAME)_$(GIT_TAG)_linux.zip $(CMD_DIR)$(PROJECT_NAME)/$(PROJECT_NAME)
	@rm -f $(CMD_DIR)$(PROJECT_NAME)/$(PROJECT_NAME)
	@rm -f $(PROJECT_NAME).tar.gz

	# Windows (64bit)
	CC=x86_64-w64-mingw32-gcc fyne package -appBuild 1 -appID s77rt.hashcat.launcher -appVersion $(GIT_TAG) -icon $(CMD_DIR)$(PROJECT_NAME)/../../Icon.png -os windows -sourceDir $(CMD_DIR)$(PROJECT_NAME)
	@zip -j $(BIN_DIR)$(PROJECT_NAME)_$(GIT_TAG)_windows.zip $(CMD_DIR)$(PROJECT_NAME)/$(PROJECT_NAME).exe
	@rm -f $(CMD_DIR)$(PROJECT_NAME)/$(PROJECT_NAME).exe
	@rm -f $(CMD_DIR)$(PROJECT_NAME)/*.syso
	@rm -f $(PROJECT_NAME).tar.gz

	@echo "[OK]"
