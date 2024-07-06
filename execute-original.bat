@echo OFF
cls

echo ### Stop previous rabbitmq
docker stop myrabbit
docker rm myrabbit

echo ### Remove images
echo y| docker volume prune
echo y| docker image prune
echo y| docker container prune

echo ### Execute new instance of rabbitmq
docker run -d --memory=6g --cpus=5.0 --name myrabbit -p 5672:5672 rabbitmq

echo ### Configure paths
set PATH=C:\WINDOWS\system32;C:\WINDOWS;C:\WINDOWS\System32\Wbem;C:\WINDOWS\System32\WindowsPowerShell\v1.0\;C:\WINDOWS\System32\OpenSSH\;C:\Program Files\Go\bin;C:\Program Files\Git\cmd;C:\Program Files\Docker\Docker\resources\bin;C:\Program Files\MATLAB\R2023b\bin;C:\Users\user\AppData\Local\Microsoft\WindowsApps;C:\Program Files\JetBrains\GoLand 2022.1.2\bin;C:\MinGW\bin;C:\Users\user\AppData\Local\Google\Cloud SDK\google-cloud-sdk\bin;C:\Program Files\Docker\Docker
set PATH=%PATH%;C:\Users\user\go;C:\Users\user\go\adaptive-moms\subscriber

echo ### Configure Go
set GO111MODULE=on
set GOPATH=C:\Users\user\go\adaptive-moms
set GOROOT=C:\Program Files\Go

timeout /t 20

cd C:\Users\user\go\adaptive-moms

docker build --tag subscriber .
START /B docker run --rm --name some-subscriber --memory="1g" --cpus="1.0" -v C:\Users\user\go\adaptive-moms\data:/app/data subscriber

cd C:\Users\user\go\adaptive-moms