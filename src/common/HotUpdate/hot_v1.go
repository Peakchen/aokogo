package HotUpdate

//add by stefan

import (
	"common/Log"
	"os"
	"os/signal"
	"syscall"
)

type THotUpdate struct {
	HUInfo *TServerHotUpdateInfo
}

var (
	_hotU = &THotUpdate{
		HUInfo: nil,
	}
)

func RunHotUpdateCheck(svrSignal *TServerHotUpdateInfo) {
	chsignal := make(chan os.Signal, 1)
	//listen sign: ctrl+c, kill, user1, user2...
	//SIGUSR1,SIGUSR2 for linux.
	signal.Notify(chsignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT /*, syscall.SIGUSR1, syscall.SIGUSR2*/)
	_hotU.HUInfo = svrSignal
	go _hotU.checkloop(chsignal)
}

func (this *THotUpdate) checkloop(chsignal chan os.Signal) {
	if this.HUInfo == nil {
		return
	}

	for {
		select {
		case s := <-chsignal:
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				Log.FmtPrintln("signal exit:", s)
			/*
				case syscall.SIGUSR1:
					Log.FmtPrintln("signal usr1:", s)
				case syscall.SIGUSR2:
					Log.FmtPrintln("signal usr2:", s)
			*/
			default:
				Log.FmtPrintln("other signal:", s)
			}

			if this.HUInfo.Recvsignal == s && this.HUInfo.HUCallback != nil {
				this.HUInfo.HUCallback()
			}
		}
	}
}
