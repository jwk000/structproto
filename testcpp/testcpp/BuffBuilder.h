#pragma once
#include <assert.h>
#include <stdio.h>
#include <memory.h>
#include <stdint.h>
#include <string>
#include "StructProto.h"

class BuffBuilder
{
public:
	BuffBuilder(void* buff, int maxlen) :m_buff((char*)buff), m_pushpos(0), m_maxlen(maxlen)
	{

	}

	BuffBuilder& reset()
	{
		m_pushpos = 0;
		return *this;
	}
	template<typename T>
	void Push(const T& obj)
	{
		assert(m_pushpos + sizeof(obj) <= m_maxlen);
		if (m_pushpos + sizeof(obj) > m_maxlen)
		{
			printf("buffshell写越界");
			return ;
		}
		memcpy(m_buff + m_pushpos, &obj, sizeof(obj));
		m_pushpos += sizeof(obj);
	}

	void PushInt8(int8_t n) { Push(n); }
	void PushUInt8(uint8_t n) { Push(n); }
	void PushInt16(int16_t n) { Push(n); }
	void PushUInt16(uint16_t n) { Push(n); }
	void PushInt32(int32_t n) { Push(n); }
	void PushUInt32(uint32_t n) { Push(n); }
	void PushFloat32(float n) { Push(n); }
	void PushFloat64(double n) { Push(n); }
	void PushString(std::string s)
	{
		uint16_t slen = s.length();
		Push(slen);

		if (s.length() == 0) return;
		assert(m_pushpos + slen <= m_maxlen);
		if (m_pushpos + slen > m_maxlen)
		{
			printf("buffshell写越界");
			return;
		}
		memcpy(m_buff + m_pushpos, s.c_str(), slen);
		m_pushpos += slen;
	}

	//push之后的size
	char* Data() { return m_buff; }
	char* Cursor() { return m_buff + m_pushpos; }
	int Size() { return m_pushpos; }
	int SizeToPush() { return m_maxlen - m_pushpos; }
	int MaxSize() { return m_maxlen; }
	void SkipSize(int len) { m_pushpos += len; }
private:
	char* m_buff = nullptr;
	int m_pushpos = 0;
	int m_maxlen = 0;

};