#include <physfs.h>

#include "_cgo_export.h"

void getCdRomDirsCallback(void *d)
{
	PHYSFS_getCdRomDirsCallback(&wrapStringCallback, d);
}

void getSearchPathCallback(void *d)
{
	PHYSFS_getSearchPathCallback(&wrapStringCallback, d);
}

void enumerateFilesCallback(char *dir, void *d)
{
	PHYSFS_enumerateFilesCallback(dir, &wrapEnumFilesCallback, d);
}
