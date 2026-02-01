FROM golang:1.24-alpine

# Устанавливаем системный компилятор
RUN apk add --no-cache protobuf protobuf-dev

# Устанавливаем плагины Go и gRPC
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# В образе golang бинарники скачиваются в /go/bin.
# Мы принудительно добавляем эту папку в системный PATH.
ENV PATH="/go/bin:${PATH}"

WORKDIR /workspace

ENTRYPOINT ["protoc"]