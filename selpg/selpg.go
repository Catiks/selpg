package main
import (
  "fmt"
  "os"
  "bufio"
  "io"
  "log"
  "strconv"
  //"io/ioutil"
)
func self_logger(myerr interface{}) {
	logfile, err := os.OpenFile("./ErrorLog", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(-1)
	}
	defer logfile.Close()
	logger := log.New(logfile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)
	logger.Println(myerr)
}
func printUsage() {
  fmt.Printf("Usage: ./selpg -s开始页号 -e结束页号 [ -l每页行号 | -f | 输入文件名]\n")
}
func main()  {
  defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			self_logger(err) // 这里的err其实就是panic传入的内容，55
		}
	}()

  var b = os.Args
  //length := len(b)
  var lnum int = 12//default line number
  var err error
  var inputFile *os.File = os.Stdin
  var outputFile *os.File = os.Stdout
  var infile, outfile bool = false, false
  var flag bool = false
  var snum int = -1
  var enum int = -1
  for i, a := range b[1:] {
    if a[0:2] == "-s" {
      snum,err = strconv.Atoi(a[2:])
      if (err != nil) {
        printUsage()
        panic(err)
      }
    } else if a[0:2] == "-e" {
      enum, err = strconv.Atoi(a[2:])
      if (err!= nil) {
        printUsage()
        panic(err)
      }
    } else if a[0:2] == "-l" {
      lnum, err = strconv.Atoi(a[2:])
      if (err!= nil) {
        printUsage()
        panic(err)
      }
    } else if a == "-f" {
      flag = true
    } else if infile == false {
      inputFile, err = os.OpenFile(a, os.O_RDONLY, 0)
      if err != nil {
        errstr:= fmt.Sprintf("Can not open file:%s", a)
        printUsage()
        panic(errstr)
      }
      infile = true
      defer inputFile.Close()
    } else if outfile == false {
      outputFile, err = os.OpenFile(a, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0)
      if err != nil {
        errstr:= fmt.Sprintf("Can not open file:%s", a)
        printUsage()
        panic(errstr)
      }
      outfile = true
      defer outputFile.Close()
    } else {
      fmt.Printf("Undefined option [%s]\n", a)
      errstr := fmt.Sprintf("Undefined option [%s]", a)
      printUsage()
      panic(errstr)
    }
    i++
  }
  if snum < 0||enum < 0 {
    printUsage()
    os.Exit(-1)
  }
  if (snum > enum) {
    errstr := fmt.Sprintf("start num(%d) is gratter than end number(%d)", snum, enum)
    panic(errstr)
  }
  if flag {
    lnum = 1
  }
  rd := bufio.NewReader(inputFile)
  wt := bufio.NewWriter(outputFile)
  var lc int = 0
  var pc int = 0
  for true {
      line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
      if err != nil || io.EOF == err {
          break
      }
      lc++
      if lc == lnum {
        lc = 0;
        pc++
      }
      if pc >= snum&&pc<=enum {
        _, err := wt.WriteString(line)
        if err != nil {
          panic(err)
        }
        wt.Flush()
      } else if pc > enum {
        break
      }
  }
  if pc < snum {
    errstr := fmt.Sprintf("start number (and end number) is grater than the whole pages in the file")
    panic(errstr)
  }
  if pc < enum {
    errstr := fmt.Sprintf("end number is grater than the whole pages in the file")
    panic(errstr)
  }
}
