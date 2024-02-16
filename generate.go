//go:generate protoc  -I protobufs\ --go_out=. --go_opt=paths=source_relative --go-vtproto_out=paths=source_relative:. protobufs\meshtastic\*.proto
package meshtastic
