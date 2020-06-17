package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var seconds int
var verbose bool

func getFileLength(f string) int {
	cmd := exec.Command(
		"ffprobe",
		"-loglevel", "error", "-show_entries", "format=duration", "-print_format", "default=noprint_wrappers=1:nokey=1", f,
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return -1
	}

	// ffprobe output is SS.MICROSECONDS
	s := strings.Split(out.String(), ".")

	l, err := strconv.Atoi(s[0])
	if err != nil {
		return -1
	}

	return l
}

func splitFileByTime(file string, length int, parts int) error {
	extension := filepath.Ext(file)
	name := strings.TrimSuffix(file, extension)

	for i := 0; i <= parts; i++ {
		p := fmt.Sprintf("%v-%v%v", name, i+1, extension)
		s := seconds * i

		// Avoid generating empty file
		if s == length {
			continue
		}

		cmd := exec.Command(
			"ffmpeg",
			"-i", file, "-acodec", "copy", "-vcodec", "copy", "-ss", strconv.Itoa(s), "-t", strconv.Itoa(seconds), p,
		)

		if verbose {
			fmt.Println(cmd.String())
		}

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			if strings.Contains(stderr.String(), "already exists") {
				fmt.Println(fmt.Sprintf("error: the file %v already exists", p))
			}
			return err
		}
	}

	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vsplit [FILE]",
	Short: "Split a video or audio file into multiple files",
	Long:  "Split a video or audio file into multiple files",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// File
		f := args[0]

		if seconds <= 0 {
			fmt.Println("error: invalid seconds parameter")
			return
		}

		length := getFileLength(f)
		if length <= 0 {
			fmt.Println("error: file length 0 seconds")
			return
		}

		parts := length / seconds

		fmt.Println(fmt.Sprintf("Splitting %v", f))

		err := splitFileByTime(f, length, parts)
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(&seconds, "seconds", "s", 0, "Length of split files in seconds")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Display FFmpeg commands executed")

	rootCmd.MarkFlagRequired("seconds")
}
