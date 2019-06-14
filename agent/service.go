// Copyright 2018-2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/log"
	runc "github.com/containerd/containerd/runtime/v2/runc/v2"
	"github.com/containerd/containerd/runtime/v2/shim"
	taskAPI "github.com/containerd/containerd/runtime/v2/task"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/firecracker-microvm/firecracker-containerd/internal/bundle"
	"github.com/firecracker-microvm/firecracker-containerd/internal/vm"
	"github.com/firecracker-microvm/firecracker-containerd/proto"
)

const (
	containerRootDir = "/container"
)

// TaskService represents inner shim wrapper over runc in order to:
// - Add default namespace to ctx as it's not passed by ttrpc over vsock
// - Add debug logging to simplify debugging
// - Make place for future extensions as needed
type TaskService struct {
	taskManager vm.TaskManager
	runcService taskAPI.TaskService

	publisher shim.Publisher

	// Normally, it's ill-advised to store a context object in a struct. However,
	// in this case there appears to be little choice.
	//
	// The problem is that sometimes service methods, such as CreateVM, require
	// a context scoped to the lifetime of the agent process but they are only provided
	// by the TTRPC server a context scoped to the lifetime of the individual request,
	// not the shimCtx.
	//
	// shimCtx should thus only be used by service methods that need to provide
	// a context that will be canceled only when agent is shutting down and
	// cleanup should happen.
	//
	// This approach is also taken by containerd's current reference runc shim
	// v2 implementation
	shimCtx      context.Context
	shimCancel   context.CancelFunc
	driveHandler *driveHandler
}

// NewTaskService creates new runc shim wrapper
func NewTaskService(shimCtx context.Context, shimCancel context.CancelFunc, publisher shim.Publisher) (taskAPI.TaskService, error) {
	// We provide an empty string for "id" as the service manages multiple tasks; there is no single
	// "id" being managed. As noted in the comments of the called code, the "id" arg is only used by
	// the Cleanup function, so it will never be invoked as part of the task service API, which is all
	// we need.
	runcService, err := runc.New(shimCtx, "", publisher, shimCancel)
	if err != nil {
		return nil, err
	}

	dh, err := newDriveHandler(blockPath, drivePath)
	if err != nil {
		return nil, err
	}

	return &TaskService{
		taskManager: vm.NewTaskManager(shimCtx, log.G(shimCtx)),
		runcService: runcService,

		publisher:    publisher,
		shimCtx:      shimCtx,
		shimCancel:   shimCancel,
		driveHandler: dh,
	}, nil
}

func logPanicAndDie(logger *logrus.Entry) {
	if err := recover(); err != nil {
		logger.WithError(err.(error)).Fatalf("panic: %s", string(debug.Stack()))
	}
}

func unmarshalExtraData(marshalled *types.Any) (*proto.ExtraData, error) {
	// get json bytes from task request
	extraData := &proto.ExtraData{}
	err := types.UnmarshalAny(marshalled, extraData)
	if err != nil {
		return nil, err
	}

	return extraData, nil
}

