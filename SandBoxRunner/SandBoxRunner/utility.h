#ifndef UTILITY_H
#define UTILITY_H

#include <sstream>
#include <stdexcept>
#include <Windows.h>

std::runtime_error createException(const char *errorPlace, DWORD errorNum);

#endif