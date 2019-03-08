package test

import (
	pb "fabric-sdk-go/protos"
	"testing"
)

func TestCreateChannel(t *testing.T) {
	status, err := CreateChannel("mychannel")
	if status != pb.StatusCode_SUCCESS || err != nil {
		t.Error("Create channel failed")
	}
}

func TestJoinChannel(t *testing.T) {
	status, err := JoinChannel("mychannel")
	if status != pb.StatusCode_SUCCESS || err != nil {
		t.Error("Join channel failed")
	}
}
