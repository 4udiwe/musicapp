mocks:
	mockgen -source=internal/service/albums/contracts.go -destination=internal/repo/mock_albums/mock_repository.go
	mockgen -source=internal/service/genres/contracts.go -destination=internal/repo/mock_genres/mock_repository.go

cover: 
	go test -count=1 -coverprofile=coverage.out ./internal/service/...
	go tool cover -html=coverage.out
	rm coverage.out

test-albums: 
	go test ./internal/service/albums 

test-genres: 
	go test ./internal/service/genres 