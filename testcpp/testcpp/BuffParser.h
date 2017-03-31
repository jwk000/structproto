#pragma once
#include <assert.h>
#include <stdio.h>
#include <memory.h>
#include <stdint.h>
#include <string>
#include "StructProto.h"

class BuffParser
{
public:
	BuffParser(void* buff, int maxlen) :m_buff((char*)buff), m_poppos(0), m_maxlen(maxlen)
	{

	}

	BuffParser& Reset()
	{
		m_poppos = 0;
		return *this;
	}

	template<typename T>
	T Pop()
	{
		assert(m_poppos + (int)sizeof(T) <= m_maxlen);
		if (m_poppos + (int)sizeof(T) > m_maxlen)
		{
			printf("pop读越界");
		}

		T obj{};
		memcpy(&obj, m_buff + m_poppos, sizeof(obj));
		m_poppos += sizeof(obj);
		return obj;
	}

	int8_t PopInt8() { return Pop<int8_t>(); }
	uint8_t PopUInt8() { return Pop<uint8_t>(); }
	int16_t PopInt16() { return Pop<int16_t>(); }
	uint16_t PopUInt16() { return Pop<uint16_t>(); }
	int32_t PopInt32() { return Pop<int32_t>(); }
	uint32_t PopUInt32() { return Pop<uint32_t>(); }
	float PopFloat32() { return Pop<float>(); }
	double PopFloat64() { return Pop<double>(); }

	std::string PopString()
	{
		uint16_t slen = PopUInt16();

		assert(m_poppos + slen <= m_maxlen);
		if (m_poppos + slen > m_maxlen)
		{
			printf("buffshell读越界");
		}

		std::string s = std::string(m_buff + m_poppos, m_buff + m_poppos + slen);
		m_poppos += slen;
		return s;

	}

	int SizeToPop() { return m_maxlen - m_poppos; }
	
	char* Data() { return m_buff; }
	
	char* Cursor() { return m_buff + m_poppos; }

	int MaxSize() { return m_maxlen; }

	void SkipSize(int x) { m_poppos += x; }

	void SkipType(uint8_t t)
	{
		switch (t)
		{
		case    eTypeInt8:
		case 	eTypeUint8:
			SkipSize(1); break;
		case 	eTypeInt16:
		case 	eTypeUint16:
			SkipSize(2); break;
		case 	eTypeInt32:
		case 	eTypeUint32:
		case 	eTypeFloat32:
			SkipSize(4); break;
		case 	eTypeInt64:
		case 	eTypeUint64:
		case 	eTypeFloat64:
			SkipSize(8); break;
		case 	eTypeString:
		case 	eTypeStruct:
			SkipSize(PopUInt16()); break;
		case    eTypeInt8Array:
		case 	eTypeUint8Array:
		{
			uint16_t c = PopUInt16();
			for (uint16_t i = 0; i < c; i++) {
				SkipSize(1);
			}
			break;
		}
		case 	eTypeInt16Array:
		case 	eTypeUint16Array:
		{
			uint16_t c = PopUInt16();
			for (uint16_t i = 0; i < c; i++) {
				SkipSize(2);
			}
			break;
		}
		case 	eTypeInt32Array:
		case 	eTypeUint32Array:
		case 	eTypeFloat32Array:
		{
			uint16_t c = PopUInt16();
			for (uint16_t i = 0; i < c; i++) {
				SkipSize(4);
			}
			break;
		}
		case 	eTypeInt64Array:
		case 	eTypeUint64Array:
		case 	eTypeFloat64Array:
		{
			uint16_t c = PopUInt16();
			for (uint16_t i = 0; i < c; i++) {
				SkipSize(8);
			}
			break;
		}
		case 	eTypeStringArray:
		case 	eTypeStructArray:
		{
			uint16_t c = PopUInt16();
			for (uint16_t i = 0; i < c; i++) {
				SkipType(eTypeString);
			}
			break;
		}
		default:
			assert(false);
			printf("发现了不支持的数据类型:%u", t);
		}
	}
private:
	char* m_buff = nullptr;
	int m_poppos = 0;
	int m_maxlen = 0;

};