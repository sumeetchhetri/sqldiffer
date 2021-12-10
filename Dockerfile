FROM oraclelinux:8

RUN  dnf -y install make git unzip oracle-instantclient-release-el8 && \
     dnf -y install oracle-instantclient-basic oracle-instantclient-devel oracle-instantclient-sqlplus && \
	 dnf -y module install go-toolset && \
     rm -rf /var/cache/dnf

ENV PROTOC_ZIP_VERSION=3.19.1
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_ZIP_VERSION/protoc-$PROTOC_ZIP_VERSION-linux-x86_64.zip
RUN unzip -o protoc-$PROTOC_ZIP_VERSION-linux-x86_64 -d /usr/local bin/protoc
RUN unzip -o protoc-$PROTOC_ZIP_VERSION-linux-x86_64 -d /usr/local 'include/*'
RUN rm -f protoc-$PROTOC_ZIP_VERSION-linux-x86_64

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

#ENV GO111MODULE=off
WORKDIR /root/go/src
RUN git clone https://github.com/sumeetchhetri/sqldiffer
COPY oci8.pc /usr/lib64/pkgconfig/
ENV PATH="$PATH:/root/go/bin"
WORKDIR /root/go/src/sqldiffer
COPY Makefile /root/go/src/sqldiffer/
COPY protos/differ.proto /root/go/src/sqldiffer/protos/
RUN GO111MODULE=off make all

RUN ./bin/schema_gen --help
RUN ./bin/diff_gen --help
