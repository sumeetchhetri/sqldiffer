deps:
	go get google.golang.org/protobuf
	go get github.com/jessevdk/go-flags
	go get github.com/denisenkom/go-mssqldb
	go get github.com/go-sql-driver/mysql
	go get github.com/lib/pq
	go get github.com/mattn/go-oci8
	cd protos && protoc --go_out=. *.proto && cd ..

schema_gen: clean deps
	env GO111MODULE=off go build -gcflags="-l=4" -ldflags="-s -w" -o bin/schema_gen ./schema_gen

diff_gen: clean deps
	env GO111MODULE=off go build -gcflags="-l=4" -ldflags="-s -w" -o bin/diff_gen ./diff_gen

win_version:
	env GO111MODULE=off GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc CGO_CFLAGS="-I/Users/sumeetc/Projects/home/instantclient_12_2_winx/sdk/include" CGO_LDFLAGS="-L/Users/sumeetc/Projects/home/instantclient_12_2_winx/ -L/Users/sumeetc/Projects/home/instantclient_12_2_winx/sdk/msvc  -lstdc++ -loci" go build -v -tags noPkgConfig -gcflags="-l=4" -ldflags="-s -w" -o bin/schema_gen.exe ./schema_gen
	env GO111MODULE=off GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc CGO_CFLAGS="-I/Users/sumeetc/Projects/home/instantclient_12_2_winx/sdk/include" CGO_LDFLAGS="-L/Users/sumeetc/Projects/home/instantclient_12_2_winx/ -L/Users/sumeetc/Projects/home/instantclient_12_2_winx/sdk/msvc  -lstdc++ -loci" go build -v -tags noPkgConfig -gcflags="-l=4" -ldflags="-s -w" -o bin/diff_gen.exe ./diff_gen

all: schema_gen diff_gen

clean:
	rm -rf bin
	mkdir bin