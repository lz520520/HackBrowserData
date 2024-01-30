#include "run.h"


void(*internalPrint)(int, int, char*, int);
size_t gTaskId = 0;


void initArgs(void * args) {
    if (args) {
        install_args_t * args_list = (install_args_t *) args;
        if (internalPrint == NULL) {
            internalPrint = (void(*)(int, int, char*, int))args_list->log_ptr;
        }
        gTaskId = args_list->task_id;
    }
}
void cPrint(int logType , char* msg, int size) {
    if (internalPrint) {
        internalPrint(gTaskId, logType, msg, size);
    }
}