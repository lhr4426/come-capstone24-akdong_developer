FROM golang:1.22.2

WORKDIR /server_login

COPY ./login_http_go /server_login

RUN go mod download && go mod verify

RUN go build -o main

CMD ["./main"]

# # 컨테이너 내에서 실행할 명령어 설정
# CMD ["go","run","/server/main.go"]
