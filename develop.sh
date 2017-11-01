function kill() {
	pkill -e gopherjs
}
trap kill EXIT
cd $GOPATH/src/github.com/SkyrisBactera/StudentVIEW/public
GOOS=linux gopherjs build -w github.com/SkyrisBactera/StudentVIEW/public &
go run server.go
