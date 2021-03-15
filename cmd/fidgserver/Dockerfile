ARG BUILD_PLATFORM=linux/amd64
ARG TARGET_PLATFORM=linux/amd64

FROM --platform=${BUILD_PLATFORM} golang:1.16-alpine AS builder
ARG COMPILE_GOOS=linux
ARG COMPILE_GOARCH=amd64
COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=${COMPILE_GOOS} GOARCH=${COMPILE_GOARCH} go build -ldflags '-extldflags "-static"' -o fidgserver ./cmd/fidgserver/main.go
RUN mv /build/fidgserver /

FROM --platform=${TARGET_PLATFORM} alpine
WORKDIR /
COPY --from=builder /fidgserver /fidgserver
ENTRYPOINT [ "/fidgserver" ]
CMD ["/fidgserver"]