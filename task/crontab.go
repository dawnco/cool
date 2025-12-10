package task

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

type Server struct {
	scheduler gocron.Scheduler
	stopSig   chan int
}

func (s *Server) Start() {

	if s.scheduler == nil {
		logx.Error("启动任务失败")
		return
	}

	s.scheduler.Start()
	logx.Infof("定时任务启动成功")
	<-s.stopSig

}
func (s *Server) Stop() {
	s.stopSig <- 1
	err := s.scheduler.Shutdown()
	if err != nil {
		logx.Errorf("关闭定时任务失败 %s", err)
	}
}

// AddJob rule 格式 "8 8 * * *" task 定时执行的任务
func (s *Server) AddJob(rule string, task func()) {

	_, err := s.scheduler.NewJob(
		gocron.CronJob(rule, false),
		gocron.NewTask(task),
	)

	if err != nil {
		logx.Errorf("添加定时任务失败 %s", err.Error())
	}
}

func MustScheduler(loc string) service.Service {

	location, err := time.LoadLocation(loc)
	if err != nil {
		logx.Errorf("启动定时任务失败 时区错误 %s", err.Error())
	}

	scheduler, err := gocron.NewScheduler(gocron.WithLocation(location))
	if err != nil {
		logx.Error("启动定时任务服务失败")
	}
	return &Server{
		scheduler: scheduler,
	}
}
