# sqldiffer - The database schema generation and compare-diff tools

Currently supports the following databases
1. Oracle
2. PostgreSQL
3. MySQL
4. SQL Server

**schema_gen** - Tool for generating schema file representing a database (tables, sequences, stored procedures, views, indexes, constraints, triggers)<br/>
It generates protobuf binary files representing the structure of the database entity, for eg, name, columns + their types etc for a table
Please ensure that the path for oracle install client libraries is in the %PATH% (Windows) & LD_LIBRARY_PATH (for macos, linux)
```
set PATH=%PATH%;C:\path\to\instantclient_21_3
export LD_LIBRARY_PATH=/path/to/instantclient_19_8
```

```
Apples-MacBook-Pro:sqldiffer sumeetc$ schema_gen --help
Usage:
  schema_gen [OPTIONS]

Application Options:
  -t, --type=[postgres|oracle|mysql|sqlserver] The database type
  -n, --name=                                  The database name
  -i, --host=                                  The database host
  -p, --port=                                  The database port
  -u, --user=                                  The database user
  -w, --password=                              The database password
  -s, --sch-nam=                               The database schema name
  -f, --fil-nam=                               The generated schema/diff file name
      --gen-proc                               Generate stored procedure files

Help Options:
  -h, --help                                   Show this help message
```

**diff_gen** - Tool for generating diff sql when comparing 2 databases<br/>
It compares 2 previously generated protobuf schema files and provides diff sql statements which can be safely applied to the target database to get it closer to the source database

```
Apples-MacBook-Pro:sqldiffer sumeetc$ diff_gen --help
Usage:
  diff_gen [OPTIONS]

Application Options:
  -f, --src-scf=                                  The source schema file name
  -d, --tgt-scf=                                  The target schema file name
  -n, --one-file                                  Whether a single diff file is needed or multiple?
  -t, --tgt-dbt=[postgres|oracle|mysql|sqlserver] The target database type
  -m, --tgt-dbn=                                  The target database name
  -s, --tgt-dbs=                                  The target database schema name
  -r, --rdiff                                     Whether a reverse diff is needed or not?
  -o, --out-file=                                 Diff File Name
  -p, --options=                                  Extra Diff options

Help Options:
  -h, --help                                      Show this help message
```


Usage Examples
==============

**Postgresql**
```
schema_gen -i host1 -n db1 -t postgres -u user -p 5432 -w 'pwd' -f schema_fp.json
schema_gen -i host2 -n db2 -t postgres -u user -p 5432 -w 'pwd' -f schema_tp.json
diff_gen -f /path/to/file/schema_fp.json -d /path/to/file/schema_tp.json -m db2 -t postgres
```

**Oracle**
```
schema_gen -i host1 -n orcl -t oracle -u db1 -p 1521 -w 'pwd' -f schema_fo.json
schema_gen -i host2 -n orcl -t oracle -u db2 -p 1521 -w 'pwd' -f schema_to.json
diff_gen -f /path/to/file/schema_fo.json -d /path/to/file/schema_to.json -m db2 -t oracle
```

**Mysql**
```
schema_gen -i host1 -n db1 -t mysql -u user -p 3306 -w 'pwd' -f schema_fm.json
schema_gen -i host2 -n db2 -t mysql -u user -p 3306 -w 'pwd' -f schema_tm.json
diff_gen -f /path/to/file/schema_fm.json -d /path/to/file/schema_tm.json -m db2 -t mysql
```

**SQL Server**
```
schema_gen -i host1 -n db1 -t sqlserver -u user -p 1433 -w 'pwd' -f schema_fs.json
schema_gen -i host2 -n db2 -t sqlserver -u user -p 1433 -w 'pwd' -f schema_ts.json
diff_gen -f /path/to/file/schema_fs.json -d /path/to/file/schema_ts.json -m db2 -t sqlserver
```


Development/Build Commands
==========================

**Generating protobuf files**
```
cd protos && protoc --go_out=. *.proto && cd ..
```

**Creating executable**
```
go install ./...
OR
make all
```

**Cross compiling and building for windows x64**
```
env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc CGO_CFLAGS="-I/path/to/windows_oracle_install_client/instantclient_12_2_winx/sdk/include" CGO_LDFLAGS="-L/path/to/windows_oracle_install_client/instantclient_12_2_winx/ -L/path/to/windows_oracle_install_client/instantclient_12_2_winx/sdk/msvc -lstdc++ -loci" go install -v -tags noPkgConfig ./...
OR
make win_version
```
