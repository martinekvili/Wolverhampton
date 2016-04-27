#include "pipe.h"

void Pipe::freeResources() {
	CloseHandle(pipeReadHandle);

	if (!writeHandleClosed) {
		CloseHandle(pipeWriteHandle);
	}
}

Pipe::Pipe() : writeHandleClosed{ false } {
	SECURITY_ATTRIBUTES saAttr;
	saAttr.nLength = sizeof(SECURITY_ATTRIBUTES);
	saAttr.bInheritHandle = TRUE;
	saAttr.lpSecurityDescriptor = NULL;

	if (!CreatePipe(&pipeReadHandle, &pipeWriteHandle, &saAttr, 0)) {
		throw createException("Pipe creation failed", GetLastError());
	}

	if (!SetHandleInformation(pipeReadHandle, HANDLE_FLAG_INHERIT, 0)) {
		freeResources();
		throw createException("The read handle to the pipe is inherited", GetLastError());
	}
}

Pipe::~Pipe() {
	freeResources();
}

void Pipe::freeWriteHandle() {
	CloseHandle(pipeWriteHandle);
	writeHandleClosed = true;
}
