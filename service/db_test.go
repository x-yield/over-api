package service

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/go-pg/pg"
	"github.com/stretchr/testify/assert"

	"github.com/x-yield/over-api/internal/config"
	"github.com/x-yield/over-api/internal/models"
	"github.com/x-yield/over-api/pkg/overload-service"
	"github.com/x-yield/over-api/tools"
)

var overloadService *OverloadService
var jobId int32
var aggrId int32
var jobCreate overload.CreateJobRequest
var aggregateCreate overload.CreateAggregatesRequest
var jobUpdate overload.UpdateJobRequest

//Creates a connection to staging database
func init() {
	// create env provider
	envProvider := env.NewProvider("OVERLOAD_API")

	// create config client
	// "gitlab.ozon.ru/platform/realtime-config-go/client"
	client, _ := client.New("", "", envProvider)

	// waiting config ready
	<-client.Ready(context.TODO())
	configClient := config.NewClient(client)
	configClient.GetValue(context.Background(), config.DbAddr)

	db := pg.Connect(&pg.Options{
		User:     "loadtest_user",
		Addr:     "loadtestdb432z20.h.o3.ru:6432",
		Database: "loadtest_staging",
		Password: "heeMahw7vienoobi",
	})
	stageConfig := map[string]string{
		"host":     "overload.o3.ru",
		"port":     "8086",
		"database": "yandextank_stg",
		"username": "overload_stg",
		"password": "4t3CCfqF",
	}
	influxdb := tools.NewCustomInfluxConnector(stageConfig)
	s3 := tools.NewS3Service()
	overloadService = NewOverloadService(db, influxdb, s3)
	ReturnExpectedValues()
}

//Fills some data for test
func ReturnExpectedValues() {
	var collections []*overload.Collection
	collections = append(collections, &overload.Collection{
		Env:     "stg",
		Project: "test-project",
		Ref:     "master",
		Name:    "test",
		Author:  "autotest",
		Service: "test-service",
	}, &overload.Collection{
		Env:     "prod",
		Project: "test-project",
		Ref:     "master",
		Name:    "test",
		Author:  "autotest",
		Service: "test-service",
	})
	jobCreate = overload.CreateJobRequest{
		TestStart:          1546214400,
		TestStop:           0,
		Config:             "empty",
		Author:             "autotest",
		RegressionId:       "",
		Description:        "This is a test record",
		Tank:               "no tank used",
		Target:             "test.ru",
		EnvironmentDetails: "test details",
		Collections:        collections,
	}
	aggregateCreate = overload.CreateAggregatesRequest{
		Label:        "__TEST__",
		Q50:          20,
		Q75:          20,
		Q80:          40,
		Q85:          42,
		Q90:          45,
		Q95:          45,
		Q98:          47,
		Q99:          50,
		Q100:         50,
		Avg:          36,
		ResponseCode: "200",
		OkCount:      100,
		ErrCount:     2,
		NetRecv:      0,
		NetSend:      0,
		JobId:        jobId,
	}
	jobUpdate = overload.UpdateJobRequest{
		Id:           jobId,
		TestStop:     1546224400,
		Config:       "full",
		Description:  "This is an updated record",
		Status:       "finished autotest",
		AutostopTime: 1546224300,
	}
}

