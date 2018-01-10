FROM debian:jessie

# install curl
RUN apt-get update && apt-get install -qy curl git

# install go runtime
RUN curl -s https://dl.google.com/go/go1.9.2.linux-amd64.tar.gz | tar -C /usr/local -xz

# prepare go environment
ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV PATH $PATH:/usr/local/go/bin:/go/bin

# add the current build context
ADD . /go/src/github.com/ramukima/vc-frekon

# compile the binary
RUN cd /go/src/github.com/ramukima/vc-frekon/ && \ 
    go get -d -v github.com/matryer/way && \
    go get -d -v github.com/machinebox/sdk-go/facebox && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

EXPOSE 80

ENTRYPOINT ["/go/src/github.com/ramukima/vc-frekon/main"]
