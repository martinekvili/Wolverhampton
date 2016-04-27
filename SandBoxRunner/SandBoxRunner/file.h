#ifndef FILE_H
#define FILE_H

#include <Windows.h>

#include "unicodestring.h"
#include "utility.h"

class File {
	HANDLE fileHandle;

public:
	File(const File&) = delete;
	File& operator=(const File&) = delete;

	File(const char* fileNameAnsi);
	~File();

	HANDLE getFileHandle() const {
		return fileHandle;
	}
};

#endif