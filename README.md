# StudentVIEW
### Try it out at http://skyrisbactera.com/studentview
#### An unofficial client for StudentVUE, because I was disatisifed with the GUI and user experience.

## Build instructions

Download this project and it's dependencies with:
```bash
go get -u github.com/gopherjs/gopherjs
go get -u github.com/SkyrisBactera/govue
go get -u github.com/SkyrisBactera/StudentVIEW
go get -u github.com/fabioberger/cookie
go get -u github.com/go-humble/locstor
go get -u github.com/gopherjs/jquery
```
To get started developing, you should use:
```bash
sudo setcap CAP_NET_BIND_SERVICE=+eip $GOROOT/bin/go #Allows access to port 80 without root
cd $GOPATH/src/github.com/SkyrisBactera/StudentVIEW
bash start.sh
```
## TODO
* Make password saving method more secure
* Optimize code to get grades faster. (Goroutines are you friend here!)
* Make the UI look better

## Donate!
If you appreciate this project, feel free to send me some Monero at:
```424kUkn7pU7SVF4qmbzVt6fLPRFw6Q1v7RmVq4qE9sE2V5dyhZoGdWHDUFui85SfhVTfmDN3CzcYDQSwo41z3AuLDzU3qQt```

or just use PayPal:
```davis@skyrisbactera.com```