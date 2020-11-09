# datagears 

Tools for next-gen cloud based data processing based on RedisGears.
## Goals

  * [x] Support for runtime and secret values
  * [x] Keep track of deployed gears
  * [x] Provide versioning support for a deployment
  * [x] Attach additional metadata for each gear
  * [ ] Provide gear rollbacks
  * [ ] Provide gear based workflows for reproducible data processing

## Install

```
make install
```

This command will compile the binary and install it into `$GOPATH/bin`.

### Add a trigger

First we need to write a project manifest. To write a project manifest following structure can be used:
```
version: 1

remotes:
  origin:
    host: localhost
    port: 6379
    database: 0

gears:
  mygear1:
    entrypoint: ./gear.py
    version: 0.1.0
    type: trigger
    description: "write to database"
    values:
      DatabaseHost: localhost
      DatabasePort: 6543
    requirements:
      - numpy==1.19.4

  calc2:
    entrypoint: ./gear.py
    type: run
    description: "calculation"
    blocking: true

```
And save it in your project root under the name `datagears.yml`

```
dg trigger add mygear1 
```

This will deploy the gear `mygear1` to the `origin` instance and track it. 

### List triggers
To list all registration of triggers use:

```
dg trigger list
```

Which should give similar result:
```
ID         Name     Description        Type     Version  
XmW3QAFDp  mygear1  write to database  trigger  0.1.0   
```

### Remove trigger

To remove the trigger:
```
dg trigger remove XmW3QAFDp
```