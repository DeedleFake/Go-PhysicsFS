#include <physfs.h>

#include "_cgo_export.h"

void getCdRomDirsCallback(void *d)
{
	PHYSFS_getCdRomDirsCallback((PHYSFS_StringCallback)&wrapStringCallback, d);
}

void getSearchPathCallback(void *d)
{
	PHYSFS_getSearchPathCallback((PHYSFS_StringCallback)&wrapStringCallback, d);
}

void enumerateFilesCallback(char *dir, void *d)
{
	PHYSFS_enumerateFilesCallback(dir, (PHYSFS_EnumFilesCallback)&wrapEnumFilesCallback, d);
}
