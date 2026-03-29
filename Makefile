.PHONY: help integration-test integration-test-local integration-test-local-connector integration-test-local-backend integration-test-docker local-db-init local-up local-down

help:
	@echo "Targets:"
	@echo "  make local-db-init                 Apply DB schema to local PostgreSQL"
	@echo "  make local-up                      Start jiraConnector + backend + frontend locally"
	@echo "  make local-down                    Stop locally started services"
	@echo "  make integration-test-local        Run all local integration tests"
	@echo "  make integration-test-docker       Run integration tests in docker compose"

# Main entrypoint: local integration tests
integration-test: integration-test-local

integration-test-local: integration-test-local-connector integration-test-local-backend

# Mock-based connector integration tests (does not require started services)
integration-test-local-connector:
	cd backend/jiraConnector && go test ./tests/integration/jiraApi/... -tags=integration -v

# Full backend integration tests against running local services (localhost:8000)
integration-test-local-backend:
	powershell -NoProfile -Command "Set-Location backend/endpointHandler; $$env:BACKEND_BASE_URL='http://localhost:8000/api/v1'; go test ./tests/integration/... -tags=integration -v"

