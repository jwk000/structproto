#include "c2s.struct.h"
#include "StructProto.h"

int main()
{

	s2c_all_player_info msg;
	msg.infos.push_back(s2c_player_info{ 1,"aaa",offline });
	msg.infos.push_back(s2c_player_info{ 2,"bbb",offline });
	msg.infos.push_back(s2c_player_info{ 3,"ccc",offline });
	msg.infos.push_back(s2c_player_info{ 4,"ddd",offline });

	char buff[1024];
	int len = StructProto::Serialize(msg, buff, 1024);
	printf("len = %d\n", len);

	s2c_all_player_info xx = StructProto::Deserialize<s2c_all_player_info>(buff, len);
	printf("xx.infos[3].playername: %s\n", xx.infos[3].playername.c_str());

}