// Create creates a new initial process and container using runc
func (ts *TaskService) Create(requestCtx context.Context, req *taskAPI.CreateTaskRequest) (*taskAPI.CreateTaskResponse, error) {
	defer logPanicAndDie(log.G(requestCtx))

	logger := log.G(requestCtx).WithFields(logrus.Fields{"id": req.ID, "bundle": req.Bundle})
	logger.Info("create")

	extraData, err := unmarshalExtraData(req.Options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal extra data")
	}

	// Just provide runc the options it knows about, not our wrapper
	req.Options = extraData.RuncOptions

	// Override the bundle dir and rootfs paths, which were set on the Host and thus not valid here in the Guest
	bundleDir := bundle.Dir(filepath.Join(containerRootDir, req.ID))
	err = bundleDir.Create()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create bundle dir")
	}

	defer func() {
		if err != nil {
			removeErr := os.RemoveAll(bundleDir.RootPath())
			if removeErr != nil {
				logger.WithError(removeErr).Error("failed to cleanup bundle dir")
			}
		}
	}()

	err = bundleDir.OCIConfig().Write(extraData.JsonSpec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to write oci config file")
	}

	driveID := strings.TrimSpace(extraData.DriveID)
	drive, ok := ts.driveHandler.GetDrive(driveID)
	if !ok {
		return nil, fmt.Errorf("Drive %q could not be found", driveID)
	}

	err = bundleDir.MountRootfs(drive.Path(), "ext4", nil)
	if err != nil {
		return nil, err
	}

	req.Bundle = bundleDir.RootPath()
	req.Rootfs = nil

	// Override the incoming stdio FIFOs, which have paths from the host that we can't use
	fifoSet, err := cio.NewFIFOSetInDir(bundleDir.RootPath(), req.ID, req.Terminal)
	if err != nil {
		logger.WithError(err).Error("failed opening stdio FIFOs")
		return nil, errors.Wrap(err, "failed to open stdio FIFOs")
	}

	var stdinConnectorPair *vm.IOConnectorPair
	if req.Stdin != "" {
		req.Stdin = fifoSet.Stdin
		stdinConnectorPair = &vm.IOConnectorPair{
			ReadConnector:  vm.VSockAcceptConnector(extraData.StdinPort),
			WriteConnector: vm.FIFOConnector(fifoSet.Stdin),
		}
	}

	var stdoutConnectorPair *vm.IOConnectorPair
	if req.Stdout != "" {
		req.Stdout = fifoSet.Stdout
		stdoutConnectorPair = &vm.IOConnectorPair{
			ReadConnector:  vm.FIFOConnector(fifoSet.Stdout),
			WriteConnector: vm.VSockAcceptConnector(extraData.StdoutPort),
		}
	}

	var stderrConnectorPair *vm.IOConnectorPair
	if req.Stderr != "" {
		req.Stderr = fifoSet.Stderr
		stderrConnectorPair = &vm.IOConnectorPair{
			ReadConnector:  vm.FIFOConnector(fifoSet.Stderr),
			WriteConnector: vm.VSockAcceptConnector(extraData.StderrPort),
		}
	}

	ioConnectorSet := vm.NewIOConnectorProxy(stdinConnectorPair, stdoutConnectorPair, stderrConnectorPair)
	resp, err := ts.taskManager.CreateTask(requestCtx, req, ts.runcService, ioConnectorSet)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create runc shim for task")
	}

	logger.WithField("pid", resp.Pid).Debugf("create succeeded")
	return resp, nil
}

// State returns process state information
func (ts *TaskService) State(requestCtx context.Context, req *taskAPI.StateRequest) (*taskAPI.StateResponse, error) {
	defer logPanicAndDie(log.G(requestCtx))
	log.G(requestCtx).WithFields(logrus.Fields{"id": req.ID, "exec_id": req.ExecID}).Debug("state")

	resp, err := ts.runcService.State(requestCtx, req)
	if err != nil {
		log.G(requestCtx).WithError(err).Error("state failed")
		return nil, err
	}

	log.G(requestCtx).WithFields(logrus.Fields{
		"id":     resp.ID,
		"bundle": resp.Bundle,
		"pid":    resp.Pid,
		"status": resp.Status,
	}).Debug("state succeeded")
	return resp, nil
}

// Start starts a process
func (ts *TaskService) Start(requestCtx context.Context, req *taskAPI.StartRequest) (*taskAPI.StartResponse, error) {
	defer logPanicAndDie(log.G(requestCtx))
	log.G(requestCtx).WithFields(logrus.Fields{"id": req.ID, "exec_id": req.ExecID}).Debug("start")

	resp, err := ts.runcService.Start(requestCtx, req)
	if err != nil {
		log.G(requestCtx).WithError(err).Error("start failed")
		return nil, err
	}

	log.G(requestCtx).WithField("pid", resp.Pid).Debug("start succeeded")
	return resp, nil
}

// Delete deletes the initial process and container
func (ts *TaskService) Delete(requestCtx context.Context, req *taskAPI.DeleteRequest) (*taskAPI.DeleteResponse, error) {
	defer logPanicAndDie(log.G(requestCtx))
	log.G(requestCtx).WithFields(logrus.Fields{"id": req.ID, "exec_id": req.ExecID}).Debug("delete")

	resp, err := ts.taskManager.DeleteTask(requestCtx, req, ts.runcService)
	if err != nil {
		return nil, err
	}

	log.G(requestCtx).WithFields(logrus.Fields{
		"pid":         resp.Pid,
		"exit_status": resp.ExitStatus,
	}).Debug("delete succeeded")
	return resp, nil
}

