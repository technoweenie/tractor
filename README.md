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
- inspector renderer ... 

- virtual components (get, set, fields, call)

- spreadsheet ideas

- !project twitch command
- adding component doesn't trigger save
- removing component should have a hook
- new components won't get added to global registry
- registry will populate field with children components

added ValueTo to registry,
now use it to get a value out of the registry of a particular type