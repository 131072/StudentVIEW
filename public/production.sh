cd ~/go/src/github.com/SkyrisBactera/StudentVIEW/public
gopherjs build -w github.com/SkyrisBactera/StudentVIEW/public &
gopherjs serve --localmap ../ &
while [ true ]
do
    git pull    
    sleep 60
done
