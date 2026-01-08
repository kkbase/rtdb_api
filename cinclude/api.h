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

typedef rtdb_error (RTDBAPI_CALLRULE *rtdb_get_api_version_fn)(
    rtdb_int32 *major,
    rtdb_int32 *minor,
    rtdb_int32 *beta
);
RTDBAPI rtdb_error RTDBAPI_CALLRULE
rtdb_get_api_version_warp(
    rtdb_int32 *major,
    rtdb_int32 *minor,
    rtdb_int32 *beta
) {
    rtdb_get_api_version_fn fn = (rtdb_get_api_version_fn)get_function("rtdb_get_api_version");
    return fn(major, minor, beta);
}


#ifdef __cplusplus
}
#endif

#endif // _C_API_H_
