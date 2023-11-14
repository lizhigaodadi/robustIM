package trace

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc/codes"
)

const (
	GRPCStatusCodeKey             = attribute.Key("rpc.grpc.status_code")
	RPCNameKey                    = attribute.Key("name")
	RPCMessageTypeKey             = attribute.Key("message.type")
	RPCMessageIDKey               = attribute.Key("message.id")
	RPCMessageCompressedSizeKey   = attribute.Key("message.compressed_size")
	RPCMessageUnCompressedSizeKey = attribute.Key("message.uncompressed_size")
	ServerEnvironment             = attribute.Key("environment")
)

var (
	RPCSystemGRPC          = semconv.RPCSystemKey.String("grpc")
	RPCNameMessage         = RPCNameKey.String("message")
	RPCMessageTypeSent     = RPCMessageTypeKey.String("SENT")
	RPCMessageTypeReceived = RPCMessageTypeKey.String("RECEIVED")
)

func StatusCodeAttr(code codes.Code) attribute.KeyValue {
	return GRPCStatusCodeKey.Int64(int64(code))
}
