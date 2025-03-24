
mkdir -p build

echo "Compilando para macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -o build/sstighost_mac

echo "Compilando para macOS (M1)..."
GOOS=darwin GOARCH=arm64 go build -o build/sstighost_macm1

echo "Compilando para Windows..."
GOOS=windows GOARCH=amd64 go build -o build/sstighost.exe

echo "Compilando para Linux..."
GOOS=linux GOARCH=amd64 go build -o build/sstighost_linux

echo "Build concluído!"
