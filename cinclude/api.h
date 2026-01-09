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


/**
* \brief 列出 RTDB 服务器的所有连接句柄
* \param [in] handle       连接句柄
* \param [in] node_number   双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP
* \param [out] sockets    整形数组，所有连接的套接字句柄
* \param [in,out]  count   输入时表示sockets的长度，输出时表示返回的连接个数
* \return rtdb_error
* \remark 用户须保证分配给 sockets 的空间与 count 相符。如果输入的 count 小于输出的 count，则只返回部分连接
*/
rtdb_error RTDBAPI_CALLRULE rtdb_get_connections_warp(
    rtdb_int32 handle,
    rtdb_int32 node_number,
    rtdb_int32 *sockets,
    rtdb_int32 *count
) {
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_get_connections_fn)(
        rtdb_int32 handle,
        rtdb_int32 node_number,
        rtdb_int32 *sockets,
        rtdb_int32 *count
    );
    rtdb_get_connections_fn fn = (rtdb_get_connections_fn)get_function("rtdb_get_connections");
    return fn(handle, node_number, sockets, count);
}

/**
* 命名：rtdb_get_own_connection
* 功能：获取当前连接的socket句柄
* 参数：
* \param [in] handle       连接句柄
* \param [in] node_number   双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP
* \param [out] sockets    整形数组，所有连接的套接字句柄
*/
rtdb_error RTDBAPI_CALLRULE rtdb_get_own_connection_warp(
    rtdb_int32 handle,
    rtdb_int32 node_number,
    rtdb_int32* socket
) {
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_get_own_connection_fn)(
        rtdb_int32 handle,
        rtdb_int32 node_number,
        rtdb_int32* socket
    );
    rtdb_get_own_connection_fn fn = (rtdb_get_own_connection_fn)get_function("rtdb_get_own_connection");
    return fn(handle, node_number, socket);
}

/**
* \brief 获取 RTDB 服务器指定连接的信息
* \param [in] handle          连接句柄，参见 \ref rtdb_connect
* \param [in] node_number   双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP
* \param [in] socket          指定的连接
* \param [out] info          与连接相关的信息，参见 \ref RTDB_HOST_CONNECT_INFO
* \return rtdb_error
*/
rtdb_error RTDBAPI_CALLRULE rtdb_get_connection_info_warp(rtdb_int32 handle, rtdb_int32 node_number, rtdb_int32 socket, RTDB_HOST_CONNECT_INFO *info)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_get_connection_info_fn)(rtdb_int32 handle, rtdb_int32 node_number, rtdb_int32 socket, RTDB_HOST_CONNECT_INFO *info);
    rtdb_get_connection_info_fn fn = (rtdb_get_connection_info_fn)get_function("rtdb_get_connection_info");
    return fn(handle, node_number, socket, info);
}

/**
* \brief 获取 RTDB 服务器指定连接的ipv6版本
* \param [in] handle          连接句柄，参见 \ref rtdb_connect
* \param [in] node_number     双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP，双活模式仅支持ipv4
* \param [in] socket          指定的连接
* \param [out] info           与连接相关的信息，参见 \ref RTDB_HOST_CONNECT_INFO_IPV6
* \return rtdb_error
*/
rtdb_error RTDBAPI_CALLRULE rtdb_get_connection_info_ipv6_warp(rtdb_int32 handle, rtdb_int32 node_number, rtdb_int32 socket, RTDB_HOST_CONNECT_INFO_IPV6* info)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_get_connection_info_ipv6_fn)(rtdb_int32 handle, rtdb_int32 node_number, rtdb_int32 socket, RTDB_HOST_CONNECT_INFO_IPV6* info);
    rtdb_get_connection_info_ipv6_fn fn = (rtdb_get_connection_info_ipv6_fn)get_function("rtdb_get_connection_info_ipv6");
    return fn(handle, node_number, socket, info);
}

/**
* \brief 断开同 RTDB 数据平台的连接
* \param handle  连接句柄
* \return rtdb_error
* \remark 完成对 RTDB 的访问后调用本函数断开连接。连接一旦断开，则需要重新连接后才能调用其他的接口函数。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_disconnect_warp(rtdb_int32 handle)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_disconnect_fn)(rtdb_int32 handle);
    rtdb_disconnect_fn fn = (rtdb_disconnect_fn)get_function("rtdb_disconnect");
    return fn(handle);
}

/**
* \brief 以有效帐户登录
* \param handle          连接句柄
* \param user            登录帐户
* \param password        帐户口令
* \param [out] priv     账户权限， 枚举 \ref RTDB_PRIV_GROUP 的值之一
* \return rtdb_error
*/
rtdb_error RTDBAPI_CALLRULE rtdb_login_warp(rtdb_int32 handle, const char *user, const char *password, rtdb_int32 *priv)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_login_fn)(rtdb_int32 handle, const char *user, const char *password, rtdb_int32 *priv);
    rtdb_login_fn fn = (rtdb_login_fn)get_function("rtdb_login");
    return fn(handle, user, password, priv);
}

