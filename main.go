package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	serveAddr = flag.String("serve_addr", ":80", "HTTP Port")
	modelCP   = flag.String("ckeckpoint", "", "Checkpoint file of TensorFlow model")
	modelWL   = flag.String("wordlist", "", "Word list of TensorFlow model")
)

type apiHandler struct {
}

type DataFile struct {
	Path string
	FD   *os.File
}

func generateDataFile() (f *DataFile, err error) {
	fn := uuid.New().String() + ".jpg"

	f = &DataFile{
		Path: filepath.Join("/tmp", fn),
	}

	f.FD, err = os.Create(f.Path)
	return
}

func (h *apiHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	f, err := generateDataFile()
	if err != nil {
		fmt.Printf("Fail to create file cuz %s", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer f.FD.Close()
	defer req.Body.Close()

	_, err = io.Copy(f.FD, req.Body)
	if err != nil {
		fmt.Printf("Fail to write file to %s cuz %s", f.Path, err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	var cmd *exec.Cmd
	cmd = exec.Command("bazel-bin/im2txt/run_inference", fmt.Sprintf("--checkpoint_path=%s", *modelCP),
		fmt.Sprintf("--vocab_file=%s", *modelWL), fmt.Sprintf("--input_files=%s", f.Path))
	cmd.Stderr = os.Stderr
	bins, err := cmd.Output()
	if err != nil {
		fmt.Printf("Fail to write file to %s cuz %s", f.Path, err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = resp.Write(bins)
	if err != nil {
		fmt.Printf("Fail to write response cuz %s", err)
	}
}

func main() {
	flag.Parse()

	s := &http.Server{
		Addr:    *serveAddr,
		Handler: &apiHandler{},
	}

	s.SetKeepAlivesEnabled(false)
	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
