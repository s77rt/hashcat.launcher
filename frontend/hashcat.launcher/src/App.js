import './App.css';
import 'antd/dist/antd.dark.css';

import React, { Component } from "react";
import { Alert, Tooltip, message, Row, Col, Statistic, Button, Layout, Tabs, Typography, Menu } from 'antd';
import { QuestionCircleOutlined, InfoCircleOutlined, DeploymentUnitOutlined, RocketOutlined, ToolOutlined, UserOutlined, HomeOutlined, PlusOutlined, SettingOutlined, UnorderedListOutlined, QuestionOutlined } from '@ant-design/icons';

import data from "./jsx/data/data";

import { getAlgorithms } from './jsx/data/algorithms';

import EventBus from "./jsx/eventbus/EventBus";

import NewTask from './jsx/newtask';
import Settings from './jsx/settings';
import Tools from './jsx/tools';
import Tasks from './jsx/tasks';
import Help from './jsx/help';
import About from './jsx/about';
import TasksStats from './jsx/stats/tasks';

const { TabPane } = Tabs;
const { Header, Footer, Sider } = Layout;
const { Text, Title } = Typography;

class App extends Component {
	constructor(props) {
		super(props);

		data.callback = () => {
			EventBus.dispatch("dataUpdate");
		};
		data.getHashes();
		data.getAlgorithms();
		data.getDictionaries();
		data.getRules();
		data.getMasks();

		this.onSelectMenu = this.onSelectMenu.bind(this);

		this.newTaskView = <NewTask />;
		this.tasksView = <Tasks />;
		this.settingsView = <Settings />;
		this.toolsView = <Tools />;
		this.helpView = <Help />;
		this.aboutView = <About />;

		this.state = {
			version: undefined,
			currentView: "New Task",

			isLoadedHashcat: undefined
		}
	}

	init() {
		if (typeof window.GOgetVersion === "function") {
			window.GOgetVersion().then(
				response => {
					this.setState({
						version: response
					});
				},
				error => {
					message.warning("Failed to get version"  + " " + error);
				}
			);
		}
		if (typeof window.GOrestoreTasks === "function") {
			window.GOrestoreTasks().then(
				() => null,
				error => {
					message.warning("Failed to restore tasks" + " " + error);
				}
			);
		}
	}

	setView(view) {
		this.setState({
			currentView: view
		});
	}

	onSelectMenu(e) {
		this.setView(e.key);
	}

	componentDidMount() {
		EventBus.on("taskUpdate", "App", (taskUpdate) => {
			TasksStats._update(taskUpdate);
			EventBus.dispatch("tasksUpdate");
		});
		EventBus.on("taskDelete", "App", (taskID) => {
			TasksStats._delete(taskID);
			EventBus.dispatch("tasksUpdate");
		});
		EventBus.on("dataUpdate", "App", () => {
			this.setState({
				isLoadedHashcat: Object.keys(getAlgorithms()).length > 0
			});
		});
		this.init();
	}

	componentWillUnmount() {
		EventBus.remove("taskUpdate", "App");
		EventBus.remove("taskDelete", "App");
		EventBus.remove("dataUpdate", "App");
	}

	render() {
		return (
			<Layout>
				<Sider
					style={{
						overflow: 'auto',
						height: '100vh',
						position: 'fixed',
						left: 0
					}}
					collapsed
				>
					<Menu theme="dark" onSelect={this.onSelectMenu} defaultSelectedKeys={[this.state.currentView]} mode="inline">
						<Menu.Item key="New Task" icon={<PlusOutlined />}>
							New Task
						</Menu.Item>
						<Menu.Item key="Tasks" icon={<UnorderedListOutlined />}>
							Tasks
						</Menu.Item>
						<Menu.Item key="Settings" icon={<SettingOutlined />}>
							Settings
						</Menu.Item>
						<Menu.Divider />
						<Menu.Item key="Tools" icon={<DeploymentUnitOutlined />}>
							Tools
						</Menu.Item>
						<Menu.Divider />
						<Menu.Item key="Help" icon={<QuestionCircleOutlined />}>
							Help
						</Menu.Item>
						<Menu.Item key="About" icon={<InfoCircleOutlined />}>
							About
						</Menu.Item>
					</Menu>
				</Sider>

				<div style={{ marginLeft: '80px'}}></div>

				<Layout style={{ height: "100vh" }}>
					<Header
						style={{
							display: 'flex',
							alignItems: 'center',
							position: 'fixed',
							zIndex: 1,
							width: '100%',
							backgroundColor: '#000',
							borderBottom: '1px #1d1d1d solid'
						}}
					>
						<img style={{ height: '100%'}} src={require('./images/Icon.png').default} />
						<Title level={3} style={{ margin: '0 10px', color: '#fff' }}>
							hashcat.launcher
						</Title>
						<span>
							{this.state.version ? (
								this.state.version === "dev" ? (
									"dev"
								) : (
									"v" + this.state.version
								)
							) : "dev"}
						</span>
					</Header>

					<div style={{ marginTop: '64px'}}></div>

					{this.state.isLoadedHashcat === false && (
						<Alert
							style={{ maxHeight: "38px" }}
							type="warning"
							message={
								<Tooltip
									title={
										<>
											hashcat is expected to be in the same directory as hashcat.launcher
											inside a subfolder <Text code>/hashcat</Text>
										</>
									}
								>
									hashcat not found
								</Tooltip>
							}
							banner
						/>
					)}

					<div
						style={{ display: this.state.currentView === "New Task" ? "block" : "none" }}
					>
						{this.newTaskView}
					</div>

					<div
						style={{
							display: this.state.currentView === "Tasks" ? "flex" : "none",
							flexDirection: "column",
							flex: "1 0 auto",
							maxHeight: this.state.isLoadedHashcat === false ? "calc(100% - 64px - 38px)" : "calc(100% - 64px)"
						}}
					>
						{this.tasksView}
					</div>

					<div
						style={{ display: this.state.currentView === "Settings" ? "block" : "none" }}
					>
						{this.settingsView}
					</div>

					<div
						style={{ display: this.state.currentView === "Tools" ? "block" : "none" }}
					>
						{this.toolsView}
					</div>

					<div
						style={{ display: this.state.currentView === "Help" ? "block" : "none" }}
					>
						{this.helpView}
					</div>

					<div
						style={{ display: this.state.currentView === "About" ? "block" : "none" }}
					>
						{this.aboutView}
					</div>
				</Layout>
			</Layout>
		)
	}
}

export default App;
