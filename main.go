// make_http_request.go
package main

import (
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
func setImageAsWallpaper() {
	url, _ := getImageURL()
	now := time.Now()
	file := downloadImg(url, now.Format("2006-01-02")+".jpg")
	var text = []string{"易怒的人", "缺乏就事论事的能力"}
	waterMakerFile := addWaterMark(file, text)
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

func main() {

	setImageAsWallpaper()
}
