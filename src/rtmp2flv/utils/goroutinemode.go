package utils

func OrDone(done <-chan interface{}, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok { // 外界关闭数据流
					return
				}
				select { // 防止写入阻塞
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func OrDoneInt(done <-chan int, c <-chan int) <-chan int {
	valStream := make(chan int)
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok { // 外界关闭数据流
					return
				}
				select { // 防止写入阻塞
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func Tee(done <-chan interface{}, in <-chan interface{}) (<-chan interface{}, <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})
	go func() {
		defer close(out1)
		defer close(out2)
		for val := range in {
			var out1, out2 = out1, out2 // 私有变量覆盖
			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case out1 <- val:
					out1 = nil // 置空阻塞机制完成select轮询
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()
	return out1, out2
}

func ToString(done <-chan interface{}, valueStream <-chan interface{}) <-chan string {
	stringStream := make(chan string)
	go func() {
		defer close(stringStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- v.(string):
			}
		}
	}()
	return stringStream
}
