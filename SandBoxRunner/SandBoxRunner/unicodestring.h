#ifndef UNICODESTRING_H
#define UNICODESTRING_H

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