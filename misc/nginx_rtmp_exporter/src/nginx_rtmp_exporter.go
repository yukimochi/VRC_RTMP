package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type rtmp struct {
	NginxVersion     string   `xml:"nginx_version"`
	NginxRtmpVersion string   `xml:"nginx_rtmp_version"`
	Compiler         string   `xml:"compiler"`
	Built            string   `xml:"built"`
	Pid              int      `xml:"pid"`
	Uptime           int      `xml:"uptime"`
	Naccepted        int      `xml:"naccepted"`
	BwIn             int      `xml:"bw_in"`
	BytesIn          int      `xml:"bytes_in"`
	BwOut            int      `xml:"bw_out"`
	BytesOut         int      `xml:"bytes_out"`
	Server           []server `xml:"server"`
}

type server struct {
	Application []application `xml:"application"`
}

type application struct {
	Name string `xml:"name"`
	Live []live `xml:"live"`
}

type live struct {
	Stream   []stream `xml:"stream"`
	Nclients int      `xml:"nclients"`
}

type stream struct {
	Name     string   `xml:"name"`
	Time     int      `xml:"time"`
	BwIn     int      `xml:"bw_in"`
	BytesIn  int      `xml:"bytes_in"`
	BwOut    int      `xml:"bw_out"`
	BytesOut int      `xml:"bytes_out"`
	BwAudio  int      `xml:"bw_audio"`
	BwVideo  int      `xml:"bw_video"`
	Client   []client `xml:"client"`
	Meta     []meta   `xml:"meta"`
	Nclients int      `xml:"nclients"`
}

type client struct {
	ID        int    `xml:"id"`
	Address   string `xml:"address"`
	Time      int    `xml:"time"`
	Flashver  string `xml:"flashver"`
	Pageurl   string `xml:"pageurl"`
	Swfurl    string `xml:"swfurl"`
	Dropped   int    `xml:"dropped"`
	Avsync    int    `xml:"avsync"`
	Timestamp int    `xml:"timestamp"`
}

type meta struct {
	Video []video `xml:"video"`
	Audio []audio `xml:"audio"`
}

type video struct {
	Width     int    `xml:"width"`
	Height    int    `xml:"height"`
	FrameRate int    `xml:"frame_rate"`
	Codec     string `xml:"codec"`
	Profile   string `xml:"profile"`
	Compat    int    `xml:"compat"`
	Level     string `xml:"level"`
}

type audio struct {
	Codec      string `xml:"codec"`
	Profile    string `xml:"profile"`
	Channels   int    `xml:"channels"`
	SampleRate int    `xml:"sample_rate"`
}

var url string

func main() {
	f := flag.String("listen", "9100", "listen port.")
	g := flag.String("url", "", "url of nginx-rtmp-module statistics xml.")
	flag.Parse()

	url = *g

	fmt.Println("Server listening at :" + *f)
	http.HandleFunc("/metrics", handler)
	http.ListenAndServe(":"+*f, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data := fetch(url)
		fmt.Fprintf(w, build(data))
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
}

func fetch(url string) *rtmp {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "HTTP Get error:", err)
		return nil
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	raw := string(byteArray)

	data := new(rtmp)
	if err := xml.Unmarshal([]byte(raw), data); err != nil {
		fmt.Fprintln(os.Stderr, "XML Unmarshal error:", err)
		return nil
	}
	return data
}

func build(data *rtmp) string {
	var text = ""

	if data != nil {
		// Server Info
		prefix := make([]string, 0)
		prefix = append(prefix, "nginx_rtmp")

		prefix = append(prefix, "accepted")
		text += makeText(prefix, "accepted", data.Naccepted)
		prefix = prefix[:len(prefix)-1]

		prefix = append(prefix, "bandwidth")
		text += makeText(prefix, "in", data.BwIn)
		text += makeText(prefix, "out", data.BwOut)
		prefix = prefix[:len(prefix)-1]

		prefix = append(prefix, "transmited")
		text += makeText(prefix, "in", data.BytesIn)
		text += makeText(prefix, "out", data.BytesOut)
		prefix = prefix[:len(prefix)-1]

		prefix = append(prefix, "live")

		prefix = append(prefix, "clients")
		text += makeText(prefix, "clients", data.Server[0].Application[0].Live[0].Nclients)
		prefix = prefix[:len(prefix)-1]

		prefix = append(prefix, "stream")
		streams := data.Server[0].Application[0].Live[0].Stream
		for key := range streams {
			prefix = append(prefix, streams[key].Name)

			prefix = append(prefix, "media_bandwidth")
			text += makeText(prefix, "audio", streams[key].BwAudio)
			text += makeText(prefix, "video", streams[key].BwVideo)
			prefix = prefix[:len(prefix)-1]

			prefix = append(prefix, "clients")
			text += makeText(prefix, "clients", streams[key].Nclients)
			prefix = prefix[:len(prefix)-1]

			prefix = append(prefix, "bandwidth")
			text += makeText(prefix, "in", streams[key].BwIn)
			text += makeText(prefix, "out", streams[key].BwOut)
			prefix = prefix[:len(prefix)-1]

			prefix = append(prefix, "transmited")
			text += makeText(prefix, "in", streams[key].BytesIn)
			text += makeText(prefix, "out", streams[key].BytesOut)
			prefix = prefix[:len(prefix)-1]

			prefix = prefix[:len(prefix)-1]
		}
	}

	return text
}

func makeText(prefix []string, title string, value int) string {
	var text = ""

	for i := 0; i < len(prefix); i++ {
		text += prefix[i]
		text += "."
	}
	text += title
	text += "=" + strconv.Itoa(value)
	text += "\n"

	return text
}
