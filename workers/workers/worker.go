package workers

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
)

// Matchar din Job-struktur på serversidan
type MsgFromWorker struct {
	RespType string
	Res      Result
}
type Result struct {
	Job      Job    `json:"job"`      // jobbet
	Result   string `json:"result"`   // valfritt: kan vara text, hash etc.
	Duration int64  `json:"duration"` // hur lång tid jobbet tog i ms
	Success  bool   `json:"success"`
}
type Job struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Duration int    `json:"duration,omitempty"`
	Input    string `json:"input,omitempty"`
	Attempts int    `json:"attempts"`
}

func Start(wsURL string) {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("WebSocket connection failed:", err)
	}
	defer conn.Close()

	log.Println("Connected to server!")
	success := true

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading:", err)
			break
		}

		var job Job
		err = json.Unmarshal(message, &job)
		if err != nil {
			log.Printf("Text message from server: %s\n", string(message))
			continue
		}

		log.Printf("Received job: %+v\n", job)

		start := time.Now()
		var resultText string

		// Jobbtyper
		switch job.Type {
		case "execurl":
			resultText, err = handleExecURLJob(job)
			if err != nil {
				success = false
			}

		case "program_zip":
			resultText, err = handleProgramZip(job)
			if err != nil {
				success = false
			}

		case "sleep":
			time.Sleep(time.Duration(job.Duration) * time.Second)
			resultText = fmt.Sprintf("Slept for %d seconds", job.Duration)

		case "hash":
			// Simulera en hash
			resultText = fmt.Sprintf("fakehash(%s)", job.Input)

		default:
			resultText = "Unknown job type"
		}

		elapsed := time.Since(start).Milliseconds()

		// Skicka tillbaka resultatet
		result := Result{
			Job:      job,
			Result:   resultText,
			Duration: elapsed,
			Success:  success,
		}

		msgFromWorker := MsgFromWorker{
			RespType: "result",
			Res:      result,
		}

		msg, _ := json.Marshal(msgFromWorker)
		conn.WriteMessage(websocket.TextMessage, msg)
		log.Println("Sent result:", string(msg))
	}
}

// TODO: kolla på och förstå
func handleExecURLJob(job Job) (string, error) {
	filename := fmt.Sprintf("job_%s_exec.sh", job.Id)

	// Ladda ner från URL
	resp, err := http.Get(job.Input)
	if err != nil {
		return fmt.Sprintf("Download error: %s", err), err
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return fmt.Sprintf("File creation error: %s", err), err
	}
	io.Copy(out, resp.Body)
	out.Close()

	os.Chmod(filename, 0o755)

	// Kör filen
	cmd := exec.Command("./" + filename)
	output, err := cmd.CombinedOutput()
	os.Remove(filename)

	if err != nil {
		return fmt.Sprintf("Exec error: %s\nOutput: %s", err, output), err
	}

	return string(output), nil
}

func handleProgramZip(job Job) (string, error) {
	os.MkdirAll("jobs", 0o755)
	workDir := fmt.Sprintf("jobs/job-%s", job.Id)

	// Ta bort ev. gammal katalog innan vi skapar ny
	os.RemoveAll(workDir)

	// 1. Dekoda Base64
	zipData, err := base64.StdEncoding.DecodeString(job.Input)
	if err != nil {
		return "failed to decode base64 input", err
	}

	// 2. Skapa en arbetsmapp för detta jobb
	err = os.Mkdir(workDir, 0o755)
	if err != nil {
		return "failed to create work directory", err
	}

	zipPath := filepath.Join(workDir, "source.zip")

	// 3. Spara zipfilen till arbetsmappen
	err = os.WriteFile(zipPath, zipData, 0o644)
	if err != nil {
		return "failed to write zip file", err
	}

	// 4. Extrahera zip-filen
	err = unzip(zipPath, workDir)
	if err != nil {
		return "failed to unzip project", err
	}

	// 4b. Gå in i första undermappen i workDir
	entries, err := os.ReadDir(workDir)
	if err != nil {
		return "failed to read job directory", err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			workDir = filepath.Join(workDir, entry.Name())
			break
		}
	}
	fmt.Println("workDir:", workDir)

	// 5. Kör "make run" inne i arbetsmappen
	cmd := exec.Command("make", "run")
	cmd.Dir = workDir // kör kommandot i arbetsmappen

	// Fånga både stdout och stderr
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	// Lägg till timeout (max 15 sekunder)
	done := make(chan error)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Sprintf("make run error: %s\nOutput: %s", err, output.String()), err
		}
	case <-time.After(15 * time.Second):
		cmd.Process.Kill()
		return "make run timed out", fmt.Errorf("timeout")
	}

	return output.String(), nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
