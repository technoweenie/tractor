'use strict';
Object.defineProperty(exports, "__esModule", { value: true });
const vscode = require("vscode");
const path = require("path");
const fileExplorer_1 = require("./fileExplorer");
function activate(context) {
    // Samples of `window.registerTreeDataProvider`
    // const nodeDependenciesProvider = new DepNodeProvider(vscode.workspace.workspaceFolders[0].uri.path);
    // vscode.window.registerTreeDataProvider('nodeDependencies', nodeDependenciesProvider);
    // vscode.commands.registerCommand('nodeDependencies.refreshEntry', () => nodeDependenciesProvider.refresh());
    // vscode.commands.registerCommand('extension.openPackageOnNpm', moduleName => vscode.commands.executeCommand('vscode.open', vscode.Uri.parse(`https://www.npmjs.com/package/${moduleName}`)));
    // vscode.commands.registerCommand('nodeDependencies.addEntry', () => vscode.window.showInformationMessage(`Successfully called add entry.`));
    // vscode.commands.registerCommand('nodeDependencies.editEntry', (node: Dependency) => vscode.window.showInformationMessage(`Successfully called edit entry on ${node.label}.`));
    // vscode.commands.registerCommand('nodeDependencies.deleteEntry', (node: Dependency) => vscode.window.showInformationMessage(`Successfully called delete entry on ${node.label}.`));
    // const jsonOutlineProvider = new JsonOutlineProvider(context);
    // vscode.window.registerTreeDataProvider('jsonOutline', jsonOutlineProvider);
    // vscode.commands.registerCommand('jsonOutline.refresh', () => jsonOutlineProvider.refresh());
    // vscode.commands.registerCommand('jsonOutline.refreshNode', offset => jsonOutlineProvider.refresh(offset));
    // vscode.commands.registerCommand('jsonOutline.renameNode', offset => jsonOutlineProvider.rename(offset));
    // vscode.commands.registerCommand('extension.openJsonSelection', range => jsonOutlineProvider.select(range));
    var tree = new fileExplorer_1.FileExplorer(context);
    vscode.commands.registerCommand('fileExplorer.refreshEntry', () => vscode.window.showInformationMessage(`Refreshed.`));
    vscode.commands.registerCommand('fileExplorer.addEntry', () => vscode.window.showInformationMessage(`Successfully called add entry.`));
    vscode.commands.registerCommand('fileExplorer.editEntry', (node) => {
        tree.incr();
    });
    vscode.commands.registerCommand('fileExplorer.deleteEntry', (node) => vscode.window.showInformationMessage(`Successfully called delete entry on ${node.label}.`));
    vscode.window.createTerminal("Tractor", path.join(context.extensionPath, 'repl.js'));
}
exports.activate = activate;
//# sourceMappingURL=extension.js.map