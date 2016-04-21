#include "unicodestring.h"

UnicodeString::UnicodeString(const char *ansiString) {
	int length = lstrlenA(ansiString);
	unicodeString = SysAllocStringLen(NULL, length);
	MultiByteToWideChar(CP_ACP, 0, ansiString, length, unicodeString, length);
}

UnicodeString::~UnicodeString() {
	SysFreeString(unicodeString);
}