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
		"example_cc/go", "AND ('Org1MSP.member','Org2MSP.member')", [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})
	if status != pb.StatusCode_SUCCESS || err != nil {
		t.Error("Instantiate cc failed")
	}
}

func TestUpgradeCC(t *testing.T) {
	status, err := UpgradeCC("mychannel", "example_cc", "v1",
		"example_cc/go", "OutOf (1, 'Org1MSP.member')", [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})
	if status != pb.StatusCode_SUCCESS || err != nil {
		t.Error("Upgrade cc failed")
	}
}

func TestInvokeCC(t *testing.T) {
	status, err := InvokeCC("mychannel", "example_cc", "move", [][]byte{[]byte("a"), []byte("b"), []byte("10")})
	if status != pb.StatusCode_SUCCESS || err != nil {
		t.Error("Invoke cc failed")
	}
}

func TestQueryCC(t *testing.T) {
	status, err := QueryCC("mychannel", "example_cc", "query", [][]byte{[]byte("a")})
	if status != pb.StatusCode_SUCCESS || err != nil {
		t.Error("Query cc failed")
	}
}