/**
* \brief 获取连接句柄所连接的服务器操作系统类型
* \param     handle          连接句柄
* \param     ostype   操作系统类型 枚举 \ref RTDB_OS_TYPE 的值之一
* \return    rtdb_error
* \remark 如句柄未链接任何服务器，返回RTDB_OS_INVALID(当前支持操作系统类型：windows、linux)。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_get_linked_ostype_warp(rtdb_int32 handle, RTDB_OS_TYPE* ostype)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_get_linked_ostype_fn)(rtdb_int32 handle, RTDB_OS_TYPE* ostype);
    rtdb_get_linked_ostype_fn fn = (rtdb_get_linked_ostype_fn)get_function("rtdb_get_linked_ostype");
    return fn(handle, ostype);
}

/**
* \brief 获取连接句柄所连接的服务器相关信息
* \param     handle          连接句柄
* \param     handle_info   服务器相关信息
* \return    rtdb_error
* \remark 如句柄未链接任何服务器，返回RTDB_OS_INVALID(当前支持操作系统类型：windows、linux)。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_get_handle_info_warp(rtdb_int32 handle, RTDB_HANDLE_INFO* info)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_get_handle_info_fn)(rtdb_int32 handle, RTDB_HANDLE_INFO* info);
    rtdb_get_handle_info_fn fn = (rtdb_get_handle_info_fn)get_function("rtdb_get_handle_info");
    return fn(handle, info);
}

/**
* \brief 修改用户帐户口令
* \param handle    连接句柄
* \param user      已有帐户
* \param password  帐户新口令
* \return rtdb_error
* \remark 只有系统管理员可以修改其它用户的密码
*/
rtdb_error RTDBAPI_CALLRULE rtdb_change_password_warp(rtdb_int32 handle, const char *user, const char *password)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_change_password_fn)(rtdb_int32 handle, const char *user, const char *password);
    rtdb_change_password_fn fn = (rtdb_change_password_fn)get_function("rtdb_change_password");
    return fn(handle, user, password);
}

/**
* \brief 用户修改自己帐户口令
* \param handle  连接句柄
* \param old_pwd 帐户原口令
* \param new_pwd 帐户新口令
* \return rtdb_error
*/
rtdb_error RTDBAPI_CALLRULE rtdb_change_my_password_warp(rtdb_int32 handle, const char *old_pwd, const char *new_pwd)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_change_my_password_fn)(rtdb_int32 handle, const char *old_pwd, const char *new_pwd);
    rtdb_change_my_password_fn fn = (rtdb_change_my_password_fn)get_function("rtdb_change_my_password");
    return fn(handle, old_pwd, new_pwd);
}

/**
* \brief 获取连接权限
* \param handle          连接句柄
* \param [out] priv  帐户权限， 枚举 \ref RTDB_PRIV_GROUP 的值之一
* \return rtdb_error
* \remark 如果还未登陆或不在服务器信任连接中，对应权限为-1，表示没有任何权限
*/
rtdb_error RTDBAPI_CALLRULE rtdb_get_priv_warp(rtdb_int32 handle, rtdb_int32 *priv)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_get_priv_fn)(rtdb_int32 handle, rtdb_int32 *priv);
    rtdb_get_priv_fn fn = (rtdb_get_priv_fn)get_function("rtdb_get_priv");
    return fn(handle, priv);
}

/**
* \brief 修改用户帐户权限
* \param handle  连接句柄
* \param user    已有帐户
* \param priv    帐户权限， 枚举 \ref RTDB_PRIV_GROUP 的值之一
* \return rtdb_error
* \remark 只有管理员有修改权限
*/
rtdb_error RTDBAPI_CALLRULE rtdb_change_priv_warp(rtdb_int32 handle, const char *user, rtdb_int32 priv)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_change_priv_fn)(rtdb_int32 handle, const char *user, rtdb_int32 priv);
    rtdb_change_priv_fn fn = (rtdb_change_priv_fn)get_function("rtdb_change_priv");
    return fn(handle, user, priv);
}

/**
* \brief 添加用户帐户
* \param handle    连接句柄
* \param user      帐户
* \param password  帐户初始口令
* \param priv      帐户权限， 枚举 \ref RTDB_PRIV_GROUP 的值之一
* \return rtdb_error
* \remark 只有管理员有添加用户权限
*/
rtdb_error RTDBAPI_CALLRULE rtdb_add_user_warp(rtdb_int32 handle, const char *user, const char *password, rtdb_int32 priv)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_add_user_fn)(rtdb_int32 handle, const char *user, const char *password, rtdb_int32 priv);
    rtdb_add_user_fn fn = (rtdb_add_user_fn)get_function("rtdb_add_user");
    return fn(handle, user, password, priv);
}

/**
* \brief 删除用户帐户
* \param handle  连接句柄
* \param user    帐户
* \return rtdb_error
* \remark 只有管理员有删除用户权限
*/
rtdb_error RTDBAPI_CALLRULE rtdb_remove_user_warp(rtdb_int32 handle, const char *user)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_remove_user_fn)(rtdb_int32 handle, const char *user);
    rtdb_remove_user_fn fn = (rtdb_remove_user_fn)get_function("rtdb_remove_user");
    return fn(handle, user);
}

/**
* \brief 启用或禁用用户
* \param     handle    连接句柄
* \param     user      字符串，输入，帐户名
* \param     lock      布尔，输入，是否禁用
* \return    rtdb_error
* \remark 只有管理员有启用禁用权限
*/
rtdb_error RTDBAPI_CALLRULE rtdb_lock_user_warp(rtdb_int32 handle, const char *user, rtdb_int8 lock)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_lock_user_fn)(rtdb_int32 handle, const char *user, rtdb_int8 lock);
    rtdb_lock_user_fn fn = (rtdb_lock_user_fn)get_function("rtdb_lock_user");
    return fn(handle, user, lock);
}

