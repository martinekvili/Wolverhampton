#include "sandboxrunner.h"

SandBoxRunner::SandBoxRunner(int memorySizeinMB, int maxTimeInSec) : jobObject{ "SandBoxJob" } {
	jobObject.setLimitInformation(memorySizeinMB, maxTimeInSec);

	JOBOBJECT_ASSOCIATE_COMPLETION_PORT Port;
	Port.CompletionKey = jobObject.getJobHandle();
	Port.CompletionPort = ioCompletionPort.getIOCompletionPortHandle();

	if (!SetInformationJobObject(jobObject.getJobHandle(), JobObjectAssociateCompletionPortInformation, &Port, sizeof(Port)))
	{
		throw createException("Could not associate job with IO completion port", GetLastError());
	}
}

SandBoxRunner::RunResult SandBoxRunner::runProcessWithName(const char *processNameAnsi) {
	Process process{ processNameAnsi };

	if (!AssignProcessToJobObject(jobObject.getJobHandle(), process.getProcessHandle())) {
		throw createException("Assigning Process to Job failed", GetLastError());
	}

	ResumeThread(process.getThreadHanlde());

	// Wait until child process exits.
	WaitForSingleObject(process.getProcessHandle(), INFINITE);

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
		return Success;

	default:
		return Unknown;
	}
}