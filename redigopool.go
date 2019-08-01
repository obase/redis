package redis

import (
	redigo "github.com/gomodule/redigo/redis"
	"net"
	"runtime"
	"sync"
	"time"
)

/*
链表基础操作有4种:
1. 借: 将idle头的元素指向used尾
2. 还: 将conn指向idle尾的位置
3. 加: 将conn指向idle尾的位置
4. 删: 将conn的前后对接起来,并将自尾置空
*/

type redigoConn struct {
	P   *redigoPool
	C   redigo.Conn // 原生连接
	T   int64       // 空闲开始时间
	Prv *redigoConn //前驱
	Nxt *redigoConn //后继
}

type redigoPool struct {
	*Option
	testIdleSeconds int64

	Mutex *sync.Mutex
	Cond  *sync.Cond //等待信号

	// 链表容器
	Nalls int         //所有链表的数量
	Nfree int         // 空闲链表的数量
	Hused *redigoConn //在有链表的头, 头部都是最新的元素
	Hfree *redigoConn //空闲链表的头, 头部都是最新的元素
	Tfree *redigoConn //空闲链表的尾, 每次都从尾部取用
}

func newRedigoPool(opt *Option) (*redigoPool, error) {
	p := &redigoPool{
		Option:          opt,
		testIdleSeconds: int64(opt.TestIdleTimeout.Seconds()),
		Mutex:           new(sync.Mutex),
		Hused:           new(redigoConn),
		Hfree:           new(redigoConn),
	}
	p.Cond = sync.NewCond(p.Mutex)
	p.Tfree = p.Hfree //头尾指向相同表示空表

	for i := 0; i < opt.InitConns; i++ {
		_, err := p.Add(false)
		if err != nil {
			// 清理已经建好的链接
			p.Scan(closeRegigoConn)
			return nil, err
		}
	}
	return p, nil
}

/*************************START: 链表操作*******************************/
func (l *redigoPool) Scan(f func(conn *redigoConn)) {
	for e := l.Hused.Nxt; e != nil; e = e.Nxt {
		f(e)
	}
	for e := l.Hfree.Nxt; e != nil; e = e.Nxt {
		f(e)
	}
}

// 添加新元素
func (l *redigoPool) Add(used bool) (*redigoConn, error) {

	conn, err := l.newRedigoConn()
	if err != nil {
		return nil, err
	}
	if used {
		//生成直接使用. 在used头插入,默认conn.Nxt为nil
		if n := l.Hused.Nxt; n != nil {
			n.Prv = conn
			conn.Nxt = n
		}
		l.Hused.Nxt = conn
		conn.Prv = l.Hused
	} else {
		//生成放在空闲列表
		if n := l.Hfree.Nxt; n != nil {
			n.Prv = conn
			conn.Nxt = n
		}
		l.Hfree.Nxt = conn
		conn.Prv = l.Hfree

		if l.Tfree == l.Hfree {
			l.Tfree = conn
			l.Tfree.Nxt = nil
		}
		l.Nfree++
	}
	l.Nalls++

	return conn, nil
}

// 删除已有元素
func (l *redigoPool) Del(conn *redigoConn) {
	if n := conn.Prv; n != nil {
		n.Nxt = conn.Nxt
	}
	if n := conn.Nxt; n != nil {
		n.Prv = conn.Prv
	}
	l.Nalls--
}

//借出元素,没有返回nil
func (l *redigoPool) Take() (conn *redigoConn) {
	// 非空链表
	if l.Hfree != l.Tfree {
		conn = l.Tfree
		// 从free尾移除
		l.Tfree = conn.Prv
		l.Tfree.Nxt = nil //确保尾的Nxt为nil
		//在used头插入,默认conn.Nxt为nil
		if n := l.Hused.Nxt; n != nil {
			n.Prv = conn
			conn.Nxt = n
		}
		l.Hused.Nxt = conn
		conn.Prv = l.Hused

		l.Nfree--
	}
	return
}

