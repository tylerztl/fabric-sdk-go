FROM library/golang AS build

MAINTAINER tailinzhang1993@gmail.com

ENV APP_DIR /go/src/fabric-sdk-go
RUN mkdir -p $APP_DIR
WORKDIR $APP_DIR
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o fabric-server .
ENTRYPOINT ./fabric-server start

# Create a minimized Docker mirror
FROM scratch AS prod

COPY --from=build /go/src/fabric-sdk-go/fabric-server /fabric-server
EXPOSE 8080
CMD ["/fabric-server", "start"]
