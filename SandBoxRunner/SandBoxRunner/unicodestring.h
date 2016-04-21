#ifndef UNICODESTRING
#define UNICODESTRING

#include <Windows.h>

class UnicodeString {
	BSTR unicodeString;

public:
	UnicodeString(const UnicodeString&) = delete;
	UnicodeString& operator=(const UnicodeString&) = delete;

	UnicodeString(const char *ansiString);
	~UnicodeString();

	BSTR getUnicodeString() {
		return unicodeString;
	}
	
};

#endif