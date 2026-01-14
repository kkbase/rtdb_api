/*
    这里的C函数声明，均为Go进行实现，主要是C中的回调函数
*/

#ifndef __GOFN_H__
#define __GOFN_H__

#ifdef __cplusplus
extern "C" {
#endif

#include "rtdb.h"
#include "rtdbapi.h"
#include "rtdb_error.h"

// extern rtdb_error goSubscribeTagsEx(rtdb_uint32 event_type, rtdb_int32 handle, void* param, rtdb_int32 count, rtdb_int32 *ids, rtdb_int32 what);

#ifdef __cplusplus
}
#endif

#endif // __GOFN_H__
