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


/**
* \brief 创建订阅连接
* \param [in] handle 连接句柄
* \param [in] options 选项
* \param [in] param 参数
* \param [in] callback 回调函数
* \return rtdb_error
* \remark 创建订阅连接
*/
rtdb_error RTDBAPI_CALLRULE rtdb_subscribe_connect_ex_warp(rtdb_int32 handle, rtdb_uint32 options, void* param, rtdb_connect_event_ex callback)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_subscribe_connect_ex_fn)(rtdb_int32 handle, rtdb_uint32 options, void* param, rtdb_connect_event_ex callback);
    rtdb_subscribe_connect_ex_fn fn = (rtdb_subscribe_connect_ex_fn)get_function("rtdb_subscribe_connect_ex");
    return fn(handle, options, param, callback);
}

/**
* \brief 关闭订阅链接
* \param [in] handle 连接句柄
* \return rtdb_error
* \remark 关闭订阅链接
*/
rtdb_error RTDBAPI_CALLRULE rtdb_cancel_subscribe_connect_warp(rtdb_int32 handle)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_cancel_subscribe_connect_fn)(rtdb_int32 handle);
    rtdb_cancel_subscribe_connect_fn fn = (rtdb_cancel_subscribe_connect_fn)get_function("rtdb_cancel_subscribe_connect");
    return fn(handle);
}


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
    return fn(handle, change_connection, current_ip_addr, size);
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

