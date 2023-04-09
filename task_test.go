/*
 * @Description:
 * @version:
 * @Author: Adxiong
 * @Date: 2023-04-09 00:41:17
 * @LastEditors: Adxiong
 * @LastEditTime: 2023-04-09 17:49:37
 */
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
			var caseNmae = c.caseName
			var ctx context.Context
			var cancel func() = nil
			ctx, cancel = context.WithCancel(context.Background())

			Task := NewTask(ctx, cancel, c.allowPartFailed)

			for i := 0; i < 10; i++ {
				var count = i
				Task.Do(func() error {
					var err error
					time.Sleep(time.Duration(count) * time.Second)
					if count == 4 {
						err = fmt.Errorf("failed")
					} else {
						fmt.Printf("%s hello %d\n", caseNmae, count)
					}
					return err
				})
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
