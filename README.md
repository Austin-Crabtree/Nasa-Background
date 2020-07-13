# Nasa Background
A CLI tool to query the Nasa Open APOD (Astronomy Picture of the Day) to get, the current, picture of the day. 

Plans to: 
- Handle video files when found
- Add command line arugements
    - -s only save don't change background
    - -l change location to save to from default
    - -o only output the response json
    - -d give a date to query and get that days picture (will change background)
    - -dr date range to query and save that range of pictures (will not change background)
- Check for operating system, then desktop operating system to allow for this to be used on Windows, Mac and Linux (Gnome, KDE, XFCE, etc)

Notes on features:
- Handle video files and images seperately
    - Check media type to see if it is and image or not. 