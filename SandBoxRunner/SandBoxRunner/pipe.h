#ifndef PIPE_H
#define PIPE_H

#include <Windows.h>

#include "utility.h"

class Pipe {
	HANDLE pipeReadHandle;
	HANDLE pipeWriteHandle;

	bool writeHandleClosed;

	void freeResources();

public:
	Pipe(const Pipe&) = delete;
	Pipe& operator=(const Pipe&) = delete;

	Pipe();
	~Pipe();

	void freeWriteHandle();

	HANDLE getPipeReadHandle() const {
		return pipeReadHandle;
	}

	HANDLE getPipeWriteHandle() const {
		return pipeWriteHandle;
	}
};

#endif