package jobs

import (
    "context"
    "log"
    "time"

    "github.com/anpanovv/planter/internal/services"
)

// WateringNotificationsJob handles checking and creating watering notifications
type WateringNotificationsJob struct {
    notificationService *services.NotificationService
    interval           time.Duration
    stopChan           chan struct{}
}

// NewWateringNotificationsJob creates a new watering notifications job
func NewWateringNotificationsJob(notificationService *services.NotificationService, interval time.Duration) *WateringNotificationsJob {
    return &WateringNotificationsJob{
        notificationService: notificationService,
        interval:           interval,
        stopChan:           make(chan struct{}),
    }
}

// Start starts the watering notifications job
func (j *WateringNotificationsJob) Start() {
    ticker := time.NewTicker(j.interval)
    go func() {
        for {
            select {
            case <-ticker.C:
                if err := j.checkAndCreateNotifications(); err != nil {
                    log.Printf("Error checking watering notifications: %v", err)
                }
            case <-j.stopChan:
                ticker.Stop()
                return
            }
        }
    }()
}

// Stop stops the watering notifications job
func (j *WateringNotificationsJob) Stop() {
    close(j.stopChan)
}

// checkAndCreateNotifications checks for plants that need watering and creates notifications
func (j *WateringNotificationsJob) checkAndCreateNotifications() error {
    ctx := context.Background()
    return j.notificationService.CheckAndCreateWateringNotifications(ctx)
} 