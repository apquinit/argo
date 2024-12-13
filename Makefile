# Define the binary names
ARGO_BINARY=argo
CREATE_ARGO_APP_BINARY=create-argo-app

# Define the source paths
ARGO_SOURCE=main.go
CREATE_ARGO_APP_SOURCE=main.go

# Embed the argo binary and build the create-argo-app binary
build-create-argo-app:
	@echo "Building argo-cli binary..."
	@cd argo-cli && go build -o ../$(ARGO_BINARY) $(ARGO_SOURCE)
	@echo "Embedding argo-cli binary..."
	@mv $(ARGO_BINARY) create-argo-app/
	@echo "Building create-argo-app binary..."
	@cd create-argo-app && go build -o ./release/$(CREATE_ARGO_APP_BINARY) $(CREATE_ARGO_APP_SOURCE)

# Build the argo binary
build-argo-cli:
	@echo "Building argo-cli binary..."
	@cd argo-cli && go build -o ./release/argo-cli $(ARGO_SOURCE)

# Clean up binaries
clean:
	@echo "Cleaning up binaries..."
	@rm -rf ./argo-cli/release ./create-argo-app/release
