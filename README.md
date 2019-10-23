# Tractor

Programmable computing environment

### Prerequisites
 * golang 1.13+
 * vscode
 * typescript `yarn global add typescript`
 * clone [manifold/qtalk](https://github.com/manifold/qtalk) and run `make link`

### Setup
```
$ make setup
```

### Running
Open the tractor directory in VS Code and run `Debug > Start Debugging`. 
This brings up a new instance of VS Code running the Tractor extension.

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
