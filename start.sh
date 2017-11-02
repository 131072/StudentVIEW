kill () {
	pkill -e gopherjs
	rm ./public/StudentVIEW.js
	rm ./public/StudentVIEW.js.map
}
trap kill EXIT
cd $GOPATH/src/github.com/SkyrisBactera/StudentVIEW/public
GOOS=linux gopherjs build -w github.com/SkyrisBactera/StudentVIEW/public -o StudentVIEW.js &
cd $GOPATH/src/github.com/SkyrisBactera/StudentVIEW
sudo env "PATH=$PATH" go run server.go
