package gtask

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTask(t *testing.T) {

	CaseList := []struct {
		caseName        string
		allowPartFailed bool
		wantError       string
	}{
		{
			caseName:        "TestTask NotAllowPartFailed",
			allowPartFailed: false,
			wantError:       "failed",
		},
		{
			caseName:        "TestTask AllowPartFailed",
			allowPartFailed: true,
			wantError:       "failed",
		},
	}

	for _, c := range CaseList {
		Convey(c.caseName, t, func() {
			var ctx context.Context
			var cancel func() = nil
			// if !c.allowPartFailed {
			ctx, cancel = context.WithCancel(context.Background())
			// } else {
			// 	ctx = context.Background()
			// }

			Task := NewTask(ctx, cancel, c.allowPartFailed)

			for i := 0; i < 10; i++ {
				var count = i
				func(count int) {
					Task.Do(func() error {
						var err error
						time.Sleep(time.Duration(count) * time.Second)
						if count == 4 {
							err = fmt.Errorf("failed")
						} else {
							fmt.Printf("%s hello %d\n", c.caseName, count)
						}
						return err
					})
				}(count)
			}

			errWait := Task.Wait()

			if errWait != nil {
				So(errWait.Error(), ShouldEqual, c.wantError)
			}

			if c.wantError == "" {
				So(errWait, ShouldBeNil)
			}

		})
	}

}
