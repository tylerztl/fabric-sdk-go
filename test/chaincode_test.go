package test

import (
	pb "fabric-sdk-go/protos"
	"testing"
)

func TestInstallCC(t *testing.T) {
	status, err := InstallCC("example_cc", "v0", "example_cc/go")
	if status != pb.StatusCode_SUCCESS || err != nil {
		t.Error("Install cc failed")
	}
}
