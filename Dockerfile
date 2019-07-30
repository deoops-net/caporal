# mutilstage 1 for compiling golang add zone datas
FROM golang:1.12.1-alpine3.9 as build-env
# All these steps will be cached
RUN mkdir /app
WORKDIR /app
# <- COPY go.mod and go.sum files to the workspace
COPY go.mod . 
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN apk update
RUN apk add git tzdata zip
ENV GOPROXY https://goproxy.io
RUN go mod download
# COPY the source code as the last step
COPY . .
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/server

# multistage 2 for run binary with TZ
# <- Second step to build minimal image
FROM scratch 
# timezone
COPY --from=build-env /usr/share/zoneinfo /usr/share/zoneinfo
# ENV ZONEINFO /zoneinfo.zip
ENV TZ Asia/Shanghai

COPY --from=build-env /go/bin/server /go/bin/server
ENV PROJECT_LEVEL production
ENTRYPOINT ["/go/bin/server"]
