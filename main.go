// make_http_request.go
package main

import (
	"bufio"
	"container/list"
	"fmt"
	"github.com/golang/freetype"
	"github.com/reujab/wallpaper"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// getImageLink parses Bing and gets img url
// https://www.devdungeon.com/content/web-scraping-go
func getImageURL() (imgURL string, imgFilename string) {
	domain, url := "https://www.bing.com/", "https://cn.bing.com/"
	// Make HTTP GET request
	response, err := http.Get(domain)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Turn HTML string to document
	re := regexp.MustCompile("data-ultra-definition-src=\".+?\\.jpg")
	htmlData, _ := ioutil.ReadAll(response.Body)
	imgURL = re.FindString(string(htmlData))
	imgURL = imgURL[28:]
	imgFilename = imgURL[6:]
	//document, err := goquery.NewDocumentFromReader(response.Body)
	//if err != nil {
	//	log.Fatal("Error loading HTTP response body.\n", err)
	//}

	// Find background url
	//document.Find("link#bgLink").Each(func(index int, element *goquery.Selection) {
	//	imgSrc, exists := element.Attr("href")
	//	if exists {
	//		url = domain + imgSrc
	//	}
	//})
	if 0 < len(imgURL) {
		log.Println("Image URL found: " + imgURL)
		return url + imgURL, imgFilename
	} else {
		return "", ""
	}
}

// exists returns whether the given file or directory exists
// https://stackoverflow.com/a/10510783
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// downloadImg saves image to ./.data directory
// and returns file path
// https://stackoverflow.com/questions/22417283/save-an-image-from-url-to-file
func downloadImg(url string, imgFileName string) string {
	dir := "./img_data/"

	noExistToCreate(dir)

	imageFilePath := dir + imgFileName
	fileExists, err := exists(imageFilePath)

	if !fileExists {
		// get image
		log.Println("Downloading image file to: " + imageFilePath)
		response, err := http.Get(url)
		if err != nil {
			log.Fatal("Couln't Download Image\n", err)
		}
		defer response.Body.Close()

		// create and open file
		i := len(url) - 1
		for i >= 0 && url[i] != '.' {
			i--
		}

		file, err := os.Create(imageFilePath)
		if err != nil {
			log.Fatal("Couldn't Create File\n", err)
		}
		defer file.Close()

		// Copy to file
		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatal("Couldn't Save Image\n", err)
		}

	} else {
		log.Println("Image file already existed. skip download.")
	}
	// Get current directory
	// https://gist.github.com/arxdsilva/4f73d6b89c9eac93d4ac887521121120
	dir, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir + imageFilePath[1:]
}

//不存在，则进行创建该文件件
func noExistToCreate(dir string) {
	dirExists, err := exists(dir)

	if err != nil {
		log.Fatal("Error finding directory:\n", err)
	}

	if !dirExists {
		err = os.Mkdir(dir, 0777)

		if err != nil {
			log.Fatal("Couldn't Create Directory\n", err)
		}
	}
}

