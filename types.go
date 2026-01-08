package rtdb_api

// #cgo CFLAGS: -DPNG_DEBUG=1 -I./cinclude
// #cgo CXXFLAGS: -std=c++11
// #include "api.h"
import "C"

type RtdbApiOption uint32

const (
	// RtdbApiAutoReconn api 在连接中断后是否自动重连, 0 不重连；1 重连。默认为 0 不重连
	RtdbApiAutoReconn = RtdbApiOption(C.RTDB_API_AUTO_RECONN)

	// RtdbApiConnTimeout api 连接超时值设置（单位：毫秒）,0 阻塞模式，无限等待，默认为1000
	RtdbApiConnTimeout = RtdbApiOption(C.RTDB_API_CONN_TIMEOUT)

	// RtdbApiSendTimeout api 发送超时值设置（单位：毫秒）,0 阻塞模式，无限等待，默认为1000
	RtdbApiSendTimeout = RtdbApiOption(C.RTDB_API_SEND_TIMEOUT)

	// RtdbApiRecvTimeout api 接收超时值设置（单位：毫秒）,0 阻塞模式，无限等待，默认为60000
	RtdbApiRecvTimeout = RtdbApiOption(C.RTDB_API_RECV_TIMEOUT)

	// RtdbApiUserTimeout api TCP_USER_TIMEOUT超时值设置（单位：毫秒），默认为10000，Linux内核2.6.37以上有效
	RtdbApiUserTimeout = RtdbApiOption(C.RTDB_API_USER_TIMEOUT)

	// RtdbApiDefaultPrecision api 默认的时间戳精度，当使用旧版相关的api，以及新版api中未设置时间戳精度时，则使用此默认时间戳精度。 默认为毫秒精度
	RtdbApiDefaultPrecision = RtdbApiOption(C.RTDB_API_DEFAULT_PRECISION)

	// RtdbApiServerPrecision api 连接3.0数据库时，设置3.0数据库的时间戳精度，0表示毫秒精度，非0表示纳秒精度，默认为毫秒精度
	RtdbApiServerPrecision = RtdbApiOption(C.RTDB_API_SERVER_PRECISION)
)
