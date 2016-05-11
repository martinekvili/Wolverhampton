#include "process.h"

const int Process::bufferSize = 4096;

Process::Process(const char* processNameAnsi) : pipe{} {
	STARTUPINFO si;
	PROCESS_INFORMATION pi;

	ZeroMemory(&si, sizeof(si));
	si.cb = sizeof(si);
	si.hStdError = pipe.getPipeWriteHandle();
	si.hStdOutput = pipe.getPipeWriteHandle();
	si.dwFlags |= STARTF_USESTDHANDLES;

	ZeroMemory(&pi, sizeof(pi));

	UnicodeString processName{ processNameAnsi };

	// Start the child process. 
	if (!CreateProcess(
		processName.getUnicodeString(),
		NULL,					// No command line
		NULL,					// Process handle not inheritable
		NULL,					// Thread handle not inheritable
		TRUE,					// Set handle inheritance to FALSE
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

void Process::writeStdOutToFile(const File& file) {
	DWORD dwRead, dwWritten;
	CHAR chBuf[bufferSize];
	BOOL bSuccess = FALSE;

	while (true) {
		PeekNamedPipe(pipe.getPipeReadHandle(), NULL, 0, NULL, &dwRead, NULL);
		if (dwRead == 0) {
			break;
		}

		bSuccess = ReadFile(pipe.getPipeReadHandle(), chBuf, bufferSize, &dwRead, NULL);
		if (!bSuccess || dwRead == 0) {
			break;
		}

		bSuccess = WriteFile(file.getFileHandle(), chBuf, dwRead, &dwWritten, NULL);
		if (!bSuccess) {
			break;
		}
	}
}