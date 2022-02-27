FROM golang AS build
ENV APP_PORT "50051"
ENV DB_URL "postgres://productmanager:productmanager@docker.internal/productstore"
ENV DB_DRIVER "postgres"
ENV location /go/src/productmanagement
WORKDIR ${location}/service

ADD ./service ${location}/service

RUN go get -d ./
RUN go install ./

RUN CGO_ENABLED=0 go build -o /bin/grpc-productmanagement-service

FROM scratch
COPY --from=build /bin/grpc-productmanagement-service /bin/grpc-productmanagement-service
ENTRYPOINT ["/bin/grpc-productmanagement-service"]
EXPOSE ${APP_PORT}