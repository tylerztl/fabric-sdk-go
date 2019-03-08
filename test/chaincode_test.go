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

func TestInstantiateCC(t *testing.T) {
	status, err := InstantiateCC("mychannel", "example_cc", "v0",
		"example_cc/go", [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})
	if status != pb.StatusCode_SUCCESS || err != nil {
		t.Error("Instantiate cc failed")
	}
}
