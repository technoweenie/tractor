# Setting up tractor
### Clone the Tractor repo
`go get gihtub.com/manifold/tractor`

### Install Tractor dependencies
```
yarn install 
yarn run compile
```

### Clone manifold/qtalk
In the folders **qrpc/node** and **qmux/node** run:
`yarn install`

### Run tractor extension
Using VSCode debugger (F5 on Windows) to start a VSCode environment with a Tractor tree view.


# TODO:
- rest of inspector actions (values)
    - number, bool, ref, maps/lists?
- expressions?
- components: auth, cron, etc
- digital ocean / terraform example
- inspector renderer ... virtual components

- spreadsheet ideas

- agent
  - extension start daemon command
  - rpc call to agent to "connect" (start and show logs, or just show logs)
    - run locally if unable?
  - rpc call to start/stop
  - listen on unix socket if specified
  - show status based on sockets directory
  - process manager + shutdown
