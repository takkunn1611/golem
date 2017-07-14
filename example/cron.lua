local cron = require("cron").new()

cron:addFunc(config.spec, function()
    parent:send(config.command)
end)

cron:start()

function worker(msg)
end
