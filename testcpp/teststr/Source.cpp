#include <iostream>
#include <string>
#include "string.h"


int main()
{
	std::string s = "我是中国人";
	std::cout << s.length() << std::endl;
	std::cout << strlen(s.c_str()) << std::endl;
}