FROM ubuntu:20.04

ENV REPRO_NAME  blaze
ENV REPRO_MNT   /mnt/${REPRO_NAME}
ENV REPRO_USER  repro
ENV REPRO_UID   1000
ENV REPRO_GID   1000

RUN echo '***** Update packages *****'                                      \
    && apt-get -y update

RUN echo '***** Set timezone noninteractively install JDK8 *****'
RUN DEBIAN_FRONTEND="noninteractive" TZ="America/Los_Angeles"               \
    apt -y install tzdata openjdk-8-jdk

RUN echo '***** Install packages required for creating this image *****'    \
    && apt -y install apt-utils wget curl makepasswd gcc make git           \
                                                                            \
    && echo '***** Install command-line utility packages *****'             \
    && apt -y install sudo man less file tree jq graphviz libxml2-utils

ENV GO_VERSION       1.16
ENV GO_DOWNLOADS_URL https://dl.google.com/go
ENV GO_ARCHIVE       go${GO_VERSION}.linux-amd64.tar.gz

RUN echo '****** Install Go development tools *****'                        \
    && wget ${GO_DOWNLOADS_URL}/${GO_ARCHIVE} -O /tmp/${GO_ARCHIVE}         \
    && tar -xzf /tmp/${GO_ARCHIVE} -C /usr/local

RUN echo '***** Add the REPRO user and group *****'                         \
    && groupadd ${REPRO_USER} --gid ${REPRO_GID}                            \
    && useradd ${REPRO_USER} --uid ${REPRO_UID} --gid ${REPRO_GID}          \
        --shell /bin/bash                                                   \
        --create-home                                                       \
        -p `echo repro | makepasswd --crypt-md5 --clearfrom - | cut -b8-`   \
    && echo "${REPRO_USER} ALL=(ALL) NOPASSWD: ALL"                         \
            > /etc/sudoers.d/${REPRO_USER}                                  \
    && chmod 0440 /etc/sudoers.d/repro

ENV HOME /home/${REPRO_USER}
ENV HOME_BIN ${HOME}/bin
ENV PATH ${HOME_BIN}:${PATH}
ENV BASHRC ${HOME}/.bashrc
USER  ${REPRO_USER}
WORKDIR $HOME

# (Next version: BLAZEGRAPH_VER 2_1_6_RC)
ENV BLAZEGRAPH_VER RELEASE_2_1_5
ENV BLAZEGRAPH_RELEASES https://github.com/blazegraph/database/releases
ENV BLAZEGRAPH_DOWNLOAD_DIR ${BLAZEGRAPH_RELEASES}/download/BLAZEGRAPH_${BLAZEGRAPH_VER}
ENV BLAZEGRAPH_DOWNLOAD_JAR ${BLAZEGRAPH_DOWNLOAD_DIR}/blazegraph.jar
ENV BLAZEGRAPH_JAR $HOME/blazegraph-${BLAZEGRAPH_VER}.jar

RUN echo '***** Download Blazegraph jar *****'                              \
    && wget -O ${BLAZEGRAPH_JAR} ${BLAZEGRAPH_DOWNLOAD_JAR}

ENV BLAZEGRAPH_DOT_DIR ${REPRO_MNT}/.blazegraph
ENV BLAZEGRAPH_PROPERTY_FILE=${BLAZEGRAPH_DOT_DIR}/.properties
ENV BLAZEGRAPH_OPTIONS "-server -Xmx4g -Dbigdata.propertyFile=${BLAZEGRAPH_PROPERTY_FILE}"
ENV BLAZEGRAPH_CMD "java ${BLAZEGRAPH_OPTIONS} -jar ${BLAZEGRAPH_JAR}"
ENV BLAZEGRAPH_LOG ${BLAZEGRAPH_DOT_DIR}/blazegraph_`date +%s`.log

RUN echo '***** Start Blazegraph at login *****'                            \
    && echo "source .bashrc" >> .bash_profile                               \
    && echo "cd ${BLAZEGRAPH_DOT_DIR} && ${BLAZEGRAPH_CMD} 2>&1 > ${BLAZEGRAPH_LOG} &" >> .bash_profile

ENV GOPATH ${REPRO_MNT}/.gopath

RUN echo "export PATH=${REPRO_MNT}/binaries:${GOPATH}/bin:/usr/local/go/bin:$PATH" >> ${BASHRC}
RUN echo "export IN_RUNNING_REPRO=${REPRO_NAME}" >> ${BASHRC}
RUN echo "cd ${REPRO_MNT}" >> ${BASHRC}

CMD  /bin/bash -il

