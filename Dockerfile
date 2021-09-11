FROM golang

ENV CLUSTERIP "100.25.134.115"
ENV PORT ":8081"

RUN mkdir /build

COPY . /build
COPY . /build

WORKDIR /build

RUN go mod download

RUN go build -o cassandraclient

EXPOSE 8081

ENTRYPOINT ["./cassandraclient"]

~
~
