# bson2mgoindex
bson2mgoindex is a cmd tool to generate mongodb shell for index from golang struct 

## usage
```
go run bson2mgoindex.go -f models/bson_test.go
```

```
db.getCollection("tb_agent").createIndex({"uuid":1},{background: true});
db.getCollection("tb_agent").createIndex({"relate_id":-1},{background: true});
db.getCollection("tb_user").createIndex({"name":1},{background: true});
```