#include <iostream>
#include <Windows.h>

#include "sandboxrunner.h"

int main(int argc, char** argv) {
	SetErrorMode(SEM_FAILCRITICALERRORS | SEM_NOGPFAULTERRORBOX);
	_set_abort_behavior(0, _WRITE_ABORT_MSG);

	try {
		if (argc != 5) {
			std::cout << "Usage: " << argv[0] << " [memSizeInMB] [timeInSec] [executable] [outFile]" << std::endl;
			return 1;
		}

		int memSize = atoi(argv[1]);
		int time = atoi(argv[2]);
		SandBoxRunner sandBoxRunner{ memSize, time };

		auto result = sandBoxRunner.runProcessWithName(argv[3], argv[4]);

		std::cout << std::endl;
		std::cout << "RESULT: ";
		switch (result) {
		case SandBoxRunner::NotEnoughMemory:
			std::cout << "NOT_ENOUGH_MEMORY" << std::endl;
			break;

		case SandBoxRunner::NotEnoughTime:
			std::cout << "NOT_ENOUGH_TIME" << std::endl;
			break;

		case SandBoxRunner::Success:
			std::cout << "SUCCESS" << std::endl;
			break;

		case SandBoxRunner::Unknown:
			std::cout << "UNKNOWN_ERROR" << std::endl;
			break;
		}
	}
	catch (std::runtime_error error) {
		std::cout << std::endl;
		std::cout << "ERROR: " <<  error.what() << std::endl;
		return 1;
	}

	return 0;
}