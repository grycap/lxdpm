package apilxd

import (
	"fmt"
	"net/http"
	"time"
	"encoding/json"
	"database/sql"
	"math/rand"
	"errors"
	"strconv"
)

const PLANNER_RANDOM_STATE string = "random"
const PLANNER_GATHER_STATE string = "gather"
const PLANNER_SCATTER_STATE string = "scatter"



type Planner struct {
	db 	*sql.DB
}

func (planner * Planner) State() string {
	state,_ := getPlannerStateDB(planner.db)

	return state
}

func (planner * Planner) HostToDeploy() (string,error) {
	state := planner.State()
	containersInfo,err := getContainersHostInfoDB(planner.db)
	if err != nil {
		return "",err
	}
	switch state {
	case PLANNER_GATHER_STATE:
		fmt.Println("Gather strategy")
		hostname,err := calculateGatherStrategy(containersInfo)
		if err != nil {
			return "",err
		}
		return hostname,nil
		

	case PLANNER_SCATTER_STATE:
		fmt.Println("Scatter strategy")
		hostname,err := calculateScatterStrategy(containersInfo)
		if err != nil {
			return "",err
		}
		return hostname,nil

	case PLANNER_RANDOM_STATE:
		return calculateRandomStrategy(containersInfo),nil

	default:
		return "",errors.New("Strategy not set.")
	}
	return "",errors.New("Strategy not implemented.")
}
func calculateGatherStrategy(info [][]interface{}) (string,error) {
	score := make([]float64,len(info))
	for i := 0; i < len(info); i++ {
		memory,err := doGetAvailableMemory(info[i][1].(string))
		if err != nil {
			return "",err
		}
		cores,err := doGetCores(info[i][1].(string))
		if err != nil {
			return "",err
		}
		containernumber,err := strconv.ParseInt(info[i][2].(string),10,0)
		if err != nil {
			return "",err
		}
		score[i] = float64(memory/int(containernumber))
		score[i] = score[i] + score[i]*(float64(cores/int(containernumber)))
	}
	hostIndex := findMinIndex(score)
	hostname := info[hostIndex][1].(string)
	fmt.Printf("Selecting host %v with score %v",hostname,score[hostIndex])
	return hostname,nil
}

func calculateScatterStrategy(info [][]interface{}) (string,error) {
	score := make([]float64,len(info))
	for i := 0; i < len(info); i++ {
		memory,err := doGetAvailableMemory(info[i][1].(string))
		if err != nil {
			return "",err
		}
		cores,err := doGetCores(info[i][1].(string))
		if err != nil {
			return "",err
		}
		containernumber,err := strconv.ParseInt(info[i][2].(string),10,0)
		if err != nil {
			return "",err
		}
		score[i] = float64(memory/int(containernumber))
		score[i] = score[i] + score[i]*(float64(cores/int(containernumber)))
	}
	hostIndex := findMaxIndex(score)
	hostname := info[hostIndex][1].(string)
	fmt.Printf("Selecting host %v with score %v",hostname,score[hostIndex])
	return hostname,nil
}

func calculateRandomStrategy(info [][]interface{}) string {
	rand.Seed(time.Now().Unix())
	randomInt := rand.Intn(len(info))
	hostname := info[randomInt][1].(string)
	fmt.Printf("Selecting host %v",hostname,"with random id %v",randomInt)
	return hostname
}
var plannerCmd = Command{
	name: "planner",
	get:  plannerGet,
	post: plannerPost,
	put: plannerPut,
}
func plannerGet(lx *LxdpmApi,  r *http.Request) Response {
	state,err := getPlannerStateDB(lx.db)
	if err != nil {
		return BadRequest(err)
	}
	return SyncResponse(true,state)
}


type PlannerPost struct {

	State   string          `json:"state" yaml:"state"`
}

func plannerPost(lx *LxdpmApi,  r *http.Request) Response {
	req := PlannerPost{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}

	err := putPlannerDB(lx,req.State,time.Now())
	if err != nil {
		return BadRequest(err)
	}
	return SyncResponse(true,"")
}

func plannerPut(lx *LxdpmApi,  r *http.Request) Response {
	req := PlannerPost{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}

	err := updatePlannerDB(lx,req.State,time.Now())
	if err != nil {
		return BadRequest(err)
	}
	return SyncResponse(true,"")
}

func getContainersHostInfoDB(db *sql.DB) ([][]interface{},error) {
	inargs := []interface{}{}
	outargs := []interface{}{"id","name","container_number"}
	//cash, err := lx.db.Query(`SELECT * FROM hosts`)
	result, err := dbQueryScan(db, `SELECT H.id,H.name,count(*) as container_number FROM hosts H LEFT JOIN containers C ON ( H.id = C.host_id ) GROUP BY H.id,H.name`,inargs,outargs )
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	if len(result) == 0 {
		return nil , nil
	}

	return result,nil
}

func getPlannerStateDB(db *sql.DB) (string,error) {
	inargs := []interface{}{}
	outargs := []interface{}{"state"}
	//cash, err := lx.db.Query(`SELECT * FROM hosts`)
	result, err := dbQueryScan(db, `SELECT state FROM planner`,inargs,outargs )
	if err != nil {
		fmt.Println(err)
		return "",err
	}
	if len(result) == 0 {
		return "" , nil
	}

	return result[0][0].(string) ,nil
}

func putPlannerDB(lx *LxdpmApi,state string,timestamp time.Time) error{
	q := `INSERT INTO planner (state,updated_at) VALUES (?,?)`
	_,err := dbExec(lx.db,q,state,timestamp.Format(time.RFC3339))
	return err
}

func updatePlannerDB(lx *LxdpmApi,state string,timestamp time.Time) error{
	q := `UPDATE planner SET state=?,updated_at=? WHERE id=1;`
	_,err := dbExec(lx.db,q,state,timestamp.Format(time.RFC3339))
	return err
}

func findMaxIndex(array []float64) int {
	var max float64 = 0.0
	var index int
	for i := 0; i < len(array); i++ {
		if array[i] > max {
			max = array[i]
			index = i
		}
	}
	return index
}

func findMinIndex(array []float64) int {
	var min float64 = 10000000000000000000000000.0
	var index int
	for i := 0; i < len(array); i++ {
		if array[i] < min {
			min = array[i]
			index = i
		}
	}
	return index
}  