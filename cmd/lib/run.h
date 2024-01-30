#include <windows.h>

typedef struct _install_cache {
    LPVOID memory;
    size_t mem_len;
    size_t rand_len;

    LPVOID module_base;
    HMODULE* lib_handles;
    size_t lib_count;
    LPVOID entry_ptr;
    BOOL gc;
} install_cache_t;


typedef struct _install_args {
    char* install_name;
    char* entry_name;
    LPVOID log_ptr;
    int task_id;
    install_cache_t* cache;
} install_args_t;


typedef struct _args_list {
    char* params;
    int params_length;
    int task_id;
} args_list_t;


#ifdef __cplusplus
extern "C" {
#endif
    void cPrint(int logType , char* msg, int size);
    void initArgs(void * args);
#ifdef __cplusplus
}
#endif