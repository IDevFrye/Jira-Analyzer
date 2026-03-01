.PHONY: build build-backend build-frontend unit-test unit-test-backend unit-test-frontend

build: build-backend build-frontend

build-backend: unit-test-backend
	cd backend/endpointHandler && go build -o ../../bin/backend ./cmd/service
	cd backend/jiraConnector && go build -o ../../bin/jiraconnector ./cmd/service

build-frontend:
	cd frontend && npm run build

unit-test: unit-test-backend unit-test-frontend

unit-test-backend:
	cd backend/endpointHandler && go test ./... -cover
	cd backend/jiraConnector && go test ./... -cover

unit-test-frontend:
	cd frontend && npm test
