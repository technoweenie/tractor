"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
const vscode = require("vscode");
const path = require("path");
const fs = require("fs");
const qmux = require("qmux");
const qrpc = require("qrpc");
class TreeExplorer {
    constructor(context) {
        const treeDataProvider = new NodeProvider(this);
        this.explorer = vscode.window.createTreeView('treeExplorer', { treeDataProvider });
        vscode.commands.registerCommand('treeExplorer.inspectNode', (nodeId) => this.inspectNode(nodeId, context));
        this.api = new qrpc.API();
        this.api.handle("state", treeDataProvider);
        this.connect();
    }
    connect() {
        return __awaiter(this, void 0, void 0, function* () {
            try {
                var conn = yield qmux.DialWebsocket("ws://localhost:4243");
            }
            catch (e) {
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
            yield this.client.call("subscribe");
        });
    }
    addNode(name, parentId) {
        this.client.call("appendNode", { "ID": parentId || "", "Name": name });
    }
    updateNode(id, name, active) {
        let params = {
            "ID": id,
            "Name": name
        };
        if (active !== undefined) {
            params["Active"] = active;
        }
        this.client.call("updateNode", params);
    }
    deleteNode(id) {
        this.client.call("deleteNode", id);
    }
    moveNode(id, index) {
        this.client.call("moveNode", { "ID": id, "Index": index });
    }
    incr() {
        this.client.call("incr");
        vscode.window.showInformationMessage(`incremented`);
    }
    inspectNode(nodeId, context) {
        this.selectedNodeId = nodeId;
        const sendState = () => {
            this.inspectorPanel.webview.postMessage({ "event": "state", "state": this.remoteState });
            this.inspectorPanel.webview.postMessage({ "event": "select", "nodeId": this.selectedNodeId });
        };
        if (this.inspectorPanel === undefined) {
            // TODO: make another if this one is closed!
            this.inspectorPanel = vscode.window.createWebviewPanel('inspector', "Inspector", vscode.ViewColumn.One, {
                localResourceRoots: [vscode.Uri.file(path.join(context.extensionPath, 'resources'))],
                enableScripts: true
            });
            fs.readFile(path.join(context.extensionPath, 'resources', 'inspector', 'inspector.html'), 'utf8', (err, contents) => {
                this.inspectorPanel.webview.html = contents.replace(new RegExp("vscode-resource://", "g"), "vscode-resource://" + path.join(context.extensionPath, 'resources'));
            });
            this.inspectorPanel.webview.onDidReceiveMessage(message => {
                switch (message.event) {
                    case 'ready':
                        sendState();
                        return;
                    case 'rpc':
                        this.client.call(message.method, message.params);
                        return;
                    case 'edit':
                        if (message.Filepath !== undefined) {
                            vscode.window.showTextDocument(vscode.Uri.file(message.Filepath), {
                                viewColumn: vscode.ViewColumn.Two
                            });
                            return;
                        }
                        if (message.params.Component === "Delegate") {
                            vscode.window.showTextDocument(vscode.Uri.file(path.join(vscode.workspace.workspaceFolders[0].uri.path, 'delegates', message.params.ID, 'delegate.go')), {
                                viewColumn: vscode.ViewColumn.Two
                            });
                        }
                        else {
                            vscode.window.showTextDocument(vscode.Uri.file(this.remoteState.componentPaths[message.params.Component]), {
                                viewColumn: vscode.ViewColumn.Two
                            });
                        }
                        return;
                }
            }, undefined, context.subscriptions);
        }
        sendState();
    }
}
exports.TreeExplorer = TreeExplorer;
class NodeProvider {
    constructor(explorer) {
        this._onDidChangeTreeData = new vscode.EventEmitter();
        this.onDidChangeTreeData = this._onDidChangeTreeData.event;
        this.explorer = explorer;
    }
    serveRPC(r, c) {
        return __awaiter(this, void 0, void 0, function* () {
            var msg = yield c.decode();
            if (this.explorer.inspectorPanel !== undefined) {
                this.explorer.inspectorPanel.webview.postMessage({ "event": "state", "state": msg });
            }
            this.explorer.remoteState = msg;
            this.refresh();
            // output.appendLine(JSON.stringify(msg));
            r.return();
        });
    }
    refresh() {
        this._onDidChangeTreeData.fire();
    }
    getTreeItem(element) {
        return element;
    }
    getChildren(element) {
        if (this.explorer.remoteState === undefined) {
            return Promise.resolve([]);
        }
        if (element) {
            let n = this.explorer.remoteState.nodes[element.id];
            let childrenPaths = this.explorer.remoteState.hierarchy.filter((p) => {
                let basePath = element.abspath + "/";
                if (p.startsWith(basePath)) {
                    return (p.replace(basePath, "").lastIndexOf("/") === -1);
                }
                else {
                    return false;
                }
            });
            return Promise.resolve(childrenPaths.map((p) => {
                return { id: this.explorer.remoteState.nodePaths[p], path: p };
            }).map((obj) => {
                let n = this.explorer.remoteState.nodes[obj.id];
                let collapse = vscode.TreeItemCollapsibleState.None;
                if (this.explorer.remoteState.hierarchy.filter((p) => p.startsWith(obj.path + "/")).length > 0) {
                    collapse = vscode.TreeItemCollapsibleState.Collapsed;
                }
                return new Node(n.name, obj.path, obj.id, n.index, collapse, { command: 'treeExplorer.inspectNode', title: "Inspect", arguments: [obj.id], });
            }));
        }
        else {
            let rootPaths = this.explorer.remoteState.hierarchy.filter((p) => {
                return (p.lastIndexOf("/") === 0);
            });
            return Promise.resolve(rootPaths.map((p) => {
                return { id: this.explorer.remoteState.nodePaths[p], path: p };
            }).map((obj) => {
                let n = this.explorer.remoteState.nodes[obj.id];
                let collapse = vscode.TreeItemCollapsibleState.None;
                if (this.explorer.remoteState.hierarchy.filter((p) => p.startsWith(obj.path + "/")).length > 0) {
                    collapse = vscode.TreeItemCollapsibleState.Collapsed;
                }
                return new Node(n.name, obj.path, obj.id, n.index, collapse, { command: 'treeExplorer.inspectNode', title: "Inspect", arguments: [obj.id], });
            }));
        }
    }
}
exports.NodeProvider = NodeProvider;
class Node extends vscode.TreeItem {
    constructor(label, abspath, id, index, collapsibleState, command) {
        super(label, collapsibleState);
        this.label = label;
        this.abspath = abspath;
        this.id = id;
        this.index = index;
        this.collapsibleState = collapsibleState;
        this.command = command;
        // get description(): string {
        // 	return "$(alert)";
        // }
        this.iconPath = {
            light: path.join(__filename, '..', '..', 'resources', 'icons', 'light', 'document.svg'),
            dark: path.join(__filename, '..', '..', 'resources', 'icons', 'dark', 'document.svg')
        };
        this.contextValue = 'node';
    }
    get tooltip() {
        return `${this.label} (${this.id})`;
    }
}
exports.Node = Node;
//# sourceMappingURL=manifold.js.map