//归还元素
func (l *redigoPool) Back(conn *redigoConn) {

	if n := conn.Prv; n != nil {
		n.Nxt = conn.Nxt
	}
	if n := conn.Nxt; n != nil {
		n.Prv = conn.Prv
	}

	if n := l.Hfree.Nxt; n != nil {
		n.Prv = conn
		conn.Nxt = n
	}
	l.Hfree.Nxt = conn
	conn.Prv = l.Hfree

	if l.Tfree == l.Hfree {
		l.Tfree = conn
		l.Tfree.Nxt = nil
	}

	l.Nfree++
	return
}

func (p *redigoPool) newRedigoConn() (*redigoConn, error) {

	conn, err := net.DialTimeout(p.Option.Network, p.Option.Address[0], p.Option.ConnectTimeout)
	if err != nil {
		return nil, err
	}
	tcp := conn.(*net.TCPConn)

	if p.Option.Keepalive > 0 {
		tcp.SetKeepAlive(true)
		tcp.SetKeepAlivePeriod(p.Option.Keepalive)
	}
	if err != nil {
		return nil, err
	}
	c := redigo.NewConn(tcp, p.Option.ReadTimeout, p.Option.WriteTimeout)
	if p.Option.Password != "" {
		_, err = c.Do("AUTH", p.Option.Password)
	}

	// 创建元素
	return &redigoConn{
		P: p,
		C: c,
		T: time.Now().Unix(), // 链接以及其放入时间
	}, nil
}

func (p *redigoPool) Get() (ret *redigoConn, err error) {
	for {
		p.Mutex.Lock()
		if p.Option.MaxConns > 0 {
			for p.Nfree == 0 && p.Nalls >= p.Option.MaxConns {
				if p.ErrExceMaxConns {
					p.Mutex.Unlock()
					return nil, ErrExceedMaxConns
				} else {
					p.Cond.Wait()
				}
			}
		}

		if p.Nfree > 0 {
			ret = p.Take()
		} else {
			ret, err = p.Add(true)
			if err != nil {
				p.Mutex.Unlock()
				return
			}
		}
		p.Mutex.Unlock()

		if ret.T == 0 || ret.T+p.testIdleSeconds > time.Now().Unix() {
			//  无需检测或未超时,直接返回
			return
		} else if _, err = ret.C.Do("PING"); err == nil {
			// 检测通过, 直接返回
			return
		}
		p.Put(ret, &err) //回收销毁,再从循环获取下一个
	}
	return
}

// 使用指针主要预防defer的处理
func (p *redigoPool) Put(conn *redigoConn, perr *error) {
	var off bool //是否销毁链接池
	p.Mutex.Lock()
	if *perr == nil || *perr == redigo.ErrNil { //排除空错误的影响
		p.Back(conn)
	} else {
		p.Del(conn)
		off = true
	}
	p.Cond.Signal()
	p.Mutex.Unlock()

	if off {
		//销毁有问题的连接
		conn.C.Close()
	}

	// 避免独占
	runtime.Gosched()
}
/*************************DONE: 链表操作*******************************/

/*************************START: 接口操作*******************************/
func (p *redigoPool) Do(cmd string, keysArgs ...interface{}) (reply interface{}, err error) {
	rc, err := p.Get()
	if err != nil {
		return
	}
	if p.Keyfix != "" && len(keysArgs) > 0 {
		keysArgs = Keyfix(&p.Keyfix, keysArgs)
	}
	reply, err = rc.C.Do(cmd, keysArgs...)
	if err == redigo.ErrNil {
		reply = nil
		err = nil
	}
	p.Put(rc, &err)
	return
}

