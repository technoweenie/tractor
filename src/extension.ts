'use strict';

import * as vscode from 'vscode';
import * as path from 'path';


import { TreeExplorer } from './manifold';

let serverTask = undefined;

export function activate(context: vscode.ExtensionContext) {
	if (vscode.workspace.workspaceFolders === undefined) {
		return;
	}
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

	var tree = new TreeExplorer(context);
	vscode.commands.registerCommand('treeExplorer.addNode', () => {
		vscode.window.showInputBox({ placeHolder: 'Enter a node name' })
			.then(value => {
				if (value !== null && value !== undefined) {
					tree.addNode(value);
				}
			});
		
	});
	vscode.commands.registerCommand('treeExplorer.renameNode', (node: any) => {
		// tree.incr();
	});
	vscode.commands.registerCommand('treeExplorer.deleteNode', (node: any) => {
		tree.deleteNode(node.id);
	});

	vscode.window.createTerminal("Tractor", path.join(context.extensionPath, 'repl.js'));

	serverTask = vscode.tasks.executeTask(new vscode.Task({ type: 'server', task: 'server' }, "server", "tractor", new vscode.ShellExecution("cd _workspace && go run ./cmd/daemon/daemon.go")));
}

export function deactivate() {
	if (serverTask) {
		serverTask.terminate();
	}
}