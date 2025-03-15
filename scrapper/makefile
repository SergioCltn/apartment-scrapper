IMAGE_NAME = apartment-scrapper-builder
BINARY_NAME = apartment-scrapper
DOCKER = docker

# Default target
.PHONY: all
all: build-linux

# Build and extract for Linux
.PHONY: build-linux
build-linux:
	$(DOCKER) build --build-arg TARGET_OS=linux -t $(IMAGE_NAME) .
	$(DOCKER) create --name temp-builder $(IMAGE_NAME)
	$(DOCKER) cp temp-builder:/app/$(BINARY_NAME) ./$(BINARY_NAME)
	$(DOCKER) rm temp-builder
	@echo "Binary extracted to ./$(BINARY_NAME)"

# Build and extract for macOS
.PHONY: build-macos
build-macos:
	$(DOCKER) build --build-arg TARGET_OS=darwin -t $(IMAGE_NAME) .
	$(DOCKER) create --name temp-builder $(IMAGE_NAME)
	$(DOCKER) cp temp-builder:/app/$(BINARY_NAME) ./$(BINARY_NAME)
	$(DOCKER) rm temp-builder
	@echo "Binary extracted to ./$(BINARY_NAME)"

# Build and extract for Windows
.PHONY: build-windows
build-windows:
	$(DOCKER) build --build-arg TARGET_OS=windows -t $(IMAGE_NAME) .
	$(DOCKER) create --name temp-builder $(IMAGE_NAME)
	$(DOCKER) cp temp-builder:/app/$(BINARY_NAME).exe ./$(BINARY_NAME).exe
	$(DOCKER) rm temp-builder
	@echo "Binary extracted to ./$(BINARY_NAME).exe"

# Clean up Docker artifacts
.PHONY: clean
clean:
	-$(DOCKER) rm -f temp-builder 2>/dev/null || true
	-$(DOCKER) rmi -f $(IMAGE_NAME) 2>/dev/null || true
	-rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	@echo "Cleaned up Docker artifacts and binaries"
