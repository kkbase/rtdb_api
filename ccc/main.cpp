#include<stdio.h>
#include"rtdb.h"
#include"rtdbapi.h"
#include"rtdb_error.h"

int main() {
	rtdb_int32 handle;
	rtdb_error err;
	rtdb_int32 priv;

	err = rtdb_connect("159.75.187.68", 6327, &handle);
	if (err != RtE_OK) {
		printf("创建连接报错：%d\n", err);
	}

	err = rtdb_login(handle, "sa", "golden", &priv);
	if (err != RtE_OK) {
		printf("登录报错：%d\n", err);
	}

	/*
	RTDB_TABLE table_info;
	table_info.name[0] = 'a';
	table_info.name[1] = 'b';
	table_info.name[2] = 'c';
	table_info.name[3] = '\0';
	err = rtdbb_append_table(handle, &table_info);
	if (err != RtE_OK) {
		printf("创建表失败：%d\n", err);
	}
	printf("table_id: %d\n", table_info.id);
	*/

	// table_id = 1;
	
	/*
	rtdb_int32 point_id;
	err = rtdbb_insert_base_point(handle, "aaa", RTDB_INT32, 1, 1, &point_id);
	if (err != RtE_OK) {
		printf("创建点失败：%d\n", err);
	}
	printf("point_id:%d\n", point_id);
	*/
	// point_id = 80
	
	rtdb_int32 count = 1;
	rtdb_int32 ids[1] = { 80 };
	rtdb_timestamp_type datetimes[1] = { 1769485603 };
	rtdb_subtime_type subtimes[1] = { 0 };
	rtdb_float64 values[1] = { 0 };
	rtdb_int64 states[1] = { 1 };
	rtdb_int16 qualities[1] = { 0 };
	rtdb_error errors[1] = { 0 };
	err = rtdbs_put_snapshots64(handle, &count, ids, datetimes, subtimes, values, states, qualities, errors);
	if (err != RtE_OK) {
		printf("写数值失败：%d\n", err);
	}
	if (errors[0] != RtE_OK) {
		printf("errors[0]:%u\n", errors[0]);
	}

	err = rtdb_disconnect(handle);
	if (err != RtE_OK) {
		printf("断开连接报错：%d\n", err);
	}
	return 0;
}
