function kill() {
	pkill -e gopherjs
}
trap kill EXIT
cd $GOPATH/src/github.com/SkyrisBactera/StudentVIEW
GOOS=linux gopherjs build -w github.com/SkyrisBactera/StudentVIEW &
GOOS=linux gopherjs serve --localmap ../
