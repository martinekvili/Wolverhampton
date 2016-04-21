#ifndef UTILITY
#define UTILITY

#include <sstream>
#include <stdexcept>
#include <Windows.h>

std::runtime_error createException(const char *errorPlace, DWORD errorNum);

#endif