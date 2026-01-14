#ifndef _GOFN_H_
#define _GOFN_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "rtdb.h"
#include "rtdbapi.h"
#include "rtdb_error.h"

extern rtdb_error goSubscribeTagsEx(rtdb_uint32 event_type, rtdb_int32 handle, void* param, rtdb_int32 count, rtdb_int32 *ids, rtdb_int32 what);

#ifdef __cplusplus
}
#endif

#endif // _C_API_H_