//增加水印
func addWaterMark(file string, text []string) string {

	var dpi float64 = 72
	fontfile := "msyh.ttf"
	hinting := "none"
	var size float64 = 200
	spacing := 1.5
	wonb := false

	dir, err := GetCurrentPath()
	if err != nil {
		log.Fatal(err)
	}

	//读取字体
	fontBytes, err := ioutil.ReadFile(dir + fontfile)
	if err != nil {
		log.Println(err)
		return ""
	}
	//解析字体
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return ""
	}

	// 初始化图片背景
	fg := image.Black

	if wonb {
		fg = image.White
	}
	//初始化一张图片,生成原图
	imgB, _ := os.Open(file)
	img, _ := jpeg.Decode(imgB)
	defer imgB.Close()
	b := img.Bounds()
	rgba := image.NewNRGBA(b)
	draw.Draw(rgba, rgba.Bounds(), img, image.ZP, draw.Src)

	//在图片上面添加文字
	c := freetype.NewContext()
	c.SetDPI(dpi)
	//设置字体
	c.SetFont(f)
	//设置大小
	c.SetFontSize(size)
	//设置边界
	c.SetClip(rgba.Bounds())
	//设置背景底图
	c.SetDst(rgba)
	//设置背景图
	c.SetSrc(fg)
	//设置字体颜色(红色)
	c.SetSrc(image.NewUniform(color.RGBA{255, 0, 0, 255}))
	//设置提示
	switch hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// 画文字
	//设置水印偏移量
	pt := freetype.Pt(img.Bounds().Dx()-2500-10, img.Bounds().Dy()-int(c.PointToFixed(size)>>6)*2*len(text))
	//pt := freetype.Pt(10, 10+int(c.PointToFixed(size)>>6))
	for _, s := range text {
		_, err = c.DrawString(s, pt)
		if err != nil {
			log.Println(err)
			return ""
		}
		pt.Y += c.PointToFixed(size * spacing)
	}

	noExistToCreate(dir + "img_data/.tmp")

	imgw, _ := os.Create(dir + "img_data/.tmp/out.jpg")
	jpeg.Encode(imgw, rgba, &jpeg.Options{100})
	defer imgw.Close()
	return dir + "img_data/.tmp/out.jpg"
}
func setImageAsWallpaper(text *[]string) {
	url, _ := getImageURL()
	now := time.Now()
	file := downloadImg(url, now.Format("2006-01-02")+".jpg")

	waterMakerFile := addWaterMark(file, *text)
	fmt.Println(waterMakerFile)
	wallpaper.SetFromFile(waterMakerFile)
	log.Println("Enjoy, bye.")

	//获娶当前壁纸
	//background, err := wallpaper.Get()
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("Current wallpaper:", background)
	//wallpaper.SetFromFile("/usr/share/backgrounds/gnome/adwaita-day.jpg")
	//wallpaper.SetFromURL("https://i.imgur.com/pIwrYeM.jpg")
}

func GetMotto(fileName string, mottoList *list.List) {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	var size = stat.Size()
	fmt.Println("file size=", size)

	buf := bufio.NewReader(file)
	var text = []string{}
	for {
		line, err := buf.ReadString('\n')
		//if err == io.EOF {
		//	break  //finished reading
		//}

		line = strings.TrimSpace(line)
		//本组里面没有内容，直接遇到了空行，则直接跳过
		if strings.TrimSpace(line) == "" && len(text) == 0 && err != io.EOF {
			continue
		}
		//本组里面有内容，遇到了空行，则直接把数组存起来，新建数据，然后继续下一次循环
		if (strings.TrimSpace(line) == "" || err == io.EOF) && len(text) != 0 {
			mottoList.PushBack(text)
			text = []string{} //走到了文件结尾
			if err == io.EOF {
				break
			} else {
				continue
			}
		}
		//走到了文件结尾
		if err == io.EOF {
			break
		}
		//处理异常
		if err != nil {
			fmt.Errorf("read failed:%v", err)
			return
		}
		text = append(text, line)
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return
			}
		}
	}
}
//go build -ldflags "-H windowsgui"
func main() {

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

	select {}
}

func setNextBwall(mottoList *list.List, conf *BwallConfig) {
	var curentItem *list.Element
	curentItem = GetNextItem(curentItem, mottoList, conf)
	//还原成数组类型
	itemValue := curentItem.Value.([]string)

	setImageAsWallpaper(&itemValue)

	//保存最新的配置文件
	conf.CurrentText = strings.Join(itemValue, "$")
	SaveConfig(conf)

}

//获取下一项名言名句
func GetNextItem(curentItem *list.Element, mottoList *list.List, conf *BwallConfig) *list.Element {
	//默认设置为第一项
	curentItem = mottoList.Front()
	//配置中有值，则进行查找
	if conf.CurrentText != "" {
		//遍历查找，如果找到相同的，则设为相同的下一项
		for item := mottoList.Front(); nil != item; item = item.Next() {
			if strings.Join(item.Value.([]string), "$") == conf.CurrentText {
				if item.Next() == nil { //最后一个，则设为第一个
					curentItem = mottoList.Front()
				} else {
					curentItem = item.Next()
				}
				break
			}
		}
	}

	return curentItem
}
