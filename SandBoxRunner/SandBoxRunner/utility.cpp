#include "utility.h"

std::runtime_error createException(const char *errorPlace, DWORD errorNum)
{
	std::stringstream stringStream;
	stringStream << errorPlace << " (" << errorNum << ").";
	return std::runtime_error(stringStream.str());
}