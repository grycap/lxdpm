package apilxd

import (
	"fmt"
	"net/http"
//	"runtime"
//	"strings"
	"sync"
	"time"

//	"github.com/gorilla/mux"
	"github.com/pborman/uuid"

	"github.com/lxc/lxd/shared"
	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/version"
)

type StatusCode int

const (
	TaskCreated StatusCode = 100
)

func (o StatusCode) String() string {
	return map[StatusCode]string{
		TaskCreated: "Task created",
	}[o]
}

var tasksLock sync.Mutex
var tasks map[string]*task = make(map[string]*task)

type taskClass int

const (
	taskClassOp      taskClass = 1
	taskClassWebsocket taskClass = 2
	taskClassToken     taskClass = 3
)

func (t taskClass) String() string {
	return map[taskClass]string{
		taskClassOp:      "task",
		taskClassWebsocket: "websocket",
		taskClassToken:     "token",
	}[t]
}

type task struct {
	id        string
	class 	  taskClass
	createdAt time.Time
	updatedAt time.Time
	status    api.StatusCode
	url       string
	resources map[string][]string
	metadata  map[string]interface{}
	err       string
	readonly  bool

	// Those functions are called at various points in the task lifecycle
	onRun     func(*task) error
	onCancel  func(*task) error
	onConnect func(*task, *http.Request, http.ResponseWriter) error

	// Channels used for error reporting and state tracking of background actions
	chanDone chan error

	// Locking for concurent access to the task
	lock sync.Mutex
}

func (tk *task) done() {
	if tk.readonly {
		return
	}

	tk.lock.Lock()
	tk.readonly = true
	tk.onRun = nil
	tk.onCancel = nil
	tk.onConnect = nil
	close(tk.chanDone)
	tk.lock.Unlock()

	time.AfterFunc(time.Second*5, func() {
		tasksLock.Lock()
		_, ok := tasks[tk.id]
		if !ok {
			tasksLock.Unlock()
			return
		}
		
		delete(tasks, tk.id)
		tasksLock.Unlock()

		/*
		 * When we create a new lxc.Container, it adds a finalizer (via
		 * SetFinalizer) that frees the struct. However, it sometimes
		 * takes the go GC a while to actually free the struct,
		 * presumably since it is a small amount of memory.
		 * Unfortunately, the struct also keeps the log fd open, so if
		 * we leave too many of these around, we end up running out of
		 * fds. So, let's explicitly do a GC to collect these at the
		 * end of each request.
		 */
		//runtime.GC()
	})
}

func (tk *task) Run() (chan error, error) {
	if tk.status != api.Pending {
		return nil, fmt.Errorf("Only pending tasks can be started")
	}
	fmt.Println("VIVEEEEE!")
	chanRun := make(chan error, 1)

	tk.lock.Lock()
	tk.status = api.Running

	if tk.onRun != nil {
		go func(tk *task, chanRun chan error) {
			err := tk.onRun(tk)
			if err != nil {
				tk.lock.Lock()
				tk.status = api.Failure
				tk.err = SmartError(err).String()
				tk.lock.Unlock()
				tk.done()
				chanRun <- err

				shared.LogDebugf("Failure for %s task: %s: %s", tk.class.String(), tk.id, err)

				//_, md, _ := tk.Render()
				//eventSend("task", md)
				return
			}

			tk.lock.Lock()
			tk.status = api.Success
			tk.lock.Unlock()
			tk.done()
			chanRun <- nil

			tk.lock.Lock()
			shared.LogDebugf("Success for %s task: %s", tk.class.String(), tk.id)
			//_, md, _ := tk.Render()
			//eventSend("task", md)
			tk.lock.Unlock()
		}(tk, chanRun)
	}
	tk.lock.Unlock()

	shared.LogDebugf("Started %s task: %s", tk.class.String(), tk.id)
	//_, md, _ := tk.Render()
	//eventSend("task", md)

	return chanRun, nil
}

