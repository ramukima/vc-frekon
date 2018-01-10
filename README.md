git clone https://github.com/ramukima/vc-frekon
cd vc-frekon

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

docker build -t vc-frekon .

docker run -e FACEBOX_URL="http://192.168.1.125:8080" -e MOTION_URL="http://192.168.1.230:8082" -p 9000:9000 -it vc-frekon

