package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Rtmp struct {
	Nginx_version      string   `xml:"nginx_version"`
	Nginx_rtmp_version string   `xml:"nginx_rtmp_version"`
	Compiler           string   `xml:"compiler"`
	Built              string   `xml:"built"`
	Pid                int      `xml:"pid"`
	Uptime             int      `xml:"uptime"`
	Naccepted          int      `xml:"naccepted"`
	Bw_in              int      `xml:"bw_in"`
	Bytes_in           int      `xml:"bytes_in"`
	Bw_out             int      `xml:"bw_out"`
	Bytes_out          int      `xml:"bytes_out"`
	Server             []Server `xml:"server"`
}

type Server struct {
	Application []Application `xml:"application"`
}

type Application struct {
	Name string `xml:"name"`
	Live []Live `xml:"live"`
}

type Live struct {
	Stream   []Stream `xml:"stream"`
	Nclients int      `xml:"nclients"`
}

type Stream struct {
	Name      string   `xml:"name"`
	Time      int      `xml:"time"`
	Bw_in     int      `xml:"bw_in"`
	Bytes_in  int      `xml:"bytes_in"`
	Bw_out    int      `xml:"bw_out"`
	Bytes_out int      `xml:"bytes_out"`
	Bw_audio  int      `xml:"bw_audio"`
	Bw_video  int      `xml:"bw_video"`
	Client    []Client `xml:"client"`
	Meta      []Meta   `xml:"meta"`
	Nclients  int      `xml:"nclients"`
}

type Client struct {
	Id        int    `xml:"id"`
	Address   string `xml:"address"`
	Time      int    `xml:"time"`
	Flashver  string `xml:"flashver"`
	Pageurl   string `xml:"pageurl"`
	Swfurl    string `xml:"swfurl"`
	Dropped   int    `xml:"dropped"`
	Avsync    int    `xml:"avsync"`
	Timestamp int    `xml:"timestamp"`
}

type Meta struct {
	Video []Video `xml:"video"`
	Audio []Audio `xml:"audio"`
}

type Video struct {
	Width      int    `xml:"width"`
	Height     int    `xml:"height"`
	Frame_rate int    `xml:"frame_rate"`
	Codec      string `xml:"codec"`
	Profile    string `xml:"profile"`
	Compat     int    `xml:"compat"`
	Level      string `xml:"level"`
}

type Audio struct {
	Codec       string `xml:"codec"`
	Profile     string `xml:"profile"`
	Channels    int    `xml:"channels"`
	Sample_rate int    `xml:"sample_rate"`
}

func main() {
	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") == "1" {
		fmt.Println("# mackerel-agent-plugin\n{ \"graphs\": { \"nginx_rtmp.accepted\": { \"label\": \"NGINX_RTMP.Accepted\", \"unit\": \"integer\", \"metrics\": [ { \"name\": \"accepted\", \"label\": \"Accepted\" } ] }, \"nginx_rtmp.bandwidth\": { \"label\": \"NGINX_RTMP.Bandwidth\", \"unit\": \"bytes/sec\", \"metrics\": [ { \"name\": \"in\", \"label\": \"In\" }, { \"name\": \"out\", \"label\": \"Out\" } ] }, \"nginx_rtmp.transmited\": { \"label\": \"NGINX_RTMP.Transmited\", \"unit\": \"bytes\", \"metrics\": [ { \"name\": \"in\", \"label\": \"In\" }, { \"name\": \"out\", \"label\": \"Out\" } ] }, \"nginx_rtmp.live.clients\": { \"label\": \"NGINX_RTMP.Live.Clients\", \"unit\": \"integer\", \"metrics\": [ { \"name\": \"clients\", \"label\": \"Clients\" } ] } }}")
		return
	}

	f := flag.String("url", "", "url of nginx-rtmp-module statistics xml.")
	flag.Parse()

	url := *f
	data := fetch(url)
	time := fmt.Sprint(time.Now().Unix())

	if data != nil {
		// Server Info
		fmt.Println(build("nginx_rtmp.accepted.accepted", strconv.Itoa(data.Naccepted), time))
		fmt.Println(build("nginx_rtmp.bandwidth.in", strconv.Itoa(data.Bw_in), time))
		fmt.Println(build("nginx_rtmp.bandwidth.out", strconv.Itoa(data.Bw_out), time))
		fmt.Println(build("nginx_rtmp.transmited.in", strconv.Itoa(data.Bytes_in), time))
		fmt.Println(build("nginx_rtmp.transmited.out", strconv.Itoa(data.Bytes_out), time))
		fmt.Println(build("nginx_rtmp.live.clients.clients", strconv.Itoa(data.Server[0].Application[0].Live[0].Nclients), time))
		streams := data.Server[0].Application[0].Live[0].Stream
		for key := range streams {
			fmt.Println(build("nginx_rtmp.live.stream."+streams[key].Name+".media_bandwidth.audio", strconv.Itoa(data.Server[0].Application[0].Live[0].Stream[key].Bw_audio), time))
			fmt.Println(build("nginx_rtmp.live.stream."+streams[key].Name+".media_bandwidth.video", strconv.Itoa(data.Server[0].Application[0].Live[0].Stream[key].Bw_video), time))
			fmt.Println(build("nginx_rtmp.live.stream."+streams[key].Name+".clients.clients", strconv.Itoa(data.Server[0].Application[0].Live[0].Stream[key].Nclients), time))
			fmt.Println(build("nginx_rtmp.live.stream."+streams[key].Name+".bandwidth.in", strconv.Itoa(data.Server[0].Application[0].Live[0].Stream[key].Bw_in), time))
			fmt.Println(build("nginx_rtmp.live.stream."+streams[key].Name+".bandwidth.out", strconv.Itoa(data.Server[0].Application[0].Live[0].Stream[key].Bw_out), time))
			fmt.Println(build("nginx_rtmp.live.stream."+streams[key].Name+".transmited.in", strconv.Itoa(data.Server[0].Application[0].Live[0].Stream[key].Bytes_in), time))
			fmt.Println(build("nginx_rtmp.live.stream."+streams[key].Name+".transmited.out", strconv.Itoa(data.Server[0].Application[0].Live[0].Stream[key].Bytes_out), time))
		}
	}
}

func fetch(url string) *Rtmp {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("HTTP Get error:", err)
		return nil
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	raw := string(byteArray)

	data := new(Rtmp)
	if err := xml.Unmarshal([]byte(raw), data); err != nil {
		fmt.Println("XML Unmarshal error:", err)
		return nil
	}
	return data
}

func build(name string, value string, time string) string {
	return name + "\t" + value + "\t" + time
}
