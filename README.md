# datagears 

The missing toolkit for next-gen cloud based data processing using `redisgears`.
## Goals

  * Support for runtime and secret values
  * Keep track of deployed gears
  * Provide versioning support for a deployment
  * Attach additional metadata for each gear
  * Provide gear rollbacks
  * Provide gear based workflows for reproducible data processing

## Install

```
make install
```

This command will compile the binary and install it into `$GOPATH/bin`.
## Manifest

To write a project manifest use the following structure:
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
    name: "data-mover"
    description: "write to database"
    values:
      DatabaseHost: localhost
      DatabasePort: 6543
    requirements:
      - numpy==1.19.4

  mygear2:
    entrypoint: ./gear.py
    type: run
    name: "calc2"
    description: "calculation"
    blocking: true
```
And save it in your project root under the name `datagears.yml`

```
dg trigger add mygear1 
```

This will deploy the gear `mygear1` to the `origin` instance it will track this action locally. 
Each deployment gets assigned the deployment ID with which remote and local states are matches.

To list all registration of triggers use:

```
dg trigger list
```

Which should give similar result:
```
ID         Name     Description        Type     Version  
XmW3QAFDp  mygear1  write to database  trigger  0.1.0   
```

To remove the trigger:
```
dg trigger remove XmW3QAFDp
```