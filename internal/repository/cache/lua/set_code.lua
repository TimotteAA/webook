-- 存验证码的key
-- phone_code:login:1111111
local key = KEYS[1]
-- 存验证码验证次数的key
-- phone_code:login:1111111
-- ..用于字符串拼接
local countKey = key.."count"
-- 验证码
local code = ARGV[1]
-- code的过期时间
local expiration = tonumber(redis.call("ttl", key))
if expiration == -1 then
--     发生了奇怪的事情
    return -2;
--     540 = 600 - 60，超过了一分钟，或者key不存在
elseif expiration == -2 or expiration < 540 then
    redis.call("set", key, code)
--     10分钟有效期
    redis.call("expire", key, 600)
--     发送后的验证次数
    redis.call("set", countKey, 3)
--     验证次数的过期时间
    redis.call("expire", countKey, 600)
    return 0
else
    -- 发送太频繁
    return -1;
end