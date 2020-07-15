package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

//定义配置文件解析后的结构
type BwallConfig struct {
	//当前展示的句子,注，首字母要大写，否则无法访问
	CurrentText   string //记录当前句的第一句话
	MottoFileName string //名言文本名称
	Interval      int    //间隔时间
}

//读取配置文件
func Config() *BwallConfig {

	dir, err := GetCurrentPath()
	if err != nil {
		log.Fatal(err)
	}

	JsonParse := NewJsonStruct()
	v := BwallConfig{}
	//下面使用的是相对路径，config.json文件和main.go文件处于同一目录下
	JsonParse.Load(dir+"config.json", &v)
	return &v
}

func SaveConfig(conf *BwallConfig) {
	dir, err := GetCurrentPath()
	if err != nil {
		log.Fatal(err)
	}

	JsonParse := NewJsonStruct()
	JsonParse.Save(dir+"config.json", &conf)
}

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

func (jst *JsonStruct) Load(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

//保存文件
func (jst *JsonStruct) Save(filename string, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println("error:", err)
	}
	//os.Stdout.Write(b)

	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	//func (f *File) WriteAt(b []byte, off int64) (n int, err error)
	_, err = f.WriteAt(b, 0)
	if err != nil {
		panic(err)
	}
	//fmt.Println(length) //8
}
