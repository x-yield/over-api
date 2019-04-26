package service

import (
	"errors"
	"fmt"

	"github.com/x-yield/over-api/internal/models"
	"github.com/x-yield/over-api/pkg/overload-service"
)

func (s *OverloadService) CreateJoint(req *overload.CreateJointRequest) (*overload.CreateJointResponse, error) {
	joint := &models.Joint{
		Jobs: req.Jobs,
		Name: req.Name,
	}
	err := s.Db.Insert(joint)
	if err != nil {
		return nil, errors.New(
				fmt.Sprintf("Failed to create new joint: %v", err))
	}

	return &overload.CreateJointResponse{Id: joint.Id}, nil
}

func (s *OverloadService) ListJoints(req *overload.ListJointsRequest) (*overload.ListJointsResponse, error) {
	var joints []*models.Joint
	err := s.Db.Model(&joints).Order("id DESC").Select()

	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Failed to get joints: %v", err))
	}
	var preparedJoints []*overload.Joint
	for _, joint := range joints {
		thisJointJobs, err := s.GetJointJobs(joint.Jobs)
		if err != nil {
			return nil, err
		}
		preparedJoints = append(preparedJoints, &overload.Joint{
			Id:   joint.Id,
			Jobs: thisJointJobs,
			Name: joint.Name,
		})
	}
	response := &overload.ListJointsResponse{
		Joints: preparedJoints,
	}
	return response, nil
}

func (s *OverloadService) GetJoint(req *overload.GetJointRequest) (*overload.GetJointResponse, error) {
	joint := models.Joint{
		Id: req.Id,
	}

	err := s.Db.Select(&joint)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to select joint: %v", err))
	}

	thisJointJobs, err := s.GetJointJobs(joint.Jobs)
	if err != nil {
		return nil, err
	}

	response := overload.GetJointResponse{
		Joint: &overload.Joint{
			Id:   joint.Id,
			Jobs: thisJointJobs,
			Name: joint.Name,
		},
	}
	return &response, nil
}

func (s *OverloadService) GetJointJobs(jobIds []int32) ([]*overload.Job, error) {
	var thisJointJobs []*overload.Job
	for _, jointJobId := range jobIds {
		job := models.Job{
			Id: jointJobId,
		}
		err := s.Db.Model(&job).WherePK().Select()
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf("Failed to get job associated w/ joint: %v", err))
		}
		preparedJob := &overload.Job{
			Id:              job.Id,
			TestStart:       job.TestStart,
			TestStop:        job.TestStop,
			Description:     job.Description,
			Author:          job.Author,
			Status:          job.Status,
			Tank:            job.Tank,
			Target:          job.Target,
			AutostopTime:    job.AutostopTime,
			AutostopMessage: job.AutostopMessage,
		}
		thisJointJobs = append(thisJointJobs, preparedJob)
	}
	return thisJointJobs, nil
}
