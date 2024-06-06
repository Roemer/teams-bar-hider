package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/roemer/gotaskr"
	"github.com/roemer/gotaskr/gttools"
)

func main() {
	os.Exit(gotaskr.Execute())
}

func init() {
	gotaskr.Task("Build", func() error {
		if err := gotaskr.Tools.DotNet.DotNetBuild("src/TeamsBarHider.csproj", &gttools.DotNetBuildSettings{
			ToolSettingsBase: gttools.ToolSettingsBase{
				OutputToConsole: true,
			},
			Configuration: "Release",
		}); err != nil {
			return err
		}

		return zipFolder(`src\bin\Release\net5.0-windows`, "TeamsBarHider.zip")
	})
}

func zipFolder(srcFolder, destZip string) error {
	// Create a zip file
	zipFile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Create a new zip archive
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk the source folder
	err = filepath.Walk(srcFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path to the file or directory
		relPath, err := filepath.Rel(srcFolder, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			if relPath != "." {
				_, err := zipWriter.Create(relPath + "/")
				if err != nil {
					return err
				}
			}
			return nil
		}

		// Create a new entry in the zip file
		zipFileWriter, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// Open the file to be zipped
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Copy the file content to the zip entry
		_, err = io.Copy(zipFileWriter, file)
		return err
	})

	return err
}
