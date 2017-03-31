package main

import (
	"flag"
	"bufio"
	"os"
	"io"
	"strings"
	"path/filepath"
)

//protogen.exe -i c2s.struct -cpp ./ -cs ./
func main(){
	//解析参数
	infile := flag.String("i","","input file")
	cppdir := flag.String("cpp","./", "cpp out dir")
	csdir := flag.String("cs","./","csharp out dir")

	flag.Parse()

	//读文件
	if infile == nil {
		return
	}

	fi, err := os.Open(*infile)
	if err != nil{
		panic(err)
	}

	lex := NewLexer()

	buf := bufio.NewReader(fi)


	for{

		line, err := buf.ReadString('\n')
		if err != nil{
			if err == io.EOF{
				break
			}
			panic(err)
		}

		line = strings.TrimSpace(line)

		err = lex.ReadLine(line)
		if err!=nil{
			panic(err)
		}

	}


	_, file := filepath.Split(*infile)
	//s := strings.Split(file,".")
	s := file
	if cppdir != nil {
		hfile := *cppdir+ s+ ".h"
		cppfile := *cppdir+s+".cpp"
		lex.GenCppCode(hfile, cppfile)
	}

	if csdir != nil {
		csfile := *csdir+ s+ ".cs"
		lex.GenCsharpCode(csfile,s)
	}
}
