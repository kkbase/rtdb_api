#ifndef _C_DYLIB_H_
#define _C_DYLIB_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#include <windows.h>
#define LOAD_LIBRARY(name) LoadLibraryA(name)
#define GET_FUNCTION GetProcAddress
#define CLOSE_LIBRARY FreeLibrary
typedef HMODULE LIBRARY_HANDLE;
#else
#include <dlfcn.h>
#define LOAD_LIBRARY(name) dlopen(name, RTLD_LAZY)
#define GET_FUNCTION dlsym
#define CLOSE_LIBRARY dlclose
typedef void* LIBRARY_HANDLE;
#endif

typedef struct {
    LIBRARY_HANDLE handle;
} DylibHandle;

DylibHandle load_library(const char *name) {
    DylibHandle h;
    h.handle = LOAD_LIBRARY(name);
    return h;
}

void* get_function(DylibHandle handle, const char *name) {
    return GET_FUNCTION(handle.handle, name);
}

int free_library(DylibHandle handle) {
    return CLOSE_LIBRARY(handle.handle);
}

#ifdef __cplusplus
}
#endif

#endif // _C_DYLIB_H_