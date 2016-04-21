#ifndef SANDBOXRUNNER
#define SANDBOXRUNNER

#include "jobobject.h"
#include "iocompletionport.h"
#include "utility.h"
#include <iostream>

class SandBoxRunner {
	JobObject jobObject;
	IOCompletionPort ioCompletionPort;

public:
	enum RunResult {
		Success,
		NotEnoughTime,
		NotEnoughMemory,
		Unknown
	};

	SandBoxRunner(const SandBoxRunner&) = delete;
	SandBoxRunner& operator=(const SandBoxRunner&) = delete;

	SandBoxRunner(int memorySizeinMB, int maxTimeInSec);

	RunResult runProcessWithName(const char *processNameAnsi);
};

#endif