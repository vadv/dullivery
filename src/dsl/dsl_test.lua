commands = helpers.new()
db_engine = storage.new()

if commands.alive('localhost:3131') and commands.alive('localhost:3132') then
    commands.log("all servers up")
end

commands.copy('content://localhost:3131/my_file', 'content://localhost:3132/xxx/my_file_new')

for _, file  in pairs(commands.find('content://localhost:3132/xxx/*')) do
    print("name:", file.name, "size:", file.size)
end

commands.remove('content://localhost:3131/my_file')
commands.copy('content://localhost:3132/xxx/my_file_new', 'content://localhost:3131/yyy/my_file_2')
commands.copy('content://localhost:3131/yyy/my_file_2',   'content://localhost:3131/yyy/my_file_3')

local find_files_count = 0
for _, file  in pairs(commands.find('content://localhost:3131/yyy/*')) do
    print("name:", file.name, "size:", file.size)
    find_files_count = find_files_count + 1
end

if not find_files_count == 2 then
    error("bad count of files")
end

http = {}
function http.get(url)
    return commands.http("GET", url)
end

response = http.get("http://yandex.ru")
print(response.code)
print(response.body)

db_engine.set("db1", "key", "value")
db_engine.set("db1", "key2", "value2")
db_engine.set("db2", "key", "value2")

if not db_engine.get("db1", "key") == "value" then
    error("key may return value")
end

if db_engine.get("db1", "key") == db_engine.get("db1", "key2") then
    error("isolation not working")
end

if db_engine.get("db1", "key") == db_engine.get("db2", "key") then
    error("isolation failed")
end

db_engine.expire("db2", "key", 1)
commands.sleep(2)
if db_engine.get("db2", "key") == "value" then
    error("expire failed")
end