/**
* 命名：rtdbb_get_equation_by_file_name
* 功能：根据文件名获取方程式
* 参数：
*      [handle]   连接句柄
*      [file_name] 输入，字符串，方程式路径
*      [equation]  输出，返回的方程式长度最长为RTDB_MAX_EQUATION_SIZE-1
*备注：用户调用时为equation分配的空间不得小于RTDB_MAX_EQUATION_SIZE
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_equation_by_file_name_warp(rtdb_int32 handle, const char* file_name, char equation[RTDB_MAX_EQUATION_SIZE])
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_equation_by_file_name_fn)(rtdb_int32 handle, const char* file_name, char equation[RTDB_MAX_EQUATION_SIZE]);
    rtdbb_get_equation_by_file_name_fn fn = (rtdbb_get_equation_by_file_name_fn)get_function("rtdbb_get_equation_by_file_name");
    return fn(handle, file_name, equation);
}

/**
* 命名：rtdbb_get_equation_by_id
* 功能：根ID径获取方程式
* 参数：
*      [handle]   连接句柄
*      [id]				输入，整型，方程式ID
*      [equation] 输出，返回的方程式长度最长为RTDB_MAX_EQUATION_SIZE-1
*
*备注：用户调用时为equation分配的空间不得小于RTDB_MAX_EQUATION_SIZE
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_equation_by_id_warp(rtdb_int32 handle, rtdb_int32 id, char equation[RTDB_MAX_EQUATION_SIZE])
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_equation_by_id_fn)(rtdb_int32 handle, rtdb_int32 id, char equation[RTDB_MAX_EQUATION_SIZE]);
    rtdbb_get_equation_by_id_fn fn = (rtdbb_get_equation_by_id_fn)get_function("rtdbb_get_equation_by_id");
    return fn(handle, id, equation);
}

/**
*
* \brief 添加新表
*
* \param handle   连接句柄
* \param field    RTDB_TABLE 结构，输入/输出，表信息。
*                 在输入时，type、name、desc 字段有效；
*                 输出时，id 字段由系统自动分配并返回给用户。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_append_table_warp(rtdb_int32 handle, RTDB_TABLE *field)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_append_table_fn)(rtdb_int32 handle, RTDB_TABLE *field);
    rtdbb_append_table_fn fn = (rtdbb_append_table_fn)get_function("rtdbb_append_table");
    return fn(handle, field);
}

/**
*
* \brief 取得标签点表总数
*
* \param handle   连接句柄
* \param count    整型，输出，标签点表总数
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_tables_count_warp(rtdb_int32 handle, rtdb_int32 *count)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_tables_count_fn)(rtdb_int32 handle, rtdb_int32 *count);
    rtdbb_tables_count_fn fn = (rtdbb_tables_count_fn)get_function("rtdbb_tables_count");
    return fn(handle, count);
}

/**
*
* \brief 取得所有标签点表的ID
*
* \param handle   连接句柄
* \param ids      整型数组，输出，标签点表的id
* \param count    整型，输入/输出，
*                 输入表示 ids 的长度，输出表示标签点表个数
* \remark 用户须保证分配给 ids 的空间与 count 相符
*      如果输入的 count 小于输出的 count，则只返回部分表id
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_tables_warp(rtdb_int32 handle, rtdb_int32 *ids, rtdb_int32 *count)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_tables_fn)(rtdb_int32 handle, rtdb_int32 *ids, rtdb_int32 *count);
    rtdbb_get_tables_fn fn = (rtdbb_get_tables_fn)get_function("rtdbb_get_tables");
    return fn(handle, ids, count);
}

/**
*
* \brief 根据表 id 获取表中包含的标签点数量
*
* \param handle   连接句柄
* \param id       整型，输入，表ID
* \param size     整型，输出，表中标签点数量
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_table_size_by_id_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 *size)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_table_size_by_id_fn)(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 *size);
    rtdbb_get_table_size_by_id_fn fn = (rtdbb_get_table_size_by_id_fn)get_function("rtdbb_get_table_size_by_id");
    return fn(handle, id, size);
}

/**
*
* \brief 根据表名称获取表中包含的标签点数量
*
* \param handle   连接句柄
* \param name     字符串，输入，表名称
* \param size     整型，输出，表中标签点数量
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_table_size_by_name_warp(rtdb_int32 handle, const char *name, rtdb_int32 *size)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_table_size_by_name_fn)(rtdb_int32 handle, const char *name, rtdb_int32 *size);
    rtdbb_get_table_size_by_name_fn fn = (rtdbb_get_table_size_by_name_fn)get_function("rtdbb_get_table_size_by_name");
    return fn(handle, name, size);
}

/**
*
* \brief 根据表 id 获取表中实际包含的标签点数量
*
* \param handle   连接句柄
* \param id       整型，输入，表ID
* \param size     整型，输出，表中标签点数量
* 注意：通过此API获取标签点数量，然后搜索此表中的标签点得到的数量可能会不一致，这是由于服务内部批量建点采取了异步的方式。
*       一般情况下请使用rtdbb_get_table_size_by_id来获取表中的标签点数量。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_table_real_size_by_id_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 *size)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_table_real_size_by_id_fn)(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 *size);
    rtdbb_get_table_real_size_by_id_fn fn = (rtdbb_get_table_real_size_by_id_fn)get_function("rtdbb_get_table_real_size_by_id");
    return fn(handle, id, size);
}

/**
*
* \brief 根据标签点表 id 获取表属性
*
* \param handle 连接句柄
* \param field  RTDB_TABLE 结构，输入/输出，标签点表属性，
*               输入时指定 id 字段，输出时返回 type、name、desc 字段。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_table_property_by_id_warp(rtdb_int32 handle, RTDB_TABLE *field)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_table_property_by_id_fn)(rtdb_int32 handle, RTDB_TABLE *field);
    rtdbb_get_table_property_by_id_fn fn = (rtdbb_get_table_property_by_id_fn)get_function("rtdbb_get_table_property_by_id");
    return fn(handle, field);
}

/**
*
* \brief 根据表名获取标签点表属性
*
* \param handle 连接句柄
* \param field  RTDB_TABLE 结构，输入/输出，标签点表属性
*               输入时指定 name 字段，输出时返回 id、type、desc 字段。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_table_property_by_name_warp(rtdb_int32 handle, RTDB_TABLE *field)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_table_property_by_name_fn)(rtdb_int32 handle, RTDB_TABLE *field);
    rtdbb_get_table_property_by_name_fn fn = (rtdbb_get_table_property_by_name_fn)get_function("rtdbb_get_table_property_by_name");
    return fn(handle, field);
}

/**
*
* \brief 使用完整的属性集来创建单个标签点
*
* \param handle 连接句柄
* \param base RTDB_POINT 结构，输入/输出，
*      输入除 id, createdate, creator, changedate, changer 字段外的其它字段，输出 id 字段。
* \param scan RTDB_SCAN_POINT 结构，输入，采集标签点扩展属性集。
* \param calc RTDB_CALC_POINT 结构，输入，计算标签点扩展属性集。
* \remark 如果新建的标签点没有对应的扩展属性集，可置为空指针。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_insert_point_warp(rtdb_int32 handle, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_CALC_POINT *calc)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_insert_point_fn)(rtdb_int32 handle, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_CALC_POINT *calc);
    rtdbb_insert_point_fn fn = (rtdbb_insert_point_fn)get_function("rtdbb_insert_point");
    return fn(handle, base, scan, calc);
}

/**
* 命名：rtdbb_insert_max_point
* 功能：使用最大长度的完整属性集来创建单个标签点
* 参数：
*      [handle] 连接句柄
*      [base] RTDB_POINT 结构，输入/输出，
*      输入除 id, createdate, creator, changedate, changer 字段外的其它字段，输出 id 字段。
*      [scan] RTDB_SCAN_POINT 结构，输入，采集标签点扩展属性集。
*      [calc] RTDB_MAX_CALC_POINT 结构，输入，计算标签点扩展属性集。
* 备注：如果新建的标签点没有对应的扩展属性集，可置为空指针。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_insert_max_point_warp(rtdb_int32 handle, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_MAX_CALC_POINT *calc)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_insert_max_point_fn)(rtdb_int32 handle, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_MAX_CALC_POINT *calc);
    rtdbb_insert_max_point_fn fn = (rtdbb_insert_max_point_fn)get_function("rtdbb_insert_max_point");
    return fn(handle, base, scan, calc);
}

/**
* 命名：rtdbb_insert_max_points
* 功能：使用最大长度的完整属性集来批量创建标签点
* 参数：
*      [handle] 连接句柄
*	   [count] count, 输入，base/scan/calc数组的长度；输出，成功的个数
*      [bases] RTDB_POINT 数组，输入/输出，
*      输入除 id, createdate, creator, changedate, changer 字段外的其它字段，输出 id 字段。
*      [scans] RTDB_SCAN_POINT 数组，输入，采集标签点扩展属性集。
*      [calcs] RTDB_MAX_CALC_POINT 数组，输入，计算标签点扩展属性集。
*	   [errors] rtdb_error数组，输出，对应每个标签点的结果
* 备注：如果新建的标签点没有对应的扩展属性集，可置为空指针。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_insert_max_points_warp(rtdb_int32 handle, rtdb_int32* count, RTDB_POINT* bases, RTDB_SCAN_POINT* scans, RTDB_MAX_CALC_POINT* calcs, rtdb_error* errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_insert_max_points_fn)(rtdb_int32 handle, rtdb_int32* count, RTDB_POINT* bases, RTDB_SCAN_POINT* scans, RTDB_MAX_CALC_POINT* calcs, rtdb_error* errors);
    rtdbb_insert_max_points_fn fn = (rtdbb_insert_max_points_fn)get_function("rtdbb_insert_max_points");
    return fn(handle, count, bases, scans, calcs, errors);
}

/**
*
* 功能  使用最小的属性集来创建单个标签点
*
* \param handle     连接句柄
* \param tag        字符串，输入，标签点名称
* \param type       整型，输入，标签点数据类型，取值 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、
*                   RTDB_CHAR、RTDB_UINT16、RTDB_UINT32、RTDB_INT32、RTDB_INT64、
*                   RTDB_REAL16、RTDB_REAL32、RTDB_REAL64、RTDB_COOR、RTDB_STRING、RTDB_BLOB 之一。
* \param table_id   整型，输入，标签点所属表 id
* \param use_ms     短整型，输入，标签点时间戳精度，0 为秒；1 为纳秒。
* \param point_id   整型，输出，标签点 id
* \remark 标签点的其余属性将取默认值。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_insert_base_point_warp(rtdb_int32 handle, const char *tag, rtdb_int32 type, rtdb_int32 table_id, rtdb_int16 use_ms, rtdb_int32 *point_id)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_insert_base_point_fn)(rtdb_int32 handle, const char *tag, rtdb_int32 type, rtdb_int32 table_id, rtdb_int16 use_ms, rtdb_int32 *point_id);
    rtdbb_insert_base_point_fn fn = (rtdbb_insert_base_point_fn)get_function("rtdbb_insert_base_point");
    return fn(handle, tag, type, table_id, use_ms, point_id);
}

/**
* 命名：rtdbb_insert_named_type_point
* 功能：使用完整的属性集来创建单个自定义数据类型标签点
* 参数：
*      [handle] 连接句柄
*      [base] RTDB_POINT 结构，输入/输出，
*      输入除 id, createdate, creator, changedate, changer 字段外的其它字段，输出 id 字段。
*      [scan] RTDB_SCAN_POINT 结构，输入，采集标签点扩展属性集。
*      [name] 字符串，输入，自定义数据类型的名字。
* 备注：如果新建的标签点没有对应的扩展属性集，可置为空指针。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_insert_named_type_point_warp(rtdb_int32 handle, RTDB_POINT *base, RTDB_SCAN_POINT *scan, const char* name)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_insert_named_type_point_fn)(rtdb_int32 handle, RTDB_POINT *base, RTDB_SCAN_POINT *scan, const char* name);
    rtdbb_insert_named_type_point_fn fn = (rtdbb_insert_named_type_point_fn)get_function("rtdbb_insert_named_type_point");
    return fn(handle, base, scan, name);
}

/**
*
* \brief 根据 id 删除单个标签点
*
* \param handle 连接句柄
* \param id     整型，输入，标签点标识
* \remark 通过本接口删除的标签点为可回收标签点，
*        可以通过 rtdbb_recover_point 接口恢复。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_remove_point_by_id_warp(rtdb_int32 handle, rtdb_int32 id)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_remove_point_by_id_fn)(rtdb_int32 handle, rtdb_int32 id);
    rtdbb_remove_point_by_id_fn fn = (rtdbb_remove_point_by_id_fn)get_function("rtdbb_remove_point_by_id");
    return fn(handle, id);
}

/**
*
* \brief 根据标签点全名删除单个标签点
* \param handle        连接句柄
* \param table_dot_tag  字符串，输入，标签点全名称："表名.标签点名"
* \remark 通过本接口删除的标签点为可回收标签点，
*        可以通过 rtdbb_recover_point 接口恢复。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_remove_point_by_name_warp(rtdb_int32 handle, const char *table_dot_tag)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_remove_point_by_name_fn)(rtdb_int32 handle, const char *table_dot_tag);
    rtdbb_remove_point_by_name_fn fn = (rtdbb_remove_point_by_name_fn)get_function("rtdbb_remove_point_by_name");
    return fn(handle, table_dot_tag);
}

/**
* 命名：rtdbb_move_point_by_id
* 功能：根据 id 移动单个标签点到其他表
* 参数：
*      [handle] 连接句柄
*      [id]     整型，输入，标签点标识
*	   [dest_table_name] 字符串，输入，移动的目标表名称
* 备注：通过本接口移动标签点后不改变标签点的id，且快照
*       和历史数据都不受影响
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_move_point_by_id_warp(rtdb_int32 handle, rtdb_int32 id, const char* dest_table_name)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_move_point_by_id_fn)(rtdb_int32 handle, rtdb_int32 id, const char* dest_table_name);
    rtdbb_move_point_by_id_fn fn = (rtdbb_move_point_by_id_fn)get_function("rtdbb_move_point_by_id");
    return fn(handle, id, dest_table_name);
}

/**
*
* \brief 批量获取标签点属性
*
* \param handle 连接句柄
* \param count  整数，输入，表示标签点个数。
* \param base   RTDB_POINT 结构数组，输入/输出，标签点基本属性列表，
*                 输入时，id 字段指定需要得到属性的标签点，输出时，其它字段返回标签点属性值。
* \param scan   RTDB_SCAN_POINT 结构数组，输出，采集标签点扩展属性列表
* \param calc   RTDB_CALC_POINT 结构数组，输出，计算标签点扩展属性列表
* \param errors 无符号整型数组，输出，获取标签属性的返回值列表，参考rtdb_error.h
* \remark 用户须保证分配给 base、scan、calc、errors 的空间与 count 相符，
*        扩展属性集 scan、calc 可为空指针，此时将不返回对应的扩展属性集。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_points_property_warp(rtdb_int32 handle, rtdb_int32 count, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_CALC_POINT *calc, rtdb_error *errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_points_property_fn)(rtdb_int32 handle, rtdb_int32 count, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_CALC_POINT *calc, rtdb_error *errors);
    rtdbb_get_points_property_fn fn = (rtdbb_get_points_property_fn)get_function("rtdbb_get_points_property");
    return fn(handle, count, base, scan, calc, errors);
}

/**
* 命名：rtdbb_get_max_points_property
* 功能：按最大长度批量获取标签点属性
* 参数：
*        [handle] 连接句柄
*        [count]  整数，输入，表示标签点个数。
*        [base]   RTDB_POINT 结构数组，输入/输出，标签点基本属性列表，
*                 输入时，id 字段指定需要得到属性的标签点，输出时，其它字段返回标签点属性值。
*        [scan]   RTDB_SCAN_POINT 结构数组，输出，采集标签点扩展属性列表
*        [calc]   RTDB_MAX_CALC_POINT 结构数组，输出，计算标签点扩展属性列表
*        [errors] 无符号整型数组，输出，获取标签属性的返回值列表，参考rtdb_error.h
* 备注：用户须保证分配给 base、scan、calc、errors 的空间与 count 相符，
*        扩展属性集 scan、calc 可为空指针，此时将不返回对应的扩展属性集。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_max_points_property_warp(rtdb_int32 handle, rtdb_int32 count, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_MAX_CALC_POINT *calc, rtdb_error *errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_max_points_property_fn)(rtdb_int32 handle, rtdb_int32 count, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_MAX_CALC_POINT *calc, rtdb_error *errors);
    rtdbb_get_max_points_property_fn fn = (rtdbb_get_max_points_property_fn)get_function("rtdbb_get_max_points_property");
    return fn(handle, count, base, scan, calc, errors);
}

/**
*
* \brief 搜索符合条件的标签点，使用标签点名时支持通配符
*
* \param handle        连接句柄
* \param tagmask       字符串，输入，标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
* \param tablemask     字符串，输入，标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
* \param source        字符串，输入，数据源集合，字符串中的每个字符均表示一个数据源，
*                        空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
* \param unit          字符串，输入，标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，
*                        空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
* \param desc          字符串，输入，标签点描述的子集，描述中包含该参数的标签点均满足条件，
*                        空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
* \param instrument    字符串，输入参数，标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
* \param mode          整型，RTDB_SORT_BY_TABLE、RTDB_SORT_BY_TAG、RTDB_SORT_BY_ID 之一，
*                        搜索结果的排序模式，输入，缺省值为RTDB_SORT_BY_TABLE
* \param ids           整型数组，输出，返回搜索到的标签点标识列表
* \param count         整型，输入/输出，输入时表示 ids 的长度，输出时表示搜索到的标签点个数
* \remark 用户须保证分配给 ids 的空间与 count 相符，各参数中包含的搜索条件之间的关系为"与"的关系，
*        用包含通配符的标签点名称作搜索条件时，如果第一个字符不是通配符(如"ai67*")，会得到最快的搜索速度。
*        如果 tagmask、tablemask 为空指针，则表示使用缺省设置"*",
*        多个搜索条件可以通过空格分隔，比如"demo_*1 demo_*2"，会将满足demo_*1或者demo_*2条件的标签点搜索出来。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_search_warp(rtdb_int32 handle, const char *tagmask, const char *tablemask, const char *source, const char *unit, const char *desc, const char *instrument, rtdb_int32 mode, rtdb_int32 *ids, rtdb_int32 *count)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_search_fn)(rtdb_int32 handle, const char *tagmask, const char *tablemask, const char *source, const char *unit, const char *desc, const char *instrument, rtdb_int32 mode, rtdb_int32 *ids, rtdb_int32 *count);
    rtdbb_search_fn fn = (rtdbb_search_fn)get_function("rtdbb_search");
    return fn(handle, tagmask, tablemask, source, unit, desc, instrument, mode, ids, count);
}

/**
*
* \brief 分批继续搜索符合条件的标签点，使用标签点名时支持通配符
*
* \param handle        连接句柄
* \param start         整型，输入，搜索起始位置。
* \param tagmask       字符串，输入，标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
* \param tablemask     字符串，输入，标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
* \param source        字符串，输入，数据源集合，字符串中的每个字符均表示一个数据源，
*                        空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
* \param unit          字符串，输入，标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，
*                        空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
* \param desc          字符串，输入，标签点描述的子集，描述中包含该参数的标签点均满足条件，
*                        空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
* \param instrument    字符串，输入参数，标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
* \param mode          整型，RTDB_SORT_BY_TABLE、RTDB_SORT_BY_TAG、RTDB_SORT_BY_ID 之一，
*                        搜索结果的排序模式，输入，缺省值为RTDB_SORT_BY_TABLE
* \param ids           整型数组，输出，返回搜索到的标签点标识列表
* \param count         整型，输入/输出，输入时表示 ids 的长度，输出时表示搜索到的标签点个数
* \remark 用户须保证分配给 ids 的空间与 count 相符，各参数中包含的搜索条件之间的关系为"与"的关系，
*        用包含通配符的标签点名称作搜索条件时，如果第一个字符不是通配符(如"ai67*")，会得到最快的搜索速度。
*        如果 tagmask、tablemask 为空指针，则表示使用缺省设置"*"。
*        当搜索到的标签点数比提供的要小时，表示这是最后一批符合条件的标签点 (即全部搜索完毕),
*        多个搜索条件可以通过空格分隔，比如"demo_*1 demo_*2"，会将满足demo_*1或者demo_*2条件的标签点搜索出来。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_search_in_batches_warp(rtdb_int32 handle, rtdb_int32 start, const char *tagmask, const char *tablemask, const char *source, const char *unit, const char *desc, const char *instrument, rtdb_int32 mode, rtdb_int32 *ids, rtdb_int32 *count)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_search_in_batches_fn)(rtdb_int32 handle, rtdb_int32 start, const char *tagmask, const char *tablemask, const char *source, const char *unit, const char *desc, const char *instrument, rtdb_int32 mode, rtdb_int32 *ids, rtdb_int32 *count);
    rtdbb_search_in_batches_fn fn = (rtdbb_search_in_batches_fn)get_function("rtdbb_search_in_batches");
    return fn(handle, start, tagmask, tablemask, source, unit, desc, instrument, mode, ids, count);
}

/**
*
* \brief 搜索符合条件的标签点，使用标签点名时支持通配符
*
* \param handle        连接句柄
* \param tagmask       字符串，输入，标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
* \param tablemask     字符串，输入，标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
* \param source        字符串，输入，数据源集合，字符串中的每个字符均表示一个数据源，
*                        空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
* \param unit          字符串，输入，标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，
*                        空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
* \param desc          字符串，输入，标签点描述的子集，描述中包含该参数的标签点均满足条件，
*                        空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
* \param instrument    字符串，输入参数，标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
* \param typemask      字符串，输入参数，标签点类型名称。缺省设置为空，长度不得超过 RTDB_TYPE_NAME_SIZE,
*                        内置的普通数据类型可以使用 bool、uint8、datetime等字符串表示，不区分大小写，支持模糊搜索。
* \param classofmask   整型，输入参数，标签点的类别，缺省设置为-1，表示可以是任意类型的标签点，
*                        当使用标签点类型作为搜索条件时，必须是RTDB_CLASS枚举中的一项或者多项的组合。
* \param timeunitmask  整型，输入参数，标签点的时间戳精度，缺省设置为-1，表示可以是任意时间戳精度，
*                        当使用此时间戳精度作为搜索条件时，timeunitmask的值可以为0或1，0表示时间戳精度为秒，1表示纳秒
* \param othertypemask 整型，输入参数，使用其他标签点属性作为搜索条件，缺省设置为0，表示不作为搜索条件，
*                        当使用此参数作为搜索条件时，othertypemaskvalue作为对应的搜索值，
*                        此参数的取值可以参考rtdb.h文件中的RTDB_SEARCH_MASK。
* \param othertypemaskvalue
*                        字符串，输入参数，当使用其他标签点属性作为搜索条件时，此参数作为对应的搜索值，缺省设置为0，表示不作为搜索条件，
*                        如果othertypemask的值为0，或者RTDB_SEARCH_NULL，则此参数被忽略,
*                        当othertypemask对应的标签点属性为数值类型时，此搜索值只支持相等判断，
*                        当othertypemask对应的标签点属性为字符串类型时，此搜索值支持模糊搜索。
* \param mode          整型，RTDB_SORT_BY_TABLE、RTDB_SORT_BY_TAG、RTDB_SORT_BY_ID 之一，
*                        搜索结果的排序模式，输入，缺省值为RTDB_SORT_BY_TABLE
* \param ids           整型数组，输出，返回搜索到的标签点标识列表
* \param count         整型，输入/输出，输入时表示 ids 的长度，输出时表示搜索到的标签点个数
* \remark 用户须保证分配给 ids 的空间与 count 相符，各参数中包含的搜索条件之间的关系为"与"的关系，
*        用包含通配符的标签点名称作搜索条件时，如果第一个字符不是通配符(如"ai67*")，会得到最快的搜索速度。
*        如果 tagmask、tablemask 为空指针，则表示使用缺省设置"*",
*        多个搜索条件可以通过空格分隔，比如"demo_*1 demo_*2"，会将满足demo_*1或者demo_*2条件的标签点搜索出来。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_search_ex_warp(rtdb_int32 handle, const char *tagmask, const char *tablemask, const char *source, const char *unit, const char *desc, const char *instrument, const char *typemask, rtdb_int32 classofmask, rtdb_int32 timeunitmask, rtdb_int32 othertypemask, const char *othertypemaskvalue, rtdb_int32 mode, rtdb_int32 *ids, rtdb_int32 *count)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_search_ex_fn)(rtdb_int32 handle, const char *tagmask, const char *tablemask, const char *source, const char *unit, const char *desc, const char *instrument, const char *typemask, rtdb_int32 classofmask, rtdb_int32 timeunitmask, rtdb_int32 othertypemask, const char *othertypemaskvalue, rtdb_int32 mode, rtdb_int32 *ids, rtdb_int32 *count);
    rtdbb_search_ex_fn fn = (rtdbb_search_ex_fn)get_function("rtdbb_search_ex");
    return fn(handle, tagmask, tablemask, source, unit, desc, instrument, typemask, classofmask, timeunitmask, othertypemask, othertypemaskvalue, mode, ids, count);
}

/**
*
* \brief 搜索符合条件的标签点，使用标签点名时支持通配符
*
* \param handle        连接句柄
* \param tagmask       字符串，输入，标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
* \param tablemask     字符串，输入，标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
* \param source        字符串，输入，数据源集合，字符串中的每个字符均表示一个数据源，
*                        空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
* \param unit          字符串，输入，标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，
*                        空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
* \param desc          字符串，输入，标签点描述的子集，描述中包含该参数的标签点均满足条件，
*                        空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
* \param instrument    字符串，输入参数，标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
* \param typemask      字符串，输入参数，标签点类型名称。缺省设置为空，长度不得超过 RTDB_TYPE_NAME_SIZE,
*                        内置的普通数据类型可以使用 bool、uint8、datetime等字符串表示，不区分大小写，支持模糊搜索。
* \param classofmask   整型，输入参数，标签点的类别，缺省设置为-1，表示可以是任意类型的标签点，
*                        当使用标签点类型作为搜索条件时，必须是RTDB_CLASS枚举中的一项或者多项的组合。
* \param timeunitmask  整型，输入参数，标签点的时间戳精度，缺省设置为-1，表示可以是任意时间戳精度，
*                        当使用此时间戳精度作为搜索条件时，timeunitmask的值可以为0或1，0表示时间戳精度为秒，1表示纳秒
* \param othertypemask 整型，输入参数，使用其他标签点属性作为搜索条件，缺省设置为0，表示不作为搜索条件，
*                        当使用此参数作为搜索条件时，othertypemaskvalue作为对应的搜索值，
*                        此参数的取值可以参考rtdb.h文件中的RTDB_SEARCH_MASK。
* \param othertypemaskvalue
*                        字符串，输入参数，当使用其他标签点属性作为搜索条件时，此参数作为对应的搜索值，缺省设置为0，表示不作为搜索条件，
*                        如果othertypemask的值为0，或者RTDB_SEARCH_NULL，则此参数被忽略,
*                        当othertypemask对应的标签点属性为数值类型时，此搜索值只支持相等判断，
*                        当othertypemask对应的标签点属性为字符串类型时，此搜索值支持模糊搜索。
* \param count         整型，输出，表示搜索到的标签点个数
* \remark  各参数中包含的搜索条件之间的关系为"与"的关系，
*        用包含通配符的标签点名称作搜索条件时，如果第一个字符不是通配符(如"ai67*")，会得到最快的搜索速度。
*        如果 tagmask、tablemask 为空指针，则表示使用缺省设置"*",
*        多个搜索条件可以通过空格分隔，比如"demo_*1 demo_*2"，会将满足demo_*1或者demo_*2条件的标签点搜索出来。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_search_points_count_warp(rtdb_int32 handle, const char *tagmask, const char *tablemask, const char *source, const char *unit, const char *desc, const char *instrument, const char *typemask, rtdb_int32 classofmask, rtdb_int32 timeunitmask, rtdb_int32 othertypemask, const char *othertypemaskvalue, rtdb_int32 *count)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_search_points_count_fn)(rtdb_int32 handle, const char *tagmask, const char *tablemask, const char *source, const char *unit, const char *desc, const char *instrument, const char *typemask, rtdb_int32 classofmask, rtdb_int32 timeunitmask, rtdb_int32 othertypemask, const char *othertypemaskvalue, rtdb_int32 *count);
    rtdbb_search_points_count_fn fn = (rtdbb_search_points_count_fn)get_function("rtdbb_search_points_count");
    return fn(handle, tagmask, tablemask, source, unit, desc, instrument, typemask, classofmask, timeunitmask, othertypemask, othertypemaskvalue, count);
}

/**
* 命名：rtdbb_remove_table_by_id
* \brief 根据表 id 删除表及表中标签点
*
* \param handle        连接句柄
* \param id            整型，输入，表 id
* \remark 删除的表不可恢复，删除的标签点可以通过 rtdbb_recover_point 接口恢复。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_remove_table_by_id_warp(rtdb_int32 handle, rtdb_int32 id)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_remove_table_by_id_fn)(rtdb_int32 handle, rtdb_int32 id);
    rtdbb_remove_table_by_id_fn fn = (rtdbb_remove_table_by_id_fn)get_function("rtdbb_remove_table_by_id");
    return fn(handle, id);
}

/**
*
* \brief 根据表名删除表及表中标签点
*
* \param handle        连接句柄
* \param name          字符串，输入，表名称
* \remark 删除的表不可恢复，删除的标签点可以通过 rtdbb_recover_point 接口恢复。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_remove_table_by_name_warp(rtdb_int32 handle, const char *name)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_remove_table_by_name_fn)(rtdb_int32 handle, const char *name);
    rtdbb_remove_table_by_name_fn fn = (rtdbb_remove_table_by_name_fn)get_function("rtdbb_remove_table_by_name");
    return fn(handle, name);
}

/**
*
* \brief 更新单个标签点属性
*
* \param handle        连接句柄
* \param base RTDB_POINT 结构，输入，基本标签点属性集。
* \param scan RTDB_SCAN_POINT 结构，输入，采集标签点扩展属性集。
* \param calc RTDB_CALC_POINT 结构，输入，计算标签点扩展属性集。
* \remark 标签点由 base 参数的 id 字段指定，其中 id、table、type、millisecond 字段不能修改，
*      changedate、changer、createdate、creator 字段由系统维护，其余字段均可修改，
*      包括 classof 字段。输入参数中 scan、calc 可为空指针，对应的扩展属性将保持不变。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_update_point_property_warp(rtdb_int32 handle, const RTDB_POINT *base, const RTDB_SCAN_POINT *scan, const RTDB_CALC_POINT *calc)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_update_point_property_fn)(rtdb_int32 handle, const RTDB_POINT *base, const RTDB_SCAN_POINT *scan, const RTDB_CALC_POINT *calc);
    rtdbb_update_point_property_fn fn = (rtdbb_update_point_property_fn)get_function("rtdbb_update_point_property");
    return fn(handle, base, scan, calc);
}

/**
* 命名：rtdbb_update_max_point_property
* 功能：按最大长度更新单个标签点属性
* 参数：
*        [handle]        连接句柄
*        [base] RTDB_POINT 结构，输入，基本标签点属性集。
*        [scan] RTDB_SCAN_POINT 结构，输入，采集标签点扩展属性集。
*        [calc] RTDB_MAX_CALC_POINT 结构，输入，计算标签点扩展属性集。
* 备注：标签点由 base 参数的 id 字段指定，其中 id、table、type、millisecond 字段不能修改，
*      changedate、changer、createdate、creator 字段由系统维护，其余字段均可修改，
*      包括 classof 字段。输入参数中 scan、calc 可为空指针，对应的扩展属性将保持不变。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_update_max_point_property_warp(rtdb_int32 handle, const RTDB_POINT *base, const RTDB_SCAN_POINT *scan, const RTDB_MAX_CALC_POINT *calc)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_update_max_point_property_fn)(rtdb_int32 handle, const RTDB_POINT *base, const RTDB_SCAN_POINT *scan, const RTDB_MAX_CALC_POINT *calc);
    rtdbb_update_max_point_property_fn fn = (rtdbb_update_max_point_property_fn)get_function("rtdbb_update_max_point_property");
    return fn(handle, base, scan, calc);
}

/**
*
* \brief 根据 "表名.标签点名" 格式批量获取标签点标识
*
* \param handle           连接句柄
* \param count            整数，输入/输出，输入时表示标签点个数
*                           (即table_dot_tags、ids、types、classof、use_ms 的长度)，
*                           输出时表示找到的标签点个数
* \param table_dot_tags   字符串指针数组，输入，"表名.标签点名" 列表
* \param ids              整型数组，输出，标签点标识列表, 返回 0 表示未找到
* \param types            整型数组，输出，标签点数据类型
* \param classof          整型数组，输出，标签点类别
* \param use_ms           短整型数组，输出，时间戳精度，
*                           返回 1 表示时间戳精度为纳秒， 为 0 表示为秒。
* \remark 用户须保证分配给 table_dot_tags、ids、types、classof、use_ms 的空间与count相符，
*        其中 types、classof、use_ms 可为空指针，对应的字段将不再返回。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_find_points_warp(rtdb_int32 handle, rtdb_int32 *count, const char* const* table_dot_tags, rtdb_int32 *ids, rtdb_int32 *types, rtdb_int32 *classof, rtdb_int16 *use_ms)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_find_points_fn)(rtdb_int32 handle, rtdb_int32 *count, const char* const* table_dot_tags, rtdb_int32 *ids, rtdb_int32 *types, rtdb_int32 *classof, rtdb_int16 *use_ms);
    rtdbb_find_points_fn fn = (rtdbb_find_points_fn)get_function("rtdbb_find_points");
    return fn(handle, count, table_dot_tags, ids, types, classof, use_ms);
}

/**
*
* \brief 根据 "表名.标签点名" 格式批量获取标签点标识
*
* \param handle           连接句柄
* \param count            整数，输入/输出，输入时表示标签点个数
*                           (即table_dot_tags、ids、types、classof、use_ms 的长度)，
*                           输出时表示找到的标签点个数
* \param table_dot_tags   字符串指针数组，输入，"表名.标签点名" 列表
* \param ids              整型数组，输出，标签点标识列表, 返回 0 表示未找到
* \param types            整型数组，输出，标签点数据类型
* \param classof          整型数组，输出，标签点类别
* \param precisions       数组，输出，时间戳精度，
*                           0表示秒，1表示毫秒，2表示微秒，3纳秒。
* \param errors           无符号整型数组，输出，表示每个标签点的查询结果的错误码
* \remark 用户须保证分配给 table_dot_tags、ids、types、classof、precisions、errors 的空间与count相符，
*        其中 types、classof、precisions、errors 可为空指针，对应的字段将不再返回。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_find_points_ex_warp(rtdb_int32 handle, rtdb_int32* count, const char* const* table_dot_tags, rtdb_int32* ids, rtdb_int32* types, rtdb_int32* classof, rtdb_precision_type* precisions, rtdb_error* errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_find_points_ex_fn)(rtdb_int32 handle, rtdb_int32* count, const char* const* table_dot_tags, rtdb_int32* ids, rtdb_int32* types, rtdb_int32* classof, rtdb_precision_type* precisions, rtdb_error* errors);
    rtdbb_find_points_ex_fn fn = (rtdbb_find_points_ex_fn)get_function("rtdbb_find_points_ex");
    return fn(handle, count, table_dot_tags, ids, types, classof, precisions, errors);
}

/**
*
* \brief 根据标签属性字段对标签点标识进行排序
*
* \param handle           连接句柄
* \param count            整数，输入，表示标签点个数, 即 ids 的长度
* \param ids              整型数组，输入，标签点标识列表
* \param index            整型，输入，属性字段枚举，参见 RTDB_TAG_FIELD_INDEX，
*                           将根据该字段对 ID 进行排序。
* \param flag             整型，输入，标志位组合，参见 RTDB_TAG_SORT_FLAG 枚举，其中
*                           RTDB_SORT_FLAG_DESCEND             表示降序排序，不设置表示升序排列；
*                           RTDB_SORT_FLAG_CASE_SENSITIVE      表示进行字符串类型字段比较时大小写敏感，不设置表示不区分大小写；
*                           RTDB_SORT_FLAG_RECYCLED            表示对可回收标签进行排序，不设置表示对正常标签排序，
*                           不同的标志位可通过"或"运算连接在一起，
*                           当对可回收标签排序时，以下字段索引不可使用：
*                               RTDB_TAG_INDEX_TIMESTAMP
*                               RTDB_TAG_INDEX_VALUE
*                               RTDB_TAG_INDEX_QUALITY
* \remark 用户须保证分配给 ids 的空间与 count 相符, 如果 ID 指定的标签并不存在，
*        或标签不具备要求排序的字段 (如对非计算点进行方程式排序)，它们将被放置在数组的尾部。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_sort_points_warp(rtdb_int32 handle, rtdb_int32 count, rtdb_int32 *ids, rtdb_int32 index, rtdb_int32 flag)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_sort_points_fn)(rtdb_int32 handle, rtdb_int32 count, rtdb_int32 *ids, rtdb_int32 index, rtdb_int32 flag);
    rtdbb_sort_points_fn fn = (rtdbb_sort_points_fn)get_function("rtdbb_sort_points");
    return fn(handle, count, ids, index, flag);
}

/**
*
* \brief 根据表 ID 更新表名称。
*
* \param handle    连接句柄
* \param tab_id    整型，输入，要修改表的标识
* \param name      字符串，输入，新的标签点表名称。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_update_table_name_warp(rtdb_int32 handle, rtdb_int32 tab_id, const char *name)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_update_table_name_fn)(rtdb_int32 handle, rtdb_int32 tab_id, const char *name);
    rtdbb_update_table_name_fn fn = (rtdbb_update_table_name_fn)get_function("rtdbb_update_table_name");
    return fn(handle, tab_id, name);
}

/**
*
* \brief 根据表 ID 更新表描述。
*
* \param handle    连接句柄
* \param tab_id    整型，输入，要修改表的标识
* \param desc      字符串，输入，新的表描述。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_update_table_desc_by_id_warp(rtdb_int32 handle, rtdb_int32 tab_id, const char *desc)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_update_table_desc_by_id_fn)(rtdb_int32 handle, rtdb_int32 tab_id, const char *desc);
    rtdbb_update_table_desc_by_id_fn fn = (rtdbb_update_table_desc_by_id_fn)get_function("rtdbb_update_table_desc_by_id");
    return fn(handle, tab_id, desc);
}

/**
*
* \brief 根据表名称更新表描述。
*
* \param handle    连接句柄
* \param name      字符串，输入，要修改表的名称。
* \param desc      字符串，输入，新的表描述。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_update_table_desc_by_name_warp(rtdb_int32 handle, const char *name, const char *desc)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_update_table_desc_by_name_fn)(rtdb_int32 handle, const char *name, const char *desc);
    rtdbb_update_table_desc_by_name_fn fn = (rtdbb_update_table_desc_by_name_fn)get_function("rtdbb_update_table_desc_by_name");
    return fn(handle, name, desc);
}

/**
*
* \brief 恢复已删除标签点
*
* \param handle    连接句柄
* \param table_id  整型，输入，要将标签点恢复到的表标识
* \param point_id  整型，输入，待恢复的标签点标识
* 备注: 本接口只对可回收标签点(通过接口rtdbb_remove_point_by_id/rtdbb_remove_point_by_tag)有效，
*        对正常的标签点没有作用。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_recover_point_warp(rtdb_int32 handle, rtdb_int32 table_id, rtdb_int32 point_id)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_recover_point_fn)(rtdb_int32 handle, rtdb_int32 table_id, rtdb_int32 point_id);
    rtdbb_recover_point_fn fn = (rtdbb_recover_point_fn)get_function("rtdbb_recover_point");
    return fn(handle, table_id, point_id);
}

/**
*
* \brief 清除标签点
*
* \param handle    连接句柄
* \param id        整数，输入，要清除的标签点标识
* 备注: 本接口仅对可回收标签点(通过接口rtdbb_remove_point_by_id/rtdbb_remove_point_by_name)有效，
*      对正常的标签点没有作用。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_purge_point_warp(rtdb_int32 handle, rtdb_int32 id)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_purge_point_fn)(rtdb_int32 handle, rtdb_int32 id);
    rtdbb_purge_point_fn fn = (rtdbb_purge_point_fn)get_function("rtdbb_purge_point");
    return fn(handle, id);
}

/**
*
* \brief 获取可回收标签点数量
*
* \param handle    连接句柄
* \param count     整型，输出，可回收标签点的数量
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_recycled_points_count_warp(rtdb_int32 handle, rtdb_int32 *count)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_recycled_points_count_fn)(rtdb_int32 handle, rtdb_int32 *count);
    rtdbb_get_recycled_points_count_fn fn = (rtdbb_get_recycled_points_count_fn)get_function("rtdbb_get_recycled_points_count");
    return fn(handle, count);
}

/**
*
* \brief 获取可回收标签点 id 列表
*
* \param handle    连接句柄
* \param ids       整型数组，输出，可回收标签点 id
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids 的长度，
*                    输出时表示成功获取标签点的个数。
* \remark 用户须保证 ids 的长度与 count 一致
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_recycled_points_warp(rtdb_int32 handle, rtdb_int32 *ids, rtdb_int32 *count)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_recycled_points_fn)(rtdb_int32 handle, rtdb_int32 *ids, rtdb_int32 *count);
    rtdbb_get_recycled_points_fn fn = (rtdbb_get_recycled_points_fn)get_function("rtdbb_get_recycled_points");
    return fn(handle, ids, count);
}

/**
* 命名：rtdbb_search_recycled_points
* 功能：搜索符合条件的可回收标签点，使用标签点名时支持通配符
* 参数：
*        [handle]        连接句柄
*        [tagmask]       字符串，输入，标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE。
*        [tablemask]     字符串，输入，标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE。
*        [source]        字符串，输入，数据源集合，字符串中的每个字符均表示一个数据源，
*                        空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
*        [unit]          字符串，输入，标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，
*                        空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
*        [desc]          字符串，输入，标签点描述的子集，描述中包含该参数的标签点均满足条件，
*                        空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
*        [instrument]    字符串，输入参数，标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
*        [mode]          整型，RTDB_SORT_BY_TABLE、RTDB_SORT_BY_TAG、RTDB_SORT_BY_ID 之一，
*                        搜索结果的排序模式，输入，缺省值为RTDB_SORT_BY_TABLE
*        [ids]           整型数组，输出，返回搜索到的标签点标识列表
*        [count]         整型，输入/输出，输入时表示 ids 的长度，输出时表示搜索到的标签点个数
* 备注：用户须保证分配给 ids 的空间与 count 相符，各参数中包含的搜索条件之间的关系为"与"的关系，
*        用包含通配符的标签点名称作搜索条件时，如果第一个字符不是通配符(如"ai67*")，会得到最快的搜索速度。
*        如果 tagmask、fullmask 为空指针，则表示使用缺省设置"*"
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_search_recycled_points_warp(rtdb_int32 handle, const char *tagmask, const char *fullmask, const char *source, const char *unit, const char *desc, const char *instrument, rtdb_int32 mode, rtdb_int32 *ids, rtdb_int32 *count)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_search_recycled_points_fn)(rtdb_int32 handle, const char *tagmask, const char *fullmask, const char *source, const char *unit, const char *desc, const char *instrument, rtdb_int32 mode, rtdb_int32 *ids, rtdb_int32 *count);
    rtdbb_search_recycled_points_fn fn = (rtdbb_search_recycled_points_fn)get_function("rtdbb_search_recycled_points");
    return fn(handle, tagmask, fullmask, source, unit, desc, instrument, mode, ids, count);
}

/**
*
* \brief 分批搜索符合条件的可回收标签点，使用标签点名时支持通配符
*
* \param handle        连接句柄
* \param start         整型，输入，搜索的起始位置。
* \param tagmask       字符串，输入，标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE。
* \param tablemask     字符串，输入，标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE。
* \param source        字符串，输入，数据源集合，字符串中的每个字符均表示一个数据源，
*                        空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
* \param unit          字符串，输入，标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，
*                        空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
* \param desc          字符串，输入，标签点描述的子集，描述中包含该参数的标签点均满足条件，
*                        空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
* \param instrument    字符串，输入参数，标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
* \param mode          整型，RTDB_SORT_BY_TABLE、RTDB_SORT_BY_TAG、RTDB_SORT_BY_ID 之一，
*                        搜索结果的排序模式，输入，缺省值为RTDB_SORT_BY_TABLE
* \param ids           整型数组，输出，返回搜索到的标签点标识列表
* \param count         整型，输入/输出，输入时表示 ids 的长度，输出时表示搜索到的标签点个数
* \remark 用户须保证分配给 ids 的空间与 count 相符，各参数中包含的搜索条件之间的关系为"与"的关系，
*        用包含通配符的标签点名称作搜索条件时，如果第一个字符不是通配符(如"ai67*")，会得到最快的搜索速度。
*        如果 tagmask、fullmask 为空指针，则表示使用缺省设置"*"
*        当搜索到的标签点数比提供的要小时，表示这是最后一批符合条件的标签点 (即全部搜索完毕)。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_search_recycled_points_in_batches_warp(rtdb_int32 handle, rtdb_int32 start, const char *tagmask, const char *fullmask, const char *source, const char *unit, const char *desc, const char *instrument, rtdb_int32 mode, rtdb_int32 *ids, rtdb_int32 *count)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_search_recycled_points_in_batches_fn)(rtdb_int32 handle, rtdb_int32 start, const char *tagmask, const char *fullmask, const char *source, const char *unit, const char *desc, const char *instrument, rtdb_int32 mode, rtdb_int32 *ids, rtdb_int32 *count);
    rtdbb_search_recycled_points_in_batches_fn fn = (rtdbb_search_recycled_points_in_batches_fn)get_function("rtdbb_search_recycled_points_in_batches");
    return fn(handle, start, tagmask, fullmask, source, unit, desc, instrument, mode, ids, count);
}

/**
*
* \brief 获取可回收标签点的属性
*
* \param handle   连接句柄
* \param base     RTDB_POINT 结构，输入/输出，标签点基本属性。
输入时，由 id 字段指定要取得的可回收标签点。
* \param scan     RTDB_SCAN_POINT 结构，输出，标签点采集扩展属性
* \param calc     RTDB_CALC_POINT 结构，输出，标签点计算扩展属性
* \remark scan、calc 可为空指针，对应的扩展信息将不返回。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_recycled_point_property_warp(rtdb_int32 handle, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_CALC_POINT *calc)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_recycled_point_property_fn)(rtdb_int32 handle, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_CALC_POINT *calc);
    rtdbb_get_recycled_point_property_fn fn = (rtdbb_get_recycled_point_property_fn)get_function("rtdbb_get_recycled_point_property");
    return fn(handle, base, scan, calc);
}

/**

* 命名：rtdbb_get_recycled_max_point_property
* 功能：按最大长度获取可回收标签点的属性
* 参数：
*        [handle]   连接句柄
*        [base]     RTDB_POINT 结构，输入/输出，标签点基本属性。
                    输入时，由 id 字段指定要取得的可回收标签点。
*        [scan]     RTDB_SCAN_POINT 结构，输出，标签点采集扩展属性
*        [calc]     RTDB_MAX_CALC_POINT 结构，输出，标签点计算扩展属性
* 备注：scan、calc 可为空指针，对应的扩展信息将不返回。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_recycled_max_point_property_warp(rtdb_int32 handle, RTDB_POINT* base, RTDB_SCAN_POINT* scan, RTDB_MAX_CALC_POINT* calc)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_recycled_max_point_property_fn)(rtdb_int32 handle, RTDB_POINT* base, RTDB_SCAN_POINT* scan, RTDB_MAX_CALC_POINT* calc);
    rtdbb_get_recycled_max_point_property_fn fn = (rtdbb_get_recycled_max_point_property_fn)get_function("rtdbb_get_recycled_max_point_property");
    return fn(handle, base, scan, calc);
}

/**
*
* \brief 清空标签点回收站
*
* \param handle   连接句柄
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_clear_recycler_warp(rtdb_int32 handle)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_clear_recycler_fn)(rtdb_int32 handle);
    rtdbb_clear_recycler_fn fn = (rtdbb_clear_recycler_fn)get_function("rtdbb_clear_recycler");
    return fn(handle);
}

/**
* 命名：rtdbb_subscribe_tags_ex
* 功能：标签点属性更改通知订阅
* 参数：
*        [handle]    连接句柄
*        [options]   整型，输入，订阅选项，参见枚举RTDB_OPTION
*                    RTDB_O_AUTOCONN 订阅客户端与数据库服务器网络中断后自动重连并订阅
*        [param]     输入，用户参数，
*                    作为rtdbb_tags_change_ex的param参数
*        [callback]  rtdbb_tags_change_ex 类型回调接口，输入，当回掉函数返回非RtE_OK时退出订阅
*                    当未设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
*                    作为event_type取值调用回掉函数后退出订阅。
*                    当设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
*                    作为event_type取值调用回掉函数直到连接恢复或回掉函数返回非RtE_OK，
*                    网络中断期间回掉函数调用频率为最少3秒
*                    event_type参数值含义如下：
*                      RTDB_E_DATA        标签点属性发生更改
*                      RTDB_E_DISCONNECT  订阅客户端与数据库网络断开
*                      RTDB_E_RECOVERY    订阅客户端与数据库网络及订阅恢复
*                    handle 产生订阅回掉的连接句柄，调用rtdbb_subscribe_tags_ex时的handle参数
*                    param  用户自定义参数，调用rtdbb_subscribe_tags_ex时的param参数
*                    count  event_type为RTDB_E_DATA时表示ids的数量
*                           event_type为其它值时，count值为0
*                    ids    event_type为RTDB_E_DATA时表示属性更改的标签点ID，数量由count指定
*                           event_type为其它值时，ids值为NULL
*                    what   event_type为RTDB_E_DATA时表示属性变更原因，参考RTDB_TAG_CHANGE_REASON
*                           event_type为其它值时，what时值为0
* 备注：用于订阅测点的连接句柄必需是独立的，不能再用来调用其它 api，
*       否则返回 RtE_OTHER_SDK_DOING 错误。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_subscribe_tags_ex_warp(rtdb_int32 handle, rtdb_uint32 options, void* param, rtdbb_tags_change_event_ex callback)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_subscribe_tags_ex_fn)(rtdb_int32 handle, rtdb_uint32 options, void* param, rtdbb_tags_change_event_ex callback);
    rtdbb_subscribe_tags_ex_fn fn = (rtdbb_subscribe_tags_ex_fn)get_function("rtdbb_subscribe_tags_ex");
    return fn(handle, options, param, callback);
}

/**
*
* \brief 取消标签点属性更改通知订阅
*
* \param handle    连接句柄
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_cancel_subscribe_tags_warp(rtdb_int32 handle)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_cancel_subscribe_tags_fn)(rtdb_int32 handle);
    rtdbb_cancel_subscribe_tags_fn fn = (rtdbb_cancel_subscribe_tags_fn)get_function("rtdbb_cancel_subscribe_tags");
    return fn(handle);
}

/**
* 命名：rtdbb_create_named_type
* 功能：创建自定义类型
* 参数：
*        [handle]      连接句柄，输入参数
*        [name]        自定义类型的名称，类型的唯一标示,不能重复，长度不能超过RTDB_TYPE_NAME_SIZE，输入参数
*        [field_count]    自定义类型中包含的字段的个数,输入参数
*        [fields]      自定义类型中包含的字段的属性，RTDB_DATA_TYPE_FIELD结构的数组，个数与field_count相等，输入参数
*              RTDB_DATA_TYPE_FIELD中的length只对type为str或blob类型的数据有效。其他类型忽略
* 备注：自定义类型的大小必须要小于数据页大小(小于数据页大小的2/3，即需要合理定义字段的个数及每个字段的长度)。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_create_named_type_warp(rtdb_int32 handle, const char* name, rtdb_int32 field_count, const RTDB_DATA_TYPE_FIELD* fields, char desc[RTDB_DESC_SIZE])
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_create_named_type_fn)(rtdb_int32 handle, const char* name, rtdb_int32 field_count, const RTDB_DATA_TYPE_FIELD* fields, char desc[RTDB_DESC_SIZE]);
    rtdbb_create_named_type_fn fn = (rtdbb_create_named_type_fn)get_function("rtdbb_create_named_type");
    return fn(handle, name, field_count, fields, desc);
}

/**
* 命名：rtdbb_get_named_types_count
* 功能：获取所有的自定义类型的总数
* 参数：
*        [handle]      连接句柄，输入参数
*        [count]      返回所有的自定义类型的总数，输入/输出参数
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_named_types_count_warp(rtdb_int32 handle, rtdb_int32* count)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_named_types_count_fn)(rtdb_int32 handle, rtdb_int32* count);
    rtdbb_get_named_types_count_fn fn = (rtdbb_get_named_types_count_fn)get_function("rtdbb_get_named_types_count");
    return fn(handle, count);
}

/**
* 命名：rtdbb_get_all_named_types
* 功能：获取所有的自定义类型
* 参数：
*        [handle]      连接句柄，输入参数
*        [count]      返回所有的自定义类型的总数，输入/输出参数，输入:为name,field_counts数组的长度，输出:获取的实际自定义类型的个数
*        [name]        返回所有的自定义类型的名称的数组，每个自定义类型的名称的长度不超过RTDB_TYPE_NAME_SIZE，输入/输出参数
*              输入：name数组长度要等于count.输出：实际获取的自定义类型名称的数组
*        [field_counts]    返回所有的自定义类型所包含字段个数的数组，输入/输出参数
*              输入：field_counts数组长度要等于count。输出:实际每个自定义类型所包含的字段的个数的数组
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_all_named_types_warp(rtdb_int32 handle, rtdb_int32* count, char* name[RTDB_TYPE_NAME_SIZE], rtdb_int32* field_counts)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_all_named_types_fn)(rtdb_int32 handle, rtdb_int32* count, char* name[RTDB_TYPE_NAME_SIZE], rtdb_int32* field_counts);
    rtdbb_get_all_named_types_fn fn = (rtdbb_get_all_named_types_fn)get_function("rtdbb_get_all_named_types");
    return fn(handle, count, name, field_counts);
}

/**
* 命名：rtdbb_get_named_type
* 功能：获取自定义类型的所有字段
* 参数：
*        [handle]         连接句柄，输入参数
*        [name]           自定义类型的名称，输入参数
*        [field_count]    返回name指定的自定义类型的字段个数，输入/输出参数
*                         输入：指定fields数组长度.输出：实际的name自定义类型的字段的个数
*        [fields]         返回由name所指定的自定义类型所包含字段RTDB_DATA_TYPE_FIELD结构的数组，输入/输出参数
*                         输入：fields数组长度要等于count。输出:RTDB_DATA_TYPE_FIELD结构的数组
*        [type_size]      所有自定义类型fields结构中长度字段的累加和，输出参数
*        [desc]           自定义类型的描述，输出参数
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_named_type_warp(rtdb_int32 handle, const char* name, rtdb_int32* field_count, RTDB_DATA_TYPE_FIELD* fields, rtdb_int32* type_size, char desc[RTDB_DESC_SIZE])
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_named_type_fn)(rtdb_int32 handle, const char* name, rtdb_int32* field_count, RTDB_DATA_TYPE_FIELD* fields, rtdb_int32* type_size, char desc[RTDB_DESC_SIZE]);
    rtdbb_get_named_type_fn fn = (rtdbb_get_named_type_fn)get_function("rtdbb_get_named_type");
    return fn(handle, name, field_count, fields, type_size, desc);
}

/**
* 命名：rtdbb_remove_named_type
* 功能：删除自定义类型
* 参数：
*        [handle]      连接句柄，输入参数
*        [name]        自定义类型的名称，输入参数
*        [reserved]      保留字段,暂时不用
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_remove_named_type_warp(rtdb_int32 handle, const char* name, rtdb_int32 reserved GAPI_DEFAULT_VALUE(0))
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_remove_named_type_fn)(rtdb_int32 handle, const char* name, rtdb_int32 reserved GAPI_DEFAULT_VALUE(0));
    rtdbb_remove_named_type_fn fn = (rtdbb_remove_named_type_fn)get_function("rtdbb_remove_named_type");
    return fn(handle, name, reserved);
}

/**
* 命名：rtdbb_get_named_type_names_property
* 功能：根据标签点id查询标签点所对应的自定义类型的名字和字段总数
* 参数：
*        [handle]           连接句柄
*        [count]            输入/输出，标签点个数，
*                           输入时表示 ids、named_type_names、field_counts、errors 的长度，
*                           输出时表示成功获取自定义类型名字的标签点个数
*        [ids]              整型数组，输入，标签点标识列表
*        [named_type_names] 字符串数组，输出，标签点自定义类型的名字
*        [field_counts]     整型数组，输出，标签点自定义类型的字段个数
*        [errors]           无符号整型数组，输出，获取自定义类型名字的返回值列表，参考rtdb_error.h
* 备注：用户须保证 ids、named_type_names、field_counts、errors 的长度与 count 一致。
*        本接口只对数据类型为 RTDB_NAMED_T 的标签点有效。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_named_type_names_property_warp(rtdb_int32 handle, rtdb_int32 *count, rtdb_int32 *ids, char* const *named_type_names, rtdb_int32 *field_counts, rtdb_error *errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_named_type_names_property_fn)(rtdb_int32 handle, rtdb_int32 *count, rtdb_int32 *ids, char* const *named_type_names, rtdb_int32 *field_counts, rtdb_error *errors);
    rtdbb_get_named_type_names_property_fn fn = (rtdbb_get_named_type_names_property_fn)get_function("rtdbb_get_named_type_names_property");
    return fn(handle, count, ids, named_type_names, field_counts, errors);
}

/**
* 命名：rtdbb_get_recycled_named_type_names_property
* 功能：根据回收站标签点id查询标签点所对应的自定义类型的名字和字段总数
* 参数：
*        [handle]           连接句柄
*        [count]            输入/输出，标签点个数，
*                           输入时表示 ids、named_type_names、field_counts、errors 的长度，
*                           输出时表示成功获取自定义类型名字的标签点个数
*        [ids]              整型数组，输入，回收站标签点标识列表
*        [named_type_names] 字符串数组，输出，标签点自定义类型的名字
*        [field_counts]     整型数组，输出，标签点自定义类型的字段个数
*        [errors]           无符号整型数组，输出，获取自定义类型名字的返回值列表，参考rtdb_error.h
* 备注：用户须保证 ids、named_type_names、field_counts、errors 的长度与 count 一致。
*        本接口只对数据类型为 RTDB_NAMED_T 的标签点有效。
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_recycled_named_type_names_property_warp(rtdb_int32 handle, rtdb_int32 *count, rtdb_int32 *ids, char* const *named_type_names, rtdb_int32 *field_counts, rtdb_error *errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_recycled_named_type_names_property_fn)(rtdb_int32 handle, rtdb_int32 *count, rtdb_int32 *ids, char* const *named_type_names, rtdb_int32 *field_counts, rtdb_error *errors);
    rtdbb_get_recycled_named_type_names_property_fn fn = (rtdbb_get_recycled_named_type_names_property_fn)get_function("rtdbb_get_recycled_named_type_names_property");
    return fn(handle, count, ids, named_type_names, field_counts, errors);
}

/**
* 命名：rtdbb_get_named_type_points_count
* 功能：获取该自定义类型的所有标签点个数
* 参数：
*        [handle]           连接句柄，输入参数
*        [name]             自定义类型的名称，输入参数
*        [points_count]     返回name指定的自定义类型的标签点个数，输入参数
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_named_type_points_count_warp(rtdb_int32 handle, const char* name, rtdb_int32 *points_count)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_named_type_points_count_fn)(rtdb_int32 handle, const char* name, rtdb_int32 *points_count);
    rtdbb_get_named_type_points_count_fn fn = (rtdbb_get_named_type_points_count_fn)get_function("rtdbb_get_named_type_points_count");
    return fn(handle, name, points_count);
}

/**
*
* \brief 获取该内置的基本类型的所有标签点个数
*
* \param handle           整型，输入参数，连接句equation[RTDB_MAX_EQUATION_SIZE]柄
* \param type             整型，输入参数，内置的基本类型，参数的值可以是除RTDB_NAME_T以外的所有RTDB_TYPE枚举值
* \param points_count     整型，输入参数，返回type指定的内置基本类型的标签点个数
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_base_type_points_count_warp(rtdb_int32 handle, rtdb_int32 type, rtdb_int32 *points_count)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_base_type_points_count_fn)(rtdb_int32 handle, rtdb_int32 type, rtdb_int32 *points_count);
    rtdbb_get_base_type_points_count_fn fn = (rtdbb_get_base_type_points_count_fn)get_function("rtdbb_get_base_type_points_count");
    return fn(handle, type, points_count);
}

/**
* 命名：rtdbb_modify_named_type
* 功能：修改自定义类型名称,描述,字段名称,字段描述
* 参数：
*        [handle]             连接句柄，输入参数
*        [name]               自定义类型的名称，输入参数
*        [modify_name]        要修改的自定义类型名称，输入参数
*        [modify_desc]        要修改的自定义类型的描述，输入参数
*        [modify_field_name]  要修改的自定义类型字段的名称，输入参数
*        [modify_field_desc]  要修改的自定义类型字段的描述，输入参数
*        [field_count]        自定义类型字段的个数，输入参数
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_modify_named_type_warp(rtdb_int32 handle, const char* name, const char* modify_name, const char* modify_desc, const char* modify_field_name[RTDB_TYPE_NAME_SIZE], const char* modify_field_desc[RTDB_DESC_SIZE], rtdb_int32 field_count)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_modify_named_type_fn)(rtdb_int32 handle, const char* name, const char* modify_name, const char* modify_desc, const char* modify_field_name[RTDB_TYPE_NAME_SIZE], const char* modify_field_desc[RTDB_DESC_SIZE], rtdb_int32 field_count);
    rtdbb_modify_named_type_fn fn = (rtdbb_modify_named_type_fn)get_function("rtdbb_modify_named_type");
    return fn(handle, name, modify_name, modify_desc, modify_field_name, modify_field_desc, field_count);
}

/**
*
* \brief 获取元数据同步信息
*
* \param handle           整型，输入参数，连接句柄
* \param node_number      整型，输入参数，双活节点id，1表示第一个节点，2表示第二个节点。0表示所有节点
* \param count            整型，输入参数，sync_infos参数的数量
*                              输出参数，输出实际获取到的sync_infos的个数
* \param sync_infos       RTDB_SYNC_INFO数组，输出参数，输出实际获取到的同步信息
* \param errors           rtdb_error数组，输出参数，输出对应节点的错误信息
*/
rtdb_error RTDBAPI_CALLRULE rtdbb_get_meta_sync_info_warp(rtdb_int32 handle, rtdb_int32 node_number, rtdb_int32* count, RTDB_SYNC_INFO* sync_infos, rtdb_error* errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbb_get_meta_sync_info_fn)(rtdb_int32 handle, rtdb_int32 node_number, rtdb_int32* count, RTDB_SYNC_INFO* sync_infos, rtdb_error* errors);
    rtdbb_get_meta_sync_info_fn fn = (rtdbb_get_meta_sync_info_fn)get_function("rtdbb_get_meta_sync_info");
    return fn(handle, node_number, count, sync_infos, errors);
}

