package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/jhoonb/archivex"
	"github.com/mh-cbon/archive/uncompress"
	"github.com/urfave/cli"
)

var VERSION = "0.0.0"

func main() {
	app := cli.NewApp()
	app.Name = "archive"
	app.Version = VERSION
	app.Usage = "Command line to create and extract archive files"
	app.UsageText = "archive <cmd> <options>"
	app.Commands = []cli.Command{
		{
			Name:      "create",
			Usage:     "Create a new archive",
			UsageText: "archive create <options>",
			Action:    create,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output, o",
					Value: "",
					Usage: "Output file",
				},
				cli.StringFlag{
					Name:  "change-dir, C",
					Value: "",
					Usage: "Change directory before archiving files",
				},
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "Force overwrite",
				},
			},
		},
		{
			Name:      "extract",
			Usage:     "Extract an archive file",
			UsageText: "archive extract <options> <packages>",
			Action:    extract,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dest, d",
					Value: "",
					Usage: "Destination path",
				},
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "Force overwrite",
				},
			},
		},
	}

	app.Run(os.Args)
}

func create(c *cli.Context) error {
	output := c.String("output")
	changeDir := c.String("change-dir")
	cwd, err := os.Getwd()
	if err != nil {
    return cli.NewExitError(err.Error(), 1)
	}
	if len(changeDir) > 0 {
		err := os.Chdir(changeDir)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}
	files := make([]string, 0)
	for _, f := range c.Args() {
		sF := string(f)
		_, err := os.Stat(f)
		if !os.IsNotExist(err) {
			files = append(files, sF)
		} else {
			fmt.Println("File '" + sF + "' does not exist")
		}
	}
	if len(changeDir) > 0 {
		err := os.Chdir(cwd)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}
	if len(output) == 0 {
		return cli.NewExitError("Output argument is required", 1)
	}
	if len(files) == 0 {
		return cli.NewExitError("<file> arguments to archive are missing, or provided paths does not exists.", 1)
	}
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		if c.Bool("force") == false {
			return cli.NewExitError("Output file already exists. Use force argument.", 1)
		}
	}

	var archiver archivex.Archivex
	if strings.Index(output, ".zip") > -1 {
		archiver = new(archivex.ZipFile)
	} else if strings.Index(output, ".tar") > -1 {
		archiver = new(archivex.TarFile)
	} else if strings.Index(output, ".tgz") > -1 {
		output = strings.Replace(output, filepath.Ext(output), "", -1) + ".tar.gz"
		archiver = new(archivex.TarFile)
	} else {
		return cli.NewExitError("Cannot handle archive '"+output+"' is not zip|tar|tar.gz.", 1)
	}

	dir := filepath.Dir(output)
	if len(dir) > 0 {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return cli.NewExitError("Failed to create output directory '"+dir+"'.", 1)
		}
	}

	err = archiver.Create(output)
	if err != nil {
		fmt.Println(err)
		return cli.NewExitError("Failed to create output archive '"+output+"'.", 1)
	}
	gotErr := false
	if len(changeDir) > 0 {
		err := os.Chdir(changeDir)
		if err != nil {
			fmt.Println(err)
			gotErr = true
		}
	}
	for _, f := range files {
		s, _ := os.Stat(f)
		if s.IsDir() {
			err = archiver.AddAll(f, true)
			if err != nil {
				fmt.Println(err)
				gotErr = true
			}
		} else {
			err = archiver.AddFile(f)
			if err != nil {
				fmt.Println(err)
				gotErr = true
			}
		}
	}
	if len(changeDir) > 0 {
		err := os.Chdir(cwd)
		if err != nil {
			fmt.Println(err)
			gotErr = true
		}
	}
	err = archiver.Close()
	if err != nil {
		fmt.Println(err)
		gotErr = true
	}
	if gotErr {
		os.Remove(output)
		return cli.NewExitError("Failed to create the archive '"+output+"'.", 1)
	}
	s, _ := os.Stat(output)
	fmt.Println("✓ " + output + ": " + humanize.Bytes(uint64(s.Size())))
	return nil
}

func extract(c *cli.Context) error {
	dest := c.String("dest")
	src := string(c.Args().Get(0))
	if len(dest) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			return cli.NewExitError("Cannot determine destination path", 1)
		}
		dest = wd
	}
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		if c.Bool("force") == false {
			return cli.NewExitError("Destination already exists. Use force argument.", 1)
		}
	}
	if len(src) == 0 {
		return cli.NewExitError("<file> source arguments are required.", 1)
	}
	info := make(chan string)
	itemsCnt := 0
	go func() {
		for f := range info {
			fmt.Println(f)
			itemsCnt++
		}
	}()
	err := uncompress.Uncompress(src, dest, info)
	close(info)
	if err != nil {
		fmt.Println(err)
		return cli.NewExitError("Failed to extract the archive '"+src+"'.", 1)
	}
	fmt.Println("✓ " + dest + ": " + strconv.Itoa(itemsCnt) + " files created")
	return nil
}
