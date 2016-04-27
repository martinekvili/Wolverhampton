#ifndef JOBOBJECT_H
#define JOBOBJECT_H

#include <stdexcept>
#include <Windows.h>

#include "unicodestring.h"
#include "utility.h"
#include "process.h"

class JobObject {
	HANDLE jobHandle;

public:
	JobObject(const JobObject&) = delete;
	JobObject& operator=(const JobObject&) = delete;

	JobObject(const char *jobNameAnsi);
	~JobObject();

	void setLimitInformation(int memorySizeinMB, int maxTimeInSec);

	HANDLE getJobHandle() {
		return jobHandle;
	}
};

#endif