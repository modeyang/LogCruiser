# LogCruiser
format log, extract log, metrics from log

# denpendcies
- kafka
- yaml
- go template

# roadmap
## v1
1. input: read from stdin, kafka, file
2. format: convert to map[String]interface{}
3. metric: extract metric from log by using metric template (like `access.qps/typeName={{.TypeName}}`) with counter value.
4. output: send metric to storage
