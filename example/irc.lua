local irc = require("irc").irc(config.nick or config.user, config.user)

irc:connect(config.server)

for i, channel in ipairs(config.channels) do
    irc:join(channel)
end

local function say(params)
    if params.channel then
        irc:privmsg(params.channel, params.message)
    else
        for i, channel in ipairs(config.channels) do
            irc:privmsg(channel, params.message)
        end         
    end
end

local dispatch = {
    say=say
}

function worker(msg)
    local f = dispatch[msg.func]

    if not f then
        return
    end

    f(msg.params)
end
