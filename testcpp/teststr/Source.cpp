#include <iostream>
#include <string>
#include "string.h"


int main()
{
	std::string s = "�����й���";
	std::cout << s.length() << std::endl;
	std::cout << strlen(s.c_str()) << std::endl;
}