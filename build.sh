echo Building Linux
GOOS=linux GOARCH=amd64 go build main.go
mv main bin/discordeq-linux-x64
GOOS=linux GOARCH=386 go build main.go
mv main bin/discordeq-linux-x86
echo Building Windows
GOOS=windows GOARCH=amd64 go build main.go
mv main.exe bin/discordeq-windows-x64.exe
GOOS=windows GOARCH=386 go build main.go
mv main.exe bin/discordeq-windows-x86.exe
echo Building OSX
GOOS=darwin GOARCH=amd64 go build main.go
mv main bin/discordeq-darwin-x64
