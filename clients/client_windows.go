package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"gdback_client/db"

	"github.com/gosuri/uiprogress"
	"gopkg.in/AlecAivazis/survey.v1"
)

var VERSION = "1.0.0"
var BANNER = `    _____________
   /____________/
   ||     ______
   ||    |_____/!
   ||          ||
   ||   _||_   ||
   ||  |_||_|  ||  
   ||__________||
   |/__________|/` + "  by Siriil v" + VERSION + "\n\n"
var BATCH_MAX_SIZE = 300
var BAR *uiprogress.Bar

const MaxInt64 = 1<<63 - 1

func main() {

	logger := db.NewCustomLogger()
	fmt.Println(BANNER)

	if runtime.GOOS != "windows" {
		logger.Error("The program must be run on Windows S.O.")
	}

	if !isSuperUserWindows() {
		logger.Error("The program must be run with administrator privileges")
	}

	selectedDisk, numCPUs, path := askOptions(logger)
	fmt.Println()

	time_start := time.Now()

	database, dbfilename, err := db.CreateDatabase()
	if err != nil {
		logger.Error("Error creating the database:", err)
	}
	logger.Info("Database saved as", "'"+dbfilename+"'")

	// Create a channel for communication files - Workers
	filePathChan := make(chan string)

	// Create a BAR with the maximum length until we find out the total num of files
	BAR = uiprogress.AddBar(MaxInt64).AppendCompleted().PrependElapsed()
	// ToDo: change AppendFunct to PrependFunc and try to make this work
	BAR.AppendFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("[+] Processed Files (%d/", b.Current())
	})
	uiprogress.Start()

	var wg sync.WaitGroup
	wg.Add(numCPUs)

	for i := 1; i <= numCPUs; i++ {
		go dataProcessor(logger, filePathChan, dbfilename, i, &wg)
	}

	// s := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
	// s.Prefix = "Calculating number of files "
	// s.Start()
	nrows, err := sendFilePaths2Channel(logger, database, filePathChan, selectedDisk+path)
	if err != nil {
		logger.Error("Error get files in disk selected:", err)
	}
	// s.Stop()

	// When all files sent, write numCPUs 0s to stop go routines
	for i := 1; i <= numCPUs; i++ {
		filePathChan <- "0"
	}

	// Update the bar now that we know the number of files
	BAR.Total = nrows
	BAR.AppendFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("%d)", nrows)
	})
	wg.Wait()
	uiprogress.Stop()

	database.Close()
	logger.Info("Processed", nrows, "files")

	logger.Info("Update metadata table")
	processMetadata(logger, dbfilename, time_start)

	time_elapsed := time.Since(time_start)
	logger.Info("Time elapsed:", time_elapsed)
	fmt.Println()

	fmt.Printf("Press ENTER key to close...")
	fmt.Scanln()
	os.Exit(0)
}

func askOptions(logger *db.CustomLogger) (string, int, string) {
	disks, err := getDisksWindows()
	if err != nil {
		logger.Error("Error getting the list of disks:", err)
	}
	var selectedDisk string
	prompt_select := &survey.Select{
		Message: "Select an disk:",
		Options: disks,
	}
	survey.AskOne(prompt_select, &selectedDisk, nil)

	numAvailableCPUs := runtime.NumCPU()
	var cpus []string
	var strNumCPUs string
	for i := 1; i <= numAvailableCPUs; i++ {
		cpus = append(cpus, fmt.Sprintf("%d", i))
	}
	prompt_select = &survey.Select{
		Message: "Select a number of cpu cores:",
		Options: cpus,
	}
	survey.AskOne(prompt_select, &strNumCPUs, nil)
	numCPUs, _ := strconv.Atoi(strNumCPUs)
	runtime.GOMAXPROCS(numCPUs)

	path := ""
	prompt_input := &survey.Input{
		Message: `Specifies path in disk (Leave blank for all) (Example: \Users\User1\):`,
	}
	survey.AskOne(prompt_input, &path, nil)
	if path != "" {
		_, err = os.Stat(selectedDisk + path)
		if err != nil {
			path = ""
			logger.Warning("The path entered is wrong, so the whole disk will be scanned")
		}
	}

	if (selectedDisk == "") || (numCPUs < 1) {
		logger.Error("Error parameter in menu selection")
	}

	return selectedDisk, numCPUs, path
}

