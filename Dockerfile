FROM golang

ENV CLUSTERIP "cassandra"
ENV PORT ":8081"
ENV KAFKA_HOST "172.31.38.67:29092"
ENV TOPIC "myTopic1"

RUN mkdir /build

COPY . /build
COPY . /build

WORKDIR /build

RUN go mod download

RUN go build -o consumer

EXPOSE 8081

ENTRYPOINT ["./consumer"]

