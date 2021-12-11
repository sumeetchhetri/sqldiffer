
deps:
	cd protos && protoc --go_out=. *.proto && cd ..

schema_gen: clean deps
	go build -gcflags="-l=4" -ldflags="-s -w" -o bin/schema_gen ./schema_gen

diff_gen: clean deps
	go build -gcflags="-l=4" -ldflags="-s -w" -o bin/diff_gen ./diff_gen

win_version:
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=$(CROSS_CXX) CC=$(CROSS_CC) CGO_CFLAGS="-I$(ORCL_INSCL_PATH)/sdk/include" CGO_LDFLAGS="-L$(ORCL_INSCL_PATH) -L$(ORCL_INSCL_PATH)/sdk/msvc -lstdc++ -loci" go build -v -tags noPkgConfig -gcflags="-l=4" -ldflags="-s -w" -o bin/schema_gen.exe ./schema_gen
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=$(CROSS_CXX) CC=$(CROSS_CC) CGO_CFLAGS="-I$(ORCL_INSCL_PATH)/sdk/include" CGO_LDFLAGS="-L$(ORCL_INSCL_PATH) -L$(ORCL_INSCL_PATH)/sdk/msvc -lstdc++ -loci" go build -v -tags noPkgConfig -gcflags="-l=4" -ldflags="-s -w" -o bin/diff_gen.exe ./diff_gen

all: schema_gen diff_gen

wo_orcl:
	sed -i'' -e 's|_ "github.com/mattn/go-oci8"|//_ "github.com/mattn/go-oci8"|g' schema_gen/schema_gen.go
	sed -i'' -e 's|github.com/mattn/go-oci8|//github.com/mattn/go-oci8|g' go.mod
	sed -i'' -e 's|github.com/mattn/go-oci8|//github.com/mattn/go-oci8|g' go.sum

all_wo_orcl: wo_orcl all

all_wo_orcl_wver: wo_orcl win_version

clean:
	rm -rf bin
	mkdir bin