
DATABASE_HOST = "{{.DatabaseHost}}"
DATABASE_PORT = "{{.DatabasePort}}"
DATABASE_USER = "{{.DATABASE_USER}}"
DATABASE_PASSWD = "{{.DATABASE_PASSWORD}}"


gb = GearsBuilder('StreamReader')
gb.foreach(lambda x: execute('HMSET', x['streamId'], *x))
gb.foreach(lambda x: execute('SET', "database", f"{DATABASE_HOST}:{DATABASE_PORT}"))
gb.register('mystream')