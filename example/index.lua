local golem = require("golem")

function main()
    local server = channel.make()

    golem.worker(
        "cron.lua",
        server,
        {
            spec="@every 1s",
            command={
                func="say",
                params={
                    message="Yo!!"
                }
            }
        }
    )
    
    -- local irc = golem.worker(
    --     "irc.lua",
    --     server,
    --     {
    --         user="foo",
    --         server="127.0.0.1:6667",
    --         channels={
    --             "bar"
    --         }
    --     }
    -- )

    while true do
        local ok, req = server:receive()
        if not ok then
            break
        end

        print(req.params.message)
        -- irc:send(req)
    end
end