/**
* \brief 获得所有用户
* \param handle          连接句柄
* \param [in,out]  count 输入时表示 users、privs 的长度，即用户个数；输出时表示成功返回的用户信息个数
* \param [out] users     字符串指针数组，用户名称
* \param [out] privs    整型数组，用户权限，枚举 \ref RTDB_PRIV_GROUP 的值之一
* \return rtdb_error
* \remark 用户须保证分配给 users, privs 的空间与 count 相符，如果输入的 count 小于总的用户数，则只返回部分用户信息。且每个指针指向的字符串缓冲区尺寸不小于 \ref RTDB_USER_SIZE。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_get_users_warp(rtdb_int32 handle, rtdb_int32 *count, RTDB_USER_INFO *infos)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_get_users_fn)(rtdb_int32 handle, rtdb_int32 *count, RTDB_USER_INFO *infos);
    rtdb_get_users_fn fn = (rtdb_get_users_fn)get_function("rtdb_get_users");
    return fn(handle, count, infos);
}

/**
* \brief 添加连接黑名单项
* \param handle  连接句柄
* \param [in] addr    阻止连接段地址
* \param [in] mask    阻止连接段子网掩码
* \param [in] desc    阻止连接段的说明，超过 511 字符将被截断
* \return rtdb_error
* \remark addr 和 mask 进行与运算形成一个子网，
* 来自该子网范围内的连接都将被阻止，黑名单的优先级高于信任连接，
* 如果有连接既位于黑名单中，也位于信任连接中，则它将先被阻止。
* 有效的子网掩码的所有 1 位于 0 左侧，例如："255.255.254.0"。
* 当全部为 1 时，表示该子网中只有 addr 一个地址；但不能全部为 0。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_add_blacklist_warp(rtdb_int32 handle, const char *addr, const char *mask, const char *desc)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_add_blacklist_fn)(rtdb_int32 handle, const char *addr, const char *mask, const char *desc);
    rtdb_add_blacklist_fn fn = (rtdb_add_blacklist_fn)get_function("rtdb_add_blacklist");
    return fn(handle, addr, mask, desc);
}

/**
* \brief 更新连接连接黑名单项
* \param handle    连接句柄
* \param addr      原阻止连接段地址
* \param mask      原阻止连接段子网掩码
* \param addr_new  新的阻止连接段地址
* \param mask_new  新的阻止连接段子网掩码
* \param desc      新的阻止连接段的说明，超过 511 字符将被截断
*/
rtdb_error RTDBAPI_CALLRULE rtdb_update_blacklist_warp(rtdb_int32 handle, const char *addr, const char *mask, const char *addr_new, const char *mask_new, const char *desc)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_update_blacklist_fn)(rtdb_int32 handle, const char *addr, const char *mask, const char *addr_new, const char *mask_new, const char *desc);
    rtdb_update_blacklist_fn fn = (rtdb_update_blacklist_fn)get_function("rtdb_update_blacklist");
    return fn(handle, addr, mask, addr_new, mask_new, desc);
}

/**
* \brief 删除连接黑名单项
* \param handle  连接句柄
* \param addr    阻止连接段地址
* \param mask    阻止连接段子网掩码
* \remark 只有 addr 与 mask 完全相同才视为同一个阻止连接段
*/
rtdb_error RTDBAPI_CALLRULE rtdb_remove_blacklist_warp(rtdb_int32 handle, const char *addr, const char *mask)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_remove_blacklist_fn)(rtdb_int32 handle, const char *addr, const char *mask);
    rtdb_remove_blacklist_fn fn = (rtdb_remove_blacklist_fn)get_function("rtdb_remove_blacklist");
    return fn(handle, addr, mask);
}

