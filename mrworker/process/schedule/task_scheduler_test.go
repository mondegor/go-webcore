package schedule_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtrace"
	"github.com/stretchr/testify/suite"

	"github.com/mondegor/go-webcore/mrworker/process/schedule"
	"github.com/mondegor/go-webcore/mrworker/process/schedule/mock"
)

//go:generate mockgen -destination=./mock/mrworker.go -package=mock github.com/mondegor/go-webcore/mrworker Task
//go:generate mockgen -destination=./mock/mrtrace.go -package=mock github.com/mondegor/go-sysmess/mrtrace ContextManager

const deadlineTimeout = 3 * time.Second

type SchedulerTestSuite struct {
	suite.Suite

	ctrl *gomock.Controller
	ctx  context.Context

	mockTask           *mock.MockTask
	mockContextManager *mock.MockContextManager
}

func TestSchedulerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(SchedulerTestSuite))
}

func (ts *SchedulerTestSuite) SetupSuite() {
	ts.ctrl = gomock.NewController(ts.T())
	ts.ctx = context.Background()
}

func (ts *SchedulerTestSuite) TearDownSuite() {
	ts.ctrl.Finish()
}

func (ts *SchedulerTestSuite) SetupTest() {
	ts.mockTask = mock.NewMockTask(ts.ctrl)
	ts.mockTask.EXPECT().Caption().Return("mockTaskCaption").AnyTimes()

	ts.mockContextManager = mock.NewMockContextManager(ts.ctrl)
}

func (ts *SchedulerTestSuite) Test_StartWithNoTasks() {
	taskScheduler := schedule.NewTaskScheduler(
		errors.NopHandler(),
		mrlog.NopLogger(),
		ts.mockContextManager,
	)

	err := ts.taskSchedulerStart(taskScheduler)
	ts.ErrorIs(err, schedule.ErrInternalNoTasks)
}

func (ts *SchedulerTestSuite) Test_StartWithTaskZeroPeriod() {
	ts.mockTask.EXPECT().Period().Return(time.Duration(0))
	// ts.mockTask.EXPECT().Timeout().Return(time.Minute)

	taskScheduler := schedule.NewTaskScheduler(
		errors.NopHandler(),
		mrlog.NopLogger(),
		ts.mockContextManager,
		schedule.WithTasks(ts.mockTask),
	)

	err := ts.taskSchedulerStart(taskScheduler)
	ts.ErrorIs(err, schedule.ErrInternalZeroParam)
}

func (ts *SchedulerTestSuite) Test_StartWithTaskZeroTimeout() {
	ts.mockTask.EXPECT().Period().Return(time.Minute)
	ts.mockTask.EXPECT().Timeout().Return(time.Duration(0))

	taskScheduler := schedule.NewTaskScheduler(
		errors.NopHandler(),
		mrlog.NopLogger(),
		ts.mockContextManager,
		schedule.WithTasks(ts.mockTask),
	)

	err := ts.taskSchedulerStart(taskScheduler)
	ts.ErrorIs(err, schedule.ErrInternalZeroParam)
}

func (ts *SchedulerTestSuite) Test_StartWithStartupTask() {
	errTaskDoFinished := errors.New("TaskDoFinished")

	ts.mockContextManager.EXPECT().WithGeneratedProcessID(ts.ctx, mrtrace.KeyTaskID).Return(ts.ctx)

	ts.mockTask.EXPECT().Startup().Return(true)
	ts.mockTask.EXPECT().Period().Return(time.Nanosecond)
	ts.mockTask.EXPECT().Timeout().Return(time.Second).Times(2)

	ts.mockTask.
		EXPECT().
		Do(gomock.Any()).
		Return(errTaskDoFinished)

	taskScheduler := schedule.NewTaskScheduler(
		errors.NopHandler(),
		mrlog.NopLogger(),
		ts.mockContextManager,
		schedule.WithTasks(ts.mockTask),
	)

	err := ts.taskSchedulerStart(taskScheduler)
	ts.ErrorIs(err, errTaskDoFinished)
}

func (ts *SchedulerTestSuite) Test_StartAndShutdown() {
	const minTaskExecution = 10

	taskExecuted := make(chan struct{})
	schedulerFinished := make(chan struct{})

	ts.mockContextManager.EXPECT().WithGeneratedProcessID(ts.ctx, mrtrace.KeyWorkerID).Return(ts.ctx).Times(2)
	ts.mockContextManager.EXPECT().WithGeneratedProcessID(ts.ctx, mrtrace.KeyTaskID).Return(ts.ctx).MinTimes(minTaskExecution)

	ts.mockTask.EXPECT().Startup().Return(false).Times(2)
	ts.mockTask.EXPECT().Period().Return(time.Nanosecond).MinTimes(minTaskExecution)
	ts.mockTask.EXPECT().Timeout().Return(time.Second).MinTimes(minTaskExecution)
	ts.mockTask.EXPECT().SignalDo().Return(nil).AnyTimes()

	ts.mockTask.
		EXPECT().
		Do(gomock.Any()).
		DoAndReturn(func(ctx context.Context) error {
			<-taskExecuted

			return nil
		}).
		Return(nil).
		MinTimes(minTaskExecution)

	taskScheduler := schedule.NewTaskScheduler(
		errors.NopHandler(),
		mrlog.NopLogger(),
		ts.mockContextManager,
		schedule.WithTasks(ts.mockTask, ts.mockTask), // 2 tasks
	)

	go func() {
		err := taskScheduler.Start(ts.ctx, func() {})
		ts.NoError(err)

		<-schedulerFinished
	}()

	for i := 0; i < minTaskExecution; i++ {
		select {
		case taskExecuted <- struct{}{}:
		case <-time.After(deadlineTimeout):
			ts.T().Fatal("Test timed out: taskScheduler.Start()")
		}
	}

	close(taskExecuted)

	err := taskScheduler.Shutdown(ts.ctx)
	ts.Require().NoError(err)

	select {
	case schedulerFinished <- struct{}{}:
		close(schedulerFinished)
	case <-time.After(deadlineTimeout):
		ts.T().Fatal("Test timed out: taskScheduler.Shutdown()")
	}
}

func (ts *SchedulerTestSuite) taskSchedulerStart(taskScheduler *schedule.TaskScheduler) (err error) {
	notify := make(chan error)

	go func() {
		notify <- taskScheduler.Start(ts.ctx, func() {})
	}()

	select {
	case err = <-notify:
		close(notify)
	case <-time.After(deadlineTimeout):
		ts.T().Fatal("Test timed out: taskScheduler.Start()")
	}

	return err
}
