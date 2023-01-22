To test the API locally clone the repository and use below commands:
1. docker-compose up -d
2. migrate -path db/migrations -database "postgresql://admin:password123@localhost:6500/golang_postgres?sslmode=disable" -verbose up
3. go run  main.go

And then open Postman and import API from postman_import.json file
