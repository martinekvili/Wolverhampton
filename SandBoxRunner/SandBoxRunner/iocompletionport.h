#ifndef IOCOMPLETIONPORT
#define IOCOMPLETIONPORT

#include <Windows.h>

#include "utility.h"

class IOCompletionPort {
	HANDLE ioCompletionPortHandle;

public:
	IOCompletionPort(const IOCompletionPort&) = delete;
	IOCompletionPort& operator=(const IOCompletionPort&) = delete;

	IOCompletionPort();
	~IOCompletionPort();

	HANDLE getIOCompletionPortHandle() {
		return ioCompletionPortHandle;
	}
};

#endif