#include <windows.h>

typedef struct _args_list {
    char* entry_name;
    char* params;
    int params_length;
    LPVOID log_ptr;
    size_t taskId;
    BOOL clear;
} args_list_t;


#ifdef __cplusplus
extern "C" {
#endif
    void cPrint(int logType , char* msg, int size);
    void initArgs(void * args);
#ifdef __cplusplus
}
#endif