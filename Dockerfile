FROM golang:1-alpine as build
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh
WORKDIR /app
COPY . /app
RUN go build -o chat ./cmd/.

FROM alpine:latest
COPY --from=build /app/static /app/static
COPY --from=build /app/chat /app/chat
WORKDIR /app
CMD [ "./chat" ]
EXPOSE 3000