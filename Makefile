# Define the binary names
ARGO_BINARY=argo
CREATE_ARGO_APP_BINARY=create-argo-app

# Define the source paths
ARGO_SOURCE=argo-cli/main.go
CREATE_ARGO_APP_SOURCE=create-argo-app/main.go

# Build the argo binary
build-argo-cli:
	@echo "Building argo-cli binary..."
	@cd argo-cli && go build -o ../$(ARGO_BINARY) $(ARGO_SOURCE)

# Embed the argo binary and build the create-argo-app binary
build-create-argo-app: build-argo-cli
	@echo "Embedding argo-cli binary and building create-argo-app..."
	@cd create-argo-app && go build -o ../$(CREATE_ARGO_APP_BINARY) $(CREATE_ARGO_APP_SOURCE)

# Clean up binaries
clean:
	@echo "Cleaning up binaries..."
	@rm -f $(CREATE_ARGO_APP_BINARY) $(ARGO_BINARY)

# Default target
all: build-create-argo-app