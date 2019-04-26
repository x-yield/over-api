package tests

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"

	"github.com/x-yield/over-api/internal/models"
	"github.com/x-yield/over-api/pkg/overload-service"
)

const (
	ammoFileName      = "unittest_ammo_d8d3cee5d43d40da935ca2d1af7c1575"
	artifactFileName  = "unittest_artifact_b993d7d3fe3542fd97e17fa23cb5e58e"
	multipartBoundary = "3831c42c660c4f4a9b883ccf38991717"
	unittestJobNumber = "1000"
	unittestFile      = "testFile"
)

// Struct analog to JobResponse.Used because of unmarshaller
type Job struct {
	Job *models.Job
}

// Struct analog to GetAggregatesResponse.Used because of unmarshaller
type Aggregate struct {
	Aggregates []*models.Aggregate
}

// Global id of test job and its aggregate
var jobId int32
var aggrId int32
var collectionIds []int32
var collections []*overload.Collection
var url string
var service = flag.String("URL", "localhost:80", "service name")

func parseUrl() {
	cluster := ".qa.dev.s.o3.ru"
	if strings.Contains(*service, "localhost") {
		cluster = ""
	}
	if strings.Contains(*service, "release-") {
		cluster = ".qa.stg.s.o3.ru"
	}
	log.Print(*service)
	url = "http://" + *service + cluster
}

