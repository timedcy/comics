/* vim: se ts=2 sw=2 enc=utf-8: */
package comics

import (
	"fmt"
	"strings"
	"time"
	"net/http"
	"io/ioutil"
	"errors"
)

var BASE_URL="http://www.gocomics.com/"

type GoComics struct {
	Title string
	Prefix string
	Subfix string
	DownloadDir string
	url string
	PicUrl string
}

func (gc *GoComics) Init(){
	gc.Prefix =`class="strip" src="`
	gc.Subfix =`"`
	gc.url = BASE_URL+gc.Title+time.Now().Format("/2006/01/02")
	fmt.Println("comics base url: ",gc.url)
}

func (gc *GoComics) Parse() error{
	resp, err := http.Get(gc.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	html := string(bytes)
	idx := strings.Index(html, gc.Prefix)
	if idx == -1 {
		return errors.New("GoComics html invalid, maybe it changed.")
	}
	html = html[idx+len(gc.Prefix):]
	idx = strings.Index(html, gc.Subfix)
	if idx == -1 {
		return errors.New("GoComics html invalid, maybe it changed.")
	}
	gc.PicUrl = html[:idx]
	return nil
}

func (gc *GoComics) PicBytes() ([]byte, error) {
	if gc.PicUrl == ""{
		return nil, errors.New("PicUrl should not be empty")
	}
	fmt.Println("Got pic url: ", gc.PicUrl)
	client := &http.Client{}
	req, err :=http.NewRequest("GET",gc.PicUrl,nil)
	gc.forgeryHttpRequest(req)
	resp, err :=client.Do(req)
	if err != nil{
		return nil, err
	}
	bytes,err :=ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}
	return bytes, nil
}

func (gc *GoComics) forgeryHttpRequest(req *http.Request){
	req.Header.Add("Accept",`*/*`)
	req.Header.Add("Accept-Charset",`UTF-8,*;q=0.5`)
	req.Header.Add("Accept-Encoding",`gzip,deflate,sdch`)
	req.Header.Add("Accept-Language",`en-US,en;q=0.8`)
	req.Header.Add("Cache-Control",`max-age=0`)
	req.Header.Add("Connection",`keep-alive`)
	req.Header.Add("Referer",gc.url)
	req.Header.Add("User-Agent",`Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.31 (KHTML, like Gecko) Chrome/26.0.1410.64 Safari/537.31`)
}