// 注意: 集群模式不支持
func (p *redigoPool) Pi(bf Batch, keysArgs ...interface{}) (ret []interface{}, err error) {
	rc, err := p.Get()
	if err != nil {
		return
	}
	pi := &redigoOP{
		redigoConn: rc,
	}
	err = bf(pi, keysArgs...)
	if pi.Err != nil {
		//内部错误,需要销毁CONN
		pi.P.Put(pi.redigoConn, &pi.Err)
		return nil, pi.Err
	} else if err != nil {
		//用户错误. 需要接收完才释放链接. 但是不要用户错误!
		if pi.Rcv > 0 {
			pi.Err = pi.C.Flush()
			if pi.Err != nil {
				pi.P.Put(pi.redigoConn, &pi.Err)
				return
			}
			ret = make([]interface{}, pi.Rcv)
			for i := 0; i < pi.Rcv; i++ {
				ret[i], pi.Err = pi.C.Receive()
				if pi.Err != nil {
					pi.P.Put(pi.redigoConn, &pi.Err)
					return
				}
			}
		}
		pi.P.Put(pi.redigoConn, &pi.Err)
		return
	}

	err = pi.C.Flush()
	if err != nil {
		pi.P.Put(pi.redigoConn, &err)
		return
	}
	ret = make([]interface{}, pi.Rcv)
	for i := 0; i < pi.Rcv; i++ {
		ret[i], err = pi.C.Receive()
		if err != nil {
			pi.P.Put(pi.redigoConn, &err)
			return
		}
	}
	// 在return之前归还资源
	pi.P.Put(pi.redigoConn, &err)
	return
}

// 注意: 集群模式不支持
func (p *redigoPool) Tx(bf Batch, args ...interface{}) (ret []interface{}, err error) {
	rc, err := p.Get()
	if err != nil {
		return
	}
	if err = rc.C.Send("MULTI"); err != nil {
		rc.P.Put(rc, &err)
		return
	}

	pi := &redigoOP{
		redigoConn: rc,
	}
	err = bf(pi, args...)
	if pi.Err != nil { //内部错误,需要销毁CONN
		pi.P.Put(pi.redigoConn, &pi.Err)
		return nil, pi.Err
	} else if err != nil { //用户返回数据. 不影响链接,但要执行DISCARD
		_, xerr := rc.C.Do("DISCARD")
		rc.P.Put(rc, &xerr)
		return
	}

	ret, _, err = ValueSlice(rc.C.Do("EXEC"))
	pi.P.Put(pi.redigoConn, &err)
	return
}

// Publish
func (p *redigoPool) Pub(key string, msg interface{}) (err error) {
	rc, err := p.Get()
	if err != nil {
		return
	}
	_, err = rc.C.Do("PUBLISH", key, msg)
	rc.P.Put(rc, &err)
	return
}

// Subscribe, 阻塞执行sf直到返回stop或error才会结束
func (p *redigoPool) Sub(key string, sf SubFun) (err error) {
	rc, err := p.Get()
	if err != nil {
		return
	}
	err = rc.C.Send("SUBSCRIBE", key)
	if err != nil {
		rc.P.Put(rc, &err)
		return
	}

	var stop bool
	for {
		stop, err = sf(rc.C.Receive())
		if stop || err != nil {
			rc.P.Put(rc, &err)
			return
		}
	}
}

func (p *redigoPool) Eval(script string, keyCount int, keysArgs ...interface{}) (reply interface{}, err error) {
	rc, err := p.Get()
	if err != nil {
		return
	}
	s := redigo.NewScript(keyCount, script)
	reply, err = s.Do(rc.C, keysArgs...)
	if err == redigo.ErrNil {
		reply = nil
		err = nil
	}
	p.Put(rc, &err)
	return
}

func closeRegigoConn(conn *redigoConn) {
	conn.C.Close()
}

func (p *redigoPool) Close() {
	// 关闭无需获取锁,避免二次阻塞
	p.Scan(closeRegigoConn)
	return
}

/*************************DONE: 接口操作*******************************/