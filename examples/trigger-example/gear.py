
DATABASE_HOST = "{{.DatabaseHost}}"
DATABASE_PORT = "{{.DatabasePort}}"

gb = GearsBuilder('StreamReader')
gb.foreach(lambda x: execute('HMSET', x['streamId'], *x))
gb.foreach(lambda x: execute('SET', "database", f"{DATABASE_HOST}:{DATABASE_PORT}"))
gb.register('mystream')