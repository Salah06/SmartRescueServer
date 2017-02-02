#! /bin/bash

echo "ps aux | grep ./main | grep -v grep | awk '{print \$2}'"

cd /home/user/SI5/AL/SmartRescueServer/
go build main.go
./main


while true
do
	./main recover
done