func (tk *task) Cancel() (chan error, error) {
	if tk.status != api.Running {
		return nil, fmt.Errorf("Only running tasks can be cancelled")
	}

	if !tk.mayCancel() {
		return nil, fmt.Errorf("This task can't be cancelled")
	}

	chanCancel := make(chan error, 1)

	tk.lock.Lock()
	oldStatus := tk.status
	tk.status = api.Cancelling
	tk.lock.Unlock()

	if tk.onCancel != nil {
		go func(tk *task, oldStatus api.StatusCode, chanCancel chan error) {
			err := tk.onCancel(tk)
			if err != nil {
				tk.lock.Lock()
				tk.status = oldStatus
				tk.lock.Unlock()
				chanCancel <- err

				shared.LogDebugf("Failed to cancel %s task: %s: %s", tk.class.String(), tk.id, err)
				//_, md, _ := tk.Render()
				//eventSend("task", md)
				return
			}

			tk.lock.Lock()
			tk.status = api.Cancelled
			tk.lock.Unlock()
			tk.done()
			chanCancel <- nil

			shared.LogDebugf("Cancelled %s task: %s", tk.class.String(), tk.id)
			//_, md, _ := tk.Render()
			//eventSend("task", md)
		}(tk, oldStatus, chanCancel)
	}

	shared.LogDebugf("Cancelling %s task: %s", tk.class.String(), tk.id)
	//_, md, _ := tk.Render()
	//eventSend("task", md)

	if tk.onCancel == nil {
		tk.lock.Lock()
		tk.status = api.Cancelled
		tk.lock.Unlock()
		tk.done()
		chanCancel <- nil
	}

	shared.LogDebugf("Cancelled %s task: %s", tk.class.String(), tk.id)
	//_, md, _ = tk.Render()
	//eventSend("task", md)

	return chanCancel, nil
}

func (tk *task) Connect(r *http.Request, w http.ResponseWriter) (chan error, error) {
	if tk.class != taskClassWebsocket {
		return nil, fmt.Errorf("Only websocket tasks can be connected")
	}

	if tk.status != api.Running {
		return nil, fmt.Errorf("Only running tasks can be connected")
	}

	chanConnect := make(chan error, 1)

	tk.lock.Lock()

	go func(tk *task, chanConnect chan error) {
		err := tk.onConnect(tk, r, w)
		if err != nil {
			chanConnect <- err

			shared.LogDebugf("Failed to handle %s task: %s: %s", tk.class.String(), tk.id, err)
			return
		}

		chanConnect <- nil

		shared.LogDebugf("Handled %s task: %s", tk.class.String(), tk.id)
	}(tk, chanConnect)
	tk.lock.Unlock()

	shared.LogDebugf("Connected %s task: %s", tk.class.String(), tk.id)

	return chanConnect, nil
}

func (tk *task) mayCancel() bool {
	return tk.onCancel != nil || tk.class == taskClassToken
}

func (tk *task) Render() (string, *api.Operation, error) {
	// Setup the resource URLs
	resources := tk.resources
	if resources != nil {
		tmpResources := make(map[string][]string)
		for key, value := range resources {
			var values []string
			for _, c := range value {
				values = append(values, fmt.Sprintf("/%s/%s/%s", version.APIVersion, key, c))
			}
			tmpResources[key] = values
		}
		resources = tmpResources
	}

	return tk.url, &api.Operation{
		ID:         tk.id,
		Class:      tk.class.String(),
		CreatedAt:  tk.createdAt,
		UpdatedAt:  tk.updatedAt,
		Status:     tk.status.String(),
		StatusCode: tk.status,
		Resources:  resources,
		Metadata:   tk.metadata,
		MayCancel:  tk.mayCancel(),
		Err:        tk.err,
	}, nil
}