// Pids returns all pids inside the container
func (ts *TaskService) Pids(requestCtx context.Context, req *taskAPI.PidsRequest) (*taskAPI.PidsResponse, error) {
	defer logPanicAndDie(log.G(requestCtx))
	log.G(requestCtx).WithField("id", req.ID).Debug("pids")

	resp, err := ts.runcService.Pids(requestCtx, req)
	if err != nil {
		log.G(requestCtx).WithError(err).Error("pids failed")
		return nil, err
	}

	log.G(requestCtx).Debug("pids succeeded")
	return resp, nil
}

// Pause pauses the container
func (ts *TaskService) Pause(requestCtx context.Context, req *taskAPI.PauseRequest) (*types.Empty, error) {
	defer logPanicAndDie(log.G(requestCtx))
	log.G(requestCtx).WithField("id", req.ID).Debug("pause")

	resp, err := ts.runcService.Pause(requestCtx, req)
	if err != nil {
		log.G(requestCtx).WithError(err).Error("pause failed")
		return nil, err
	}

	log.G(requestCtx).Debug("pause succeeded")
	return resp, nil
}

// Resume resumes the container
func (ts *TaskService) Resume(requestCtx context.Context, req *taskAPI.ResumeRequest) (*types.Empty, error) {
	defer logPanicAndDie(log.G(requestCtx))
	log.G(requestCtx).WithField("id", req.ID).Debug("resume")

	resp, err := ts.runcService.Resume(requestCtx, req)
	if err != nil {
		log.G(requestCtx).WithError(err).Debug("resume failed")
		return nil, err
	}

	log.G(requestCtx).Debug("resume succeeded")
	return resp, nil
}

// Checkpoint saves the state of the container instance
func (ts *TaskService) Checkpoint(requestCtx context.Context, req *taskAPI.CheckpointTaskRequest) (*types.Empty, error) {
	defer logPanicAndDie(log.G(requestCtx))
	log.G(requestCtx).WithFields(logrus.Fields{"id": req.ID, "path": req.Path}).Info("checkpoint")

	resp, err := ts.runcService.Checkpoint(requestCtx, req)
	if err != nil {
		log.G(requestCtx).WithError(err).Error("checkout failed")
		return nil, err
	}

	log.G(requestCtx).Debug("checkpoint sutaskAPI")
	return resp, nil
}

// Kill kills a process with the provided signal
func (ts *TaskService) Kill(requestCtx context.Context, req *taskAPI.KillRequest) (*types.Empty, error) {
	defer logPanicAndDie(log.G(requestCtx))
	log.G(requestCtx).WithFields(logrus.Fields{"id": req.ID, "exec_id": req.ExecID}).Debug("kill")

	resp, err := ts.runcService.Kill(requestCtx, req)
	if err != nil {
		log.G(requestCtx).WithError(err).Error("kill failed")
		return nil, err
	}

	log.G(requestCtx).Debug("kill succeeded")
	return resp, nil
}

// Exec runs an additional process inside the container
func (ts *TaskService) Exec(requestCtx context.Context, req *taskAPI.ExecProcessRequest) (*types.Empty, error) {
	defer logPanicAndDie(log.G(requestCtx))
	log.G(requestCtx).WithFields(logrus.Fields{"id": req.ID, "exec_id": req.ExecID}).Debug("exec")

	resp, err := ts.runcService.Exec(requestCtx, req)
	if err != nil {
		log.G(requestCtx).WithError(err).Error("exec failed")
		return nil, err
	}

	log.G(requestCtx).Debug("exec succeeded")
	return resp, nil
}

