mocks:
	go tool mockgen -source=internal/service/albums/contracts.go -destination=internal/mocks/mock_albums/mock_repository.go
	go tool mockgen -source=internal/service/genres/contracts.go -destination=internal/mocks/mock_genres/mock_repository.go
	go tool mockgen -source=pkg/transactor/transactor.go -destination=internal/mocks/mock_transactor/mock_transactor.go

cover: 
	go test -count=1 -coverprofile=coverage.out ./internal/service/...
	go tool cover -html=coverage.out
	rm coverage.out

test-albums: 
	go test ./internal/service/albums 

test-genres: 
	go test ./internal/service/genres 