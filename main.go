package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func main() {
	// Declare the flags
	var wantedLanguages string
	var automatic bool

	// Define the flags
	flag.StringVar(&wantedLanguages, "langs", "und,eng,swe,jap,jpn", "comma-separated list of languages")
	flag.BoolVar(&automatic, "automatic", false, "make the program non-interactive")

	// Parse the flags
	flag.Parse()

	// Convert the comma-separated string of languages into a slice
	languages := strings.Split(wantedLanguages, ",")

	var fileNames []string

	// Loop over the arguments and check if they are valid file or directory names.
	for _, arg := range flag.Args() {
		info, err := os.Stat(arg)
		if err == nil {
			if info.IsDir() {
				// If the argument is a directory, call the findFilesInDirectory function.
				files := findFilesInDirectory(arg)

				for _, file := range files {
					fileNames = append(fileNames, file)
				}
			} else {
				fileNames = append(fileNames, arg)
			}
		}
	}

	processFiles(fileNames, languages, automatic)
}

func processFiles(fileNames []string, wantedLanguages []string, automatic bool) {
	// Count all the files
	fileCount := len(fileNames)
	fileNo := 0

	for _, fileName := range fileNames {
		fileNo++

		var ffprobe FFProbe

		// Identify the input file to index all the tracks using ffprobe.
		ffprobe.Identify(fileName)

		// Put the filtered data into ffmpeg to build and run an
		// ffmpeg command to copy the wanted tracks to a new file.
		ffmpeg := FFMpeg{
			inputFilePath:  fileName,
			audioTracks:    ffprobe.GetAudioTracks(wantedLanguages),
			subtitleTracks: ffprobe.GetSubtitleTracks(wantedLanguages),
			videoTracks:    ffprobe.GetVideoTracks(),
		}

		// Some nice output.
		fmt.Printf("[%d/%d] Preparing to process %s\n", fileNo, fileCount, fileName)

		// Merge all the tracks we want to keep to a single slice.
		allTracks := append(ffmpeg.videoTracks, append(ffmpeg.audioTracks, ffmpeg.subtitleTracks...)...)

		// Print a table for the overview of the changes
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Index", "CodecName", "CodecType", "Language", "Keep"})
		for _, item := range ffprobe.GetTracksStatus(allTracks) {
			table.Append([]string{item.Index, item.CodecName, item.CodecType, item.Language, item.Keep})
		}

		// Output the table
		table.Render()

		// Check if a file actually needs processing and ignore if it
		// doesn't need any changes.
		if !ffprobe.NeedsProcessing(allTracks) {
			fmt.Printf("[%d/%d] There's no change to apply to this file\n", fileNo, fileCount)
			continue
		}

		fmt.Printf("[%d/%d] Command to execute: %s\n", fileNo, fileCount, ffmpeg.FormatCommandParts())

		// If the doesn't run in automatic mode, we prompt the user
		// for a confirmation to apply changes.
		if !automatic {
			fmt.Printf("[%d/%d] Run ffmpeg command on %s? [Y/n] ", fileNo, fileCount, fileName)

			// Prompt user for a confirmation of the actions that will be
			// taken on said file.
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(response)

			// Check response, if not yes or empty, skip item.
			if strings.ToLower(response) != "y" && response != "" {
				fmt.Printf("[%d/%d] Skipping %s\n", fileNo, fileCount, fileName)
				continue
			}
		}

		// Run the ffmpeg command.
		err := ffmpeg.Run()

		// Check for errors, if error, skip remaining steps.
		if err != nil {
			fmt.Printf("[%d/%d] Failed at running ffmpeg command\n", fileNo, fileCount)
			fmt.Printf("[%d/%d] Error: %s\n", fileNo, fileCount, err)
			continue
		}

		fmt.Printf("[%d/%d] Sucessfully ran ffmpeg command\n", fileNo, fileCount)

		// Rename the original file to a different temporary path.
		err = os.Rename(ffmpeg.inputFilePath, ffmpeg.inputFilePath+".rename-tmp")

		if err != nil {
			fmt.Printf("[%d/%d] Failed to rename the input file, bailing out by removing the output file\n", fileNo, fileCount)

			os.Remove(ffmpeg.outputFilePath)
			continue
		}

		err = os.Rename(ffmpeg.outputFilePath, ffmpeg.inputFilePath)
		if err != nil {
			fmt.Printf("[%d/%d] Failed to rename the temporary file to the original file name", fileNo, fileCount)
			continue
		}

		// Remove the temporary storage used of the input file while
		// renaming the new file to it's new place.
		os.Remove(ffmpeg.inputFilePath + ".rename-tmp")
	}
}
