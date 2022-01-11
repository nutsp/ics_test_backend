FROM golang:1.17-alpine
RUN apk update && apk add gcc git 
WORKDIR /app
COPY . /app
RUN go mod download
RUN go build -o main cmd/api/main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
RUN apk update \
	&& apk add tzdata \
	&& cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime \
	&& echo "Asia/Bangkok" >  /etc/timezone \
	&& echo "Asia/Bangkok" >  /etc/TZ \
	&& unset TZ

WORKDIR /app
COPY --from=0 /app/main .
CMD ["./main"]  
