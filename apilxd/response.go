package apilxd

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/mattn/go-sqlite3"

	"github.com/lxc/lxd/shared"
	"github.com/lxc/lxd/shared/api"
)

type Response interface {
	Render(w http.ResponseWriter) error
	String() string
}

// Sync response
type syncResponse struct {
	success  bool
	etag     interface{}
	metadata interface{}
	location string
	headers  map[string]string
}

func (r *syncResponse) Render(w http.ResponseWriter) error {
	// Set an appropriate ETag header
	if r.etag != nil {
		etag, err := etagHash(r.etag)
		if err == nil {
			w.Header().Set("ETag", etag)
		}
	}

	// Prepare the JSON response
	status := api.Success
	if !r.success {
		status = api.Failure
	}

	if r.headers != nil {
		for h, v := range r.headers {
			w.Header().Set(h, v)
		}
	}

	if r.location != "" {
		w.Header().Set("Location", r.location)
		w.WriteHeader(201)
	}

	resp := api.ResponseRaw{
		Response: api.Response{
			Type:       api.SyncResponse,
			Status:     status.String(),
			StatusCode: int(status)},
		Metadata: r.metadata,
	}

	return WriteJSON(w, resp)
}

func (r *syncResponse) String() string {
	if r.success {
		return "success"
	}

	return "failure"
}

func SyncResponse(success bool, metadata interface{}) Response {
	return &syncResponse{success: success, metadata: metadata}
}

func SyncResponseETag(success bool, metadata interface{}, etag interface{}) Response {
	return &syncResponse{success: success, metadata: metadata, etag: etag}
}

func SyncResponseLocation(success bool, metadata interface{}, location string) Response {
	return &syncResponse{success: success, metadata: metadata, location: location}
}

func SyncResponseHeaders(success bool, metadata interface{}, headers map[string]string) Response {
	return &syncResponse{success: success, metadata: metadata, headers: headers}
}

var EmptySyncResponse = &syncResponse{success: true, metadata: make(map[string]interface{})}

// async response
type asyncResponse struct {
	success  bool
	etag     interface{}
	metadata interface{}
	location string
	headers  map[string]string
}

func (r *asyncResponse) Render(w http.ResponseWriter) error {
	// Set an appropriate ETag header
	if r.etag != nil {
		etag, err := etagHash(r.etag)
		if err == nil {
			w.Header().Set("ETag", etag)
		}
	}

	// Prepare the JSON response
	status := api.Success
	if !r.success {
		status = api.Failure
	}

	if r.headers != nil {
		for h, v := range r.headers {
			w.Header().Set(h, v)
		}
	}

	if r.location != "" {
		w.Header().Set("Location", r.location)
		w.WriteHeader(201)
	}

	resp := api.ResponseRaw{
		Response: api.Response{
			Type:       api.AsyncResponse,
			Status:     status.String(),
			StatusCode: int(status)},
		Metadata: r.metadata,
	}

	return WriteJSON(w, resp)
}

func (r *asyncResponse) String() string {
	if r.success {
		return "success"
	}

	return "failure"
}

func AsyncResponse(success bool, metadata interface{}) Response {
	return &asyncResponse{success: success, metadata: metadata}
}

func AsyncResponseETag(success bool, metadata interface{}, etag interface{}) Response {
	return &asyncResponse{success: success, metadata: metadata, etag: etag}
}

func AsyncResponseLocation(success bool, metadata interface{}, location string) Response {
	return &asyncResponse{success: success, metadata: metadata, location: location}
}

func AsyncResponseHeaders(success bool, metadata interface{}, headers map[string]string) Response {
	return &asyncResponse{success: success, metadata: metadata, headers: headers}
}

var EmptyasyncResponse = &asyncResponse{success: true, metadata: make(map[string]interface{})}

// File transfer response
type fileResponseEntry struct {
	identifier string
	path       string
	filename   string
	buffer     []byte /* either a path or a buffer must be provided */
}