/**
* \brief 获得连接黑名单
* \param handle          连接句柄
* \param addrs           字符串指针数组，输出，阻止连接段地址列表
* \param masks           字符串指针数组，输出，阻止连接段子网掩码列表
* \param descs           字符串指针数组，输出，阻止连接段的说明。
* \param [in,out]  count 整型，输入/输出，用户个数
* \remark 用户须保证分配给 addrs, masks, descs 的空间与 count 相符，
*      如果输入的 count 小于输出的 count，则只返回部分阻止连接段，
*      addrs, masks 中每个字符串指针所指缓冲区尺寸不得小于 32 字节，
*      descs 中每个字符串指针所指缓冲区尺寸不得小于 512 字节。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_get_blacklist_warp(rtdb_int32 handle, char* const* addrs, char* const* masks, char* const* descs, rtdb_int32 *count)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_get_blacklist_fn)(rtdb_int32 handle, char* const* addrs, char* const* masks, char* const* descs, rtdb_int32 *count);
    rtdb_get_blacklist_fn fn = (rtdb_get_blacklist_fn)get_function("rtdb_get_blacklist");
    return fn(handle, addrs, masks, descs, count);
}

/**
* \brief 添加信任连接段
* \param handle  连接句柄
* \param addr    字符串，输入，信任连接段地址
* \param mask    字符串，输入，信任连接段子网掩码。
* \param priv    整数，输入，信任连接段拥有的用户权限。
* \param desc    字符串，输入，信任连接段的说明，超过 511 字符将被截断。
* \remark addr 和 mask 进行与运算形成一个子网，
*        来自该子网范围内的连接都被视为可信任的，
*        可以不用登录 (rtdb_login)，就直接调用 API，
*        它所拥有的权限在 priv 中指定。
*        有效的子网掩码的所有 1 位于 0 左侧，
*        例如："255.255.254.0"。当全部为 1 时，
*        表示该子网中只有 addr 一个地址；
*        但不能全部为 0。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_add_authorization_warp(rtdb_int32 handle, const char *addr, const char *mask, rtdb_int32 priv, const char *desc)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_add_authorization_fn)(rtdb_int32 handle, const char *addr, const char *mask, rtdb_int32 priv, const char *desc);
    rtdb_add_authorization_fn fn = (rtdb_add_authorization_fn)get_function("rtdb_add_authorization");
    return fn(handle, addr, mask, priv, desc);
}

/**
* \brief 更新信任连接段
* \param handle    连接句柄
* \param addr      字符串，输入，原信任连接段地址。
* \param mask      字符串，输入，原信任连接段子网掩码。
* \param addr_new  字符串，输入，新的信任连接段地址。
* \param mask_new  字符串，输入，新的信任连接段子网掩码。
* \param priv      整数，输入，新的信任连接段拥有的用户权限。
* \param desc      字符串，输入，新的信任连接段的说明，超过 511 字符将被截断。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_update_authorization_warp(rtdb_int32 handle, const char *addr, const char *mask, const char *addr_new, const char *mask_new, rtdb_int32 priv, const char *desc)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_update_authorization_fn)(rtdb_int32 handle, const char *addr, const char *mask, const char *addr_new, const char *mask_new, rtdb_int32 priv, const char *desc);
    rtdb_update_authorization_fn fn = (rtdb_update_authorization_fn)get_function("rtdb_update_authorization");
    return fn(handle, addr, mask, addr_new, mask_new, priv, desc);
}

/**
* \brief 删除信任连接段
* \param handle  连接句柄
* \param addr    字符串，输入，信任连接段地址。
* \param mask    字符串，输入，信任连接段子网掩码。
* \remark 只有 addr 与 mask 完全相同才视为同一个信任连接段
*/
rtdb_error RTDBAPI_CALLRULE rtdb_remove_authorization_warp(rtdb_int32 handle, const char *addr, const char *mask)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_remove_authorization_fn)(rtdb_int32 handle, const char *addr, const char *mask);
    rtdb_remove_authorization_fn fn = (rtdb_remove_authorization_fn)get_function("rtdb_remove_authorization");
    return fn(handle, addr, mask);
}

