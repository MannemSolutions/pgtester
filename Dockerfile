FROM golang:alpine AS build-stage
WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine AS export-stage
#RUN mkdir /lib64 && cd /lib64 && find /lib -type l -name "libc.musl*" -exec cp -s {} . \; && ls
COPY --from=build-stage /go/bin/pgtester /usr/bin/
COPY testdata /etc/pgtester/tests
CMD pgtester -d /etc/pgtester/tests
