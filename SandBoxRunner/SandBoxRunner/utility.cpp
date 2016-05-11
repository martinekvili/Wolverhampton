#include "utility.h"

std::runtime_error createException(const char *errorPlace, DWORD errorNum)
{
	std::stringstream stringStream;
	stringStream << errorPlace << " (" << errorNum << ").";
	return std::runtime_error(stringStream.str());
}

ULARGE_INTEGER fileTimeToLargeInteger(FILETIME fileTime) {
	ULARGE_INTEGER integer;

	integer.HighPart = fileTime.dwHighDateTime;
	integer.LowPart = fileTime.dwLowDateTime;

	return integer;
}