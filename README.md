# Youtube Audio Streamer

A simple web service written in Go. As the name suggests, it's purpose is streaming audio from youtube videos. **For more info, please refer [here](https://longwintershadows.com/articles/4c92a655-4cea-43cc-8198-a7ea97b78abc).**

### Building
Compiling a binary is pretty straightforward. Run `go install` and `go build` to make a binary. Tested with go version **3.17**. 

### Usage
**Important!** The service requires FFmpeg installation to work! 
The binary allows downloading the audio from youtube without launching the service. Simply run `./{binary-name} [OPTIONS] {youtube-video-url}`. The file is downloaded to the same folder as binary by default. Please mind that download functionality is pretty limited. The details of the download process are covered in the **Downloading** section. 
Run `./{binary-name} serve [OPTIONS]` to launch the web service. 
Binary requires a config.yml file to be present in the same folder. Here're the possible configuration properties:
 - `env` - specifies environment, expected values are `prod` or `debug`. Could be overriden and set to `debug` by specifying `-d` flag when running the command. Affects log output.
 - `port` - specifies port web service will be listening on. Could be overriden by provding flag `-p` and specifying new port value.
 - `version` - specifies app version. Doesn't really affect anything
 - `output-route` - specifies URL route for statically served output m3u8 files.
 - `output-directory` - specifies directory where converted m3u8 files are saved. Those files are also statically served from this directory. Could be overriden by providing `-o` flag and specifying a new path.
 - `source-directory` - specifies directory where files downloaded from youtube are saved. Could be overriden by providing `-s` flag and specifying a new path.
 - `completion-marker` - specifies completion marker file name.

### Under the hood
Service command-line interactions and configuration are handled with [`Cobra`](https://github.com/spf13/cobra) and [`Viper`](https://github.com/spf13/viper). Both are amazing libraries that are just such a pleasure to work with. The setup is pretty conventional. Using `cobra` we check for command options and select the corresponding handler. If a user didn’t provide the serve command, we simply use the download handler to attempt to obtain the audio. If serve command is provided, we start the server. The service API is built with [`echo`](https://github.com/labstack/echo) framework. 

### Frontend
Could be found [here](https://github.com/AlexKLWS/youtube-audio-stream-client).
