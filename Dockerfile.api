FROM golang:1.13-alpine3.11 as build
WORKDIR /app
COPY . /app
RUN apk add make
RUN make api 

FROM alpine:3.11.3
WORKDIR /app
COPY --from=build  /app/api .
CMD ["./api"]
