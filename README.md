# Nasa Background
A CLI tool to query the Nasa Open APOD (Astronomy Picture of the Day) to get, the current, picture of the day. 

Requirements:
- Linux:
  - Gnome desktop environment
- Windows:
  - Have powershell installed to be able to set registry values
- API:
  - You in the example config file is the demo api key that you can use, but it
    suggested to go get your own api key from [Nasa](https://api.nasa.gov/)

Install:
Download or clone the repo and then build main.go using go. Then setup the config.json file to have your api key (optional) and then the path on your desktop to where you would like to save the pictures to. 

Plans to: 
- Handle video files when found
- Add command line arugements
    - -s only save don't change background (Do I need to do this?)
    - -o only output the response json
    - -d give a date to query and get that days picture (will change background)
    - -dr date range to query and save that range of pictures (will not change background)
    - -f force re-download the picture
- Check for operating system, then desktop operating system to allow for this to be used on Windows, Mac and Linux (Gnome, KDE, XFCE, etc)

Notes on features:
- Handle video files and images seperately
    - Check media type to see if it is and image or not. 
