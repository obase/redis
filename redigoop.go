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
	// TODO: 兼容性做法
	if !pi.P.Cluster && pi.P.Keyfix != "" && len(keysArgs) > 0 {
		keysArgs[0] = Keyfix(keysArgs[0], pi.P.Keyfix)
	}

	err = pi.C.Send(cmd, keysArgs...)
	if err != nil {
		pi.Err = err
		return
	}
	pi.Rcv++
	return
}