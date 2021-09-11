FROM golang

ENV CLUSTERIP "cassandra"
ENV PORT ":8081"

RUN mkdir /build

COPY . /build
COPY . /build

WORKDIR /build

RUN go mod download

RUN go build -o consumer

EXPOSE 8081

ENTRYPOINT ["./consumer"]

