package main

import (
	"encoding/json"
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

func checkErr(err error) {
	if os.IsNotExist(err) {
		panic(err)
	}
}

func main() {
	url := "https://api.nasa.gov/planetary/apod"
	today := time.Now().Format("2006-01-02")
	configPath := filepath.Join(".", "config.json")
	configData, err := ioutil.ReadFile(configPath)
	checkErr(err)
	config := model.Config{}
	err = json.Unmarshal(configData, &config)
	checkErr(err)
	// TODO Change this to not create a new variable and to use config directly
	rootPath := config.SavePath
	// TODO Change this to not create a new variable and to use config directly
	apiKey := config.APIKey
	req, err := http.NewRequest("GET", url, nil)
	checkErr(err)

	query := req.URL.Query()
	query.Add("api_key", apiKey)
	query.Add("hd", "True")
	query.Add("date", today)

	req.URL.RawQuery = query.Encode()

	resp, err := http.Get(req.URL.String())
	checkErr(err)

	nasa := model.NasaResp{}
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	// TODO Check media type of the response to know how to proceed
	// possible types: video, image
	// images are jpg
	// videos are youtube urls
	err = json.Unmarshal(body, &nasa)
	checkErr(err)

	imgTitle := strings.ReplaceAll(nasa.Title, " ", "-")
	imgName := fmt.Sprintf("%s_%s.jpg", imgTitle, nasa.Date)
	imgPath := filepath.Join(rootPath, imgName)

	if _, err = os.Stat(imgPath); os.IsNotExist(err) {
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
