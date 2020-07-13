package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"./model"
)

// TODO make better error handling
func checkErr(err error) {
	if os.IsNotExist(err) {
		panic(err)
	}
}

func main() {
	// Setup the commandline args
	save := flag.Bool("s", false, "Only save the picture and don't change background.")
	print := flag.Bool("o", false, "Only output the response json")
	//dateChange := flag.String("d", "", "Give date to query and get that days picture, while changing background.")
	//dateQuery := flag.String("dr", "", "Give a date range to query and save those days pictures. Will not change background")
	flag.Parse()

	// Setup request values
	url := "https://api.nasa.gov/planetary/apod"
	date := time.Now().Format("2006-01-02")

	// Read in the config data
	configPath := filepath.Join(".", "config.json") // Change to this being fallback if a path isn't given
	configData, err := ioutil.ReadFile(configPath)
	checkErr(err)
	config := model.Config{}
	err = json.Unmarshal(configData, &config)
	checkErr(err)

	// Start the http request
	req, err := http.NewRequest("GET", url, nil)
	checkErr(err)
	query := req.URL.Query()
	query.Add("api_key", config.APIKey)
	query.Add("hd", "True")
	query.Add("date", date)
	req.URL.RawQuery = query.Encode()

	// Execute the http request
	resp, err := http.Get(req.URL.String())
	checkErr(err)

	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	if *print {
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, body, "", "\t") // pretty print the json
		fmt.Println(prettyJSON.String())
		return
	}

	// Read the http response into the response model
	nasa := model.NasaResp{}
	err = json.Unmarshal(body, &nasa)
	checkErr(err)

	// TODO Check media type of the response to know how to proceed
	// possible types: video, image
	// images are jpg
	// videos are youtube urls
	imgTitle := strings.ReplaceAll(nasa.Title, " ", "-")
	imgName := fmt.Sprintf("%s_%s.jpg", imgTitle, nasa.Date)
	imgPath := filepath.Join(config.SavePath, imgName)

	// If the file exists then skip this part
	if _, err = os.Stat(imgPath); os.IsNotExist(err) { // Add a -f flage to force re-download
		if nasa.Hdurl != "" {
			resp, err = http.Get(nasa.Hdurl)
		} else {
			resp, err = http.Get(nasa.URL)
		}
		checkErr(err)

		body, err = ioutil.ReadAll(resp.Body)
		checkErr(err)

		err = ioutil.WriteFile(imgPath, body, 755)
		checkErr(err)

		if *save {
			return
		}

		if runtime.GOOS == "linux" {
			// DO a better check on wether they are in GNOME, XFCE or another DE
			fileURI := fmt.Sprintf("file://%s", imgPath)
			cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", fileURI)
			err = cmd.Run()
			checkErr(err)
		} else if runtime.GOOS == "windows" {
			cmd := exec.Command("powershell.exe", "Set-ItemProperty", "-path", "'HKCU:\\Control Panel\\Desktop\\'", "-name", "wallpaper", "-value", imgPath)
			err = cmd.Run()
			checkErr(err)
			cmd = exec.Command("RUNDLL32.EXE", "user32.dll,UpdatePerUserSystemParameters")
			err = cmd.Run()
			checkErr(err)
		}
	}
}
