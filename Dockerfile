FROM golang


RUN mkdir /build

COPY . /build
COPY . /build

WORKDIR /build

RUN go mod download

RUN go build -o consumer

EXPOSE 8081

ENTRYPOINT ["./consumer"]

