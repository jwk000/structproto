#include "c2s.struct.h"
#include "BuffBuilder.h"
#include "BuffParser.h"


int s2c_refresh_stone::serialize(void* buff, int len) const
{
	BuffBuilder bb(buff, len);
	uint16_t _member_count =  1 ;
	bb.Push(_member_count);
	uint16_t _member_code = 0;
	uint8_t _member_type=0;
	_member_code =  1 ;
	bb.Push(_member_code);
	_member_type =  24 ;
	bb.Push(_member_type);
	bb.PushUInt16((uint16_t) stones .size());
	for(int n=0; n <  stones .size(); n++){
		uint16_t& elemlen = *(uint16_t*)bb.Cursor();
		elemlen = 0;
		bb.Push(elemlen);
		elemlen =  stones [n].serialize(bb.Cursor(), bb.SizeToPush());
		bb.SkipSize(elemlen);
	}

	return bb.Size();
}
void s2c_refresh_stone::deserialize(void* buff, int len)
{
	BuffParser bp(buff, len);
	uint16_t _member_count = bp.PopUInt16();
	uint16_t _member_code = 0;
	uint8_t _member_type = 0;
	for(uint16_t i=0;i<_member_count;i++)
	{
		_member_code = bp.PopUInt16();
		_member_type = bp.PopUInt8();
		switch(_member_code){
			case  1 :
			{
				uint16_t arraylen = bp.PopUInt16();
				for (;arraylen>0;arraylen--){
					uint16_t elemlen = bp.PopUInt16();
					msg_stone element{};
					element.deserialize(bp.Cursor(), bp.SizeToPop());
					stones.push_back(element);
					bp.SkipSize(elemlen);
				}
			}
			break;
			default:
			bp.SkipType((eMemberType)_member_type);break;
		}
	}
}

int c2s_login::serialize(void* buff, int len) const
{
	BuffBuilder bb(buff, len);
	uint16_t _member_count =  1 ;
	bb.Push(_member_count);
	uint16_t _member_code = 0;
	uint8_t _member_type=0;
	_member_code =  1 ;
	bb.Push(_member_code);
	_member_type =  4 ;
	bb.Push(_member_type);
	bb.PushInt32(userid);
	return bb.Size();
}
void c2s_login::deserialize(void* buff, int len)
{
	BuffParser bp(buff, len);
	uint16_t _member_count = bp.PopUInt16();
	uint16_t _member_code = 0;
	uint8_t _member_type = 0;
	for(uint16_t i=0;i<_member_count;i++)
	{
		_member_code = bp.PopUInt16();
		_member_type = bp.PopUInt8();
		switch(_member_code){
			case  1 :
			userid = bp.PopInt32();break;
			default:
			bp.SkipType((eMemberType)_member_type);break;
		}
	}
}

int c2s_enter_room::serialize(void* buff, int len) const
{
	BuffBuilder bb(buff, len);
	uint16_t _member_count =  2 ;
	bb.Push(_member_count);
	uint16_t _member_code = 0;
	uint8_t _member_type=0;
	_member_code =  1 ;
	bb.Push(_member_code);
	_member_type =  4 ;
	bb.Push(_member_type);
	bb.PushInt32(userid);
	_member_code =  2 ;
	bb.Push(_member_code);
	_member_type =  4 ;
	bb.Push(_member_type);
	bb.PushInt32(roomtype);
	return bb.Size();
}
void c2s_enter_room::deserialize(void* buff, int len)
{
	BuffParser bp(buff, len);
	uint16_t _member_count = bp.PopUInt16();
	uint16_t _member_code = 0;
	uint8_t _member_type = 0;
	for(uint16_t i=0;i<_member_count;i++)
	{
		_member_code = bp.PopUInt16();
		_member_type = bp.PopUInt8();
		switch(_member_code){
			case  1 :
			userid = bp.PopInt32();break;
			case  2 :
			roomtype = bp.PopInt32();break;
			default:
			bp.SkipType((eMemberType)_member_type);break;
		}
	}
}

int c2s_get_player_info::serialize(void* buff, int len) const
{
	BuffBuilder bb(buff, len);
	uint16_t _member_count =  0 ;
	bb.Push(_member_count);
	uint16_t _member_code = 0;
	uint8_t _member_type=0;
	return bb.Size();
}
void c2s_get_player_info::deserialize(void* buff, int len)
{
	BuffParser bp(buff, len);
	uint16_t _member_count = bp.PopUInt16();
	uint16_t _member_code = 0;
	uint8_t _member_type = 0;
	for(uint16_t i=0;i<_member_count;i++)
	{
		_member_code = bp.PopUInt16();
		_member_type = bp.PopUInt8();
		switch(_member_code){
			default:
			bp.SkipType((eMemberType)_member_type);break;
		}
	}
}