func sendFilePaths2Channel(logger *db.CustomLogger, database *db.Database, filePathChan chan string, root string) (int, error) {
	var num_files int = 0

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if (err == nil) && (!info.IsDir()) {
			// send path to a worker through the channel
			filePathChan <- path
			num_files++
		}

		return nil
	})

	if err != nil {
		return num_files, fmt.Errorf("Failed to walk through files: %v", err)
	}

	return num_files, nil
}

func dataProcessor(logger *db.CustomLogger, fileChan <-chan string, dbpath string, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	files := make([]*db.Data, 0)

	database, err := db.Connect(dbpath)
	if err != nil {
		logger.Error("Error connecting the database in goroutine", id, ":", err)
	}
	defer database.Close()

	for {
		path := <-fileChan
		if path == "0" {
			break
		}

		fileInfo, err := os.Stat(path)
		if err == nil {
			DateCreation, _ := getFileCreationTimeWindows(path)
			HashMD5, _ := getFileMD5(path)
			fileData := &db.Data{
				FullPath:             path,
				FileName:             fileInfo.Name(),
				FileExtension:        filepath.Ext(path),
				HashMD5:              HashMD5,
				SizeBytes:            int(fileInfo.Size()),
				DateCreation:         DateCreation,
				DateLastModification: fileInfo.ModTime().Format("2006-01-02 15:04:05"),
			}

			files = append(files, fileData)
			BAR.Incr()
		}

		if len(files) >= BATCH_MAX_SIZE {
			// *QUESTION: No need to use mutex to write in the database?
			err := database.InsertDatas(files)
			if err != nil {
				logger.Error("Failed to insert batch of files: %v", err)
			}
			files = files[:0]
		}

	}

	if len(files) > 0 {
		// *QUESTION: No need to use mutex to write in the database?
		err := database.InsertDatas(files)
		if err != nil {
			logger.Error("Failed to insert remaining files: %v", err)
		}
	}
}

func processMetadata(logger *db.CustomLogger, dbpath string, time_start time.Time) {
	database, err := db.Connect(dbpath)
	if err != nil {
		logger.Error("Error connecting the database in processMetadata:", err)
	}
	defer database.Close()
	metadatas := make([]*db.Metadata, 0)

	randomString := "510D5B0B6245A77B40B52C60DF3E0F85480D323C184B3BF6673044E249E12B3F"
	datecreation := time_start.Format("2006-01-02 15:04:05")
	sigmd5, _ := database.GetTableMD5("data", randomString)

	metadata := &db.Metadata{
		SignatureMD5:   sigmd5,
		Challenge:      randomString,
		SO:             runtime.GOOS,
		Architecture:   runtime.GOARCH,
		DateDBCreation: datecreation,
	}
	metadatas = append(metadatas, metadata)
	err = database.InsertMetadatas(metadatas)
	if err != nil {
		logger.Info("Failed to insert batch of metadatas: %v", err)
	}
}

func isSuperUserWindows() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

func getDisksWindows() ([]string, error) {
	var disks []string

	driveLetters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for _, letter := range driveLetters {
		path := string(letter) + ":\\"
		_, err := os.Open(path)
		if err == nil {
			disks = append(disks, path)
		}
	}

	return disks, nil
}

func getFileMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashBytes := hash.Sum(nil)
	md5String := hex.EncodeToString(hashBytes)

	return md5String, nil
}

func getFileCreationTimeWindows(path string) (string, error) {
	var fileInfo os.FileInfo
	var err error

	if fileInfo, err = os.Stat(path); err != nil {
		return "", err
	}

	var creationTime time.Time
	winFileSys := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	nsec := winFileSys.CreationTime.Nanoseconds()
	creationTime = time.Unix(0, nsec)

	formattedTime := creationTime.Format("2006-01-02 15:04:05")
	return formattedTime, nil
}
