#ifndef PROCESS
#define PROCESS

#include <Windows.h>

#include "unicodestring.h"
#include "utility.h"

class Process {
	HANDLE processHandle;
	HANDLE threadHandle;

public:
	Process(const Process&) = delete;
	Process& operator=(const Process&) = delete;

	Process(const char* processNameAnsi);
	~Process();

	HANDLE getProcessHandle() {
		return processHandle;
	}

	HANDLE getThreadHanlde() {
		return threadHandle;
	}
};

#endif PROCESS