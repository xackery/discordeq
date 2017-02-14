set -e 
export VERSION="0.43"
export NAME="discordeq"
rm -rf bin/*
echo "Building Linux"
GOOS=linux GOARCH=amd64 go build main.go
mv main bin/$NAME-$VERSION-linux-x64
GOOS=linux GOARCH=386 go build main.go
mv main bin/$NAME-$VERSION-linux-x86
echo "Building Windows"
GOOS=windows GOARCH=amd64 go build main.go
mv main.exe bin/$NAME-$VERSION-windows-x64.exe
GOOS=windows GOARCH=386 go build main.go
mv main.exe bin/$NAME-$VERSION-windows-x86.exe
echo "Building OSX"
GOOS=darwin GOARCH=amd64 go build main.go
mv main bin/$NAME-$VERSION-osx-x64
