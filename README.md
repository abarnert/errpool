# errpool
Like `sync.errgroup` but with limited concurrency

Like `errgroup`, `errpool` provides synchronization, error propagation, and `Context` cancelation for groups of goroutines working on subtasks of a common task.

Unlike `errgroup`, `errpool` limits the number of concurrent goroutines to a specified max concurrency—much like the thread/process/coroutine pools found in other languages' standard libraries (although with an `errgroup`-like API, rather than the more powerful APIs provided by languages like, e.g., C# or Python).

> `type Pool`
>> `func WithContext(ctx context.Context, maxTasks int) (*Pool, context.Context)`
>> `func (p *Pool) Go(f func() error)`
>> `func (p *Pool) Wait() error`
