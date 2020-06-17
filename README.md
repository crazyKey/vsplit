# Video Split Go Command Line Tool

Split a video into multiple files of a specified length. Make sure you have [FFmpeg](https://www.ffmpeg.org/download.html) installed.

As `vsplit` uses FFmpeg, it's possible to split audio files too.

#### Example

Display help message `vsplit -h`

Split file.mp4 into many files of 60 seconds: `vsplit file.mp4 -s 60`

Split audio.mp3 into many files of 15 seconds: `vsplit file.mp3 -s 15`