// Ratelimiter code is heavily modified from Delly (Discord Extreme List)
// Modified with permission! (Thanks Ice <3)

package ratelimiter

//
//// Ratelimit is a single API ratelimit from a IP.
//type Ratelimit struct {
//	// Returns how many remaining requests there are before
//	// being ratelimited.
//	Remaining int `json:"remaining"`
//}
//
//func NewRatelimit() Ratelimit {
//	return Ratelimit{
//		Remaining: 500,
//	}
//}
//
//// Consume consumes the current remaining requests in a Ratelimit
//// and returns that object.
//func (r Ratelimit) Consume() Ratelimit {
//	r.Remaining = r.Remaining - 1
//	return r
//}
//
//func (r Ratelimit) Exceeded() bool {
//	return !r.Expired() && r.Remaining == 0
//}
//
//func (r Ratelimit) Expired() bool {
//	return r.ResetAt.UnixNano() < r.ResetAt.UnixNano()
//}
//
//// Ratelimiter is the base ratelimiter that handles all ratelimits
//type Ratelimiter struct {
//	logger     slog.Logger
//	Ratelimits map[string]Ratelimit
//	NextReset  time.Time
//	Redis      *managers.RedisManager
//	Limit      int
//	Reset      int
//}
//
//func NewRatelimiter(redis *managers.RedisManager) Ratelimiter {
//	log := slog.Make(sloghuman.Sink(os.Stdout))
//	s := time.Now()
//	count := redis.Connection.HLen(context.TODO(), "tsubaki:ratelimits").Val()
//
//	rl := Ratelimiter{
//		logger:     log,
//		Limit:      500,
//		Reset:      int((1 * time.Hour).Milliseconds()),
//		NextReset:  time.Now().Add(1 * time.Hour),
//		Ratelimits: make(map[string]Ratelimit),
//	}
//
//	log.Info(context.Background(), fmt.Sprintf("Took %s to get %d ratelimits.", time.Now().Sub(s).String(), count))
//	go rl.resetCurrentRatelimits()
//
//	return rl
//}
//
//func (r Ratelimiter) resetCurrentRatelimits() {
//	for {
//		select {
//		case <-time.After(time.Duration(r.Reset)):
//			{
//				r.NextReset = time.Now().Add(time.Duration(r.Reset))
//				for key := range r.All() {
//					r.Reset(key)
//				}
//			}
//		}
//	}
//}
//
//func (r Ratelimiter) All() map[string]*Ratelimit {
//	ratelimits := make(map[string]*Ratelimit)
//	result, err := r.Redis.Connection.HGetAll(context.TODO(), "tsubaki:ratelimits").Result()
//	if err != nil {
//		return ratelimits
//	}
//
//	for key, val := range result {
//		ratelimit := &Ratelimit{}
//		_ = json.Unmarshal([]byte(val), ratelimit)
//
//		ratelimits[key] = ratelimit
//	}
//
//	return ratelimits
//}
//
//func (r Ratelimiter) CacheRatelimit(ip string, ratelimit *Ratelimit) {
//	data, _ := json.Marshal(&ratelimit)
//	r.Redis.Connection.HMSet(context.TODO(), "tsubaki:ratelimits", ip, string(data))
//}
//
//func (r Ratelimiter) Get(ip string) *Ratelimit {
//	result, err := r.Redis.Connection.HGet(context.TODO(), "tsubaki:ratelimits", ip).Result()
//	if err != nil {
//		if err == redis.Nil {
//			ratelimit := NewRatelimit()
//
//			r.CacheRatelimit(ip, &ratelimit)
//			return &ratelimit
//		}
//
//		ratelimit := NewRatelimit()
//		return &ratelimit
//	}
//
//	var rl *Ratelimit
//	err = json.Unmarshal([]byte(result), &rl)
//	if err != nil {
//		ratelimit := NewRatelimit()
//		return &ratelimit
//	}
//
//	if rl == nil {
//		ru := NewRatelimit()
//		rl = &ru
//	}
//
//	rl.Consume()
//	return rl
//}
