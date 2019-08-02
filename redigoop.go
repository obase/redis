package redis

type redigoOP struct {
	*redigoConn
	Err error //内部错误
	Rcv int   //应用接收的数量
}

func (pi *redigoOP) Do(cmd string, keysArgs ...interface{}) (err error) {
	if pi.Err != nil {
		return pi.Err
	}
	if pi.P.Keyfix != "" && len(keysArgs) > 0 {
		FillKeyfix2(&pi.P.Keyfix, keysArgs)
	}
	err = pi.C.Send(cmd, keysArgs...)
	if err != nil {
		pi.Err = err
		return
	}
	pi.Rcv++
	return
}
