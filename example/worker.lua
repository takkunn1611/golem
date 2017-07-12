local golem = require("golem")

function worker(msg)
    print("worker: "..msg)
    parent:send(msg)
end
