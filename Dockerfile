FROM oraclelinux:9

RUN  dnf -y install make git unzip oracle-instantclient-release-el9 && \
     dnf -y install oracle-instantclient19.19-basic.x86_64 oracle-instantclient19.19-devel.x86_64 oracle-instantclient19.19-sqlplus.x86_64 && \
	 dnf -y install go-toolset && \
	 dnf config-manager --set-enabled ol9_codeready_builder && \
	 dnf -y install mingw64-cpp mingw64-gcc mingw64-gcc-c++ && \
     rm -rf /var/cache/dnf

ENV PROTOC_ZIP_VERSION=3.19.1
RUN curl --silent -OL https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_ZIP_VERSION/protoc-$PROTOC_ZIP_VERSION-linux-x86_64.zip
RUN unzip -qq -o protoc-$PROTOC_ZIP_VERSION-linux-x86_64 -d /usr/local bin/protoc
RUN unzip -qq -o protoc-$PROTOC_ZIP_VERSION-linux-x86_64 -d /usr/local 'include/*'
RUN rm -f protoc-$PROTOC_ZIP_VERSION-linux-x86_64

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

WORKDIR /root/go/src
RUN git clone https://github.com/sumeetchhetri/sqldiffer
COPY oci8.pc /usr/lib64/pkgconfig/
ENV PATH="$PATH:/root/go/bin"
WORKDIR /root/go/src/sqldiffer
#COPY Makefile /root/go/src/sqldiffer/
#COPY protos/differ.proto /root/go/src/sqldiffer/protos/

RUN make all
RUN mv bin bin_with_orcl
RUN ./bin_with_orcl/schema_gen --help
RUN ./bin_with_orcl/diff_gen --help

RUN make all_wo_orcl
RUN mv bin bin_wo
RUN ./bin_wo/schema_gen --help
RUN ./bin_wo/diff_gen --help

WORKDIR /tmp
RUN curl --silent -OL https://download.oracle.com/otn_software/nt/instantclient/213000/instantclient-basic-windows.x64-21.3.0.0.0.zip
RUN curl --silent -OL https://download.oracle.com/otn_software/nt/instantclient/213000/instantclient-sdk-windows.x64-21.3.0.0.0.zip
RUN mkdir basic && mv instantclient-basic-windows.x64-21.3.0.0.0.zip basic/ && cd basic && unzip -qq instantclient-basic-windows.x64-21.3.0.0.0.zip
RUN unzip -qq instantclient-sdk-windows.x64-21.3.0.0.0.zip && mv instantclient_21_3/sdk basic/instantclient_21_3/
RUN mv basic/instantclient_21_3/ /root/instantclient_21_3_win
WORKDIR /root/go/src/sqldiffer
RUN git checkout schema_gen/schema_gen.go  go.mod  go.sum
RUN make clean
RUN make ORCL_INSCL_PATH=/root/instantclient_21_3_win CROSS_CC=x86_64-w64-mingw32-gcc CROSS_CXX=x86_64-w64-mingw32-g++ win_version
RUN ls -ltr /root/go/src/sqldiffer/bin
RUN mv bin bin_with_orcl_wver
RUN make clean
RUN make ORCL_INSCL_PATH=/root/instantclient_21_3_win CROSS_CC=x86_64-w64-mingw32-gcc CROSS_CXX=x86_64-w64-mingw32-g++ all_wo_orcl_wver
RUN ls -ltr /root/go/src/sqldiffer/bin
RUN mv bin bin_wo_wver
