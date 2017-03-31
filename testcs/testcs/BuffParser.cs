using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace StructProtocol
{
    public class BuffParser
    {
        public BuffParser(byte[] buff)
        {
            mBuff = buff;
        }

        void SkipSize(int x) { mIndex += x; }

        public void SkipType(eMemberType t)
        {
            switch (t)
            {
                case eMemberType.eTypeInt8:
                case eMemberType.eTypeUint8:
                    SkipSize(1); break;
                case eMemberType.eTypeInt16:
                case eMemberType.eTypeUint16:
                    SkipSize(2); break;
                case eMemberType.eTypeInt32:
                case eMemberType.eTypeUint32:
                case eMemberType.eTypeFloat32:
                    SkipSize(4); break;
                case eMemberType.eTypeInt64:
                case eMemberType.eTypeUint64:
                case eMemberType.eTypeFloat64:
                    SkipSize(8); break;
                case eMemberType.eTypeString:
                case eMemberType.eTypeStruct:
                    SkipSize(PopUInt16()); break;
                case eMemberType.eTypeInt8Array:
                case eMemberType.eTypeUint8Array:
                    {
                        var c = PopUInt16();
                        SkipSize(c);
                        break;
                    }
                case eMemberType.eTypeInt16Array:
                case eMemberType.eTypeUint16Array:
                    {
                        var c = PopUInt16();
                        SkipSize(2 * c);
                        break;
                    }
                case eMemberType.eTypeInt32Array:
                case eMemberType.eTypeUint32Array:
                case eMemberType.eTypeFloat32Array:
                    {
                        var c = PopUInt16();
                        SkipSize(4 * c);
                        break;
                    }
                case eMemberType.eTypeInt64Array:
                case eMemberType.eTypeUint64Array:
                case eMemberType.eTypeFloat64Array:
                    {
                        var c = PopUInt16();
                        SkipSize(8 * c);
                        break;
                    }
                case eMemberType.eTypeStringArray:
                case eMemberType.eTypeStructArray:
                    {
                        var c = PopUInt16();
                        for (var i = 0; i < c; i++)
                        {
                            SkipType(eMemberType.eTypeString);
                        }
                        break;
                    }
                default:
                    Console.WriteLine("发现了不支持的数据类型", t);
                    break;
            }
        }

        public string PopString()
        {
            int len = PopUInt16();
            byte[] buf = PopBytes(len);
            return Encoding.UTF8.GetString(buf);
        }

        public bool PopBool()
        {
            bool b = BitConverter.ToBoolean(mBuff, mIndex);
            mIndex += sizeof(bool);
            return b;
        }

        public byte PopByte()
        {
            return mBuff[mIndex++];
        }

        public int PopInt()
        {
            int i = BitConverter.ToInt32(mBuff, mIndex);
            mIndex += sizeof(int);
            return i;
        }

        public short PopShort()
        {
            short s = BitConverter.ToInt16(mBuff, mIndex);
            mIndex += sizeof(short);
            return s;
        }

        public long PopLong()
        {
            long l = BitConverter.ToInt64(mBuff, mIndex);
            mIndex += sizeof(long);
            return l;
        }

        public byte PopUInt8()
        {
            return PopByte();
        }

        public sbyte PopInt8()
        {
            return (sbyte)PopByte();
        }
        public Int16 PopInt16()
        {
            return PopShort();
        }

        public Int32 PopInt32()
        {
            return PopInt();
        }

        public Int64 PopInt64()
        {
            return PopLong();
        }

        public uint PopUInt()
        {
            uint u = BitConverter.ToUInt32(mBuff, mIndex);
            mIndex += sizeof(uint);
            return u;
        }

        public ushort PopUShort()
        {
            ushort u = BitConverter.ToUInt16(mBuff, mIndex);
            mIndex += sizeof(ushort);
            return u;
        }

        public char PopChar()
        {
            char c = BitConverter.ToChar(mBuff, mIndex);
            mIndex += sizeof(char);
            return c;
        }

        public ulong PopULong()
        {
            ulong u = BitConverter.ToUInt64(mBuff, mIndex);
            mIndex += sizeof(ulong);
            return u;
        }

        public UInt16 PopUInt16()
        {
            return PopUShort();
        }

        public UInt32 PopUInt32()
        {
            return PopUInt();
        }

        public UInt64 PopUInt64()
        {
            return PopULong();
        }

        public float PopFloat()
        {
            float f = BitConverter.ToSingle(mBuff, mIndex);
            mIndex += sizeof(float);
            return f;
        }

        public float PopFloat32()
        {
            return PopFloat();
        }

        public double PopDouble()
        {
            double d = BitConverter.ToDouble(mBuff, mIndex);
            mIndex += sizeof(double);
            return d;
        }

        public double PopFloat64()
        {
            return PopDouble();
        }

        public byte[] PopBytes(int len)
        {
            byte[] ret = null;
            if (check_buff(len))
            {
                ret = new byte[len];
                Array.Copy(mBuff, mIndex, ret, 0, len);
                mIndex += len;
            }
            return ret;
        }


        private bool check_buff(int len)
        {
            if (mBuff == null) return false;
            if (mBuff.Length < len) return false;
            return true;
        }
        private byte[] mBuff;
        private int mIndex = 0;
    }
}
