-- code对应的key
local key = KEYS[1]
-- 验证次数的key
local cntKey = key.."count"
-- 验证的验证码
local code = ARGV[1]
-- 存在redis里的code
local codeInRedis = redis.call("get", key)
-- 验证次数
local count = tonumber(redis.call("get", cntKey))
if count == nil then
    return -1;
end
if count <= 0 then
    -- 用户一直输错、有人搞我
    -- 或者别人搞你
    return -1
elseif code == codeInRedis then
    -- 验证正确，过期次数为-1，防止再被验证
    redis.call("set", cntKey, -1)
    return 0
else
    -- 输错一次，减一
    redis.call("decr", cntKey)
    return -2
end