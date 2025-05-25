FROM golang:bullseye
WORKDIR /app

COPY . .
ENV TZ=Asia/Shanghai

RUN go env -w GO111MODULE=on
RUN go env -w  GOPROXY=https://goproxy.cn,direct

RUN go mod tidy
RUN go build -o server .

EXPOSE 3000

CMD ["./server"]