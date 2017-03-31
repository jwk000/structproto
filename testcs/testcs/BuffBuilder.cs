using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace StructProtocol
{

    class BuffBuilder
    {
        public BuffBuilder(byte[] buff)
        {
            mBuff = buff;
        }

        public void Destroy()
        {
            mBuff = null;
            mIndex = 0;
        }

        public void Reset()
        {
            mIndex = 0;
        }

        public int Size()
        {
            return mIndex;
        }
        public byte[] GetBuff()
        {
            if (mIndex == mBuff.Length)
                return mBuff;

            byte[] ret = new byte[mIndex];
            Array.Copy(mBuff, ret, mIndex);
            return ret;
        }

        //不带结尾0
        public BuffBuilder PushString(string s)
        {
            if (string.IsNullOrEmpty(s))
                return this;

            byte[] pAr = Encoding.UTF8.GetBytes(s);
            PushUInt16((ushort)pAr.Length);
            check_buff(pAr.Length);
            Array.Copy(pAr, 0, mBuff, mIndex, pAr.Length);
            mIndex += pAr.Length;
            return this;

        }

        public BuffBuilder PushBool(bool b)
        {
            check_buff(sizeof(bool));
            Array.Copy(BitConverter.GetBytes(b), 0, mBuff, mIndex, sizeof(bool));
            mIndex += sizeof(bool);
            return this;
        }

        public BuffBuilder PushChar(char c)
        {
            check_buff(sizeof(char));
            Array.Copy(BitConverter.GetBytes(c), 0, mBuff, mIndex, sizeof(char));
            mIndex += sizeof(char);
            return this;
        }

        public BuffBuilder PushFloat(float f)
        {
            check_buff(sizeof(float));
            Array.Copy(BitConverter.GetBytes(f), 0, mBuff, mIndex, sizeof(float));
            mIndex += sizeof(float);
            return this;
        }
        public BuffBuilder PushFloat32(float f)
        {
            return PushFloat(f);
        }
        public BuffBuilder PushDouble(double d)
        {
            check_buff(sizeof(double));
            Array.Copy(BitConverter.GetBytes(d), 0, mBuff, mIndex, sizeof(double));
            mIndex += sizeof(double);
            return this;
        }

        public BuffBuilder PushFloat64(double d)
        {
            return PushDouble(d);
        }
        public BuffBuilder PushByte(byte b)
        {
            check_buff(1);
            mBuff[mIndex++] = b;
            return this;
        }

        public BuffBuilder PushSByte(sbyte b)
        {
            check_buff(1);
            mBuff[mIndex++] = (byte)b;
            return this;
        }

        public BuffBuilder PushInt8(sbyte b)
        {
            return PushSByte(b);
        }

        public BuffBuilder PushUInt8(byte b)
        {
            return PushByte(b);
        }

        public BuffBuilder PushUInt16(UInt16 u)
        {
            check_buff(sizeof(UInt16));
            Array.Copy(BitConverter.GetBytes(u), 0, mBuff, mIndex, sizeof(UInt16));
            mIndex += sizeof(UInt16);
            return this;
        }

        public BuffBuilder PushUInt32(UInt32 u)
        {
            check_buff(sizeof(UInt32));
            Array.Copy(BitConverter.GetBytes(u), 0, mBuff, mIndex, sizeof(UInt32));
            mIndex += sizeof(UInt32);
            return this;
        }

        public BuffBuilder PushUInt64(UInt64 u)
        {
            check_buff(sizeof(UInt64));
            Array.Copy(BitConverter.GetBytes(u), 0, mBuff, mIndex, sizeof(UInt64));
            mIndex += sizeof(UInt64);
            return this;
        }

        public BuffBuilder PushUInt(uint u)
        {
            return PushUInt32(u);
        }

        public BuffBuilder PushUShort(ushort u)
        {
            return PushUInt16(u);
        }

        public BuffBuilder PushULong(ulong u)
        {
            return PushUInt64(u);
        }

        public BuffBuilder PushShort(short s)
        {
            check_buff(sizeof(short));
            Array.Copy(BitConverter.GetBytes(s), 0, mBuff, mIndex, sizeof(short));
            mIndex += sizeof(short);
            return this;
        }

        public BuffBuilder PushInt(int i)
        {
            check_buff(sizeof(int));
            Array.Copy(BitConverter.GetBytes(i), 0, mBuff, mIndex, sizeof(int));
            mIndex += sizeof(int);
            return this;
        }

        public BuffBuilder PushLong(long l)
        {
            check_buff(sizeof(long));
            Array.Copy(BitConverter.GetBytes(l), 0, mBuff, mIndex, sizeof(long));
            mIndex += sizeof(long);
            return this;
        }

        public BuffBuilder PushInt16(Int16 i)
        {
            return PushShort(i);
        }

        public BuffBuilder PushInt32(Int32 i)
        {
            return PushInt(i);
        }

        public BuffBuilder PushInt64(Int64 i)
        {
            return PushLong(i);
        }

        //不含长度
        public BuffBuilder PushBytes(byte[] buff)
        {
            if (buff == null) return this;
            check_buff(buff.Length);
            Array.Copy(buff, 0, mBuff, mIndex, buff.Length);
            mIndex += buff.Length;
            return this;
        }

        public BuffBuilder PushRangeBuff(byte[] buff, int len)
        {
            if (buff == null || len == 0) return this;
            check_buff(len);
            Array.Copy(buff, 0, mBuff, mIndex, len);
            mIndex += len;
            return this;
        }

        private void check_buff(int len)
        {
            if (null == mBuff)
                mBuff = new byte[len];
            if (mBuff.Length - mIndex < len)
            {
                byte[] t = new byte[mBuff.Length + len - (mBuff.Length - mIndex) + 1];
                Array.Copy(mBuff, t, mBuff.Length);
                mBuff = t;
            }
        }
        private byte[] mBuff;
        private int mIndex = 0;
    }

}
