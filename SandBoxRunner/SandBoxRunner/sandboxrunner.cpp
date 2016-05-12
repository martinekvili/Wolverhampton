#include "sandboxrunner.h"

SandBoxRunner::SandBoxRunner(int memorySizeinMB, int maxTimeInSec) : jobObject{ "SandBoxJob" } {
	maxTime = jobObject.setLimitInformation(memorySizeinMB, maxTimeInSec);

	JOBOBJECT_ASSOCIATE_COMPLETION_PORT Port;
	Port.CompletionKey = jobObject.getJobHandle();
	Port.CompletionPort = ioCompletionPort.getIOCompletionPortHandle();

	if (!SetInformationJobObject(jobObject.getJobHandle(), JobObjectAssociateCompletionPortInformation, &Port, sizeof(Port)))
	{
		throw createException("Could not associate job with IO completion port", GetLastError());
	}
}

bool SandBoxRunner::hasMoreTimeToRun(const Process& process) {
	FILETIME ftCreationTime;
	FILETIME ftExitTime;
	FILETIME ftKernelTime;
	FILETIME ftUserTime;

	if (!GetProcessTimes(process.getProcessHandle(), &ftCreationTime, &ftExitTime, &ftKernelTime, &ftUserTime)) {
		throw createException("Could not get information about the running process", GetLastError());
	}

	SYSTEMTIME stSystemTime;
	GetSystemTime(&stSystemTime);
	FILETIME ftSystemTime;
	SystemTimeToFileTime(&stSystemTime, &ftSystemTime);

	ULARGE_INTEGER creationTime = fileTimeToLargeInteger(ftCreationTime);
	ULARGE_INTEGER systemTime = fileTimeToLargeInteger(ftSystemTime);

	return (systemTime.QuadPart - creationTime.QuadPart < maxTime.QuadPart);
}

SandBoxRunner::RunResult SandBoxRunner::runProcessWithName(const char *processNameAnsi, const char *outFileNameAnsi) {
	Process process{ processNameAnsi };
	File outFile{ outFileNameAnsi };

	if (!AssignProcessToJobObject(jobObject.getJobHandle(), process.getProcessHandle())) {
		throw createException("Assigning Process to Job failed", GetLastError());
	}

	ResumeThread(process.getThreadHandle());

	// Wait until child process exits.
	while (WaitForSingleObject(process.getProcessHandle(), 100) == WAIT_TIMEOUT) {
		if (!hasMoreTimeToRun(process)) {
			TerminateJobObject(jobObject.getJobHandle(), 1);
			return NotEnoughTime;
		}

		process.writeStdOutToFile(outFile);
	}

	DWORD CompletionCode;
	ULONG_PTR CompletionKey;
	LPOVERLAPPED Overlapped;

	// First it is just the information about strating the process, we skip this
	GetQueuedCompletionStatus(ioCompletionPort.getIOCompletionPortHandle(), &CompletionCode, &CompletionKey, &Overlapped, INFINITE);

	GetQueuedCompletionStatus(ioCompletionPort.getIOCompletionPortHandle(), &CompletionCode, &CompletionKey, &Overlapped, INFINITE);
	switch (CompletionCode) {
	case JOB_OBJECT_MSG_END_OF_JOB_TIME:
		return NotEnoughTime;

	case JOB_OBJECT_MSG_JOB_MEMORY_LIMIT:
		return NotEnoughMemory;

	case JOB_OBJECT_MSG_EXIT_PROCESS:
		process.writeStdOutToFile(outFile);
		return Success;	

	default:
		return Unknown;
	}
}