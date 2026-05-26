.PHONY: build-web clean-web

build-web: ## Build frontend
	@echo "Building frontend..."
	cd web && npm ci && npm run build

clean-web: ## Clean build artifacts
	rm -rf web/node_modules internal/web/static