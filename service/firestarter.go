package service

import (
	"github.com/x-yield/over-api/internal/models"
	"github.com/x-yield/over-api/pkg/overload-service"
	"github.com/x-yield/over-api/tools/tankapi"
	"log"
)

// TODO: bulk
func (s *OverloadService) FirestarterSessionSyncDB(session *overload.TankSession) {
	dbSession := &models.TankSession{
		Tank:          session.Tank,
		Conf:          session.Conf,
		Name:          session.Name,
		Failures:      session.Failures,
		Status:        session.Status,
		Stage:         session.Stage,
		ExternalId:    session.ExternalId,
		OverloadId:    session.OverloadId,
		ExternalJoint: session.ExternalJoint,
		OverloadJoint: session.OverloadJoint,
		Author:        session.Author,
	}
	log.Println(dbSession)
	_, err := s.Db.Model(dbSession).
		OnConflict("(tank, name) DO UPDATE").
		Set("conf = COALESCE(EXCLUDED.conf, ts.conf)").
		Set("name = COALESCE(EXCLUDED.name, ts.name)").
		Set("failures = COALESCE(EXCLUDED.failures, ts.failures)").
		Set("status = COALESCE(EXCLUDED.status, ts.status)").
		Set("stage = COALESCE(EXCLUDED.stage, ts.stage)").
		Set("external_id = COALESCE(EXCLUDED.external_id, ts.external_id)").
		Set("overload_id = COALESCE(EXCLUDED.overload_id, ts.overload_id)").
		Set("external_joint = COALESCE(EXCLUDED.external_joint, ts.external_joint)").
		Set("overload_joint = COALESCE(EXCLUDED.overload_joint, ts.overload_joint)").
		Set("author = COALESCE(EXCLUDED.author, ts.author)").
		Insert()
	if err != nil {
		log.Println(err)
	}
	return
}

// FirestarterValidate - Does not create or update db sessions
func (s *OverloadService) FirestarterValidate(req *overload.FirestarterValidateRequest) (*overload.FirestarterValidateResponse, error) {
	reqSessions := req.Sessions
	var sessions []*tankapi.TankSession
	for _, s := range reqSessions {
		session := tankapi.NewSession(s.Tank, s.Conf, s.ExternalId, s.ExternalJoint)
		sessions = append(sessions, session)

	}
	firestarter := tankapi.NewFirestarter()

	sessions = firestarter.Validate(sessions)

	var validatedSessions []*overload.TankSession
	for _, session := range sessions {
		validatedSessions = append(validatedSessions, &overload.TankSession{
			Tank:          session.Tank.Url,
			Conf:          session.Config.Contents,
			Name:          session.Name,
			Failures:      session.Failures,
			Stage:         session.Stage,
			Status:        session.Status,
			ExternalId:    session.ExternalId,
			OverloadId:    session.OverloadId,
			ExternalJoint: session.ExternalJoint,
			OverloadJoint: session.OverloadJoint,
			Author:        session.Author,
		})
	}
	return &overload.FirestarterValidateResponse{Sessions: validatedSessions}, nil
}

func (s *OverloadService) FirestarterPrepare(req *overload.FirestarterPrepareRequest) (*overload.FirestarterPrepareResponse, error) {
	reqSessions := req.Sessions
	var sessions []*tankapi.TankSession

	for _, s := range reqSessions {
		session := tankapi.NewSession(s.Tank, s.Conf, s.ExternalId, s.ExternalJoint)
		sessions = append(sessions, session)

	}
	firestarter := tankapi.NewFirestarter()

	sessions = firestarter.Prepare(sessions)

	var preparedSessions []*overload.TankSession
	for _, session := range sessions {
		preparedSession := &overload.TankSession{
			Tank:          session.Tank.Url,
			Conf:          session.Config.Contents,
			Name:          session.Name,
			Failures:      session.Failures,
			Stage:         session.Stage,
			Status:        session.Status,
			ExternalId:    session.ExternalId,
			OverloadId:    session.OverloadId,
			ExternalJoint: session.ExternalJoint,
			OverloadJoint: session.OverloadJoint,
			Author:        session.Author,
		}
		preparedSessions = append(preparedSessions, preparedSession)
		go s.FirestarterSessionSyncDB(preparedSession)
	}
	return &overload.FirestarterPrepareResponse{Sessions: preparedSessions}, nil
}

