package dict

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	cacheDict *CachedDict[int64, *TestInfo]
	load      bool
	result    *TestInfo
	err       error
	todo      = context.TODO()
)

func TestNewCachedDict(t *testing.T) {
	cacheDict = NewCachedDict[int64, *TestInfo](func(ctx context.Context, key int64) (info *TestInfo, err error) {
		load = true
		return &TestInfo{A: key}, nil
	}, 10)
}

type TestInfo struct {
	A int64
	B bool
}

/*
TestCachedDict_Get 测试获取
参数:
*	t	*testing.T	参数1
返回值:

测试场景
 1. key==1 第一次获取时，成功，且真实发起加载(通过load==false来判断)
 2. key==1 第二次获取时，成功，且未真实记载
 3. key==2 第一次加载，成功，真实发起加载
 4. 清理key==2, 再次加载key==2,
*/
func TestCachedDict_Get(t *testing.T) {
	TestNewCachedDict(t)

	firstKey := int64(1)
	secondKey := firstKey + 1
	result, err = cacheDict.Get(todo, firstKey)
	require.NoError(t, err, `获取`)
	require.EqualValues(t, firstKey, result.A)
	require.True(t, load, `第一次会加载`)

	load = false

	// 第二次
	result, err = cacheDict.Get(todo, firstKey)
	require.NoError(t, err, `获取`)
	require.EqualValues(t, firstKey, result.A)
	require.False(t, load, `第二次不会加载`)

	load = false

	result, err = cacheDict.Get(todo, secondKey)
	require.NoError(t, err, `获取`)
	require.EqualValues(t, secondKey, result.A)
	require.True(t, load, `第一次会加载`)
	load = false

	cacheDict.Clean(1)

	result, err = cacheDict.Get(todo, firstKey)
	require.NoError(t, err, `获取`)
	require.EqualValues(t, firstKey, result.A)
	require.True(t, load, `清理后会加载`)
}
