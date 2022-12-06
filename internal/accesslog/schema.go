package accesslog

//go:generate mkdir -p ./avro
//go:generate go run github.com/actgardner/gogen-avro/v10/cmd/gogen-avro --containers ./avro ../../terraform/access_log_schema/v1.avsc
