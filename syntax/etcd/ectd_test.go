package etcd

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestName(t *testing.T) {
	fmt.Println("this is etcd test")
}

type EtcdTestSuite struct {
	suite.Suite //测试套件
}
