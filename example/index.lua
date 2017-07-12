local golem = require("golem")

function main()
    local req = channel.make()
    local worker = golem.worker("worker.lua", req)

    worker:send("test")

    local ok, msg = req:receive()
    req:close()

    print("main: "..msg)
end
