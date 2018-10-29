package rpanic

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/params"
)

/*
rpanic 是为了设计给手机用户使用
实现在进程不退出的情况下,重新初始化资源,仍然可以使用.

发生错误,处理,然后创建新的 RaidenService使用.
*/
/*
 *	repository rpanic is designed for mobile users, which implements resource reallocation while processes are executing.
 *	If faults occurs, handle them and create a new RaidenService to continue.
 */

/*
errChan should never be closed, 否则有可能会引起崩溃
*/
/*
 *	errChan should never be closed, or mobile app might crash.
 */
var errChan chan error
var notifier []string

//永不关闭.
var notifyChan chan error

//InitRaidenPanic init my panic system
func InitRaidenPanic() {
	errChan = make(chan error, 20)
	notifyChan = make(chan error, 20)
	startNotify()
}

/*
PanicRecover 用于所有 go routine 错误通知,但是不一定都会被处理,
只会处理第一个错误.
*/
/*
 *	PanicRecover : function to feed fault notifications for all go routine,
 *	but those might not be processed, only the first one will be processed.
 */
func PanicRecover(ctx string) {
	if err := recover(); err != nil {
		err2 := fmt.Errorf("%s occured err %s", ctx, err)
		log.Error(err2.Error())
		if params.MobileMode {
			errChan <- err2
		} else {
			panic(err2)
		}

	}
}

//RegisterErrorNotifier who wants to know error
func RegisterErrorNotifier(name string) {
	log.Trace(fmt.Sprintf("RegisterErrorNotifier %s ", name))
	notifier = append(notifier, name)
}

//startNotify  start notify system,只针对每个 RaidenService 实例启动一次.
func startNotify() {
	if params.MobileMode {
		go func() {
			err := <-errChan
			for i := 0; i < len(notifier); i++ {
				notifyChan <- err
			}
		}()
	}

	log.Info(fmt.Sprintf("startNotify complete..."))
}

//GetNotify 返回 通知.
func GetNotify() <-chan error {
	return notifyChan
}
