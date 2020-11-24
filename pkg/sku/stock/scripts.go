package stock

import "github.com/go-redis/redis/v8"

var (
	CounterScript = redis.NewScript(
		`
			local amount = tonumber(ARGV[1])
			local current = tonumber(redis.call('GET', KEYS[1]) or 0)
			
			if current < amount then
				return -1   -- 不够扣减
			elseif amount == 0 then
				return current
			else
				return redis.call('DECRBY', KEYS[1], amount)
			end
			`,
	)
)
