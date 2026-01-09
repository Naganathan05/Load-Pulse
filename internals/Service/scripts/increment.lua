local key = KEYS[1]
local limit = tonumber(ARGV[1])
local current = tonumber(redis.call("GET", key) or "0")
if current < limit then
	return redis.call("INCR", key)
else
	return -1
end
