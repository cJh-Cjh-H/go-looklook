```Mermaid
graph LR
    A[Scheduler] -->|每分钟触发| B[Redis队列]
    B -->|消费者拉取| C[Job Worker]
    C --> D[关闭订单]
    C --> E[微信通知]
    C --> F[结算任务]
    
    G[业务系统] -->|支付成功后| H[异步任务]
    H -->|写入| B
```

### 1. **mqueue-job（任务执行器/Worker）**

- **作用**：执行具体的异步任务
- **位置**：`cmd/job/`
- **技术栈**：Asynq + Redis
- **功能**：
  - `CloseHomestayOrderHandler`：关闭未支付的民宿订单
  - `PaySuccessNotifyUserHandler`：支付成功后微信小程序通知用户
  - `SettleRecordHandler`：给商家结算（演示任务）

### 2. **mqueue-scheduler（任务调度器）**

- **作用**：定时触发任务到 Redis 队列
- **位置**：`cmd/scheduler/`
- **技术栈**：Asynq Scheduler + Redis
- **功能**：
  - `settleRecordScheduler`：每分钟触发结算任务

## 📋 三种任务类型

### 1. **定时任务（Scheduler）**

go

```
// 调度器注册：每分钟执行一次
l.svcCtx.Scheduler.Register("*/1 * * * *", task)
```



- 示例：`ScheduleSettleRecord`（每分钟结算一次）

### 2. **延迟任务（Deferred）**

go

```
// 延迟关闭订单
jobtype.DeferCloseHomestayOrder
```



- 用途：订单创建后，延迟关闭未支付的订单

### 3. **即时消息任务（Message）**

go

```
// 支付成功立即通知用户
jobtype.MsgPaySuccessNotifyUser
```



- 用途：支付成功后立即发送微信小程序通知