// ResizePty resizes pty
func (ts *TaskService) ResizePty(requestCtx context.Context, req *taskAPI.ResizePtyRequest) (*types.Empty, error) {
	defer logPanicAndDie(log.G(requestCtx))
	log.G(requestCtx).WithFields(logrus.Fields{"id": req.ID, "exec_id": req.ExecID}).Debug("resize_pty")

	resp, err := ts.runcService.ResizePty(requestCtx, req)
	if err != nil {
		log.G(requestCtx).WithError(err).Error("resize_pty failed")
		return nil, err
	}

	log.G(requestCtx).Debug("resize_pty succeeded")
	return resp, nil
}

// CloseIO closes all IO inside container
func (ts *TaskService) CloseIO(requestCtx context.Context, req *taskAPI.CloseIORequest) (*types.Empty, error) {
	defer logPanicAndDie(log.G(requestCtx))
	log.G(requestCtx).WithFields(logrus.Fields{"id": req.ID, "exec_id": req.ExecID}).Debug("close_io")

	resp, err := ts.runcService.CloseIO(requestCtx, req)
	if err != nil {
		log.G(requestCtx).WithError(err).Error("close io failed")
		return nil, err
	}

	log.G(requestCtx).Debug("close io succeeded")
	return resp, nil
}

// Update updates running container
func (ts *TaskService) Update(requestCtx context.Context, req *taskAPI.UpdateTaskRequest) (*types.Empty, error) {
	defer logPanicAndDie(log.G(requestCtx))
	log.G(requestCtx).WithField("id", req.ID).Debug("update")

	resp, err := ts.runcService.Update(requestCtx, req)
	if err != nil {
		log.G(requestCtx).WithError(err).Error("update failed")
		return nil, err
	}

	log.G(requestCtx).Debug("update succeeded")
	return resp, nil
}

// Wait waits for a process to exit
func (ts *TaskService) Wait(requestCtx context.Context, req *taskAPI.WaitRequest) (*taskAPI.WaitResponse, error) {
	defer logPanicAndDie(log.G(requestCtx))
	log.G(requestCtx).WithFields(logrus.Fields{"id": req.ID, "exec_id": req.ExecID}).Debug("wait")

	resp, err := ts.runcService.Wait(requestCtx, req)
	if err != nil {
		log.G(requestCtx).WithError(err).Error("wait failed")
		return nil, err
	}

	log.G(requestCtx).WithField("exit_status", resp.ExitStatus).Debug("wait succeeded")
	return resp, nil
}

// Stats returns a process stats
func (ts *TaskService) Stats(requestCtx context.Context, req *taskAPI.StatsRequest) (*taskAPI.StatsResponse, error) {
	defer logPanicAndDie(log.G(requestCtx))
	log.G(requestCtx).WithField("id", req.ID).Debug("stats")

	resp, err := ts.runcService.Stats(requestCtx, req)
	if err != nil {
		log.G(requestCtx).WithError(err).Error("stats failed")
		return nil, err
	}

	log.G(requestCtx).Debug("stats succeeded")
	return resp, nil
}

// Connect returns shim information such as the shim's pid
func (ts *TaskService) Connect(requestCtx context.Context, req *taskAPI.ConnectRequest) (*taskAPI.ConnectResponse, error) {
	defer logPanicAndDie(log.G(requestCtx))
	log.G(requestCtx).WithField("id", req.ID).Debug("connect")

	resp, err := ts.runcService.Connect(requestCtx, req)
	if err != nil {
		log.G(requestCtx).WithError(err).Error("connect failed")
		return nil, err
	}

	log.G(requestCtx).WithFields(logrus.Fields{
		"shim_pid": resp.ShimPid,
		"task_pid": resp.TaskPid,
		"version":  resp.Version,
	}).Error("connect succeeded")
	return resp, nil
}

// Shutdown cleanups runc resources ans gracefully shutdowns ttrpc server
func (ts *TaskService) Shutdown(requestCtx context.Context, req *taskAPI.ShutdownRequest) (*types.Empty, error) {
	defer logPanicAndDie(log.G(requestCtx))
	log.G(requestCtx).WithFields(logrus.Fields{"id": req.ID, "now": req.Now}).Debug("shutdown")

	// shimCancel will result in the runc shim to be canceled in addition to unblocking agent's
	// main func, which will allow it to exit gracefully.
	defer ts.shimCancel()

	log.G(requestCtx).Debug("going to gracefully shutdown agent")
	return &types.Empty{}, nil
}
