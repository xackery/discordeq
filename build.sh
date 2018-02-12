set -e 
export VERSION="0.50"
export NAME="discordeq"
echo "Preparing discordeq v${VERSION}"
rm -rf bin/*
echo "Building Linux"
echo "x64"
GOOS=linux GOARCH=amd64 go build main.go
mv main bin/$NAME-$VERSION-linux-x64
echo "x86"
GOOS=linux GOARCH=386 go build main.go
mv main bin/$NAME-$VERSION-linux-x86
echo "Building Windows"
echo "x64"
GOOS=windows GOARCH=amd64 go build main.go
mv main.exe bin/$NAME-$VERSION-windows-x64.exe
echo "x86"
GOOS=windows GOARCH=386 go build main.go
mv main.exe bin/$NAME-$VERSION-windows-x86.exe
echo "Building OSX"
echo "x64"
GOOS=darwin GOARCH=amd64 go build main.go
mv main bin/$NAME-$VERSION-osx-x64
echo "Completed."
