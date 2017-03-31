using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace StructProtocol
{
    public enum eMemberType
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

    interface IStruct
    {
        int serialize(byte[] buff);
        void deserialize(byte[] buff);
    }

    class StructProto
    {
        public static int Serialize<T>(T obj, byte[] buff) where T : IStruct
        {
            return obj.serialize(buff);
        }

        public static T Deserialize<T>(byte[] buff) where T : IStruct, new()
        {
            T obj = new T();
            obj.deserialize(buff);
            return obj;
        }
    };
}