func (tk *task) WaitFinal(timeout int) (bool, error) {
	// Check current state
	if tk.status.IsFinal() {
		return true, nil
	}

	// Wait indefinitely
	if timeout == -1 {
		for {
			<-tk.chanDone
			return true, nil
		}
	}

	// Wait until timeout
	if timeout > 0 {
		timer := time.NewTimer(time.Duration(timeout) * time.Second)
		for {
			select {
			case <-tk.chanDone:
				return false, nil

			case <-timer.C:
				return false, nil
			}
		}
	}

	return false, nil
}

func (tk *task) UpdateResources(tkResources map[string][]string) error {
	if tk.status != api.Pending && tk.status != api.Running {
		return fmt.Errorf("Only pending or running operations can be updated")
	}

	if tk.readonly {
		return fmt.Errorf("Read-only operations can't be updated")
	}

	tk.lock.Lock()
	tk.updatedAt = time.Now()
	tk.resources = tkResources
	tk.lock.Unlock()

	shared.LogDebugf("Updated resources for %s operation: %s", tk.class.String(), tk.id)
	//_, md, _ := tk.Render()
	//eventSend("task", md)

	return nil
}

func (tk *task) UpdateMetadata(tkMetadata interface{}) error {
	if tk.status != api.Pending && tk.status != api.Running {
		return fmt.Errorf("Only pending or running tasks can be updated")
	}

	if tk.readonly {
		return fmt.Errorf("Read-only tasks can't be updated")
	}

	newMetadata, err := shared.ParseMetadata(tkMetadata)
	if err != nil {
		return err
	}

	tk.lock.Lock()
	tk.updatedAt = time.Now()
	tk.metadata = newMetadata
	tk.lock.Unlock()

	shared.LogDebugf("Updated metadata for %s task: %s", tk.class.String(), tk.id)
	//_, md, _ := tk.Render()
	//eventSend("task", md)

	return nil
}

func taskCreate(tkClass taskClass, tkResources map[string][]string, tkMetadata interface{},
	onRun func(*task) error,
	onCancel func(*task) error,
	onConnect func(*task, *http.Request, http.ResponseWriter) error) (*task, error) {

	// Main attributes
	tk := task{}
	tk.id = uuid.NewRandom().String()
	tk.class = tkClass
	tk.createdAt = time.Now()
	tk.updatedAt = tk.createdAt
	tk.status = api.Pending
	tk.url = fmt.Sprintf("/%s/operations/%s", version.APIVersion, tk.id)
	tk.resources = tkResources
	tk.chanDone = make(chan error)

	newMetadata, err := shared.ParseMetadata(tkMetadata)
	if err != nil {
		return nil, err
	}
	tk.metadata = newMetadata

	// Callback functions
	tk.onRun = onRun
	tk.onCancel = onCancel
	tk.onConnect = onConnect

	// Sanity check
	if tk.class != taskClassWebsocket && tk.onConnect != nil {
		return nil, fmt.Errorf("Only websocket tasks can have a Connect hook")
	}

	if tk.class == taskClassWebsocket && tk.onConnect == nil {
		return nil, fmt.Errorf("Websocket tasks must have a Connect hook")
	}

	if tk.class == taskClassToken && tk.onRun != nil {
		return nil, fmt.Errorf("Token tasks can't have a Run hook")
	}

	if tk.class == taskClassToken && tk.onCancel != nil {
		return nil, fmt.Errorf("Token tasks can't have a Cancel hook")
	}

	tasksLock.Lock()
	tasks[tk.id] = &tk
	tasksLock.Unlock()

	shared.LogDebugf("New %s task: %s", tk.class.String(), tk.id)
	//_, md, _ := tk.Render()
	//eventSend("task", md)

	return &tk, nil
}

func taskGet(id string) (*task, error) {
	tasksLock.Lock()
	tk, ok := tasks[id]
	tasksLock.Unlock()

	if !ok {
		return nil, fmt.Errorf("Task '%s' doesn't exist", id)
	}

	return tk, nil
}