type fileResponse struct {
	req              *http.Request
	files            []fileResponseEntry
	headers          map[string]string
	removeAfterServe bool
}

func (r *fileResponse) Render(w http.ResponseWriter) error {
	if r.headers != nil {
		for k, v := range r.headers {
			w.Header().Set(k, v)
		}
	}

	// No file, well, it's easy then
	if len(r.files) == 0 {
		return nil
	}

	// For a single file, return it inline
	if len(r.files) == 1 {
		var rs io.ReadSeeker
		var mt time.Time
		var sz int64

		if r.files[0].path == "" {
			rs = bytes.NewReader(r.files[0].buffer)
			mt = time.Now()
			sz = int64(len(r.files[0].buffer))
		} else {
			f, err := os.Open(r.files[0].path)
			if err != nil {
				return err
			}
			defer f.Close()

			fi, err := f.Stat()
			if err != nil {
				return err
			}

			mt = fi.ModTime()
			sz = fi.Size()
			rs = f
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", sz))
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline;filename=%s", r.files[0].filename))

		http.ServeContent(w, r.req, r.files[0].filename, mt, rs)
		if r.files[0].path != "" && r.removeAfterServe {
			err := os.Remove(r.files[0].path)
			if err != nil {
				return err
			}
		}

		return nil
	}

	// Now the complex multipart answer
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)

	for _, entry := range r.files {
		var rd io.Reader
		if entry.path != "" {
			fd, err := os.Open(entry.path)
			if err != nil {
				return err
			}
			defer fd.Close()
			rd = fd
		} else {
			rd = bytes.NewReader(entry.buffer)
		}

		fw, err := mw.CreateFormFile(entry.identifier, entry.filename)
		if err != nil {
			return err
		}

		_, err = io.Copy(fw, rd)
		if err != nil {
			return err
		}
	}
	mw.Close()

	w.Header().Set("Content-Type", mw.FormDataContentType())
	w.Header().Set("Content-Length", fmt.Sprintf("%d", body.Len()))

	_, err := io.Copy(w, body)
	return err
}

func (r *fileResponse) String() string {
	return fmt.Sprintf("%d files", len(r.files))
}

func FileResponse(r *http.Request, files []fileResponseEntry, headers map[string]string, removeAfterServe bool) Response {
	return &fileResponse{r, files, headers, removeAfterServe}
}

// Operation response
type operationResponse struct {
	op *api.Operation
	url string
}

func (r *operationResponse) Render(w http.ResponseWriter) error {

	url, md := r.url,r.op

	body := api.ResponseRaw{
		Response: api.Response{
			Type:       api.AsyncResponse,
			Status:     api.OperationCreated.String(),
			StatusCode: int(api.OperationCreated),
			Operation:  url,
		},
		Metadata: md,
	}

	w.Header().Set("Location", url)
	w.WriteHeader(202)

	return WriteJSON(w, body)
}

func (r *operationResponse) String() string {
	md := r.op

	return md.ID
}

func OperationResponse(url string,op *api.Operation) Response {
	return &operationResponse{op,url}
}

// Error response
type errorResponse struct {
	code int
	msg  string
}

func (r *errorResponse) String() string {
	return r.msg
}

