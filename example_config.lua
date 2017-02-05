commands = helpers.new()
db_enging = storage.new()
db = {name: "default"}
http = {}

function db.set(key, val)
    db_engine.set(db.name, key, val)
end

function db.get(key)
    db_engine.set(db.name, key)
end

function db.expire(key, ttl)
    db_engine.expire(db.name, key, ttl)
end

-- response = http.get("http://ya.ru")
-- response.code - код ответа
-- response.body - содержимое ответа
function http.get(url)
    return commands.http("GET", url)
end

function log(str)
    commands.log(str)
end

function sleep(t)
    commands.sleep(t)
end


-- копирует с src на dst, может вызвать exception
function copy_unsafe(src, dst)
    log("[info] starting copy from "..src.." to "..dst)
    commands.copy(src, dst)
    log("[info] completed copy from "..src.." to "..dst)
 end

-- копирует с src на dst, игнорирует exception
function copy_safe(src, dst)
    if not pcall(copy_unsafe(src, dst)) then
        log("[error] failed copy from "..src.." to "..dst..", skip error")
    end
end

-- обертка над copy_safe и copy_unsafe
function copy(src, dst, params)
    params = params or {skip_error=true}
    if params.skip_error then
        copy_safe(src, dst)
    else
        copy_unsafe(src, dst)
    end
end

-- функция возвращает true если машина доступна
function alive(srv)
    -- проверяем, есть ли 'content://' впереди, если надо, подставляем
    if string.sub(srv, 1, string.len('content://')) == 'content://' then
        addr = srv
    else
        addr = 'content://'..srv
    end
    -- выполняем собственно саму проверку
    if commands.alive(addr) then
        log("[info] host "..addr.." is up")
    else
        log("[info] host "..addr.." is down")
    end
end

log('[info] config loaded')
