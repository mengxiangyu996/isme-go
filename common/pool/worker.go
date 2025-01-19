package pool

type worker struct {
	taskQueue chan Task
	quit      chan struct{}
}

// 初始化协程
func newWorker() *worker {
	return &worker{
		taskQueue: make(chan Task),
		quit:      make(chan struct{}),
	}
}

// 启动协程
func (w *worker) start(pool *Pool) {

	// 防止阻塞，开启协程处理
	go func() {
		for {
			// 注册协程加入到协程池
			pool.workers <- w
			select {
			case <-w.quit: // 如果接收到停止信号，退出循环
				return
			case task := <-w.taskQueue: // 从任务队列中取出任务
				task.Execute() // 执行任务
				pool.wg.Done()
			}
		}
	}()
}

// 停止协程
func (w *worker) stop() {
	close(w.quit)
}