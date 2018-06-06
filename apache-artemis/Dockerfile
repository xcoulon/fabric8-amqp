FROM openjdk:9

LABEL maintainer="xcoulon@redhat.com"
LABEL product="Apache Active MQ Artemis server"

ENV ARTEMIS_VERSION=2.6.0
LABEL version="${VERSION}"

# see install docs in http://activemq.apache.org/artemis/docs/latest/using-server.html
# install the Apache ActiveMQ Artemis at the root 
WORKDIR /
ENV ARTEMIS_HOME=/var/lib/apache-artemis
RUN wget -v http://mirrors.ircam.fr/pub/apache/activemq/activemq-artemis/${ARTEMIS_VERSION}/apache-artemis-${ARTEMIS_VERSION}-bin.tar.gz && \
    tar zxvf apache-artemis-${ARTEMIS_VERSION}-bin.tar.gz && \
    mv /apache-artemis-${ARTEMIS_VERSION} ${ARTEMIS_HOME}

# create the broker structure in the image
ENV ADMIN_USERNAME=admin
ENV ADMIN_PASSWORD=password
ENV BROKER_DIR=${ARTEMIS_HOME}/data/fabric8-broker
RUN mkdir -p ${BROKER_DIR} && \
    ${ARTEMIS_HOME}/bin/artemis create \
    --user=${ADMIN_USERNAME} --password=${ADMIN_PASSWORD} \
    --allow-anonymous --http-host=0.0.0.0 \
    ${BROKER_DIR}
# modify the etc/jolokia-access.xml (TODO: avoid using hard-coded values)
COPY apache-artemis/jolokia-access.xml ${BROKER_DIR}/etc
COPY apache-artemis/broker.xml ${BROKER_DIR}/etc




# see https://docs.openshift.com/container-platform/3.9/creating_images/guidelines.html#openshift-specific-guidelines
RUN chgrp -R 0 ${ARTEMIS_HOME} && \
    chmod -R g=u ${ARTEMIS_HOME}

WORKDIR ${BROKER_DIR}

# Web console
EXPOSE 8161

# AMPQ, etc.
EXPOSE 61616
