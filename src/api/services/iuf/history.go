/*
 *
 *  MIT License
 *
 *  (C) Copyright 2022 Hewlett Packard Enterprise Development LP
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a
 *  copy of this software and associated documentation files (the "Software"),
 *  to deal in the Software without restriction, including without limitation
 *  the rights to use, copy, modify, merge, publish, distribute, sublicense,
 *  and/or sell copies of the Software, and to permit persons to whom the
 *  Software is furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included
 *  in all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 *  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 *  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 *  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 *  OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package services_iuf

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"sort"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s iufService) ListActivityHistory(activityName string) ([]iuf.History, error) {
	rawConfigMapList, err := s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		List(
			context.TODO(),
			v1.ListOptions{
				LabelSelector: fmt.Sprintf("type=%s,%s=%s", LABEL_HISTORY, LABEL_ACTIVITY_REF, activityName),
			},
		)
	if err != nil {
		s.logger.Error(err)
		return []iuf.History{}, err
	}

	sort.Slice(rawConfigMapList.Items, func(i, j int) bool {
		return rawConfigMapList.Items[i].CreationTimestamp.Before(&rawConfigMapList.Items[j].CreationTimestamp)
	})

	var res []iuf.History
	for _, rawConfigMap := range rawConfigMapList.Items {
		tmp, err := s.configMapDataToHistory(rawConfigMap.Data[LABEL_HISTORY])
		if err != nil {
			return []iuf.History{}, err
		}
		res = append(res, tmp)
	}
	return res, nil
}

func (s iufService) GetActivityHistory(activityName string, startTime int32) (iuf.History, error) {
	rawConfigMapList, err := s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		List(
			context.TODO(),
			v1.ListOptions{
				LabelSelector: fmt.Sprintf("type=%s,%s=%s", LABEL_HISTORY, LABEL_ACTIVITY_REF, activityName),
			},
		)
	if err != nil {
		s.logger.Error(err)
		return iuf.History{}, err
	}
	var res iuf.History
	for _, rawConfigMap := range rawConfigMapList.Items {
		tmp, err := s.configMapDataToHistory(rawConfigMap.Data[LABEL_HISTORY])
		if err != nil {
			s.logger.Error(err)
			return iuf.History{}, err
		}
		if tmp.StartTime == startTime {
			res = tmp
			break
		}
	}
	return res, nil
}

func (s iufService) ReplaceHistoryComment(activityName string, startTime int32, req iuf.ReplaceHistoryCommentRequest) (iuf.History, error) {
	history, err := s.GetActivityHistory(activityName, startTime)
	if err != nil {
		s.logger.Error(err)
		return iuf.History{}, err
	}
	history.Comment = req.Comment

	// update history
	configmap, err := s.iufObjectToConfigMapData(history, history.Name, LABEL_HISTORY)
	if err != nil {
		s.logger.Error(err)
		return iuf.History{}, err
	}
	configmap.Labels[LABEL_ACTIVITY_REF] = activityName
	_, err = s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		Update(
			context.TODO(),
			&configmap,
			v1.UpdateOptions{},
		)
	if err != nil {
		s.logger.Error(err)
		return iuf.History{}, err
	}
	return history, nil
}

func (s iufService) HistoryRunAction(activityName string, req iuf.HistoryRunActionRequest) (iuf.Session, error) {
	activity, err := s.GetActivity(activityName)
	if err != nil {
		s.logger.Errorf("HistoryRunAction.1: an error occurred while creating a new session for activity %s: %v", activityName, err)
		return iuf.Session{}, err
	}

	sessions, err := s.ListSessions(activityName)
	if err != nil {
		s.logger.Errorf("HistoryRunAction.2: an error occurred while creating a new session for activity %s: %v", activityName, err)
		return iuf.Session{}, err
	}

	if len(sessions) > 0 && sessions[len(sessions)-1].CurrentState == iuf.SessionStateInProgress {
		s.logger.Errorf("HistoryRunAction.3: for activity %s, cannot create a new session while a previous one is running. Run abort first.", activityName)
		return iuf.Session{}, err
	}

	inputParamsForPatch := iuf.InputParametersPatch{}
	jsonInputParams, err := json.Marshal(req.InputParameters)
	if err != nil {
		s.logger.Errorf("HistoryRunAction.4: for activity %s, error while parsing input parameters: %#v", activityName, req.InputParameters)
		return iuf.Session{}, err
	}
	err = json.Unmarshal(jsonInputParams, &inputParamsForPatch)
	if err != nil {
		s.logger.Errorf("HistoryRunAction.5: for activity %s, error while parsing input parameters: %#v", activityName, req.InputParameters)
		return iuf.Session{}, err
	}

	activity, err = s.PatchActivity(activity, iuf.PatchActivityRequest{
		InputParameters: inputParamsForPatch,
		SiteParameters:  req.SiteParameters,
	})

	if err != nil {
		s.logger.Error(err)
		return iuf.Session{}, err
	}

	// store session
	name := utils.GenerateName(activity.Name)
	session := iuf.Session{
		InputParameters: activity.InputParameters,
		SiteParameters:  activity.SiteParameters,
		Products:        activity.Products,
		Name:            name,
		ActivityRef:     activityName,
	}
	return s.CreateSession(session, name, activity)
}

func (s iufService) HistoryAbortAction(activityName string, req iuf.HistoryAbortRequest) (iuf.Session, error) {
	// go through the sessions and if there is any session that is not completed or aborted, then mark it as aborted
	// and terminate its workflows.
	sessions, err := s.ListSessions(activityName)
	if err != nil {
		s.logger.Errorf("HistoryAbortAction: An error occurred while listing sessions for activity %s: %v", activityName, err)
		return iuf.Session{}, err
	}

	var errors []error

	for i := len(sessions) - 1; i >= 0; i-- {
		session := sessions[i]

		// if this session still has workflows running, then this is a good session for abort irrespective of its
		//  current state (i.e. even if it was aborted in the past).

		abortable, workflows, err := s.isSessionAbortable(session)

		if err != nil {
			return iuf.Session{}, err
		}

		if abortable {
			// add a history entry for aborted sessions
			comment := req.Comment
			if comment == "" {
				comment = fmt.Sprintf("Aborted at stage %s", session.CurrentStage)
			}

			err := s.AbortSession(&session, comment, req.Force, workflows)
			if err != nil {
				errors = append(errors, err)
			}
		}
	}

	if len(errors) > 0 {
		s.logger.Errorf("HistoryAbortAction: An error(s) occurred while aborting sessions for activity %s: %v", activityName, errors)
		return iuf.Session{}, err
	}

	if len(sessions) > 0 {
		return sessions[len(sessions)-1], nil
	} else {
		return iuf.Session{}, nil
	}
}

// Check whether the workflows associated with the session / activity is abortable. Returns all the workflows related to
//  the session/activity irrespective.
func (s iufService) isSessionAbortable(session iuf.Session) (bool, *v1alpha1.WorkflowList, error) {
	isAbortable := session.CurrentState != iuf.SessionStateCompleted && session.CurrentState != iuf.SessionStateAborted

	workflows, err := s.workflowClient.ListWorkflows(context.TODO(), &workflow.WorkflowListRequest{
		Namespace: "argo",
		ListOptions: &v1.ListOptions{
			// note that we do not use iuf=true label selector here because we also want to include the IUF node worker
			//  rebuild workflows as well.
			// Also note that we don't have session=%s specified here because abort is really called on an activity not
			//  a session. It would be difficult for CLI to terminate workflows across sessions.
			LabelSelector: fmt.Sprintf("activity=%s", session.ActivityRef),
		},
		Fields: "-items.status.nodes,-items.spec",
	})

	if err != nil {
		return false, nil, err
	}

	for _, workflowObj := range workflows.Items {
		if workflowObj.Status.Phase == v1alpha1.WorkflowRunning || workflowObj.Status.Phase == v1alpha1.WorkflowPending {
			isAbortable = true
			break
		}
	}
	return isAbortable, workflows, nil
}

func (s iufService) HistoryPausedAction(activityName string, req iuf.HistoryActionRequest) (iuf.Session, error) {
	// go through the sessions and if there is any session that is in_progress state, then mark it as paused
	sessions, err := s.ListSessions(activityName)
	if err != nil {
		s.logger.Errorf("HistoryPausedAction: An error occurred while listing sessions for activity %s: %v", activityName, err)
		return iuf.Session{}, err
	}

	var errors []error
	for i := len(sessions) - 1; i >= 0; i-- {
		session := sessions[i]
		if session.CurrentState == iuf.SessionStateInProgress {
			// add a history entry for paused sessions
			comment := req.Comment
			if comment == "" {
				comment = fmt.Sprintf("Paused at stage %s", session.CurrentStage)
			}

			err := s.PauseSession(&session, comment)
			if err != nil {
				errors = append(errors, err)
			}

			break
		}
	}

	if len(errors) > 0 {
		s.logger.Errorf("HistoryPausedAction: An error(s) occurred while aborting sessions for activity %s: %v", activityName, errors)
		return iuf.Session{}, err
	}

	if len(sessions) > 0 {
		return sessions[len(sessions)-1], nil
	} else {
		return iuf.Session{}, nil
	}
}

func (s iufService) HistoryResumeAction(activityName string, req iuf.HistoryActionRequest) (iuf.Session, error) {
	// go through the sessions and if there is any session that is in_progress state, then mark it as paused
	sessions, err := s.ListSessions(activityName)
	if err != nil {
		return iuf.Session{}, err
	}

	// when resuming, we only look at the very last session so that we don't accidentally run more than one session
	//  from the history
	if len(sessions) == 0 {
		err := utils.GenericError{Message: fmt.Sprintf("HistoryResumeAction.1: There are no sessions in activity %s", activityName)}
		s.logger.Error(err)
		return iuf.Session{}, err
	}

	session := sessions[len(sessions)-1]

	// add a history entry for paused sessions
	comment := req.Comment
	if comment == "" {
		comment = fmt.Sprintf("Resuming stage %s", session.CurrentStage)
	}

	if session.CurrentState == iuf.SessionStateCompleted {
		err := utils.GenericError{Message: fmt.Sprintf("HistoryResumeAction.3: The session %s in activity %s cannot be resumed because it is either Completed or Aborted. Try restarting or running a new session.", session.Name, activityName)}
		s.logger.Error(err)
		return session, err
	} else if session.CurrentState == iuf.SessionStateTransitioning {
		// do nothing, we don't want to overlap when the session is transitioning to the next stage.
		return session, err
	} else {
		err := s.ResumeSession(&session, comment)
		if err != nil {
			return iuf.Session{}, err
		}
		return session, nil

	}
}

func (s iufService) HistoryRestartAction(activityName string, req iuf.HistoryRestartRequest) (iuf.Session, error) {
	// go through the sessions and if there is any session that is abortable, abort it first.
	sessions, err := s.ListSessions(activityName)
	if err != nil {
		return iuf.Session{}, err
	}

	// when restarting, we only look at the very last session so that we don't accidentally run more than one session
	//  from the history
	if len(sessions) == 0 {
		err := utils.GenericError{Message: fmt.Sprintf("HistoryRestartAction.1: There are no sessions in activity %s", activityName)}
		s.logger.Error(err)
		return iuf.Session{}, err
	}

	session := sessions[len(sessions)-1]

	abortable, workflows, err := s.isSessionAbortable(session)

	if err != nil {
		return iuf.Session{}, err
	}

	if abortable {
		// first abort this session
		err := s.AbortSession(&session, "Aborting before restart", true, workflows)
		if err != nil {
			// just print the warning. We don't care if it doesn't abort
			s.logger.Warnf("HistoryRestartAction.2: There was an error aborting the current session before restarting. Session name %s in activity %s. Error: %v", session.Name, session.ActivityRef, err)
		}
	}

	// add a history entry for restart
	comment := req.Comment
	if comment == "" {
		comment = fmt.Sprintf("Restarting from stage %s to %s", session.InputParameters.Stages[0], session.InputParameters.Stages[len(session.InputParameters.Stages)-1])
	}

	// now modify the session so the current stage is blank (so it can be restarted)
	session.CurrentStage = ""
	session.CurrentState = ""
	session.InputParameters.Force = req.Force
	err = s.UpdateSessionAndActivity(session, comment)

	return session, err
}

func (s iufService) HistoryBlockedAction(activityName string, req iuf.HistoryActionRequest) (iuf.Session, error) {
	// this is only allowed when activity is in debug, paused, or wait_for_admin state.
	activity, err := s.GetActivity(activityName)
	if err != nil {
		return iuf.Session{}, err
	}

	sessions, err := s.ListSessions(activityName)
	if err != nil {
		return iuf.Session{}, err
	}

	// there shouldn't be any running sessions
	var lastSession iuf.Session
	for _, session := range sessions {
		if session.CurrentState == iuf.SessionStateInProgress || session.CurrentState == iuf.SessionStatePaused {
			err = utils.GenericError{
				Message: fmt.Sprintf("HistoryBlockedAction.1: For the activity %s, there is currently an session %s that is in state %s.", activityName, session.Name, session.CurrentStage),
			}
			s.logger.Error(err)
			return iuf.Session{}, err
		}

		lastSession = session
	}

	switch activity.ActivityState {
	case iuf.ActivityStateWaitForAdmin, iuf.ActivityStateDebug:
		activity.ActivityState = iuf.ActivityStateBlocked
		_, err := s.updateActivity(activity)
		if err != nil {
			return iuf.Session{}, err
		}
	case iuf.ActivityStateBlocked:
		// noop
		return lastSession, nil
	default:
		err = utils.GenericError{
			Message: fmt.Sprintf("HistoryBlockedAction.2: The activity %s must be in debug or wait_for_admin state for it to be marked as blocked. Currently, it is in %s: %v", activityName, activity.ActivityState, activity.ActivityState),
		}
		s.logger.Error(err)
		return iuf.Session{}, err
	}

	comment := req.Comment
	if comment == "" {
		comment = fmt.Sprintf("Blocked at stage %s", lastSession.CurrentStage)
	}

	// add a history entry for blocked activity
	err = s.CreateHistoryEntry(activityName, iuf.ActivityStateBlocked, comment)
	if err != nil {
		return iuf.Session{}, err
	}

	return lastSession, nil
}

func (s iufService) configMapDataToHistory(data string) (iuf.History, error) {
	var res iuf.History
	err := json.Unmarshal([]byte(data), &res)
	if err != nil {
		s.logger.Errorf("configMapDataToHistory: error while parsing JSON data for history %s: %v", data, err)
		return res, err
	}
	return res, err
}
