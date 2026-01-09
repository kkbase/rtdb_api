#ifndef _C_DYLIB_H_
#define _C_DYLIB_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#include <windows.h>
#define LOAD_LIBRARY(name) LoadLibraryA(name)
#define GET_FUNCTION GetProcAddress
#define FREE_LIBRARY FreeLibrary
typedef HMODULE LIBRARY_HANDLE;
#else
#include <dlfcn.h>
#define LOAD_LIBRARY(name) dlopen(name, RTLD_LAZY)
#define GET_FUNCTION dlsym
#define FREE_LIBRARY dlclose
typedef void* LIBRARY_HANDLE;
#endif

#include "rtdb.h"
#include "rtdbapi.h"
#include "rtdb_error.h"

LIBRARY_HANDLE LIB;

// 加载动态库
void load_library(char *path) {
    LIB = LOAD_LIBRARY(path);
}


// 释放动态库
void free_library() {
    FREE_LIBRARY(LIB);
}


// 从动态库获取函数
void* get_function(char *name) {
    return GET_FUNCTION(LIB, name);
}


/**
* \brief   取得 rtdbapi 库的版本号
* \param [out]  major   主版本号
* \param [out]  minor   次版本号
* \param [out]  beta    发布版本号
* \return rtdb_error
* \remark 如果返回的版本号与 rtdb.h 中定义的不匹配(RTDB_API_XXX_VERSION)，则应用程序使用了错误的库。
*      应输出一条错误信息并退出，否则可能在调用某些 api 时会导致崩溃
*/
rtdb_error RTDBAPI_CALLRULE rtdb_get_api_version_warp(
    rtdb_int32 *major,
    rtdb_int32 *minor,
    rtdb_int32 *beta
) {
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_get_api_version_fn)(
        rtdb_int32 *major,
        rtdb_int32 *minor,
        rtdb_int32 *beta
    );
    rtdb_get_api_version_fn fn = (rtdb_get_api_version_fn)get_function("rtdb_get_api_version");
    return fn(major, minor, beta);
}


/**
* \brief 配置 api 行为参数，参见枚举 \ref RTDB_API_OPTION
* \param [in] type  选项类别
* \param [in] value 选项值
* \return rtdb_error
* \remark 选项设置后在下一次调用 api 时才生效
*/
rtdb_error RTDBAPI_CALLRULE rtdb_set_option_warp(
    rtdb_int32 type,
    rtdb_int32 value
) {
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_set_option_fn)(
        rtdb_int32 type,
        rtdb_int32 value
    );
    rtdb_set_option_fn fn = (rtdb_set_option_fn)get_function("rtdb_set_option");
    return fn(type, value);
}


/**
* \brief 创建数据流
* \param [in] in 端口
* \param [out] remotehost 对端地址
* \param [out] handle 数据流句柄
* \return rtdb_error
* \remark 创建数据流
*/
rtdb_error RTDBAPI_CALLRULE rtdb_create_datagram_handle_warp(
    rtdb_int32 port,
    const char* remotehost,
    rtdb_datagram_handle* handle
) {
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_create_datagram_handle_fn)(
        rtdb_int32 port,
        const char* remotehost,
        rtdb_datagram_handle* handle
    );
    rtdb_create_datagram_handle_fn fn = (rtdb_create_datagram_handle_fn)get_function("rtdb_create_datagram_handle");
    return fn(port, remotehost, handle);
}


/**
* \brief 删除数据流
* \param [in] handle 数据流句柄
* \return rtdb_error
* \remark 删除数据流
*/
rtdb_error RTDBAPI_CALLRULE rtdb_remove_datagram_handle_warp(
    rtdb_datagram_handle handle
) {
    typedef  rtdb_error (RTDBAPI_CALLRULE *rtdb_remove_datagram_handle_fn)(
        rtdb_datagram_handle handle
    );
    rtdb_remove_datagram_handle_fn fn = (rtdb_remove_datagram_handle_fn)get_function("rtdb_remove_datagram_handle");
    return fn(handle);
}


/**
* \brief 接收数据流
* \param [in] message 消息
* \param [in] message_len 消息长度
* \param [in] handle 数据流句柄
* \param [in] remote_addr 对端地址
* \param [in] timeout 超时时间
* \return rtdb_error
* \remark 接收数据流
*/
rtdb_error RTDBAPI_CALLRULE rtdb_recv_datagram_warp(
    char* message,
    rtdb_int32* message_len,
    rtdb_datagram_handle handle,
    char* remote_addr,
    rtdb_int32 timeout
) {
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_recv_datagram_fn)(
        char* message,
        rtdb_int32* message_len,
        rtdb_datagram_handle handle,
        char* remote_addr,
        rtdb_int32 timeout
    );
    rtdb_recv_datagram_fn fn = (rtdb_recv_datagram_fn)get_function("rtdb_recv_datagram");
    return fn(message, message_len, handle, remote_addr, timeout);
}


// TODO
// RTDBAPI
// rtdb_error
// RTDBAPI_CALLRULE
// rtdb_subscribe_connect_ex(
// rtdb_int32 handle,
// rtdb_uint32 options,
// void* param,
// rtdb_connect_event_ex callback
// );
//
// RTDBAPI
// rtdb_error
// RTDBAPI_CALLRULE
// rtdb_cancel_subscribe_connect(
// rtdb_int32 handle
// );


/**
* \brief 建立同 RTDB 数据库的网络连接
* \param [in] hostname     RTDB 数据平台服务器的网络地址或机器名
* \param [in] port         连接断开，缺省值 6327
* \param [out]  handle  连接句柄
* \return rtdb_error
* \remark 在调用所有的接口函数之前，必须先调用本函数建立同Rtdb服务器的连接
 */
rtdb_error RTDBAPI_CALLRULE rtdb_connect_warp(
    const char *hostname,
    rtdb_int32 port,
    rtdb_int32 *handle
) {
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_connect_fn)(
        const char *hostname,
        rtdb_int32 port,
        rtdb_int32 *handle
    );
    rtdb_connect_fn fn = (rtdb_connect_fn)get_function("rtdb_connect");
    return fn(hostname, port, handle);
}


/**
* \brief 获取 RTDB 服务器当前连接个数
* \param [in] handle   连接句柄 参见 \ref rtdb_connect
* \param [in] node_number   双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP
* \param [out]  count 返回当前主机的连接个数
* \return rtdb_error
*/
rtdb_error RTDBAPI_CALLRULE rtdb_connection_count_warp(
    rtdb_int32 handle,
    rtdb_int32 node_number,
    rtdb_int32 *count
) {
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_connection_count_fn)(
        rtdb_int32 handle,
        rtdb_int32 node_number,
        rtdb_int32 *count
    );
    rtdb_connection_count_fn fn = (rtdb_connection_count_fn)get_function("rtdb_connection_count");
    return fn(handle, node_number, count);
}


#ifdef __cplusplus
}
#endif

#endif // _C_API_H_
