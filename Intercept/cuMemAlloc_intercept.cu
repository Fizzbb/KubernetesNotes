#include <dlfcn.h>
#include <string.h>
#include <iostream>
#include <cuda.h>

extern "C" {void *__libc_dlsym(void *map, const char *name);}
extern "C" {void *__libc_dlopen_mode(const char *name, int maddArgumentode);}

typedef void *(*fnDlsym)(void *, const char *);
static void *real_dlsym(void *handle, const char *symbol)
{
    static fnDlsym internal_dlsym = (fnDlsym)__libc_dlsym(__libc_dlopen_mode("libdl.so.2", RTLD_LAZY), "dlsym");
    return (*internal_dlsym)(handle, symbol);
}

static void *realFunctions;

CUresult cuMemAlloc (CUdeviceptr* dptr, size_t bytesize)
{
    std::cout << "@@@@==cuMemAlloc hooked====" << std::endl;
    if (realFunctions == NULL) realFunctions = real_dlsym(RTLD_NEXT, "cuMemAlloc");
    return  ((CUresult (*)(CUdeviceptr*, size_t))realFunctions)(dptr, bytesize);
}

void *dlsym(void *handle, const char *symbol)
{
    if (strcmp(symbol, "cuMemAlloc") == 0) {
        if(realFunctions == NULL) realFunctions = real_dlsym(handle, symbol);
        return (void*)(&cuMemAlloc);
    }

    return (real_dlsym(handle, symbol));
}
