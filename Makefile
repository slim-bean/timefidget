


build:
	GOARM=7 GOARCH=arm GOOS=linux go build -o rgbled main.go
	ssh pi@rgbled.edjusted.com "sudo pkill --signal SIGINT rgbled || true"
	scp ./rgbled pi@rgbled.edjusted.com:
	ssh pi@rgbled.edjusted.com "sudo nohup ./rgbled > nohup1.out 2>&1 </dev/null &"