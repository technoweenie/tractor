import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';
import * as qmux from 'qmux';
import * as qrpc from 'qrpc';

export class TreeExplorer {

	explorer: vscode.TreeView<Node>;

	inspectorPanel: vscode.WebviewPanel;
	client: any;
	api: any;
    remoteState: any;

    selectedNodeId: any;

	constructor(context: vscode.ExtensionContext) {
		const treeDataProvider = new NodeProvider(this);
		this.explorer = vscode.window.createTreeView('treeExplorer', { treeDataProvider });
		vscode.commands.registerCommand('treeExplorer.inspectNode', (nodeId) => this.inspectNode(nodeId, context));
		
		this.api = new qrpc.API();
		this.api.handle("state", treeDataProvider);

		this.connect();
	}

	async connect() {
		try {
			var conn = await qmux.DialWebsocket("ws://localhost:4243");
		} catch (e) {
			setTimeout(() => {
				this.connect();
			}, 200);
			return;
		}
		conn.socket.onclose = () => {
			conn.close();
			setTimeout(() => {
				this.connect();
			}, 200);
		};
		var session = new qmux.Session(conn);
		this.client = new qrpc.Client(session, this.api);
		this.client.serveAPI();
		//window.rpc = client;
		await this.client.call("subscribe");
	}

	addNode(name: string) {
		this.client.call("appendNode", {"ID": "", "Name": name});
	}

	deleteNode(id: string) {
		this.client.call("deleteNode", id);
	}

	incr() {
		this.client.call("incr");
		vscode.window.showInformationMessage(`incremented`);
	}

	inspectNode(nodeId: any, context: vscode.ExtensionContext): void {
        this.selectedNodeId = nodeId;
        const sendState = () => {
            this.inspectorPanel.webview.postMessage({"event": "state", "state": this.remoteState});
            this.inspectorPanel.webview.postMessage({"event": "select", "nodeId": this.selectedNodeId});
        };
		if (this.inspectorPanel === undefined) {
			// TODO: make another if this one is closed!
			this.inspectorPanel = vscode.window.createWebviewPanel(
				'inspector',
				"Inspector",
				vscode.ViewColumn.One,
				{
					localResourceRoots: [vscode.Uri.file(path.join(context.extensionPath, 'resources'))],
					enableScripts: true
				}
			);
			fs.readFile(path.join(context.extensionPath, 'resources', 'inspector', 'inspector.html'), 'utf8', (err, contents) => {
				this.inspectorPanel.webview.html = contents.replace(new RegExp("vscode-resource://", "g"), "vscode-resource://"+path.join(context.extensionPath, 'resources'));
            });
            this.inspectorPanel.webview.onDidReceiveMessage(
                message => {
                  switch (message.event) {
                    case 'ready':
                      	sendState();
					  	return;
					case 'rpc':
						this.client.call(message.method, message.params);
						return;
                  }
                },
                undefined,
                context.subscriptions
              );
		}
        sendState();
	}
}


export class NodeProvider implements vscode.TreeDataProvider<Node> {

	private _onDidChangeTreeData: vscode.EventEmitter<Node | undefined> = new vscode.EventEmitter<Node | undefined>();
	readonly onDidChangeTreeData: vscode.Event<Node | undefined> = this._onDidChangeTreeData.event;

    private explorer: TreeExplorer;

	constructor(explorer: TreeExplorer) {
        this.explorer = explorer;
    }
    
    async serveRPC(r, c) {
        var msg = await c.decode();
        if (this.explorer.inspectorPanel !== undefined) {
            this.explorer.inspectorPanel.webview.postMessage({"event": "state", "state": msg});
        }
        this.explorer.remoteState = msg;
        this.refresh();
        // output.appendLine(JSON.stringify(msg));
        r.return();
    }

	refresh(): void {
		this._onDidChangeTreeData.fire();
	}

	getTreeItem(element: Node): vscode.TreeItem {
		return element;
	}

	getChildren(element?: Node): Thenable<Node[]> {
        if (this.explorer.remoteState === undefined) {
            return Promise.resolve([]);
        }
		if (element) {
            let n = this.explorer.remoteState.nodes[element.id];
            let childrenPaths = this.explorer.remoteState.hierarchy.filter((p) => {
                let basePath = element.abspath+"/";
                if (p.startsWith(basePath)) {
                    return (p.replace(basePath, "").lastIndexOf("/") === -1);
                    
                } else {
                    return false;
                }
            });
			return Promise.resolve(childrenPaths.map((p) => {
                return {id: this.explorer.remoteState.nodePaths[p], path: p};
            }).map((obj) => {
                let n = this.explorer.remoteState.nodes[obj.id];
                let collapse = vscode.TreeItemCollapsibleState.None;
                if (this.explorer.remoteState.hierarchy.filter((p) => p.startsWith(obj.path+"/")).length > 0) {
                    collapse = vscode.TreeItemCollapsibleState.Collapsed;
                }
                return new Node(n.name, obj.path, obj.id, collapse, { command: 'treeExplorer.inspectNode', title: "Inspect", arguments: [obj.id], });
            }));
		} else {
			let rootPaths = this.explorer.remoteState.hierarchy.filter((p) => {
                return (p.lastIndexOf("/") === 0);
            });
            return Promise.resolve(rootPaths.map((p) => {
                return {id: this.explorer.remoteState.nodePaths[p], path: p};
            }).map((obj) => {
                let n = this.explorer.remoteState.nodes[obj.id];
                let collapse = vscode.TreeItemCollapsibleState.None;
                if (this.explorer.remoteState.hierarchy.filter((p) => p.startsWith(obj.path+"/")).length > 0) {
                    collapse = vscode.TreeItemCollapsibleState.Collapsed;
                }
                return new Node(n.name, obj.path, obj.id, collapse, { command: 'treeExplorer.inspectNode', title: "Inspect", arguments: [obj.id], });
            }));
		}

	}

}

export class Node extends vscode.TreeItem {

	constructor(
        public readonly label: string,
        public readonly abspath: string,
		public readonly id: string,
		public readonly collapsibleState: vscode.TreeItemCollapsibleState,
		public readonly command?: vscode.Command
	) {
		super(label, collapsibleState);
	}

	get tooltip(): string {
		return `${this.label} (${this.id})`;
	}

	// get description(): string {
	// 	return this.label;
	// }

	iconPath = {
		light: path.join(__filename, '..', '..', 'resources', 'icons', 'light', 'document.svg'),
		dark: path.join(__filename, '..', '..', 'resources', 'icons', 'dark', 'document.svg')
	};

	contextValue = 'node';

}