func (r *errorResponse) Render(w http.ResponseWriter) error {
	var output io.Writer

	buf := &bytes.Buffer{}
	output = buf
	var captured *bytes.Buffer
	if debug {
		captured = &bytes.Buffer{}
		output = io.MultiWriter(buf, captured)
	}

	err := json.NewEncoder(output).Encode(shared.Jmap{"type": api.ErrorResponse, "error": r.msg, "error_code": r.code})

	if err != nil {
		return err
	}

	if debug {
		shared.DebugJson(captured)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(r.code)
	fmt.Fprintln(w, buf.String())

	return nil
}

/* Some standard responses */
var NotImplemented = &errorResponse{http.StatusNotImplemented, "not implemented"}
var NotFound = &errorResponse{http.StatusNotFound, "not found"}
var Forbidden = &errorResponse{http.StatusForbidden, "not authorized"}
var Conflict = &errorResponse{http.StatusConflict, "already exists"}

func BadRequest(err error) Response {
	return &errorResponse{http.StatusBadRequest, err.Error()}
}

func InternalError(err error) Response {
	return &errorResponse{http.StatusInternalServerError, err.Error()}
}

func PreconditionFailed(err error) Response {
	return &errorResponse{http.StatusPreconditionFailed, err.Error()}
}

/*
 * SmartError returns the right error message based on err.
 */
func SmartError(err error) Response {
	switch err {
	case nil:
		return EmptySyncResponse
	case os.ErrNotExist:
		return NotFound
	case sql.ErrNoRows:
		return NotFound
	case os.ErrPermission:
		return Forbidden
	case sqlite3.ErrConstraintUnique:
		return Conflict
	default:
		return InternalError(err)
	}
}


type LxdResponseRaw struct {
	LxdResponse `yaml:",inline"`

	Metadata interface{} `json:"metadata" yaml:"metadata"`
}
// Response represents a LXD operation
type LxdResponse struct {
	Type lxdResponseType `json:"type" yaml:"type"`

	// Valid only for Sync responses
	Status     string `json:"status" yaml:"status"`
	StatusCode int    `json:"status_code" yaml:"status_code"`

	// Valid only for Async responses
	Operation 	string `json:"operation" yaml:"operation"`
	Task 		string `json:"task" yaml:"task"`

	// Valid only for Error responses
	Code  int    `json:"error_code" yaml:"error_code"`
	Error string `json:"error" yaml:"error"`

	// Valid for Sync and Error responses
	Metadata json.RawMessage `json:"metadata" yaml:"metadata"`
}

// MetadataAsMap parses the Response metadata into a map
func (r *LxdResponse) MetadataAsMap() (map[string]interface{}, error) {
	ret := map[string]interface{}{}
	err := r.MetadataAsStruct(&ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// MetadataAsOperation turns the Response metadata into an Operation
func (r *LxdResponse) MetadataAsOperation() (*api.Operation, error) {
	op := api.Operation{}
	err := r.MetadataAsStruct(&op)
	if err != nil {
		return nil, err
	}

	return &op, nil
}

// MetadataAsStringSlice parses the Response metadata into a slice of string
func (r *LxdResponse) MetadataAsStringSlice() ([]string, error) {
	sl := []string{}
	err := r.MetadataAsStruct(&sl)
	if err != nil {
		return nil, err
	}

	return sl, nil
}

// MetadataAsStruct parses the Response metadata into a provided struct
func (r *LxdResponse) MetadataAsStruct(target interface{}) error {
	if err := json.Unmarshal(r.Metadata, &target); err != nil {
		return err
	}

	return nil
}

// ResponseType represents a valid LXD response type
type lxdResponseType string

// LXD response types
const (
	lxdSyncResponse  lxdResponseType = "sync"
	lxdAsyncResponse lxdResponseType = "async"
	lxdErrorResponse lxdResponseType = "error"
)

// Task response
type taskResponse struct {
	tk *task
}

func (r *taskResponse) Render(w http.ResponseWriter) error {
	_, err := r.tk.Run()
	if err != nil {
		return err
	}

	url, md, err := r.tk.Render()
	if err != nil {
		return err
	}

	body := LxdResponseRaw{
		LxdResponse: LxdResponse{
			Type:       lxdAsyncResponse,
			Status:     TaskCreated.String(),
			StatusCode: int(TaskCreated),
			Task:  url,
		},
		Metadata: md,
	}

	w.Header().Set("Location", url)
	w.WriteHeader(202)

	return WriteJSON(w, body)
}

func (r *taskResponse) String() string {
	_, md, err := r.tk.Render()
	if err != nil {
		return fmt.Sprintf("error: %s", err)
	}

	return md.ID
}

func TaskResponse(tk *task) Response {
	return &taskResponse{tk}
}