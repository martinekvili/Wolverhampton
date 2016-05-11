#include "jobobject.h"

JobObject::JobObject(const char *jobNameAnsi) {
	UnicodeString jobName{ jobNameAnsi };

	jobHandle = CreateJobObject(NULL, jobName.getUnicodeString());
	if (!jobHandle) {
		throw createException("Creating Job failed", GetLastError());
	}
}

JobObject::~JobObject() {
	// Close job handle
	CloseHandle(jobHandle);
}

LARGE_INTEGER JobObject::setLimitInformation(int memorySizeinMB, int maxTimeInSec) {
	JOBOBJECT_BASIC_LIMIT_INFORMATION jobBasicLimitInfo;

	LARGE_INTEGER maxTime;
	maxTime.QuadPart = (long)10000000 * maxTimeInSec;			// The time the job can use, in 100ns ticks
	jobBasicLimitInfo.PerJobUserTimeLimit = maxTime;

	jobBasicLimitInfo.LimitFlags = JOB_OBJECT_LIMIT_JOB_TIME | JOB_OBJECT_LIMIT_JOB_MEMORY;

	JOBOBJECT_EXTENDED_LIMIT_INFORMATION jobLimitInfo;
	jobLimitInfo.BasicLimitInformation = jobBasicLimitInfo;
	jobLimitInfo.JobMemoryLimit = memorySizeinMB * 1024 * 1024;	// The memory the job can use, in bytes

	if (!SetInformationJobObject(jobHandle, JobObjectExtendedLimitInformation, &jobLimitInfo, sizeof(jobLimitInfo))) {
		throw createException("Setting limit information failed", GetLastError());
	}

	return maxTime;
}