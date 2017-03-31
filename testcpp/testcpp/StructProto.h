#pragma once

enum eMemberType
{
	eTypeInt8,
	eTypeUint8,
	eTypeInt16,
	eTypeUint16,
	eTypeInt32,
	eTypeUint32,
	eTypeInt64,
	eTypeUint64,
	eTypeFloat32,
	eTypeFloat64,
	eTypeString,
	eTypeStruct,
	eTypeInt8Array,
	eTypeUint8Array,
	eTypeInt16Array,
	eTypeUint16Array,
	eTypeInt32Array,
	eTypeUint32Array,
	eTypeInt64Array,
	eTypeUint64Array,
	eTypeFloat32Array,
	eTypeFloat64Array,
	eTypeStringArray,
	eTypeStructArray,
	eTypeEnum,
	eTypeEnumArray,
};


class StructProto
{
public:
	template<typename T>
	static int Serialize(const T& obj, void* buff, int maxlen)
	{
		return obj.serialize(buff, maxlen);
	}

	template<typename T>
	static T Deserialize(void* buff, int maxlen)
	{
		T obj{};
		obj.deserialize(buff, maxlen);
		return obj;
	}
};