from alpine

workdir /app

copy ../client/client /app

cmd ["./client"]