/**
*
* \brief 批量读取开关量、模拟量快照数值
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表
* \param datetimes 整型数组，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
* \param values    双精度浮点型数组，输出，实时浮点型数值列表，
*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，返回相应的快照值；否则为 0
* \param states    64 位整型数组，输出，实时整型数值列表，
*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，返回相应的快照值；否则为 0
* \param qualities 短整型数组，输出，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
*        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_get_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities, rtdb_error* errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_get_snapshots64_fn)(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities, rtdb_error* errors);
    rtdbs_get_snapshots64_fn fn = (rtdbs_get_snapshots64_fn)get_function("rtdbs_get_snapshots64");
    return fn(handle, count, ids, datetimes, subtimes, values, states, qualities, errors);
}

/**
*
* \brief 批量写入开关量、模拟量快照数值
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
*                    输出时表示成功写入实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表，同一个标签点标识可以出现多次，
*                    但它们的时间戳必需是递增的。
* \param datetimes 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
* \param values    双精度浮点型数组，输入，实时浮点型数值列表，
*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的快照值；否则忽略
* \param states    64 位整型数组，输入，实时整型数值列表，
*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的快照值；否则忽略
* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，写入实时数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
*        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_put_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float64* values, const rtdb_int64* states, const rtdb_int16* qualities, rtdb_error* errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_put_snapshots64_fn)(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float64* values, const rtdb_int64* states, const rtdb_int16* qualities, rtdb_error* errors);
    rtdbs_put_snapshots64_fn fn = (rtdbs_put_snapshots64_fn)get_function("rtdbs_put_snapshots64");
    return fn(handle, count, ids, datetimes, subtimes, values, states, qualities, errors);
}

/**
*
* \brief 批量写入开关量、模拟量快照数值
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
*                    输出时表示成功写入实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表，同一个标签点标识可以出现多次，
*                    但它们的时间戳必需是递增的。
* \param datetimes 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
* \param values    双精度浮点型数组，输入，实时浮点型数值列表，
*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的快照值；否则忽略
* \param states    64 位整型数组，输入，实时整型数值列表，
*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的快照值；否则忽略
* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，写入实时数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
*        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
*        仅当输入时间戳与当前快照时间戳完全相等时，会替换当前快照的值和质量；
*        其余情况下会调用 rtdbs_put_snapshots()
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_fix_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float64* values, const rtdb_int64* states, const rtdb_int16* qualities, rtdb_error* errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_fix_snapshots64_fn)(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float64* values, const rtdb_int64* states, const rtdb_int16* qualities, rtdb_error* errors);
    rtdbs_fix_snapshots64_fn fn = (rtdbs_fix_snapshots64_fn)get_function("rtdbs_fix_snapshots64");
    return fn(handle, count, ids, datetimes, subtimes, values, states, qualities, errors);
}

/**
*
* \brief 批量回溯快照
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
*                    输出时表示成功写入实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表，同一个标签点标识可以出现多次，
*
* \param datetimes 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
* \param values    双精度浮点型数组，输入，实时浮点型数值列表，
*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的快照值；否则忽略
* \param states    64 位整型数组，输入，实时整型数值列表，
*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的快照值；否则忽略
* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，写入实时数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
*        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
* 功能说明：
*       批量将标签点的快照值vtmq改成传入的vtmq，如果传入的时间戳早于当前快照，会删除传入时间戳到当前快照的历史存储值。
*       如果传入的时间戳等于或者晚于当前快照，什么也不做。
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_back_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float64* values, const rtdb_int64* states, const rtdb_int16* qualities, rtdb_error* errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_back_snapshots64_fn)(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float64* values, const rtdb_int64* states, const rtdb_int16* qualities, rtdb_error* errors);
    rtdbs_back_snapshots64_fn fn = (rtdbs_back_snapshots64_fn)get_function("rtdbs_back_snapshots64");
    return fn(handle, count, ids, datetimes, subtimes, values, states, qualities, errors);
}

/**
*
* \brief 批量读取坐标实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、x、y、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表
* \param datetimes 整型数组，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
* \param x         单精度浮点型数组，输出，实时浮点型横坐标数值列表
* \param y         单精度浮点型数组，输出，实时浮点型纵坐标数值列表
* \param qualities 短整型数组，输出，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、x、y、qualities、errors 的长度与 count 一致。
*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_get_coor_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float32* x, rtdb_float32* y, rtdb_int16* qualities, rtdb_error* errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_get_coor_snapshots64_fn)(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float32* x, rtdb_float32* y, rtdb_int16* qualities, rtdb_error* errors);
    rtdbs_get_coor_snapshots64_fn fn = (rtdbs_get_coor_snapshots64_fn)get_function("rtdbs_get_coor_snapshots64");
    return fn(handle, count, ids, datetimes, subtimes, x, y, qualities, errors);
}

/**
*
* \brief 批量写入坐标实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、x、y、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表
* \param datetimes 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
* \param x         单精度浮点型数组，输入，实时浮点型横坐标数值列表
* \param y         单精度浮点型数组，输入，实时浮点型纵坐标数值列表
* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，写入实时坐标数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、x、y、qualities、errors 的长度与 count 一致。
*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_put_coor_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float32* x, const rtdb_float32* y, const rtdb_int16* qualities, rtdb_error* errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_put_coor_snapshots64_fn)(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float32* x, const rtdb_float32* y, const rtdb_int16* qualities, rtdb_error* errors);
    rtdbs_put_coor_snapshots64_fn fn = (rtdbs_put_coor_snapshots64_fn)get_function("rtdbs_put_coor_snapshots64");
    return fn(handle, count, ids, datetimes, subtimes, x, y, qualities, errors);
}

/**
*
* \brief 批量写入坐标实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、x、y、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表
* \param datetimes 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
* \param x         单精度浮点型数组，输入，实时浮点型横坐标数值列表
* \param y         单精度浮点型数组，输入，实时浮点型纵坐标数值列表
* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，写入实时坐标数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、x、y、qualities、errors 的长度与 count 一致。
*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
*        仅当输入时间戳与当前快照时间戳完全相等时，会替换当前快照的值和质量；
*        其余情况下会调用 rtdbs_put_coor_snapshots()
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_fix_coor_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float32* x, const rtdb_float32* y, const rtdb_int16* qualities, rtdb_error* errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_fix_coor_snapshots64_fn)(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float32* x, const rtdb_float32* y, const rtdb_int16* qualities, rtdb_error* errors);
    rtdbs_fix_coor_snapshots64_fn fn = (rtdbs_fix_coor_snapshots64_fn)get_function("rtdbs_fix_coor_snapshots64");
    return fn(handle, count, ids, datetimes, subtimes, x, y, qualities, errors);
}

/**
*
* \brief 读取二进制/字符串实时数据
*
* \param handle    连接句柄
* \param id        整型，输入，标签点标识
* \param datetime  整型，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
* \param blob      字节型数组，输出，实时二进制/字符串数值
* \param len       短整型，输出，二进制/字符串数值长度
* \param quality   短整型，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_get_blob_snapshot64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type* datetime, rtdb_subtime_type* subtime, rtdb_byte* blob, rtdb_length_type* len, rtdb_int16* quality)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_get_blob_snapshot64_fn)(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type* datetime, rtdb_subtime_type* subtime, rtdb_byte* blob, rtdb_length_type* len, rtdb_int16* quality);
    rtdbs_get_blob_snapshot64_fn fn = (rtdbs_get_blob_snapshot64_fn)get_function("rtdbs_get_blob_snapshot64");
    return fn(handle, id, datetime, subtime, blob, len, quality);
}

/**
*
* \brief 批量读取二进制/字符串实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、blobs、lens、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识
* \param datetimes 整型数组，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
* \param blobs     字节型指针数组，输出，实时二进制/字符串数值
* \param lens      短整型数组，输入/输出，二进制/字符串数值长度，
*                    输入时表示对应的 blobs 指针指向的缓冲区长度，
*                    输出时表示实际得到的 blob 长度，如果 blob 的长度大于缓冲区长度，会被截断。
* \param qualities 短整型数组，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_get_blob_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_byte* const* blobs, rtdb_length_type* lens, rtdb_int16* qualities, rtdb_error* errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_get_blob_snapshots64_fn)(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_byte* const* blobs, rtdb_length_type* lens, rtdb_int16* qualities, rtdb_error* errors);
    rtdbs_get_blob_snapshots64_fn fn = (rtdbs_get_blob_snapshots64_fn)get_function("rtdbs_get_blob_snapshots64");
    return fn(handle, count, ids, datetimes, subtimes, blobs, lens, qualities, errors);
}

/**
*
* \brief 写入二进制/字符串实时数据
*
* \param handle    连接句柄
* \param id        整型，输入，标签点标识
* \param datetime  整型，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
* \param blob      字节型数组，输入，实时二进制/字符串数值
* \param len       短整型，输入，二进制/字符串数值长度，超过一个页大小数据将被截断。
* \param quality   短整型，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_put_blob_snapshot64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime, const rtdb_byte* blob, rtdb_length_type len, rtdb_int16 quality)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_put_blob_snapshot64_fn)(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime, const rtdb_byte* blob, rtdb_length_type len, rtdb_int16 quality);
    rtdbs_put_blob_snapshot64_fn fn = (rtdbs_put_blob_snapshot64_fn)get_function("rtdbs_put_blob_snapshot64");
    return fn(handle, id, datetime, subtime, blob, len, quality);
}

/**
*
* \brief 批量写入二进制/字符串实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、blobs、lens、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识
* \param datetimes 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
* \param blobs     字节型指针数组，输入，实时二进制/字符串数值
* \param lens      短整型数组，输入，二进制/字符串数值长度，
*                    表示对应的 blobs 指针指向的缓冲区长度，超过一个页大小数据将被截断。
* \param qualities 短整型数组，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_put_blob_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_byte* const* blobs, const rtdb_length_type* lens, const rtdb_int16* qualities, rtdb_error* errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_put_blob_snapshots64_fn)(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_byte* const* blobs, const rtdb_length_type* lens, const rtdb_int16* qualities, rtdb_error* errors);
    rtdbs_put_blob_snapshots64_fn fn = (rtdbs_put_blob_snapshots64_fn)get_function("rtdbs_put_blob_snapshots64");
    return fn(handle, count, ids, datetimes, subtimes, blobs, lens, qualities, errors);
}

/**
*
* \brief 批量读取datetime类型标签点实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、dtvalues、dtlens、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识
* \param datetimes 整型数组，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
* \param dtvalues  字节型指针数组，输出，实时datetime数值
* \param dtlens    短整型数组，输入/输出，datetime数值长度，
*                    输入时表示对应的 dtvalues 指针指向的缓冲区长度，
* \param qualities 短整型数组，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \param type      短整型，输入，所有标签点的显示类型，如“yyyy-mm-dd hh:mm:ss.000”的type为1，默认类型1，
*                    “yyyy/mm/dd hh:mm:ss.000”的type为2
*                    如果不传type，则按照标签点属性显示，否则按照type类型显示
* \remark 本接口只对数据类型为 RTDB_DATETIME 的标签点有效。
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_get_datetime_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_byte* const* dtvalues, rtdb_length_type* dtlens, rtdb_int16* qualities, rtdb_error* errors, rtdb_int16 type)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_get_datetime_snapshots64_fn)(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_byte* const* dtvalues, rtdb_length_type* dtlens, rtdb_int16* qualities, rtdb_error* errors, rtdb_int16 type);
    rtdbs_get_datetime_snapshots64_fn fn = (rtdbs_get_datetime_snapshots64_fn)get_function("rtdbs_get_datetime_snapshots64");
    return fn(handle, count, ids, datetimes, subtimes, dtvalues, dtlens, qualities, errors, type);
}

/**
*
* \brief 批量插入datetime类型标签点数据
*
* \param handle      连接句柄
* \param count       整型，输入/输出，标签点个数，
*                      输入时表示 ids、datetimes、ms、dtvalues、dtlens、qualities、errors的长度，
*                      输出时表示成功写入的标签点个数
* \param ids         整型数组，输入，标签点标识
* \param datetimes   整型数组，输入，实时值时间列表
*                      表示距离1970年1月1日08:00:00的秒数
* \param ms          短整型数组，输入，实时数值时间列表，
*                      对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
* \param dtvalues    字节型指针数组，输入，datetime标签点的值
* \param dtlens      短整型数组，输入，数值长度
* \param qualities   短整型数组，输入，实时数值品质，，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors      无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \remark 被接口只对数据类型 RTDB_DATETIME 的标签点有效。
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_put_datetime_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_byte* const* dtvalues, const rtdb_length_type* dtlens, const rtdb_int16* qualities, rtdb_error* errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_put_datetime_snapshots64_fn)(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_byte* const* dtvalues, const rtdb_length_type* dtlens, const rtdb_int16* qualities, rtdb_error* errors);
    rtdbs_put_datetime_snapshots64_fn fn = (rtdbs_put_datetime_snapshots64_fn)get_function("rtdbs_put_datetime_snapshots64");
    return fn(handle, count, ids, datetimes, subtimes, dtvalues, dtlens, qualities, errors);
}

/**
*
* \brief 批量标签点快照改变的通知订阅
*
* \param handle         连接句柄
* \param count          整型，输入/输出，标签点个数，输入时表示 ids、errors 的长度，
*                           输出时表示成功订阅的标签点个数，不得超过 RTDB_MAX_SUBSCRIBE_SNAPSHOTS。
* \param ids            整型数组，输入，标签点标识列表。
* \param options        订阅选项
*                           RTDB_O_AUTOCONN 自动重连
* \param param          用户自定义参数
* \param callback       rtdbs_snaps_event_ex 类型回调接口，输入，当回掉函数返回非RtE_OK时退出订阅
*                       当未设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
*                       作为event_type取值调用回掉函数后退出订阅。
*                       当设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
*                       作为event_type取值调用回掉函数直到连接恢复或回掉函数返回非RtE_OK，
*                       网络中断期间回掉函数调用频率为最少3秒
*                       event_type参数值含义如下：
*                         RTDB_E_DATA        标签点快照改变
*                         RTDB_E_DISCONNECT  订阅客户端与数据库网络断开
*                         RTDB_E_RECOVERY    订阅客户端与数据库网络及订阅恢复
*                         RTDB_E_CHANGED     客户端修改订阅标签点信息，即通过rtdbs_change_subscribe_snapshots修改订阅信息的结果
*                       handle 产生订阅回掉的连接句柄，调用rtdbs_subscribe_snapshots_ex时的handle参数
*                       param  用户自定义参数，调用rtdbs_subscribe_snapshots_ex时的param参数
*                       count  event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示ids，datetimes,values等的数量
*                              event_type为其它值时，count值为0
*                       ids    event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示快照改变的标签点ID，数量由count指定
*                              event_type为其它值时，ids值为NULL
*                       datetimes event_type为RTDB_E_DATA时表示快照时间，数量由count指定
*                                 event_type为RTDB_E_CHANGED时表示changed_types，即通过rtdbs_change_subscribe_snapshots传入的changed_types
*                                 event_type为其它值时，datetimes值为NULL
*                       ms     event_type为RTDB_E_DATA时表示快照的毫秒，数量由count指定
*                              event_type为其它值时，ms值为NULL
*                       values event_type为RTDB_E_DATA时表示浮点数据类型快照值，数量由count指定
*                              event_type为RTDB_E_CHANGED时，表示delta_values，即通过rtdbs_change_subscribe_snapshots传入的delta_values
*                              event_type为其它值时，values值为NULL
*                       states event_type为RTDB_E_DATA时表示整形数据类型快照值，数量由count指定
*                              event_type为RTDB_E_CHANGED时，表示delta_states，通过rtdbs_change_subscribe_snapshots传入的delta_states
*                              event_type为其它值时，states值为NULL
*                       qualities event_type为RTDB_E_DATA时表示快照质量码，数量由count指定
*                              event_type为其它值时，qualities值为NULL
*                       errors event_type为RTDB_E_DATA时表示快照错误码，数量由count指定
*                              event_type为RTDB_E_CHANGED时，表示修改结果对应的错误码，数量由count指定
*                              event_type为其它值时，errors值为NULL
* \param errors         无符号整型数组，输出，
*                           写入实时数据的返回值列表，参考rtdb_error.h
* \remark   用户须保证 ids、errors 的长度与 count 一致。
*        用于订阅快照的连接句柄必需是独立的，不能再用来调用其它 api，
*        否则返回 RtE_OTHER_SDK_DOING 错误。
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_subscribe_snapshots_ex64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_uint32 options, void* param, rtdbs_snaps_event_ex64 callback, rtdb_error* errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_subscribe_snapshots_ex64_fn)(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_uint32 options, void* param, rtdbs_snaps_event_ex64 callback, rtdb_error* errors);
    rtdbs_subscribe_snapshots_ex64_fn fn = (rtdbs_subscribe_snapshots_ex64_fn)get_function("rtdbs_subscribe_snapshots_ex64");
    return fn(handle, count, ids, options, param, callback, errors);
}

/**
*
* \brief 批量标签点快照改变的通知订阅
*
* \param handle         连接句柄
* \param count          整型，输入/输出，标签点个数，输入时表示 ids、errors 的长度，
*                           输出时表示成功订阅的标签点个数，不得超过 RTDB_MAX_SUBSCRIBE_SNAPSHOTS。
* \param ids            整型数组，输入，标签点标识列表。
* \param delta_values   double型数组，输入，订阅浮点类型标签点的容差值，变化超过设置的容差值才会推送
* \param delta_values   整型数组，输入，订阅整型标签点的容差值，变化超过设置的容差值才会推送
* \param options        订阅选项
*                           RTDB_O_AUTOCONN 自动重连
* \param param          用户自定义参数
* \param callback       rtdbs_snaps_event_ex 类型回调接口，输入，当回掉函数返回非RtE_OK时退出订阅
*                         当未设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
*                         作为event_type取值调用回掉函数后退出订阅。
*                         当设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
*                         作为event_type取值调用回掉函数直到连接恢复或回掉函数返回非RtE_OK，
*                         网络中断期间回掉函数调用频率为最少3秒
*                         event_type参数值含义如下：
*                           RTDB_E_DATA        标签点快照改变
*                           RTDB_E_DISCONNECT  订阅客户端与数据库网络断开
*                           RTDB_E_RECOVERY    订阅客户端与数据库网络及订阅恢复
*                           RTDB_E_CHANGED     客户端修改订阅标签点信息，即通过rtdbs_change_subscribe_snapshots修改订阅信息的结果
*                         handle 产生订阅回掉的连接句柄，调用rtdbs_subscribe_snapshots_ex时的handle参数
*                         param  用户自定义参数，调用rtdbs_subscribe_snapshots_ex时的param参数
*                         count  event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示ids，datetimes,values等的数量
*                                event_type为其它值时，count值为0
*                         ids    event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示快照改变的标签点ID，数量由count指定
*                                event_type为其它值时，ids值为NULL
*                         datetimes event_type为RTDB_E_DATA时表示快照时间，数量由count指定
*                                   event_type为RTDB_E_CHANGED时表示changed_types，即通过rtdbs_change_subscribe_snapshots传入的changed_types
*                                   event_type为其它值时，datetimes值为NULL
*                         ms     event_type为RTDB_E_DATA时表示快照的毫秒，数量由count指定
*                                event_type为其它值时，ms值为NULL
*                         values event_type为RTDB_E_DATA时表示浮点数据类型快照值，数量由count指定
*                                event_type为RTDB_E_CHANGED时，表示delta_values，即通过rtdbs_change_subscribe_snapshots传入的delta_values
*                                event_type为其它值时，values值为NULL
*                         states event_type为RTDB_E_DATA时表示整形数据类型快照值，数量由count指定
*                                event_type为RTDB_E_CHANGED时，表示delta_states，通过rtdbs_change_subscribe_snapshots传入的delta_states
*                                event_type为其它值时，states值为NULL
*                         qualities event_type为RTDB_E_DATA时表示快照质量码，数量由count指定
*                                event_type为其它值时，qualities值为NULL
*                         errors event_type为RTDB_E_DATA时表示快照错误码，数量由count指定
*                                event_type为RTDB_E_CHANGED时，表示修改结果对应的错误码，数量由count指定
*                                event_type为其它值时，errors值为NULL
* \param errors         无符号整型数组，输出，
*                           写入实时数据的返回值列表，参考rtdb_error.h
* \remark delta_values和delta_states可以为空指针，表示不设置容差值。 只有两个参数都不为空时，设置容差值才会生效。
            用户须保证 ids、delta_values、delta_states、errors 的长度与 count 一致
*           用于订阅快照的连接句柄必需是独立的，不能再用来调用其它 api，
*           否则返回 RtE_OTHER_SDK_DOING 错误。
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_subscribe_delta_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_float64* delta_values, const rtdb_int64* delta_states, rtdb_uint32 options, void* param, rtdbs_snaps_event_ex64 callback, rtdb_error* errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_subscribe_delta_snapshots64_fn)(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_float64* delta_values, const rtdb_int64* delta_states, rtdb_uint32 options, void* param, rtdbs_snaps_event_ex64 callback, rtdb_error* errors);
    rtdbs_subscribe_delta_snapshots64_fn fn = (rtdbs_subscribe_delta_snapshots64_fn)get_function("rtdbs_subscribe_delta_snapshots64");
    return fn(handle, count, ids, delta_values, delta_states, options, param, callback, errors);
}

/**
*
* \brief 批量修改订阅标签点信息
*
* \param handle         连接句柄
* \param count          整型，输入/输出，标签点个数，输入时表示 ids、errors 的长度，
*                           输出时表示成功订阅的标签点个数，不得超过 RTDB_MAX_SUBSCRIBE_SNAPSHOTS。
* \param ids            整型数组，输入，标签点标识列表。
* \param delta_values   double型数组，输入，订阅浮点类型标签点的容差值，变化超过设置的容差值才会推送
* \param delta_values   整型数组，输入，订阅整型标签点的容差值，变化超过设置的容差值才会推送
* \param changed_types  整型数组，输入，修改类型，参考RTDB_SUBSCRIBE_CHANGE_TYPE
* \param errors         异步调用，保留参数，暂时不启用
* \remark   用户须保证 ids、delta_values、delta_states、errors 的长度与 count 一致。
*               可以同时添加、修改、删除订阅的标签点信息，
*               delta_values和delta_states，可以为空指针，为空，则表示不设置容差值，即写入新数据即推送
*               只有delta_values和delta_states都不为空时，设置的容差值才有效。
*               用于订阅快照的连接句柄必需是独立的，不能再用来调用其它 api，
*               否则返回 RtE_OTHER_SDK_DOING 错误。
*               此方法是异步方法，当网络中断等异常情况时，会通过方法的返回值返回错误，参考rtdb_error.h。
*               当方法返回值为RtE_OK时，表示已经成功发送给数据库，但是并没有等待修改结果。
*               数据库的修改结果，会异步通知给api的回调函数，通过rtdbs_snaps_event_ex的RTDB_E_CHANGED事件通知修改结果
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_change_subscribe_snapshots_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_float64* delta_values, const rtdb_int64* delta_states, const rtdb_int32* changed_types, rtdb_error* errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_change_subscribe_snapshots_fn)(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_float64* delta_values, const rtdb_int64* delta_states, const rtdb_int32* changed_types, rtdb_error* errors);
    rtdbs_change_subscribe_snapshots_fn fn = (rtdbs_change_subscribe_snapshots_fn)get_function("rtdbs_change_subscribe_snapshots");
    return fn(handle, count, ids, delta_values, delta_states, changed_types, errors);
}

/**
*
* \brief 取消标签点快照更改通知订阅
*
* \param handle    连接句柄
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_cancel_subscribe_snapshots_warp(rtdb_int32 handle)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_cancel_subscribe_snapshots_fn)(rtdb_int32 handle);
    rtdbs_cancel_subscribe_snapshots_fn fn = (rtdbs_cancel_subscribe_snapshots_fn)get_function("rtdbs_cancel_subscribe_snapshots");
    return fn(handle);
}

/**
* 命名：rtdbs_get_named_type_snapshot32
* 功能：获取自定义类型测点的单个快照
* 参数：
*        [handle]    连接句柄
*        [id]        整型，输入，标签点标识
*        [datetime]  整型，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
*        [object]    字节型数组，输出，实时自定义类型标签点的数值
*        [length]    短整型，输入/输出，自定义类型标签点的数值长度
*        [quality]   短整型，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_get_named_type_snapshot64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type* datetime, rtdb_subtime_type* subtime, void* object, rtdb_length_type* length, rtdb_int16* quality)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_get_named_type_snapshot64_fn)(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type* datetime, rtdb_subtime_type* subtime, void* object, rtdb_length_type* length, rtdb_int16* quality);
    rtdbs_get_named_type_snapshot64_fn fn = (rtdbs_get_named_type_snapshot64_fn)get_function("rtdbs_get_named_type_snapshot64");
    return fn(handle, id, datetime, subtime, object, length, quality);
}

/**
* 命名：rtdbs_get_named_type_snapshots32
* 功能：批量获取自定义类型测点的快照
* 参数：
*        [handle]    连接句柄
*        [count]     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、objects、lengths、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
*        [ids]       整型数组，输入，标签点标识
*        [datetimes] 整型数组，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型数组，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
*        [objects]   指针数组，输出，自定义类型标签点数值
*        [lengths]   短整型数组，输入/输出，自定义类型标签点数值长度，
*                    输入时表示对应的 objects 指针指向的缓冲区长度，
*                    输出时表示实际得到的 objects 长度，如果 objects 的长度大于缓冲区长度，会被截断。
*        [qualities] 短整型数组，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*        [errors]    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_get_named_type_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, void* const* objects, rtdb_length_type* lengths, rtdb_int16* qualities, rtdb_error* errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_get_named_type_snapshots64_fn)(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, void* const* objects, rtdb_length_type* lengths, rtdb_int16* qualities, rtdb_error* errors);
    rtdbs_get_named_type_snapshots64_fn fn = (rtdbs_get_named_type_snapshots64_fn)get_function("rtdbs_get_named_type_snapshots64");
    return fn(handle, count, ids, datetimes, subtimes, objects, lengths, qualities, errors);
}

/**
* 命名：rtdbs_put_named_type_snapshot32
* 功能：写入单个自定义类型标签点的快照
* 参数：
*        [handle]    连接句柄
*        [id]        整型，输入，标签点标识
*        [datetime]  整型，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
*        [object]    void类型数组，输入，自定义类型标签点数值
*        [length]    短整型，输入，自定义类型标签点数值长度，超过一个页大小数据将被截断。
*        [quality]   短整型，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_put_named_type_snapshot64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime, const void* object, rtdb_length_type length, rtdb_int16 quality)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_put_named_type_snapshot64_fn)(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime, const void* object, rtdb_length_type length, rtdb_int16 quality);
    rtdbs_put_named_type_snapshot64_fn fn = (rtdbs_put_named_type_snapshot64_fn)get_function("rtdbs_put_named_type_snapshot64");
    return fn(handle, id, datetime, subtime, object, length, quality);
}

/**
* 命名：rtdbs_put_named_type_snapshots32
* 功能：批量写入自定义类型标签点的快照
* 参数：
*        [handle]    连接句柄
*        [count]     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、objects、lengths、qualities、errors 的长度，
*                    输出时表示成功写入实时值的标签点个数
*        [ids]       整型数组，输入，标签点标识
*        [datetimes] 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
*        [objects]   void类型指针数组，输入，自定义类型标签点数值
*        [lengths]   短整型数组，输入，自定义类型标签点数值长度，
*                    表示对应的 objects 指针指向的缓冲区长度，超过一个页大小数据将被截断。
*        [qualities] 短整型数组，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*        [errors]    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
*/
rtdb_error RTDBAPI_CALLRULE rtdbs_put_named_type_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const void* const* objects, const rtdb_length_type* lengths, const rtdb_int16* qualities, rtdb_error* errors)
{
    typedef rtdb_error (RTDBAPI_CALLRULE *rtdbs_put_named_type_snapshots64_fn)(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const void* const* objects, const rtdb_length_type* lengths, const rtdb_int16* qualities, rtdb_error* errors);
    rtdbs_put_named_type_snapshots64_fn fn = (rtdbs_put_named_type_snapshots64_fn)get_function("rtdbs_put_named_type_snapshots64");
    return fn(handle, count, ids, datetimes, subtimes, objects, lengths, qualities, errors);
}

#ifdef __cplusplus
}
#endif

#endif // _C_API_H_