/**
* \brief 获得所有信任连接段
* \param handle          连接句柄
* \param addrs           字符串指针数组，输出，信任连接段地址列表
* \param masks           字符串指针数组，输出，信任连接段子网掩码列表
* \param [in,out]  privs 整型数组，输出，信任连接段权限列表
* \param descs           字符串指针数组，输出，信任连接段的说明。
* \param [in,out]  count 整型，输入/输出，用户个数
* \remark 用户须保证分配给 addrs, masks, privs, descs 的空间与 count 相符，
*        如果输入的 count 小于输出的 count，则只返回部分信任连接段，
*        addrs, masks 中每个字符串指针所指缓冲区尺寸不得小于 32 字节，
*        descs 中每个字符串指针所指缓冲区尺寸不得小于 512 字节。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_get_authorizations_warp(rtdb_int32 handle, char* const* addrs, char* const* masks, rtdb_int32 *privs, char* const* descs, rtdb_int32 *count)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_get_authorizations_fn)(rtdb_int32 handle, char* const* addrs, char* const* masks, rtdb_int32 *privs, char* const* descs, rtdb_int32 *count);
    rtdb_get_authorizations_fn fn = (rtdb_get_authorizations_fn)get_function("rtdb_get_authorizations");
    return fn(handle, addrs, masks, privs, descs, count);
}

/**
* \brief 获取 RTDB 服务器当前UTC时间
*
* \param handle       连接句柄
* \param hosttime     整型，输出，Rtdb服务器的当前UTC时间，
*                     表示距离1970年1月1日08:00:00的秒数。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_host_time_warp(rtdb_int32 handle, rtdb_int32 *hosttime)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_host_time_fn)(rtdb_int32 handle, rtdb_int32 *hosttime);
    rtdb_host_time_fn fn = (rtdb_host_time_fn)get_function("rtdb_host_time");
    return fn(handle, hosttime);
}

/**
* \brief 获取 RTDB 服务器当前UTC时间
*
* \param handle       连接句柄
* \param hosttime     整型，输出，Rtdb服务器的当前UTC时间，
*                     表示距离1970年1月1日08:00:00的秒数。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_host_time64_warp(rtdb_int32 handle, rtdb_timestamp_type* hosttime)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_host_time64_fn)(rtdb_int32 handle, rtdb_timestamp_type* hosttime);
    rtdb_host_time64_fn fn = (rtdb_host_time64_fn)get_function("rtdb_host_time64");
    return fn(handle, hosttime);
}

/**
* \brief 根据时间跨度值生成时间格式字符串
*
* \param str          字符串，输出，时间格式字符串，形如:
*                     "1d" 表示时间跨度为24小时。
*                     具体含义参见 rtdb_parse_timespan 注释。
* \param timespan     整型，输入，要处理的时间跨度秒数。
* \remark 字符串缓冲区大小不应小于 32 字节。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_format_timespan_warp(char *str, rtdb_int32 timespan)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_format_timespan_fn)(char *str, rtdb_int32 timespan);
    rtdb_format_timespan_fn fn = (rtdb_format_timespan_fn)get_function("rtdb_format_timespan");
    return fn(str, timespan);
}

/**
* \brief 根据时间格式字符串解析时间跨度值
*
* \param str          字符串，输入，时间格式字符串，规则如下：
*                     [time_span]
*
*                     时间跨度部分可以出现多次，
*                     可用的时间跨度代码及含义如下：
*                     ?y            ?年, 1年 = 365日
*                     ?m            ?月, 1月 = 30 日
*                     ?d            ?日
*                     ?h            ?小时
*                     ?n            ?分钟
*                     ?s            ?秒
*                     例如："1d" 表示时间跨度为24小时。
*
* \param timespan     整型，输出，返回解析得到的时间跨度秒数。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_parse_timespan_warp(const char *str, rtdb_int32 *timespan)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_parse_timespan_fn)(const char *str, rtdb_int32 *timespan);
    rtdb_parse_timespan_fn fn = (rtdb_parse_timespan_fn)get_function("rtdb_parse_timespan");
    return fn(str, timespan);
}

/**
* \brief 根据时间格式字符串解析时间值
*
* \param str          字符串，输入，时间格式字符串，规则如下：
*                     base_time [+|- offset_time]
*
*                     其中 base_time 表示基本时间，有三种形式：
*                     1. 时间字符串，如 "2010-1-1" 及 "2010-1-1 8:00:00"；
*                     2. 时间代码，表示客户端相对时间；
*                     可用的时间代码及含义如下：
*                     td             当天零点
*                     yd             昨天零点
*                     tm             明天零点
*                     mon            本周一零点
*                     tue            本周二零点
*                     wed            本周三零点
*                     thu            本周四零点
*                     fri            本周五零点
*                     sat            本周六零点
*                     sun            本周日零点
*                     3. 星号('*')，表示客户端当前时间。
*
*                     offset_time 是可选的，可以出现多次，
*                     可用的时间偏移代码及含义如下：
*                     [+|-] ?y            偏移?年, 1年 = 365日
*                     [+|-] ?m            偏移?月, 1月 = 30 日
*                     [+|-] ?d            偏移?日
*                     [+|-] ?h            偏移?小时
*                     [+|-] ?n            偏移?分钟
*                     [+|-] ?s            偏移?秒
*                     [+|-] ?ms           偏移?毫秒
*                     例如："*-1d" 表示当前时刻减去24小时。
*
* \param datetime     整型，输出，返回解析得到的时间值。
* \param ms           短整型，输出，返回解析得到的时间毫秒值。
*  备注：ms 可以为空指针，相应的毫秒信息将不再返回。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_parse_time_warp(const char *str, rtdb_int64 *datetime, rtdb_int16 *ms)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_parse_time_fn)(const char *str, rtdb_int64 *datetime, rtdb_int16 *ms);
    rtdb_parse_time_fn fn = (rtdb_parse_time_fn)get_function("rtdb_parse_time");
    return fn(str, datetime, ms);
}

/**
* \brief 获取 Rtdb API 调用返回值的简短描述
*
* \param ecode        无符号整型，输入，Rtdb API调用后的返回值，详见rtdb_error.h头文件
* \param message      字符串，输出，返回错误码简短描述
* \param name         字符串，输出，返回错误码宏名称
* \param size         整型，输入，message 参数的字节长度
* \remark 用户须保证分配给 message， name 的空间与 size 相符,
*      name 或 message 可以为空指针，对应的信息将不再返回。
*/
void RTDBAPI_CALLRULE rtdb_format_message_warp(rtdb_error ecode, char *message, char *name, rtdb_int32 size)
{
    typedef void (RTDBAPI_CALLRULE *rtdb_format_message_fn)(rtdb_error ecode, char *message, char *name, rtdb_int32 size);
    rtdb_format_message_fn fn = (rtdb_format_message_fn)get_function("rtdb_format_message");
    return fn(ecode, message, name, size);
}

/**
* \brief 获取任务的简短描述
*
* \param job_id       整型，输入，RTDB_HOST_CONNECT_INFO::job 字段所表示的最近任务的描述
* \param desc         字符串，输出，返回任务描述
* \param name         字符串，输出，返回任务名称
* \param size         整型，输入，desc、name 参数的字节长度
* \remark 用户须保证分配给 desc、name 的空间与 size 相符，
*      name 或 message 可以为空指针，对应的信息将不再返回。
*/
void RTDBAPI_CALLRULE rtdb_job_message_warp(rtdb_int32 job_id, char *desc, char *name, rtdb_int32 size)
{
    typedef void (RTDBAPI_CALLRULE *rtdb_job_message_fn)(rtdb_int32 job_id, char *desc, char *name, rtdb_int32 size);
    rtdb_job_message_fn fn = (rtdb_job_message_fn)get_function("rtdb_job_message");
    return fn(job_id, desc, name, size);
}

/**
* \brief 设置连接超时时间
*
* \param handle   连接句柄
* \param socket   整型，输入，要设置超时时间的连接
* \param timeout  整型，输入，超时时间，单位为秒，0 表示始终保持
*/
rtdb_error RTDBAPI_CALLRULE rtdb_set_timeout_warp(rtdb_int32 handle, rtdb_int32 socket, rtdb_int32 timeout)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_set_timeout_fn)(rtdb_int32 handle, rtdb_int32 socket, rtdb_int32 timeout);
    rtdb_set_timeout_fn fn = (rtdb_set_timeout_fn)get_function("rtdb_set_timeout");
    return fn(handle, socket, timeout);
}

