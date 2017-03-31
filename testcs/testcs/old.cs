using System;
using System.Collections.Generic;
using StructProtocol;

namespace old
{

    enum ePlayerState
    {
        offline = 0,
        online = 1,
        inteam = 2,
        ingame = 3
    }

    class s2c_player_info : IStruct
    {
        public int playerid;
        public string playername;
        public ePlayerState playerstate;

        public int serialize(byte[] buff)
        {
            BuffBuilder bb = new BuffBuilder(buff);
            UInt16 _member_count = 3;
            bb.PushUInt16(_member_count);
            UInt16 _member_code = 0;
            byte _member_type = 0;
            _member_code = 1;
            bb.PushUInt16(_member_code);
            _member_type = 4;
            bb.PushByte(_member_type);
            bb.PushInt32(playerid);
            _member_code = 2;
            bb.PushUInt16(_member_code);
            _member_type = 10;
            bb.PushByte(_member_type);
            bb.PushString(playername);
            _member_code = 3;
            bb.PushUInt16(_member_code);
            _member_type = 24;
            bb.PushByte(_member_type);
            bb.PushInt32((int)playerstate);
            return bb.Size();
        }
        public void deserialize(byte[] buff)
        {
            BuffParser bp = new BuffParser(buff);
            UInt16 _member_count = bp.PopUInt16();
            UInt16 _member_code = 0;
            byte _member_type = 0;
            for (UInt16 i = 0; i < _member_count; i++)
            {
                _member_code = bp.PopUInt16();
                _member_type = bp.PopByte();
                switch (_member_code)
                {
                    case 1:
                        playerid = bp.PopInt32(); break;
                    case 2:
                        playername = bp.PopString(); break;
                    case 3:
                        playerstate = (ePlayerState)bp.PopInt32(); break;
                    default:
                        bp.SkipType((eMemberType)_member_type); break;
                }
            }

        }
    }

    class s2c_all_player_info : IStruct
    {
        public List<s2c_player_info> infos = new List<s2c_player_info>();

        public int serialize(byte[] buff)
        {
            BuffBuilder bb = new BuffBuilder(buff);
            ushort _member_count = 1;
            bb.PushUInt16(_member_count);
            ushort _member_code = 0;
            byte _member_type = 0;
            _member_code = 1;
            bb.PushUInt16(_member_code);
            _member_type = 23;
            bb.PushByte(_member_type);
            bb.PushUInt16((ushort)infos.Count);
            for (int n = 0; n < infos.Count; n++)
            {
                byte[] _buf = new byte[4096];
                ushort elemlen = (ushort)infos[n].serialize(_buf);

                bb.PushUInt16(elemlen);
                bb.PushRangeBuff(_buf, elemlen);
            }

            return bb.Size();
        }
        public void deserialize(byte[] buff)
        {
            BuffParser bp = new BuffParser(buff);
            ushort _member_count = bp.PopUInt16();
            ushort _member_code = 0;
            byte _member_type = 0;
            for (ushort i = 0; i < _member_count; i++)
            {
                _member_code = bp.PopUInt16();
                _member_type = bp.PopByte();
                switch (_member_code)
                {
                    case 1:
                        {
                            ushort arraylen = bp.PopUInt16();
                            for (; arraylen > 0; arraylen--)
                            {
                                ushort elemlen = bp.PopUInt16();
                                byte[] _buf = bp.PopBytes(elemlen);
                                s2c_player_info element = new s2c_player_info();
                                element.deserialize(_buf);
                                infos.Add(element);
                            }
                        }
                        break;
                    default:
                        bp.SkipType((eMemberType)_member_type); break;
                }
            }

        }

    }

}