// Tests if the job is created and new id is sent
func TestCreateJob(t *testing.T) {
	parseUrl()
	t.Logf("Service for test: %s", url)
	t.Run("http:POST/create_job", func(t *testing.T) {
		collections = append(collections, &overload.Collection{
			Env:     "prod",
			Project: "test-project",
			Ref:     "master",
			Name:    "dbtest",
			Author:  "autotest",
		}, &overload.Collection{
			Env:     "stg",
			Project: "test-project",
			Ref:     "master",
			Name:    "httptest",
			Author:  "autotest",
		})
		jobResponse := overload.CreateJobResponse{}
		jobCreate := overload.CreateJobRequest{
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
		body, err := json.Marshal(jobCreate)
		if err != nil {
			t.Fatalf("Body for request couldn't be marshalled, %s", err)
		}
		resp, err := http.Post(url+"/create_job", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("http.POST failed: %s", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("API responsed with another code: %d", resp.StatusCode)
		}
		defer resp.Body.Close()
		if response, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Fatalf("Failed to read response body: %s", err)
		} else if _ = json.Unmarshal(response, &jobResponse); err != nil {
			t.Fatalf("Failed to unmarshal json: %s", err)
		} else if jobResponse.Id == 0 {
			t.Fatalf("API didn't return id of created job, %d", jobResponse.Id)
		}
		jobId = jobResponse.Id
		t.Logf("Created job: %d", jobId)
	})
}

// Tests if an aggregate was uploaded and it's id is sent
func TestUploadAggregates(t *testing.T) {
	t.Run("http:POST/upload_aggregates", func(t *testing.T) {
		uploadAggregateResponse := overload.CreateAggregatesResponse{}
		aggregateCreate := overload.CreateAggregatesRequest{
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
		body, err := json.Marshal(aggregateCreate)
		if err != nil {
			t.Fatalf("Body for request couldn't be marshalled, %s", err)
		}
		resp, err := http.Post(url+"/upload_aggregates", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("http.POST failed: %s", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("API responsed with another code: %d", resp.StatusCode)
		}
		defer resp.Body.Close()
		if response, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Fatalf("Failed to read response body: %s", err)
		} else if _ = json.Unmarshal(response, &uploadAggregateResponse); err != nil {
			t.Fatalf("Failed to unmarshal json: %s", err)
		} else if uploadAggregateResponse.Id == 0 {
			t.Fatalf("API didn't return id of created aggregate, %d, %s", uploadAggregateResponse.Id, string(response))
		}
		aggrId = uploadAggregateResponse.Id

	})
}

// Tests if there are any aggregates, that aggregates belong to the same job, check aggregate id, label, responseCode, Q95
func TestAggregates(t *testing.T) {
	t.Run("http:GET/aggregates", func(t *testing.T) {
		aggregateResponse := Aggregate{}
		resp, err := http.Get(url + "/aggregates/" + fmt.Sprint(jobId))
		if err != nil {
			t.Fatalf("http.Get failed: %s", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("API responsed with another code: %d", resp.StatusCode)
		}
		defer resp.Body.Close()
		if response, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Fatalf("Failed to read response body: %s", err)
		} else if _ = json.Unmarshal(response, &aggregateResponse); err != nil {
			t.Fatalf("Failed to unmarshal json: %s", err)
		} else if len(aggregateResponse.Aggregates) <= 0 {
			t.Fatalf("Expected some aggregates, but there are no, %d", len(aggregateResponse.Aggregates))
		} else if jobId != aggregateResponse.Aggregates[0].JobId {
			t.Fatalf("TestIdId(%d) of aggregates should be equal to job number (%d)", aggregateResponse.Aggregates[0].JobId, jobId)
		} else if aggrId != aggregateResponse.Aggregates[0].Id {
			t.Fatalf("Expected aggregate id: %d, actual: %d", aggrId, aggregateResponse.Aggregates[0].Id)
		} else if expected := "__TEST__"; expected != aggregateResponse.Aggregates[0].Label {
			t.Fatalf("Expected label: %s, actual: %s", expected, aggregateResponse.Aggregates[0].Label)
		} else if expected := "200"; expected != aggregateResponse.Aggregates[0].ResponseCode {
			t.Fatalf("Expected responseCode: %s, actual: %s", expected, aggregateResponse.Aggregates[0].ResponseCode)
		} else if expected := 45; float32(expected) != aggregateResponse.Aggregates[0].Q95 {
			t.Fatalf("Expected Q95: %f, actual: %f", float32(expected), aggregateResponse.Aggregates[0].Q95)
		}

	})
}

// Tests if job is updated
func TestUpdateJob(t *testing.T) {
	t.Run("http:POST/update_job", func(t *testing.T) {
		jobUpdate := overload.UpdateJobRequest{
			Id:           jobId,
			TestStop:     1546224400,
			Config:       "full",
			Description:  "This is an updated record",
			Status:       "finished autotest",
			AutostopTime: 1546268310,
		}
		body, err := json.Marshal(jobUpdate)
		if err != nil {
			t.Fatalf("Body for request couldn't be marshalled, %s", err)
		}
		resp, err := http.Post(url+"/update_job", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("http.POST failed: %s", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("API responsed with another code: %d", resp.StatusCode)
		}
		defer resp.Body.Close()
	})

}

// Tests that API returns the job which it was asked for, tests if update was successful, whole fields of job
func TestGetJob(t *testing.T) {
	t.Run("http:GET/job", func(t *testing.T) {
		job := Job{}
		resp, err := http.Get(url + "/job/" + fmt.Sprint(jobId))
		if err != nil {
			t.Fatalf("http.Get failed: %s", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("API responsed with another code: %d", resp.StatusCode)
		}
		defer resp.Body.Close()
		if response, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Fatalf("Failed to read response body: %s", err)
		} else if _ = json.Unmarshal(response, &job); err != nil {
			t.Fatalf("Failed to unmarshal json: %s", err)
		} else if job.Job.Id != jobId {
			t.Fatalf("Expected job id: %d, actual: %d", jobId, job.Job.Id)
		} else if expected := "autotest"; expected != job.Job.Author {
			t.Fatalf("Expected author: %s, actual: %s", expected, job.Job.Author)
		} else if expected := "full"; expected != job.Job.Config {
			t.Fatalf("Expected config: %s, actual: %s", expected, job.Job.Config)
		} else if expected := "This is an updated record"; expected != job.Job.Description {
			t.Fatalf("Expected description: %s, actual: %s", expected, job.Job.Description)
		} else if expected := "finished autotest"; expected != job.Job.Status {
			t.Fatalf("Expected status: %s, actual: %s", expected, job.Job.Status)
		} else if expected := "no tank used"; expected != job.Job.Tank {
			t.Fatalf("Expected tank: %s, actual: %s", expected, job.Job.Tank)
		} else if expected := "test.ru"; expected != job.Job.Target {
			t.Fatalf("Expected target: %s, actual: %s", expected, job.Job.Target)
		} else if expected := "test details"; expected != job.Job.EnvironmentDetails {
			t.Fatalf("Expected environment details: %s, actual: %s", expected, job.Job.EnvironmentDetails)
		} else if expected := 1546214400; float64(expected) != job.Job.TestStart {
			t.Fatalf("Expected test start: %f, actual: %f", float64(expected), job.Job.TestStart)
		} else if expected := 1546224400; float64(expected) != job.Job.TestStop {
			t.Fatalf("Expected test stop: %f, actual: %f", float64(expected), job.Job.TestStop)
		}
	})
}

// Tests how many jobs are received, and if there are no duplicates in jobs
func TestLastJobsNumber(t *testing.T) {
	jobsResponse := overload.LastJobsResponse{}
	t.Run("http:GET/lastjobs", func(t *testing.T) {
		rand.Seed(time.Now().UnixNano())
		page := rand.Intn(50)
		limit := rand.Intn(150)
		jobNumber := rand.Intn(4000)
		resp, err := http.Get(url + "/lastjobs?page=" + strconv.Itoa(page) + "&limit=" + strconv.Itoa(limit))
		if err != nil {
			t.Fatalf("http.Get failed: %s", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("API responsed with another code: %d", resp.StatusCode)
		}
		defer resp.Body.Close()
		if response, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Fatalf("Failed to read response body: %s", err)
		} else if _ = json.Unmarshal(response, &jobsResponse); err != nil {
			t.Fatalf("Failed to unmarshal json: %s", err)
		} else if expected := limit; expected != len(jobsResponse.Jobs) {
			t.Fatalf("Expected number of jobs: %d, received: %d", expected, len(jobsResponse.Jobs))
		}
		t.Logf("Job number used to test last jobs: %d", jobNumber)
	})
	t.Run("http:lastjob==created job", func(t *testing.T) {
		resp, err := http.Get(url + "/lastjobs?page=1&limit=5")
		if err != nil {
			t.Fatalf("http.Get failed: %s", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("API responsed with another code: %d", resp.StatusCode)
		}
		defer resp.Body.Close()
		if response, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Fatalf("Failed to read response body: %s", err)
		} else if _ = json.Unmarshal(response, &jobsResponse); err != nil {
			t.Fatalf("Failed to unmarshal json: %s", err)
		} else if expected := 5; expected != len(jobsResponse.Jobs) {
			t.Fatalf("Expected number of jobs on index: %d, received: %d", expected, len(jobsResponse.Jobs))
		} else if firstJobId := jobsResponse.Jobs[0].Id; jobId != firstJobId {
			t.Fatalf("Received id (%d) should be id of the autotest job (%d)", firstJobId, jobId)
		}
	})
}

func TestGetCollections(t *testing.T) {
	t.Run("http:GET/collections", func(t *testing.T) {
		//get collection ids for job
		jobWithCollections := overload.JobResponse{}
		respJob, errJob := http.Get(url + "/job/" + fmt.Sprint(jobId))
		if errJob != nil {
			t.Fatalf("http.Get failed: %s", errJob)
		}
		if respJob.StatusCode != 200 {
			t.Fatalf("API responsed with another code: %d", respJob.StatusCode)
		}
		defer respJob.Body.Close()
		if response, err := ioutil.ReadAll(respJob.Body); err != nil {
			t.Fatalf("Failed to read response body: %s", err)
		} else if _ = json.Unmarshal(response, &jobWithCollections); err != nil {
			t.Fatalf("Failed to unmarshal json: %s", err)
		}
		for i := 0; i < len(jobWithCollections.Job.Collections); i++ {
			collectionIds = append(collectionIds, jobWithCollections.Job.Collections[i].Id)
		}

		collectionResponse := overload.GetCollectionsResponse{}
		resp, err := http.Get(url + "/collections?collection_id=" + fmt.Sprint(collectionIds[0]) + "&collection_id=" + fmt.Sprint(collectionIds[1]))
		if err != nil {
			t.Fatalf("http.Get failed: %s", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("API responsed with another code: %d", resp.StatusCode)
		}
		defer resp.Body.Close()
		if response, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Fatalf("Failed to read response body: %s", err)
		} else if _ = json.Unmarshal(response, &collectionResponse); err != nil {
			t.Fatalf("Failed to unmarshal json: %s", err)
		} else if len(collectionResponse.Collections) <= 0 {
			t.Fatalf("Expected some collections, but there are no, %d", len(collectionResponse.Collections))
		}
		for i := 0; i < len(collectionResponse.Collections); i++ {
			assert.Equal(t, collections[i].Env, collectionResponse.Collections[len(collectionResponse.Collections)-i-1].Env, "Env created and got differ from each other %s, %s")
			assert.Equal(t, collections[i].Ref, collectionResponse.Collections[len(collectionResponse.Collections)-i-1].Ref, "Ref created and got differ from each other %s, %s")
			assert.Equal(t, collections[i].Project, collectionResponse.Collections[len(collectionResponse.Collections)-i-1].Project, "Project created and got differ from each other %s, %s")
			assert.Equal(t, collections[i].Name, collectionResponse.Collections[len(collectionResponse.Collections)-i-1].Name, "Name created and got differ from each other %s, %s")
			assert.Equal(t, collections[i].Author, collectionResponse.Collections[len(collectionResponse.Collections)-i-1].Author, "Author created and got differ from each other %s, %s")
		}
	})
}

//Tests if job is deleted
func TestDeleteJob(t *testing.T) {
	t.Run("http:POST/delete_job", func(t *testing.T) {
		jobDelete := overload.DeleteJobRequest{
			Id: jobId,
		}
		body, err := json.Marshal(jobDelete)
		if err != nil {
			t.Fatalf("Body for request couldn't be marshalled, %s", err)
		}
		resp, err := http.Post(url+"/delete_job", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("http.POST failed: %s", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("Couldn't delete test job: %d", resp.StatusCode)
		}
		defer resp.Body.Close()
	})
}

//TestDownloadAmmo Uploads ammo into real life ceph # TODO: mock ceph
func TestUploadAmmo(t *testing.T) {
	t.Run("http:POST/upload_ammo", func(t *testing.T) {
		ammoFile, err := ioutil.ReadFile(unittestFile)
		if err != nil {
			t.Fatalf("Couldn't read test ammo file: %s", err)
		}
		uploadAmmo := overload.UploadAmmoRequest{ // instance test
			Name: ammoFileName,
			File: string(ammoFile),
		}
		var body = bytes.NewBuffer([]byte{})
		multipartWriter := multipart.NewWriter(body)
		err = multipartWriter.SetBoundary(multipartBoundary)
		err = multipartWriter.WriteField("name", uploadAmmo.Name)
		err = multipartWriter.WriteField("file", uploadAmmo.File)
		err = multipartWriter.Close()
		if err != nil {
			t.Fatalf("multipart failed %v", err)
		}
		resp, err := http.Post(url+"/upload_ammo", "multipart/form-data; boundary="+multipartBoundary, bytes.NewReader(body.Bytes()))
		if err != nil {
			t.Fatalf("http.POST failed: %s", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fatalf("Couldn't upload test ammo file: %d", resp.StatusCode)
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Log(err)
		}
		var jBody map[string]string
		err = json.Unmarshal(respBody, &jBody)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasSuffix(jBody["url"], "/ammo/"+ammoFileName) {
			t.Fatalf("Unexpected response %s", jBody["url"])
		}
	})
}

//TestDownloadAmmo Downloads ammo from real life ceph # TODO: mock ceph
func TestDownloadAmmo(t *testing.T) {
	t.Run("http:GET/download_ammo", func(t *testing.T) {
		resp400, err := http.Get(url + "/download_ammo") // no key must fail as bad request
		if err != nil {
			t.Fatalf("http.GET failed: %s", err)
		}
		defer resp400.Body.Close()
		if resp400.StatusCode != 400 {
			t.Fatalf("Expected 400, got: %d", resp400.StatusCode)
		}
		resp200, err := http.Get(url + "/download_ammo?key=" + ammoFileName)
		if err != nil {
			t.Fatalf("http.GET failed: %s", err)
		}
		defer resp200.Body.Close()
		if resp200.StatusCode != 200 {
			t.Fatalf("Couldn't download ammo file job: %d", resp200.StatusCode)
		}
		if resp200.Header.Get("Content-Type") != "application/octet-stream" {
			t.Fatalf("Wrong response mimetype; expected 'application/octet-stream', got '%s'", resp200.Header.Get("Content-Type"))
		}
		respBody, err := ioutil.ReadAll(resp200.Body)
		if err != nil {
			t.Log(err)
		}
		ammoFile, err := ioutil.ReadFile(unittestFile)
		if err != nil {
			t.Fatalf("Couldn't read test ammo file: %s", err)
		}
		if string(ammoFile) != string(respBody) {
			t.Fatalf("expected file body differs from actual")
		}
	})
}

//TestDownloadAmmo Downloads ammo from real life ceph # TODO: mock ceph
func TestDeleteAmmo(t *testing.T) {
	t.Run("http:GET/delete_ammo", func(t *testing.T) {
		resp400, err := http.Get(url + "/delete_ammo") // no key must fail as bad request
		if err != nil {
			t.Fatalf("http.GET failed: %s", err)
		}
		defer resp400.Body.Close()
		if resp400.StatusCode != 400 {
			t.Fatalf("Expected 400, got: %d", resp400.StatusCode)
		}
		resp200, err := http.Get(url + "/delete_ammo?key=" + ammoFileName)
		if err != nil {
			t.Fatalf("http.GET failed: %s", err)
		}
		defer resp200.Body.Close()
		if resp200.StatusCode != 200 {
			t.Fatalf("Couldn't delete ammo file job: %d", resp200.StatusCode)
		}
		// check that the file is not available anymore.
		resp404, err := http.Get(url + "/download_ammo?key=" + ammoFileName)
		if err != nil {
			t.Fatalf("http.GET failed: %s", err)
		}
		defer resp404.Body.Close()
		if resp404.StatusCode != 200 {
			t.Fatalf("Couldn't download ammo file job: %d", resp404.StatusCode)
		}
		respBody, err := ioutil.ReadAll(resp404.Body)
		if err != nil {
			t.Log(err)
		}
		if string(respBody) != "" {
			t.Fatalf("Expected empty file body")
		}

	})
}

//TestDownloadAmmo Uploads artifact into real life ceph # TODO: mock ceph
func TestUploadArtifact(t *testing.T) {
	t.Run("http:POST/upload_artifact", func(t *testing.T) {
		artifactFile, err := ioutil.ReadFile(unittestFile)
		if err != nil {
			t.Fatalf("Couldn't read test ammo file: %s", err)
		}
		uploadArtifact := overload.UploadArtifactRequest{ // instance test
			Job:  unittestJobNumber,
			Name: artifactFileName,
			File: string(artifactFile),
		}
		var body = bytes.NewBuffer([]byte{})
		multipartWriter := multipart.NewWriter(body)
		err = multipartWriter.SetBoundary(multipartBoundary)
		err = multipartWriter.WriteField("job", uploadArtifact.Job)
		err = multipartWriter.WriteField("name", uploadArtifact.Name)
		err = multipartWriter.WriteField("file", uploadArtifact.File)
		err = multipartWriter.Close()
		if err != nil {
			t.Fatalf("multipart failed %v", err)
		}
		resp, err := http.Post(url+"/upload_artifact", "multipart/form-data; boundary="+multipartBoundary, bytes.NewReader(body.Bytes()))
		if err != nil {
			t.Fatalf("http.POST failed: %s", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fatalf("Couldn't upload test ammo file: %d", resp.StatusCode)
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Log(err)
		}
		var jBody map[string]string
		err = json.Unmarshal(respBody, &jBody)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasSuffix(jBody["url"], "/artifacts/"+unittestJobNumber+"/"+artifactFileName) {
			t.Fatalf("Unexpected response filename %s", jBody["url"])
		}
	})
}

// TestListAmmo Gets ammo list from real life ceph TODO: mock ceph
func TestListAmmo(t *testing.T) {
	t.Run("http:POST/list_ammo", func(t *testing.T) {
		resp, err := http.Get(url + "/list_ammo")
		if err != nil {
			t.Fatalf("http.Get failed: %s", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fatalf("Couldn't get ammo list: %d", resp.StatusCode)
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Log(err)
		}
		var jBody map[string][]interface{}
		err = json.Unmarshal(respBody, &jBody)
		if err != nil {
			t.Fatal(err)
		}
		if len(jBody["ammo"]) == 0 {
			t.Log("got empty ammo list")
		}
	})
}

// TestListArtifacts Gets artifacts list from real life ceph TODO: mock ceph
func TestListArtifacts(t *testing.T) {
	t.Run("http:POST/list_artifacts", func(t *testing.T) {
		j, _ := strconv.Atoi(unittestJobNumber)
		listAtrifacts := overload.ListArtifactsRequest{ // instance test
			Job: int32(j),
		}
		resp, err := http.Get(url + "/list_artifacts/" + strconv.Itoa(int(listAtrifacts.Job)))
		if err != nil {
			t.Fatalf("http.Get failed: %s", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fatalf("Couldn't get artifacts list: %d", resp.StatusCode)
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Log(err)
		}
		var jBody map[string][]interface{}
		err = json.Unmarshal(respBody, &jBody)
		if err != nil {
			t.Fatal(err)
		}
		if len(jBody["artifacts"]) == 0 {
			t.Log("got empty artifacts list")
		}
	})
}
