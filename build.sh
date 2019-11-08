GOOS=darwin GOARCH=amd64 go build -o bin/diskspeed-darwin diskspeed.go
GOOS=linux GOARCH=amd64 go build -o bin/diskspeed-linux diskspeed.go
