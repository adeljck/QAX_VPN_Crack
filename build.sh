export LDFLAGS='-s -w '

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$LDFLAGS" -trimpath -o VPN_linux_amd64 main.go
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags="$LDFLAGS" -trimpath -o VPN_windows_386.exe  main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$LDFLAGS" -trimpath -o VPN_windows_amd64.exe  main.go
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags="$LDFLAGS" -trimpath -o VPN_windows_arm64.exe  main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="$LDFLAGS" -trimpath -o VPN_darwin_amd64 main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$LDFLAGS" -trimpath -o VPN_darwin_arm64 main.go

upx -9 VPN_linux_amd64
upx -9 VPN_windows_386.exe
upx -9 VPN_windows_amd64.exe
upx -9 VPN_windows_arm64.exe
upx -9 VPN_darwin_amd64
upx -9 VPN_darwin_arm64

zip VPN_linux_amd64.zip VPN_linux_amd64 config.yaml
zip VPN_windows_386.zip VPN_windows_386.exe config.yaml
zip VPN_windows_amd64.zip VPN_windows_amd64.exe config.yaml
zip VPN_windows_arm64.zip VPN_windows_arm64.exe config.yaml
zip VPN_darwin_amd64.zip VPN_darwin_amd64 config.yaml
zip VPN_darwin_arm64.zip VPN_darwin_arm64 config.yaml

rm -f VPN_linux_amd64
rm -f VPN_windows_386.exe
rm -f VPN_windows_amd64.exe
rm -f VPN_windows_arm64.exe
rm -f VPN_darwin_amd64
rm -f VPN_darwin_arm64