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
RTDBAPI rtdb_error RTDBAPI_CALLRULE
rtdb_get_api_version_warp(
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
RTDBAPI rtdb_error RTDBAPI_CALLRULE
rtdb_set_option_warp(
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

#ifdef __cplusplus
}
#endif

#endif // _C_API_H_