func TestDb(t *testing.T) {
	t.Run("CreateJob", func(t *testing.T) {
		actual, err := overloadService.CreateJob(&jobCreate)
		if err != nil {
			t.Fatalf("Couldn't create a job %s", err)
		}
		assert.NotEmpty(t, actual.Id, "Id of created job is empty, %d")
		jobId = actual.Id
		ReturnExpectedValues()
		t.Logf("Job created for test: %d", jobId)
	})
	t.Run("GetJob", func(t *testing.T) {
		expected := jobCreate
		actual, err := overloadService.GetJob(&overload.JobRequest{Id: jobId})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, actual.Job.Id, jobId, "Job from get job has another id, than created one: %d, %d")
		assert.Equal(t, expected.TestStart, actual.Job.TestStart, "Job test start is not equal to created %d, %d")
		assert.Equal(t, expected.TestStop, actual.Job.TestStop, "Job test stop is not equal to created %d, %d")
		assert.Equal(t, expected.Author, actual.Job.Author, "Job author is not equal to created %s, %s")
		assert.Equal(t, expected.Config, actual.Job.Config, "Job config is not equal to created %s, %s")
		assert.Equal(t, expected.Description, actual.Job.Description, "Job description is not equal to created %s, %s")
		assert.Equal(t, expected.Tank, actual.Job.Tank, "Job tank is not equal to created %s, %s")
		assert.Equal(t, expected.Target, actual.Job.Target, "Job target is not equal to created %s, %s")
		assert.Equal(t, expected.RegressionId, actual.Job.RegressionId, "Job regression id is not equal to created %s, %s")
		assert.Equal(t, expected.EnvironmentDetails, actual.Job.EnvironmentDetails, "Job environment details are not equal to created %s, %s")
		assert.Empty(t, actual.Job.AutostopTime, "Autostop time isn't empty %d")
		//assert.Len(t, actual.Job.Collections, 2, "Get Job doesn't return collections")
	})
	t.Run("GetLastJobId", func(t *testing.T) {
		actual, err := overloadService.GetLastJobId()
		if err != nil {
			t.Fatalf("Couldn't get last job id %s", err)
		}
		assert.Equal(t, jobId, actual, "Last job has another id, than created one: %d, %d")
	})
	t.Run("CreateAggregates", func(t *testing.T) {
		actual, err := overloadService.CreateAggregates(&aggregateCreate)
		if err != nil {
			t.Fatalf("Couldn't create aggregates %s", err)
		}
		aggrId = actual.Id
		assert.NotEmpty(t, actual.Id, "There's no aggregates id for this job %d")
		assert.NotEqual(t, jobId, actual.Id, "These aggregates belong to another job %d, %d")
	})
	t.Run("GetAggregates", func(t *testing.T) {
		actual, err := overloadService.GetAggregates(&overload.GetAggregatesRequest{TestId: jobId})
		if err != nil {
			t.Fatalf("Couldn't get aggregates %s", err)
		}
		assert.NotEmpty(t, actual.Aggregates[0].Id, "Id of created aggregate is empty, %d")
		assert.Equal(t, jobId, actual.Aggregates[0].JobId, "Aggregate JobId is not equal to created Job id, %d, %d")
	})
	t.Run("UpdateJob", func(t *testing.T) {
		expected := jobUpdate
		_, err := overloadService.UpdateJob(&expected)
		if err != nil {
			t.Fatalf("Couldn't update job %s", err)
		}
		actual, err := overloadService.GetJob(&overload.JobRequest{Id: jobId})
		assert.Equal(t, actual.Job.Id, jobId, "Job from get job has another id, than created one: %d, %d")
		assert.NotEqual(t, expected.TestStart, actual.Job.TestStart, "Job test start was overwritten %d, %d")
		assert.Equal(t, expected.TestStop, actual.Job.TestStop, "Job test stop is not equal to created %d, %d")
		assert.NotEqual(t, expected.Author, actual.Job.Author, "Job author was overwritten %s, %s")
		assert.Equal(t, expected.Config, actual.Job.Config, "Job config is not equal to created %s, %s")
		assert.Equal(t, expected.Description, actual.Job.Description, "Job description is not equal to created %s, %s")
		assert.NotEqual(t, expected.Tank, actual.Job.Tank, "Job tank was overwritten %s, %s")
		assert.NotEqual(t, expected.Target, actual.Job.Target, "Job target was overwritten %s, %s")
		assert.NotEqual(t, expected.RegressionId, actual.Job.Id, "Job regression id was overwritten %s, %s")
		assert.NotEqual(t, expected.EnvironmentDetails, actual.Job.EnvironmentDetails, "Job environment was overwritten%s, %s")
	})
	t.Run("GetLastJobs", func(t *testing.T) {
		rand.Seed(time.Now().UnixNano())
		var pages []string
		var limits []string
		limitNumber := rand.Intn(150)
		pages = append(pages, strconv.Itoa(rand.Intn(50)))
		limits = append(limits, strconv.Itoa(limitNumber))
		actual, err := overloadService.GetLastJobs(&overload.LastJobsRequest{Page: pages, Limit: limits})
		if err != nil {
			t.Fatalf("Couldn't get last jobs %s", err)
		}
		assert.Len(t, actual.Jobs, limitNumber, "Returns less/more jobs than asked %d, %d")
	})
	t.Run("GetFilteredLastJobs", func(t *testing.T) {
		var (
			pages   []string
			limits  []string
			authors []string
			targets []string
		)
		pages = append(pages, "1")
		limits = append(limits, "10")
		authors = append(authors, "kzlenko")
		targets = append(targets, "api.ozon.ru:80")
		actual, err := overloadService.GetLastJobs(
			&overload.LastJobsRequest{Page: pages, Limit: limits, Author: authors, Target: targets})
		if err != nil {
			t.Fatalf("Couldn't get last jobs %s", err)
		}
		assert.Len(t, actual.Jobs, 10, "Returns less/more jobs than asked %d, %d")
		for i := 0; i < len(actual.Jobs); i++ {
			assert.True(t, actual.Jobs[i].Author == "kzlenko" || actual.Jobs[i].Target == "api.ozon.ru:80",
				"All authors should be %s,  all targets - %s, didn't filtered")
		}
	})
	//t.Run("GetImbalance", func(t *testing.T) {
	//	actual, err := overloadService.GetJob(&overload.JobRequest{Id: jobId})
	//	if err != nil {
	//		t.Fatalf("Couldn't get job %s", err)
	//	}
	//	expected := jobUpdate
	//	_, errUpdate := overloadService.UpdateJob(&overload.UpdateJobRequest{Id: jobId, Status: "calculate imbalance"})
	//	if errUpdate != nil {
	//		t.Fatalf("Couldn't update job %s", errUpdate)
	//	}
	//	job := models.Job{AutostopTime: expected.AutostopTime}
	//	go overloadService.CalculateImbalance(&job)
	//	assert.NotEmpty(t, actual.Job.Imbalance, "Imbalance wasn't counted")
	//	assert.NotEmpty(t, actual.Job.AutostopTime, "There is no autostop, how was the imbalance counted?")
	//    assert.Equal(t, actual.Job.Imbalance, job.Imbalance, "Imbalance was counted wrong %d, %d")
	//})
	t.Run("GetCollections", func(t *testing.T) {
		job := models.Job{
			Id: jobId,
		}
		err := overloadService.Db.Select(&job)
		if err != nil {
			t.Fatalf("Couldn't parse ids of collections")
		}
		collectionIds := job.CollectionIds
		//sort.Slice(collectionIds, func(i, j int) bool { return collectionIds[i] < collectionIds[j] })
		t.Logf("Collections for this job %v", collectionIds)
		expected := jobCreate
		log.Printf("Excepted: %+v", expected)
		actual, err := overloadService.GetCollections(&overload.GetCollectionsRequest{CollectionId: collectionIds})
		if err != nil {
			t.Fatalf("Couldn't get collection %s", err)
		}
		assert.Len(t, actual.Collections, 2, "Two collections were added to job, but there are only %d")
		for i := 0; i < len(actual.Collections); i++ {
			assert.Equal(t, expected.Collections[i].Author, actual.Collections[len(actual.Collections)-i-1].Author, "Author in created collection and actual are different %s, %s")
			assert.Equal(t, expected.Collections[i].Env, actual.Collections[len(actual.Collections)-i-1].Env, "Env in created collection and actual are different %s, %s")
			assert.Equal(t, expected.Collections[i].Name, actual.Collections[len(actual.Collections)-i-1].Name, "Name in created collection and actual are different %s, %s")
			assert.Equal(t, expected.Collections[i].Project, actual.Collections[len(actual.Collections)-i-1].Project, "Project in created collection and actual are different %s, %s")
			assert.Equal(t, expected.Collections[i].Ref, actual.Collections[len(actual.Collections)-i-1].Ref, "Ref in created collection and actual are different %s, %s")
		}
	})
	t.Run("DeleteJob", func(t *testing.T) {
		_, err := overloadService.DeleteJob(&overload.DeleteJobRequest{Id: jobId})
		if err != nil {
			t.Fatalf("Couldn't delete job %s", err)
		}
		actual, err := overloadService.GetJob(&overload.JobRequest{Id: jobId})
		if err != nil {
			t.Logf("Couldn't get deleted job %s", err)
		}
		assert.Empty(t, actual, "Job wasn't deleted %s")
	})

}
