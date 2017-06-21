package apilxd

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lxc/lxd/shared"
	"github.com/lxc/lxd/shared/api"
)

func WriteJSON(w http.ResponseWriter, body interface{}) error {
	var output io.Writer
	var captured *bytes.Buffer

	output = w
	if debug {
		captured = &bytes.Buffer{}
		output = io.MultiWriter(w, captured)
	}

	err := json.NewEncoder(output).Encode(body)

	if captured != nil {
		shared.DebugJson(captured)
	}

	return err
}

func etagHash(data interface{}) (string, error) {
	etag := sha256.New()
	err := json.NewEncoder(etag).Encode(data)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", etag.Sum(nil)), nil
}

func etagCheck(r *http.Request, data interface{}) error {
	match := r.Header.Get("If-Match")
	if match == "" {
		return nil
	}

	hash, err := etagHash(data)
	if err != nil {
		return err
	}

	if hash != match {
		return fmt.Errorf("ETag doesn't match: %s vs %s", hash, match)
	}

	return nil
}

func loadModule(module string) error {
	if shared.PathExists(fmt.Sprintf("/sys/module/%s", module)) {
		return nil
	}

	return shared.RunCommand("modprobe", module)
}

func parseAsyncResponse(input []byte) Response {
	req := api.ResponseRaw{}
	fmt.Println(input)
	fmt.Println(string(input))
	if err := json.NewDecoder(bytes.NewReader(input)).Decode(&req); err != nil {
		return BadRequest(err)
	}
	fmt.Printf("\nReq: %+v",req)


	return AsyncResponse(true,req.Metadata)

}
func saveFile(input []byte,filename string) error {
	os.Mkdir(filepath.Dir(filename), 0700)
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("cannot create file: %v", err)
	}
	defer f.Close()
	defer os.Remove(filename + ".new")

	_, err = f.Write(input)
	if err != nil {
		return fmt.Errorf("cannot write file: %v", err)
	}
	f.Close()
	return nil
}

/*
func parseFileResponse(input []byte) Response {
	req := fileResponse{}
	fmt.Println(input)
	fmt.Println(string(input))
	if err := json.NewDecoder(bytes.NewReader(input)).Decode(&req); err != nil {
		return BadRequest(err)
	}
	fmt.Printf("\nReq: %+v",req)


	return req

}*/