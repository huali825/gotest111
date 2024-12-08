package connpooltmh

import (
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

// 不停机数据迁移的 双写步骤 使用connPool来操作
func TestConnPool(t *testing.T) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: DoubleWritePool{},
	}))
	require.NoError(t, err)
	t.Log(db)
}
