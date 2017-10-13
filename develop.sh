cd $GOPATH/src/github.com/SkyrisBactera/StudentVIEW
GOOS=linux gopherjs build -w github.com/SkyrisBactera/StudentVIEW &
GOOS=linux gopherjs serve --localmap ../ &