int s2c_player_info::serialize(void* buff, int len) const
{
	BuffBuilder bb(buff, len);
	uint16_t _member_count =  3 ;
	bb.Push(_member_count);
	uint16_t _member_code = 0;
	uint8_t _member_type=0;
	_member_code =  1 ;
	bb.Push(_member_code);
	_member_type =  4 ;
	bb.Push(_member_type);
	bb.PushInt32(playerid);
	_member_code =  2 ;
	bb.Push(_member_code);
	_member_type =  11 ;
	bb.Push(_member_type);
	bb.PushString(playername);
	_member_code =  3 ;
	bb.Push(_member_code);
	_member_type =  10 ;
	bb.Push(_member_type);
	bb.PushInt32((int)playerstate);
	return bb.Size();
}
void s2c_player_info::deserialize(void* buff, int len)
{
	BuffParser bp(buff, len);
	uint16_t _member_count = bp.PopUInt16();
	uint16_t _member_code = 0;
	uint8_t _member_type = 0;
	for(uint16_t i=0;i<_member_count;i++)
	{
		_member_code = bp.PopUInt16();
		_member_type = bp.PopUInt8();
		switch(_member_code){
			case  1 :
			playerid = bp.PopInt32();break;
			case  2 :
			playername = bp.PopString();break;
			case  3 :
			playerstate = ( ePlayerState )bp.PopInt32();break;
			default:
			bp.SkipType((eMemberType)_member_type);break;
		}
	}
}

int s2c_all_player_info::serialize(void* buff, int len) const
{
	BuffBuilder bb(buff, len);
	uint16_t _member_count =  1 ;
	bb.Push(_member_count);
	uint16_t _member_code = 0;
	uint8_t _member_type=0;
	_member_code =  1 ;
	bb.Push(_member_code);
	_member_type =  24 ;
	bb.Push(_member_type);
	bb.PushUInt16((uint16_t) infos .size());
	for(int n=0; n <  infos .size(); n++){
		uint16_t& elemlen = *(uint16_t*)bb.Cursor();
		elemlen = 0;
		bb.Push(elemlen);
		elemlen =  infos [n].serialize(bb.Cursor(), bb.SizeToPush());
		bb.SkipSize(elemlen);
	}

	return bb.Size();
}
void s2c_all_player_info::deserialize(void* buff, int len)
{
	BuffParser bp(buff, len);
	uint16_t _member_count = bp.PopUInt16();
	uint16_t _member_code = 0;
	uint8_t _member_type = 0;
	for(uint16_t i=0;i<_member_count;i++)
	{
		_member_code = bp.PopUInt16();
		_member_type = bp.PopUInt8();
		switch(_member_code){
			case  1 :
			{
				uint16_t arraylen = bp.PopUInt16();
				for (;arraylen>0;arraylen--){
					uint16_t elemlen = bp.PopUInt16();
					s2c_player_info element{};
					element.deserialize(bp.Cursor(), bp.SizeToPop());
					infos.push_back(element);
					bp.SkipSize(elemlen);
				}
			}
			break;
			default:
			bp.SkipType((eMemberType)_member_type);break;
		}
	}
}

int msg_stone::serialize(void* buff, int len) const
{
	BuffBuilder bb(buff, len);
	uint16_t _member_count =  2 ;
	bb.Push(_member_count);
	uint16_t _member_code = 0;
	uint8_t _member_type=0;
	_member_code =  1 ;
	bb.Push(_member_code);
	_member_type =  4 ;
	bb.Push(_member_type);
	bb.PushInt32(keyid);
	_member_code =  2 ;
	bb.Push(_member_code);
	_member_type =  4 ;
	bb.Push(_member_type);
	bb.PushInt32(stoneid);
	return bb.Size();
}
void msg_stone::deserialize(void* buff, int len)
{
	BuffParser bp(buff, len);
	uint16_t _member_count = bp.PopUInt16();
	uint16_t _member_code = 0;
	uint8_t _member_type = 0;
	for(uint16_t i=0;i<_member_count;i++)
	{
		_member_code = bp.PopUInt16();
		_member_type = bp.PopUInt8();
		switch(_member_code){
			case  1 :
			keyid = bp.PopInt32();break;
			case  2 :
			stoneid = bp.PopInt32();break;
			default:
			bp.SkipType((eMemberType)_member_type);break;
		}
	}
}

