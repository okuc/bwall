//go:generate goversioninfo -icon=./icon.ico -manifest=./main.manifest
//go generate
//go build
package main

import (
	"container/list"
	"github.com/lxn/walk"
	"log"
	"time"
)

type MyWindow struct {
	*walk.MainWindow
	ni *walk.NotifyIcon
}

func NewMyWindow() *MyWindow {
	mw := new(MyWindow)
	var err error
	mw.MainWindow, err = walk.NewMainWindow()
	checkError(err)
	return mw
}

func (mw *MyWindow) init() {
	//读取配置文件
	var conf = Config()
	//读取解析名人名言
	var mottoList = list.New()
	GetMotto(conf.MottoFileName, mottoList)
	//setNextBwall(mottoList, conf, mychan)

	// 指定的时间后执行一次
	time.AfterFunc(time.Duration(conf.Interval)*time.Minute,
		func() {
			go func() { //协程函数
				for {   //死循环，
					setNextBwall(mottoList, conf)
					//tick :=time.NewTicker(time.Duration(conf.Interval) * time.Minute)
					time.Sleep(time.Duration(conf.Interval) * time.Minute)
				}
			}()
		})
}

func (mw *MyWindow) AddNotifyIcon() {
	var err error
	mw.ni, err = walk.NewNotifyIcon(mw)
	checkError(err)
	mw.ni.SetVisible(true)
 	//生成syso资源文件时，会把icon文件包含在内，生成id为0-3之间的一个，可以在这4个数间尝试一下。肯定有一个会是对的。
	icon, err := walk.NewIconFromResourceId(2)
	checkError(err)
	mw.SetIcon(icon)
	mw.ni.SetIcon(icon)

	aboutAction := mw.addAction(nil, "关于")
	aboutAction.Triggered().Attach(func() {

		//aboutAction.SetChecked(true)
		aboutAction.SetEnabled(true)
		mw.msgbox("关于", "科技改变未来，谢谢！ \r\n              okuc 开发", walk.MsgBoxIconInformation)
	})


	helpMenu := mw.addMenu("帮助")
	mw.addAction(helpMenu, "软件说明").Triggered().Attach(func() {
		walk.MsgBox(mw, "help", "https://gitee.com/okuc/bwall", walk.MsgBoxIconInformation)
	})

	mw.addAction(helpMenu, "github地址").Triggered().Attach(func() {
		walk.MsgBox(mw, "about", "https://github.com/okuc/bwall", walk.MsgBoxIconInformation)
	})

	mw.addAction(nil, "退出").Triggered().Attach(func() {
		mw.ni.Dispose()
		mw.Dispose()
		walk.App().Exit(0)
	})

	aboutAction.SetEnabled(true)
}

func (mw *MyWindow) addMenu(name string) *walk.Menu {
	helpMenu, err := walk.NewMenu()
	checkError(err)
	help, err := mw.ni.ContextMenu().Actions().AddMenu(helpMenu)
	checkError(err)
	help.SetText(name)

	return helpMenu
}

func (mw *MyWindow) addAction(menu *walk.Menu, name string) *walk.Action {
	action := walk.NewAction()
	action.SetText(name)
	if menu != nil {
		menu.Actions().Add(action)
	} else {
		mw.ni.ContextMenu().Actions().Add(action)
	}

	return action
}

func (mw *MyWindow) msgbox(title, message string, style walk.MsgBoxStyle) {

	walk.MsgBox(mw, title, message, style)
}

//系统消息中心提示
func (mw *MyWindow) showInfo(title, message string, style walk.MsgBoxStyle) {
	mw.ni.ShowInfo(title, message)

}

func main() {
	mw := NewMyWindow()

	mw.init()
	mw.AddNotifyIcon()
	mw.showInfo("提醒","壁纸切换器已运行。",walk.MsgBoxIconInformation)

	mw.Run()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}