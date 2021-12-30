# sql-parser-go
sql-parser-go is a simple sql parser, for now, it only implements the lexer part.

# Build
```
go build -o parser main.go
```

# How to use
```
./parser nfa --sql="select col1, col2 from t01 where id <= 100 and col1 = 'abc'"
```
or
```
./parser dfa --sql="select col1, col2 from t01 where id <= 100 and col1 = 'abc'"
```

# Document
document of [lexer](docs/lexer_cn.md)