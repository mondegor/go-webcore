package schedule_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mondegor/go-sysmess/mrlog/slog/nop"
	"github.com/stretchr/testify/suite"

	mock_mrcore "github.com/mondegor/go-webcore/mrcore/mock"
	mock_mrworker "github.com/mondegor/go-webcore/mrworker/mock"
	"github.com/mondegor/go-webcore/mrworker/process/schedule"
	mock_schedule "github.com/mondegor/go-webcore/mrworker/process/schedule/mock"
)

const deadlineTimeout = 3 * time.Second

type SchedulerTestSuite struct {
	suite.Suite

	ctrl *gomock.Controller
	ctx  context.Context

	mockTask            *mock_mrworker.MockTask
	mockContextEmbedder *mock_schedule.MockcontextEmbedder
	mockErrorHandler    *mock_mrcore.MockErrorHandler
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
	ts.mockTask = mock_mrworker.NewMockTask(ts.ctrl)
	ts.mockTask.EXPECT().Caption().Return("mockTaskCaption").AnyTimes()

	ts.mockContextEmbedder = mock_schedule.NewMockcontextEmbedder(ts.ctrl)
	ts.mockErrorHandler = mock_mrcore.NewMockErrorHandler(ts.ctrl)
}

func (ts *SchedulerTestSuite) Test_StartWithNoTasks() {
	taskScheduler := schedule.NewTaskScheduler(
		ts.mockErrorHandler,
		nop.NewLoggerAdapter(),
		ts.mockContextEmbedder,
	)

	err := ts.taskSchedulerStart(taskScheduler)
	ts.ErrorIs(err, schedule.ErrInternalNoTasks)
}

func (ts *SchedulerTestSuite) Test_StartWithTaskZeroPeriod() {
	ts.mockTask.EXPECT().Period().Return(time.Duration(0))
	// ts.mockTask.EXPECT().Timeout().Return(time.Minute)

	taskScheduler := schedule.NewTaskScheduler(
		ts.mockErrorHandler,
		nop.NewLoggerAdapter(),
		ts.mockContextEmbedder,
		schedule.WithTasks(ts.mockTask),
	)

	err := ts.taskSchedulerStart(taskScheduler)
	ts.ErrorIs(err, schedule.ErrInternalZeroParam)
}

func (ts *SchedulerTestSuite) Test_StartWithTaskZeroTimeout() {
	ts.mockTask.EXPECT().Period().Return(time.Minute)
	ts.mockTask.EXPECT().Timeout().Return(time.Duration(0))

	taskScheduler := schedule.NewTaskScheduler(
		ts.mockErrorHandler,
		nop.NewLoggerAdapter(),
		ts.mockContextEmbedder,
		schedule.WithTasks(ts.mockTask),
	)

	err := ts.taskSchedulerStart(taskScheduler)
	ts.ErrorIs(err, schedule.ErrInternalZeroParam)
}

func (ts *SchedulerTestSuite) Test_StartWithStartupTask() {
	errTaskDoFinished := errors.New("TaskDoFinished")

	ts.mockTask.EXPECT().Startup().Return(true)
	ts.mockTask.EXPECT().Period().Return(time.Nanosecond)
	ts.mockTask.EXPECT().Timeout().Return(time.Second)

	ts.mockTask.
		EXPECT().
		Do(gomock.Any()).
		Return(errTaskDoFinished)

	taskScheduler := schedule.NewTaskScheduler(
		ts.mockErrorHandler,
		nop.NewLoggerAdapter(),
		ts.mockContextEmbedder,
		schedule.WithTasks(ts.mockTask),
	)

	err := ts.taskSchedulerStart(taskScheduler)
	ts.ErrorIs(err, errTaskDoFinished)
}

func (ts *SchedulerTestSuite) Test_StartAndShutdown() {
	const minTaskExecution = 10

	taskExecuted := make(chan struct{})
	schedulerFinished := make(chan struct{})

	ts.mockContextEmbedder.EXPECT().WithWorkerIDContext(ts.ctx).Return(ts.ctx).Times(2)
	ts.mockContextEmbedder.EXPECT().WithTaskIDContext(ts.ctx).Return(ts.ctx).MinTimes(minTaskExecution)

	ts.mockTask.EXPECT().Startup().Return(false).Times(2)
	ts.mockTask.EXPECT().Period().Return(time.Nanosecond).Times(4)
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
		MinTimes(1)

	taskScheduler := schedule.NewTaskScheduler(
		ts.mockErrorHandler,
		nop.NewLoggerAdapter(),
		ts.mockContextEmbedder,
		schedule.WithTasks(ts.mockTask, ts.mockTask),
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