func (s *OverloadService) FirestarterRun(req *overload.FirestarterRunRequest) (*overload.FirestarterRunResponse, error) {
	reqSessions := req.Sessions
	var sessions []*tankapi.TankSession

	for _, s := range reqSessions {
		session := tankapi.NewSession(s.Tank, s.Conf, s.ExternalId, s.ExternalJoint)
		session.Name = s.Name
		session.Failures = s.Failures
		session.Stage = s.Stage
		session.Status = s.Status
		sessions = append(sessions, session)

	}
	firestarter := tankapi.NewFirestarter()

	sessions = firestarter.Run(sessions)

	var startedSessions []*overload.TankSession
	for _, session := range sessions {
		startedSession := &overload.TankSession{
			Tank:          session.Tank.Url,
			Conf:          session.Config.Contents,
			Name:          session.Name,
			Failures:      session.Failures,
			Stage:         session.Stage,
			Status:        session.Status,
			ExternalId:    session.ExternalId,
			OverloadId:    session.OverloadId,
			ExternalJoint: session.ExternalJoint,
			OverloadJoint: session.OverloadJoint,
			Author:        session.Author,
		}
		startedSessions = append(startedSessions, startedSession)
		go s.FirestarterSessionSyncDB(startedSession)
	}
	return &overload.FirestarterRunResponse{Sessions: startedSessions}, nil
}

func (s *OverloadService) FirestarterStop(req *overload.FirestarterStopRequest) (*overload.FirestarterStopResponse, error) {
	reqSessions := req.Sessions
	var sessions []*tankapi.TankSession

	for _, s := range reqSessions {
		session := tankapi.NewSession(s.Tank, s.Conf, s.ExternalId, s.ExternalJoint)
		session.Name = s.Name
		session.Failures = s.Failures
		session.Stage = s.Stage
		session.Status = s.Status
		sessions = append(sessions, session)

	}
	firestarter := tankapi.NewFirestarter()

	sessions = firestarter.Stop(sessions)

	var stoppedSessions []*overload.TankSession
	for _, session := range sessions {
		stoppedSession := &overload.TankSession{
			Tank:          session.Tank.Url,
			Conf:          session.Config.Contents,
			Name:          session.Name,
			Failures:      session.Failures,
			Stage:         session.Stage,
			Status:        session.Status,
			ExternalId:    session.ExternalId,
			OverloadId:    session.OverloadId,
			ExternalJoint: session.ExternalJoint,
			OverloadJoint: session.OverloadJoint,
			Author:        session.Author,
		}
		stoppedSessions = append(stoppedSessions, stoppedSession)
		go s.FirestarterSessionSyncDB(stoppedSession)
	}
	return &overload.FirestarterStopResponse{Sessions: stoppedSessions}, nil
}

func (s *OverloadService) FirestarterPoll(req *overload.FirestarterPollRequest) (*overload.FirestarterPollResponse, error) {
	reqSessions := req.Sessions
	var sessions []*tankapi.TankSession

	for _, s := range reqSessions {
		session := tankapi.NewSession(s.Tank, s.Conf, s.ExternalId, s.ExternalJoint)
		session.Name = s.Name
		session.Failures = s.Failures
		session.Stage = s.Stage
		session.Status = s.Status
		sessions = append(sessions, session)

	}
	firestarter := tankapi.NewFirestarter()

	sessions = firestarter.Poll(sessions)

	var polledSessions []*overload.TankSession
	for _, session := range sessions {
		polledSession := &overload.TankSession{
			Tank:          session.Tank.Url,
			Conf:          session.Config.Contents,
			Name:          session.Name,
			Failures:      session.Failures,
			Stage:         session.Stage,
			Status:        session.Status,
			ExternalId:    session.ExternalId,
			OverloadId:    session.OverloadId,
			ExternalJoint: session.ExternalJoint,
			OverloadJoint: session.OverloadJoint,
			Author:        session.Author,
		}
		polledSessions = append(polledSessions, polledSession)
		go s.FirestarterSessionSyncDB(polledSession)
	}
	return &overload.FirestarterPollResponse{Sessions: polledSessions}, nil
}

func (s *OverloadService) FirestarterTankSessions(req *overload.FirestarterTankSessionsRequest) (*overload.FirestarterTankSessionsResponse, error) {
	tank := tankapi.Tank{Url: req.Tank}
	sessions, err := tank.Sessions()
	if err != nil {
		return nil, err
	}

	var tankSessions []*overload.TankSession
	for _, session := range sessions {
		tankSessions = append(tankSessions, &overload.TankSession{
			Tank:          session.Tank.Url,
			Conf:          session.Config.Contents,
			Name:          session.Name,
			Failures:      session.Failures,
			Stage:         session.Stage,
			Status:        session.Status,
			ExternalId:    session.ExternalId,
			OverloadId:    session.OverloadId,
			ExternalJoint: session.ExternalJoint,
			OverloadJoint: session.OverloadJoint,
			Author:        session.Author,
		})
	}
	return &overload.FirestarterTankSessionsResponse{Sessions: tankSessions}, nil
}
