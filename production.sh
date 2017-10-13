cd ~/go/src/github.com/SkyrisBactera/StudentVIEW
gopherjs build -w github.com/SkyrisBactera/StudentVIEW &
gopherjs serve --localmap ../ &
while [ true ]
do
    git pull    
    sleep 60
done
