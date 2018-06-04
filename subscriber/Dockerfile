FROM centos:7
LABEL author Xavier Coulon <xcoulon@redhat.com>

EXPOSE 8080
ARG BIN_DIR
ARG BIN_NAME
COPY ${BIN_DIR}/${BIN_NAME} /usr/local/bin/${BIN_NAME}

ENV F8_USER_NAME=fabric8
RUN useradd --no-create-home -s /bin/bash ${F8_USER_NAME}
# From here onwards, any RUN, CMD, or ENTRYPOINT will be run under the following user
USER ${F8_USER_NAME}

ENTRYPOINT [ "subscriber" ]