/**
* \brief 获得连接超时时间
*
* \param handle   连接句柄
* \param socket   整型，输入，要获取超时时间的连接
* \param timeout  整型，输出，超时时间，单位为秒，0 表示始终保持
*/
rtdb_error RTDBAPI_CALLRULE rtdb_get_timeout_warp(rtdb_int32 handle, rtdb_int32 socket, rtdb_int32 *timeout)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_get_timeout_fn)(rtdb_int32 handle, rtdb_int32 socket, rtdb_int32 *timeout);
    rtdb_get_timeout_fn fn = (rtdb_get_timeout_fn)get_function("rtdb_get_timeout");
    return fn(handle, socket, timeout);
}

/**
* \brief 断开已知连接
*
* \param handle    连接句柄
* \param socket    整型，输入，要断开的连接
*/
rtdb_error RTDBAPI_CALLRULE rtdb_kill_connection_warp(rtdb_int32 handle, rtdb_int32 socket)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_kill_connection_fn)(rtdb_int32 handle, rtdb_int32 socket);
    rtdb_kill_connection_fn fn = (rtdb_kill_connection_fn)get_function("rtdb_kill_connection");
    return fn(handle, socket);
}

/**
* \brief 获得字符串型数据库系统参数
*
* \param handle    连接句柄
* \param index     整型，输入，要取得的参数索引，参见枚举 RTDB_DB_PARAM_INDEX。
* \param str       字符串型，输出，存放取得的字符串参数值。
* \param size      整型，输入，字符串缓冲区尺寸。
* \remark 本接口只接受 [RTDB_PARAM_STR_FIRST, RTDB_PARAM_STR_LAST) 范围之内参数索引。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_get_db_info1_warp(rtdb_int32 handle, rtdb_int32 index, char *str, rtdb_int32 size)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_get_db_info1_fn)(rtdb_int32 handle, rtdb_int32 index, char *str, rtdb_int32 size);
    rtdb_get_db_info1_fn fn = (rtdb_get_db_info1_fn)get_function("rtdb_get_db_info1");
    return fn(handle, index, str, size);
}

/**
* \brief 获得整型数据库系统参数
*
* \param handle    连接句柄
* \param index     整型，输入，要取得的参数索引，参见枚举 RTDB_DB_PARAM_INDEX。
* \param value     无符号整型，输出，存放取得的整型参数值。
* \remark 本接口只接受 [RTDB_PARAM_INT_FIRST, RTDB_PARAM_INT_LAST) 范围之内参数索引。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_get_db_info2_warp(rtdb_int32 handle, rtdb_int32 index, rtdb_uint32 *value)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_get_db_info2_fn)(rtdb_int32 handle, rtdb_int32 index, rtdb_uint32 *value);
    rtdb_get_db_info2_fn fn = (rtdb_get_db_info2_fn)get_function("rtdb_get_db_info2");
    return fn(handle, index, value);
}

/**
* \brief 设置字符串型数据库系统参数
*
* \param handle    连接句柄
* \param index     整型，输入，要设置的参数索引，参见枚举 RTDB_DB_PARAM_INDEX。
*                  其中，仅以下列出的枚举值可用：
*                  RTDB_PARAM_AUTO_BACKUP_PATH,
*                  RTDB_PARAM_SERVER_SENDER_IP,
* \param str       字符串型，输入，新的参数值。
* \remark 如果修改了启动参数，将返回 RtE_DATABASE_NEED_RESTART 提示码。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_set_db_info1_warp(rtdb_int32 handle, rtdb_int32 index, const char *str)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_set_db_info1_fn)(rtdb_int32 handle, rtdb_int32 index, const char *str);
    rtdb_set_db_info1_fn fn = (rtdb_set_db_info1_fn)get_function("rtdb_set_db_info1");
    return fn(handle, index, str);
}

/**
* \brief 设置整型数据库系统参数
*
* \param handle    连接句柄
* \param index     整型，输入，要取得的参数索引，参见枚举 RTDB_DB_PARAM_INDEX。
*                  其中，仅以下列出的枚举值可用：
*                  RTDB_PARAM_SERVER_IPC_SIZE,
*                  RTDB_PARAM_EQUATION_IPC_SIZE,
*                  RTDB_PARAM_HASH_TABLE_SIZE,
*                  RTDB_PARAM_TAG_DELETE_TIMES,
*                  RTDB_PARAM_SERVER_PORT,
*                  RTDB_PARAM_SERVER_SENDER_PORT,
*                  RTDB_PARAM_SERVER_RECEIVER_PORT,
*                  RTDB_PARAM_SERVER_MODE,
*                  RTDB_PARAM_ARV_PAGES_NUMBER,
*                  RTDB_PARAM_ARVEX_PAGES_NUMBER,
*                  RTDB_PARAM_EXCEPTION_AT_SERVER,
*                  RTDB_PARAM_EX_ARCHIVE_SIZE,
*                  RTDB_PARAM_ARCHIVE_BATCH_SIZE,
*                  RTDB_PARAM_ARV_ASYNC_QUEUE_SLOWER_DOOR,
*                  RTDB_PARAM_ARV_ASYNC_QUEUE_NORMAL_DOOR,
*                  RTDB_PARAM_INDEX_ALWAYS_IN_MEMORY,
*                  RTDB_PARAM_DISK_MIN_REST_SIZE,
*                  RTDB_PARAM_DELAY_OF_AUTO_MERGE_OR_ARRANGE,
*                  RTDB_PARAM_START_OF_AUTO_MERGE_OR_ARRANGE,
*                  RTDB_PARAM_STOP_OF_AUTO_MERGE_OR_ARRANGE,
*                  RTDB_PARAM_START_OF_AUTO_BACKUP,
*                  RTDB_PARAM_STOP_OF_AUTO_BACKUP,
*                  RTDB_PARAM_MAX_LATENCY_OF_SNAPSHOT,
*                  RTDB_PARAM_PAGE_ALLOCATOR_RESERVE_SIZE,
* \param value     无符号整型，输入，新的参数值。
* \remark 如果修改了启动参数，将返回 RtE_DATABASE_NEED_RESTART 提示码。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_set_db_info2_warp(rtdb_int32 handle, rtdb_int32 index, rtdb_uint32 value)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_set_db_info2_fn)(rtdb_int32 handle, rtdb_int32 index, rtdb_uint32 value);
    rtdb_set_db_info2_fn fn = (rtdb_set_db_info2_fn)get_function("rtdb_set_db_info2");
    return fn(handle, index, value);
}

/**
* \brief 获得逻辑盘符
*
* \param handle     连接句柄
* \param drivers    字符数组，输出，
*                   返回逻辑盘符组成的字符串，每个盘符占一个字符。
* \remark drivers 的内存空间由用户负责维护，长度应不小于 32。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_get_logical_drivers_warp(rtdb_int32 handle, char *drivers)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_get_logical_drivers_fn)(rtdb_int32 handle, char *drivers);
    rtdb_get_logical_drivers_fn fn = (rtdb_get_logical_drivers_fn)get_function("rtdb_get_logical_drivers");
    return fn(handle, drivers);
}

/**
* \brief 打开目录以便遍历其中的文件和子目录。
*
* \param handle       连接句柄
* \param dir          字符串，输入，要打开的目录
*/
rtdb_error RTDBAPI_CALLRULE rtdb_open_path_warp(rtdb_int32 handle, const char *dir)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_open_path_fn)(rtdb_int32 handle, const char *dir);
    rtdb_open_path_fn fn = (rtdb_open_path_fn)get_function("rtdb_open_path");
    return fn(handle, dir);
}

