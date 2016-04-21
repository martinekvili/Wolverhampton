#include "process.h"

Process::Process(const char* processNameAnsi) {
	STARTUPINFO si;
	PROCESS_INFORMATION pi;

	ZeroMemory(&si, sizeof(si));
	si.cb = sizeof(si);
	ZeroMemory(&pi, sizeof(pi));

	UnicodeString processName{ processNameAnsi };

	// Start the child process. 
	if (!CreateProcess(
		processName.getUnicodeString(),
		NULL,					// No command line
		NULL,					// Process handle not inheritable
		NULL,					// Thread handle not inheritable
		FALSE,					// Set handle inheritance to FALSE
		CREATE_SUSPENDED,		// Create suspended, we will start it after assigning to Job
		NULL,					// Use parent's environment block
		NULL,					// Use parent's starting directory 
		&si,					// Pointer to STARTUPINFO structure
		&pi)					// Pointer to PROCESS_INFORMATION structure
		)
	{
		throw createException("Creating Process failed", GetLastError());
	}

	processHandle = pi.hProcess;
	threadHandle = pi.hThread;
}

Process::~Process() {
	// Close process and thread handles. 
	CloseHandle(processHandle);
	CloseHandle(threadHandle);
}