#include "file.h"

File::File(const char* fileNameAnsi) {
	UnicodeString fileName{ fileNameAnsi };

	fileHandle = CreateFile(
		fileName.getUnicodeString(),
		GENERIC_WRITE,
		0,
		NULL,
		CREATE_ALWAYS,
		FILE_ATTRIBUTE_NORMAL,
		NULL);

	if (fileHandle == INVALID_HANDLE_VALUE) {
		throw createException("Couldn't create output file", GetLastError());
	}
}

File::~File() {
	CloseHandle(fileHandle);
}