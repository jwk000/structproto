using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using c2s;
using StructProtocol;

namespace testcs
{
    class Program
    {
        static void Main(string[] args)
        {

            s2c_all_player_info msg = new s2c_all_player_info();
            msg.infos.Add(new s2c_player_info{playerid= 1,playername = "aaa", playerstate = ePlayerState.offline });
            msg.infos.Add(new s2c_player_info{ playerid= 2,playername ="bbb", playerstate = ePlayerState.offline });
            msg.infos.Add(new s2c_player_info{ playerid= 3, playername = "ccc", playerstate = ePlayerState.offline });
            msg.infos.Add(new s2c_player_info{ playerid = 4, playername = "ddd", playerstate = ePlayerState.offline });

            byte[] buff = new byte[1024];
            int len = StructProto.Serialize(msg, buff);
            Console.WriteLine("len = {0}", len);

            s2c_all_player_info xx = StructProto.Deserialize<s2c_all_player_info>(buff);
            Console.WriteLine("xx.infos[3].playername: {0}", xx.infos[3].playername);

        }
    }
}
