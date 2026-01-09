package rtdb_api

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

const Hostname = "127.0.0.1"
const Port = 6327
const Username = "sa"
const Password = "golden"

func TestRtdbGetApiVersion(t *testing.T) {
	major, minor, beta, err := RtdbGetApiVersionWarp()
	fmt.Println(major, minor, beta, err)
}

func TestRtdbConnectWarp(t *testing.T) {
	handle, err := RtdbConnectWarp(Hostname, Port)
	if err != nil {
		fmt.Println("创建连接失败", err.Error())
		return
	}
	fmt.Println(handle)
}

func TestRtdbConnectionCountWarp(t *testing.T) {
	handle, err := RtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err.Error())
		return
	}
	count, err := RtdbConnectionCountWarp(handle, 0)
	if err != nil {
		t.Error("获取Count失败", err.Error())
		return
	}
	fmt.Println("Connect Count: ", count)
}

func TestFunc(t *testing.T) {
	code := `
/**
*
* \brief 重算或补算批量计算标签点历史数据
*
* \param handle        连接句柄
* \param count         整型，输入/输出，
*                        输入时表示 ids、errors 的长度，
*                        即标签点个数；输出时返回成功开始计算的标签点个数
* \param flag          短整型，输入，不为 0 表示进行重算，删除时间范围内已经存在历史数据；
*                        为 0 表示补算，保留时间范围内已经存在历史数据，覆盖同时刻的计算值。
* \param datetime1     整型，输入，表示起始时间秒数。
* \param ms1           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示起始时间对应的纳秒；否则忽略
* \param datetime2     整型，输入，表示结束时间秒数。如果为 0，表示计算直至存档中数据的最后时间
* \param ms2           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示结束时间对应的纳秒；否则忽略
* \param ids           整型数组，输入，标签点标识
* \param errors        无符号整型数组，输出，计算历史数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、errors 的长度与 count 一致，本接口仅对带有计算扩展属性的标签点有效。
*        由 datetime1、ms1 表示的时间可以大于 datetime2、ms2 表示的时间，
*        此时前者表示结束时间，后者表示起始时间。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbe_compute_history64(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int16 flag,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    const rtdb_int32* ids,
    rtdb_error* errors
);
`
	sp := strings.Split(code, "/**")
	for _, fn := range sp {
		ff := "/**\n" + fn
		fmt.Println(CFunc(ff))
	}
}

func TestHello(t *testing.T) {
	fmt.Println("243")
}
