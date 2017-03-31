using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace teststr
{
    class Program
    {
        static void Main(string[] args)
        {
            string s = "我是中国人";
            Console.WriteLine(s.Length);
            Console.WriteLine(Encoding.UTF8.GetBytes(s).Length);
        }
    }
}