/**
* \brief 读取目录中的文件或子目录
*
* \param handle      连接句柄
* \param path        字符数组，输出，返回的文件、子目录全路径
* \param is_dir      短整数，输出，返回 1 为目录，0 为文件
* \param atime       整数，输出，为文件时，返回访问时间
* \param ctime       整数，输出，为文件时，返回建立时间
* \param mtime       整数，输出，为文件时，返回修改时间
* \param size        64 位整数，输出，为文件时，返回文件大小
* \remark path 的内存空间由用户负责维护，尺寸应不小于 RTDB_MAX_PATH。
*      当返回值为 RtE_BATCH_END 时表示目录下所有子目录和文件已经遍历完毕。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_read_path_warp(rtdb_int32 handle, char *path, rtdb_int16 *is_dir, rtdb_int32 *atime, rtdb_int32 *ctime, rtdb_int32 *mtime, rtdb_int64 *size)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_read_path_fn)(rtdb_int32 handle, char *path, rtdb_int16 *is_dir, rtdb_int32 *atime, rtdb_int32 *ctime, rtdb_int32 *mtime, rtdb_int64 *size);
    rtdb_read_path_fn fn = (rtdb_read_path_fn)get_function("rtdb_read_path");
    return fn(handle, path, is_dir, atime, ctime, mtime, size);
}

/**
* \brief 读取目录中的文件或子目录
*
* \param handle      连接句柄
* \param path        字符数组，输出，返回的文件、子目录全路径
* \param is_dir      短整数，输出，返回 1 为目录，0 为文件
* \param atime       整数，输出，为文件时，返回访问时间
* \param ctime       整数，输出，为文件时，返回建立时间
* \param mtime       整数，输出，为文件时，返回修改时间
* \param size        64 位整数，输出，为文件时，返回文件大小
* \remark path 的内存空间由用户负责维护，尺寸应不小于 RTDB_MAX_PATH。
*      当返回值为 RtE_BATCH_END 时表示目录下所有子目录和文件已经遍历完毕。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_read_path64_warp(rtdb_int32 handle, char* path, rtdb_int16* is_dir, rtdb_timestamp_type* atime, rtdb_timestamp_type* ctime, rtdb_timestamp_type* mtime, rtdb_int64* size)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_read_path64_fn)(rtdb_int32 handle, char* path, rtdb_int16* is_dir, rtdb_timestamp_type* atime, rtdb_timestamp_type* ctime, rtdb_timestamp_type* mtime, rtdb_int64* size);
    rtdb_read_path64_fn fn = (rtdb_read_path64_fn)get_function("rtdb_read_path64");
    return fn(handle, path, is_dir, atime, ctime, mtime, size);
}

/**
*
* \brief 关闭当前遍历的目录
*
* \param handle      连接句柄
*/
rtdb_error RTDBAPI_CALLRULE rtdb_close_path_warp(rtdb_int32 handle)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_close_path_fn)(rtdb_int32 handle);
    rtdb_close_path_fn fn = (rtdb_close_path_fn)get_function("rtdb_close_path");
    return fn(handle);
}

/**
*
* \brief 建立目录
*
* \param handle       连接句柄
* \param dir          字符串，输入，新建目录的全路径
*/
rtdb_error RTDBAPI_CALLRULE rtdb_mkdir_warp(rtdb_int32 handle, const char *dir)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_mkdir_fn)(rtdb_int32 handle, const char *dir);
    rtdb_mkdir_fn fn = (rtdb_mkdir_fn)get_function("rtdb_mkdir");
    return fn(handle, dir);
}

