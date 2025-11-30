build-web:
	cd web && npm run build	

build-backend:
	cd api && go build -o reddit-content-analyzer cmd/main.go

run-web:
	cd web && npm run dev

run-backend:
	cd api && go run cmd/main.go
