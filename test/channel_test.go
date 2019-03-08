package test

import (
	pb "fabric-sdk-go/protos"
	"testing"
)

func TestCreateChannel(t *testing.T) {
	status, err := CreateChannel("mychannel")
	if status != pb.StatusCode_SUCCESS || err != nil {
		t.Error("测试失败")
	}
}

func TestJoinChannel(t *testing.T) {
	status, err := JoinChannel("mychannel")
	if status != pb.StatusCode_SUCCESS || err != nil {
		t.Error("测试失败")
	}
}
