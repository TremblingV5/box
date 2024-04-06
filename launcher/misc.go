package launcher

import "os"

const (
	ProtoBufConflictEnvKey = "GOLANG_PROTOBUF_REGISTRATION_CONFLICT"
)

func FixProtoBufConflict() {
	_ = os.Setenv(ProtoBufConflictEnvKey, "warn")
}
