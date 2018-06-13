FROM centos:7
LABEL author Xavier Coulon <xcoulon@redhat.com>

# install the qpid-proton-c-devel lib and symlink to `libqpid-proton-core`
# so that the `#cgo LDFLAGS: -lqpid-proton-core` instruction works when building
# the qpid.apache.org/* vendor libraries
RUN yum install -y \
        git \
        curl \
        golang-1.9.4 \
        qpid-proton-c-devel && \
    ln -s /usr/lib64/libqpid-proton.so /usr/lib64/libqpid-proton-core.so

# golang ENVs
ENV GOBIN=/go/bin
ENV GOPATH=/go

# compile the binary
ARG BIN_NAME
RUN mkdir -p /go/src/github.com/fabric8-services/fabric8-amqp
COPY . /go/src/github.com/fabric8-services/fabric8-amqp
WORKDIR /go/src/github.com/fabric8-services/fabric8-amqp
# RUN ${GOBIN}/dep ensure -v && \
RUN CGO_ENABLED=1 CGO_LDFLAGS="-L/usr/lib64/" CGO_CFLAGS="-I/usr/include/" \
	go build -o /usr/local/bin/${BIN_NAME} subscriber/*.go

ENV F8_USER_NAME=fabric8
RUN useradd --no-create-home -s /bin/bash ${F8_USER_NAME}
# From here onwards, any RUN, CMD, or ENTRYPOINT will be run under the following user
USER ${F8_USER_NAME}

ENTRYPOINT [ "subscriber" ]