/**
*
* \brief 获得指定服务器端文件的大小
*
* \param handle     连接句柄
* \param file       字符串，输入，文件名
* \param size       64 位整数，输出，文件大小
*/
rtdb_error RTDBAPI_CALLRULE rtdb_get_file_size_warp(rtdb_int32 handle, const char *file, rtdb_int64 *size)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_get_file_size_fn)(rtdb_int32 handle, const char *file, rtdb_int64 *size);
    rtdb_get_file_size_fn fn = (rtdb_get_file_size_fn)get_function("rtdb_get_file_size");
    return fn(handle, file, size);
}

/**
*
* \brief 读取服务器端指定文件的内容
*
* \param handle       连接句柄
* \param file         字符串，输入，要读取内容的文件名
* \param content      字符数组，输出，文件内容
* \param pos          64 位整型，输入，读取文件的起始位置
* \param size         整型，输入/输出，
*                     输入时表示要读取文件内容的字节大小；
*                     输出时表示实际读取的字节数
* \remark 用户须保证分配给 content 的空间与 size 相符。
*/
rtdb_error RTDBAPI_CALLRULE rtdb_read_file_warp(rtdb_int32 handle, const char *file, char *content, rtdb_int64 pos, rtdb_int64 *size)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_read_file_fn)(rtdb_int32 handle, const char *file, char *content, rtdb_int64 pos, rtdb_int64 *size);
    rtdb_read_file_fn fn = (rtdb_read_file_fn)get_function("rtdb_read_file");
    return fn(handle, file, content, pos, size);
}

/**
*
* \brief 取得数据库允许的blob与str类型测点的最大长度
*
* \param handle       连接句柄
* \param len          整形，输出参数，代表数据库允许的blob、str类型测点的最大长度
*/
rtdb_error RTDBAPI_CALLRULE rtdb_get_max_blob_len_warp(rtdb_int32 handle, rtdb_int32 *len)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_get_max_blob_len_fn)(rtdb_int32 handle, rtdb_int32 *len);
    rtdb_get_max_blob_len_fn fn = (rtdb_get_max_blob_len_fn)get_function("rtdb_get_max_blob_len");
    return fn(handle, len);
}

/**
*
* \brief 取得质量码对应的定义
*
* \param handle       连接句柄
* \param count        质量码个数，输入参数，
* \param qualities    质量码，输入参数
* \param definitions  输出参数，0~255为RTDB质量码（参加rtdb.h文件），256~511为OPC质量码，大于511为用户自定义质量码
* \param lens         输出参数，每个定义对应的长度
* \remark OPC质量码把8位分3部分定义：XX XXXX XX，对应着：品质位域 分状态位域 限定位域
* 品质位域：00（Bad）01（Uncertain）10（N/A）11（Good）
* 分状态位域：不同品质位域对应各自的分状态位域
* 限定位域：00（Not limited）01（Low limited）10（high limited）11（Constant）
* 三个域之间用逗号隔开，输出到definitions参数中，前面有有RTDB，OPC或者USER标识，说明标签点类别
*/
rtdb_error RTDBAPI_CALLRULE rtdb_format_quality_warp(rtdb_int32 handle, rtdb_int32 *count, rtdb_int16 *qualities, rtdb_byte **definitions, rtdb_int32 *lens)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_format_quality_fn)(rtdb_int32 handle, rtdb_int32 *count, rtdb_int16 *qualities, rtdb_byte **definitions, rtdb_int32 *lens);
    rtdb_format_quality_fn fn = (rtdb_format_quality_fn)get_function("rtdb_format_quality");
    return fn(handle, count, qualities, definitions, lens);
}

/**
*
* \brief 判断连接是否可用
* \param handle   连接句柄
*/
rtdb_error RTDBAPI_CALLRULE rtdb_judge_connect_status_warp(rtdb_int32 handle, rtdb_int8* change_connection GAPI_DEFAULT_VALUE(0), char* current_ip_addr GAPI_DEFAULT_VALUE(0), rtdb_int32 size GAPI_DEFAULT_VALUE(0))
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_judge_connect_status_fn)(rtdb_int32 handle, rtdb_int8* change_connection GAPI_DEFAULT_VALUE(0), char* current_ip_addr GAPI_DEFAULT_VALUE(0), rtdb_int32 size GAPI_DEFAULT_VALUE(0));
    rtdb_judge_connect_status_fn fn = (rtdb_judge_connect_status_fn)get_function("rtdb_judge_connect_status");
    return fn(handle, GAPI_DEFAULT_VALUE(0), GAPI_DEFAULT_VALUE(0), GAPI_DEFAULT_VALUE(0));
}

/**
* 命名：rtdb_format_ipaddr
* 功能：将整形IP转换为字符串形式的IP
* 参数：
*      [ip]        无符号整型，输入，整形的IP地址
*      [ip_addr]      字符串，输出，字符串IP地址缓冲区
*      [size]         整型，输入，ip_addr 参数的字节长度
* 备注：用户须保证分配给 ip_addr 的空间与 size 相符
*/
void RTDBAPI_CALLRULE rtdb_format_ipaddr_warp(rtdb_uint32 ip, char* ip_addr, rtdb_int32 size)
{
    typedef void (RTDBAPI_CALLRULE *rtdb_format_ipaddr_fn)(rtdb_uint32 ip, char* ip_addr, rtdb_int32 size);
    rtdb_format_ipaddr_fn fn = (rtdb_format_ipaddr_fn)get_function("rtdb_format_ipaddr");
    return fn(ip, ip_addr, size);
}

#ifdef __cplusplus
}
#endif

#endif // _C_API_H_
