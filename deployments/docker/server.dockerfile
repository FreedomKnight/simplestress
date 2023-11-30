from alpine

workdir /app

copy ../server/server /app

expose 50051

cmd ["./server"]
