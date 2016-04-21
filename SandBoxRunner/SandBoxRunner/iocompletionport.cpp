#include "iocompletionport.h"

IOCompletionPort::IOCompletionPort() {
	ioCompletionPortHandle = CreateIoCompletionPort(INVALID_HANDLE_VALUE, NULL, 0, 1);
	if (!ioCompletionPortHandle) {
		throw createException("Creating I/O completion port failed", GetLastError());
	}
}

IOCompletionPort::~IOCompletionPort() {
	CloseHandle(ioCompletionPortHandle);
}