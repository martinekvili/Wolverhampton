package main

import (
    "archive/zip"
    "log"
    "path"
    "io"
    "os"
)

func CopyFile(sourceName string, targetName string) {
    source, err := os.Open(sourceName)
    if err != nil {
        log.Fatal(err)
        return
    }
    defer source.Close()
    
    target, err := os.Create(targetName)
    if err != nil {
        log.Fatal(err)
        return
    }
    defer target.Close()
    
    _, err = io.Copy(target, source)
    if err != nil {
        log.Fatal(err)
    }
}

func ExtractZipIntoFolder(fileName string, folderName string) {
	// Open a zip archive for reading.
	r, err := zip.OpenReader(fileName)
	if err != nil {
		log.Fatal(err)
        return
	}
	defer r.Close()

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
        if !f.FileInfo().IsDir() {
            rc, err := f.Open()
            if err != nil {
                log.Fatal(err)
            }
            
            newFileName := path.Join(folderName, f.Name)
            os.MkdirAll(path.Dir(newFileName), os.ModeDir)
            
            newFile, err := os.Create(newFileName)
            if err != nil {
                log.Fatal(err)
            }
            
            _, err = io.Copy(newFile, rc)
            if err != nil {
                log.Fatal(err)
            }
            rc.Close()
            newFile.Close();
        }	
	}
}