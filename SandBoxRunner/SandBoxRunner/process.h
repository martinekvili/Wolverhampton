#ifndef PROCESS_H
#define PROCESS_H

#include <Windows.h>

#include "pipe.h"
#include "file.h"
#include "unicodestring.h"
#include "utility.h"

class Process {
	static const int bufferSize;

	Pipe pipe;

	HANDLE processHandle;
	HANDLE threadHandle;

public:
	Process(const Process&) = delete;
	Process& operator=(const Process&) = delete;

	Process(const char* processNameAnsi);
	~Process();

	void writeStdOutToFile(const File& file);

	HANDLE getProcessHandle() {
		return processHandle;
	}

	HANDLE getThreadHandle() {
		return threadHandle;
	}
};

#endif PROCESS