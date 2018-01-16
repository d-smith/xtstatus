package workitems

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"encoding/json"
	"html/template"
)

var templates = template.Must(template.ParseFiles("./workitems/status.html"))
var awsSession = session.Must(session.NewSession())
var modelInstanceTable = os.Getenv("MODEL_INSTANCE_TABLE")
var ddb = dynamodb.New(awsSession)

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

func respondError(w http.ResponseWriter, status int, err error) {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := &ErrorResponse{Errors: make([]string, 0, 1)}
	if err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}

	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

func WorkItemsHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	workitem := vars["workitem"]

	workItemStates, err := getWorkItemStates(workitem)
	if err != nil {
		respondError(rw, http.StatusInternalServerError, err)
		return
	}

	statusDetail := workflowStatusFromStates(workItemStates)

	if err = templates.ExecuteTemplate(rw, "status.html", statusDetail); err != nil {
		respondError(rw, http.StatusInternalServerError, err)
	}
}



func getWorkItemStates(workItemNo string)(map[string]string, error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":iid": {
				S: aws.String(workItemNo),
			},
		},
		KeyConditionExpression: aws.String("instanceId = :iid"),
		TableName:              aws.String(modelInstanceTable),
	}

	qout, err := ddb.Query(input)
	if err != nil {
		return nil, err
	}

	workitemTxns := make(map[string]string)
	items := qout.Items
	for _, item := range items {
		workitemTxns[*item["state"].S] = *item["txnId"].S
	}

	return workitemTxns, nil
}

type workflowStepState struct {
	Step string
	Status string
}

func workflowStatusFromStates(workflowStates map[string]string) []workflowStepState {
	states := []workflowStepState{}

	steps := []string{"Assign Project Mgr","Review Tech Requirements","Configure Solution","Send Email","Task Complete"}
	for _, step := range steps {
		txn := workflowStates[step]
		switch txn {
		case "":
			states = append(states, workflowStepState{
				Step: step,
				Status: "",
			})
		default:
			states = append(states, workflowStepState{
				Step: step,
				Status: "is-complete",
			})
		}
	}

	return states
}