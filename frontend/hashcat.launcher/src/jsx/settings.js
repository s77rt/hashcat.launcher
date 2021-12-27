import React, { Component } from "react";
import { Popconfirm, Layout, PageHeader, message, Statistic, Row, Col, Card, Select, Typography, Upload, Button, Space, Form, Radio, Divider, Collapse, Checkbox, Tabs, Steps } from 'antd';
import { FileOutlined, AimOutlined, ToolOutlined, ExportOutlined, ExperimentOutlined, SyncOutlined } from '@ant-design/icons';

import EventBus from "./eventbus/EventBus";

import data from "./data/data";
import { getHashes } from './data/hashes';
import { getAlgorithms } from './data/algorithms';
import { getDictionaries } from './data/dictionaries';
import { getRules } from './data/rules';
import { getMasks } from './data/masks';

const { Content } = Layout;
const { Dragger } = Upload;
const { Option } = Select;
const { Title, Paragraph, Text, Link } = Typography;
const { Panel } = Collapse;
const { TabPane } = Tabs;
const { Step } = Steps;

class Settings extends Component {
	constructor(props) {
		super(props);

		this.onClickRescan = this.onClickRescan.bind(this);
		this.onClickRefreshTaskCounter = this.onClickRefreshTaskCounter.bind(this);
		this.onClickResetTaskCounter = this.onClickResetTaskCounter.bind(this);

		this.state = {
			taskCounter: "-",

			isLoadingRescan: false,
			isLoadingRefreshTaskCounter: false,
			isLoadingResetTaskCounter: false,

			_dictionaries: getDictionaries(),
			_rules: getRules(),
			_masks: getMasks(),
			_hashes: getHashes(),
			_algorithms: getAlgorithms()
		}
	}

	onClickRescan() {
		if (typeof window.GOscan !== "function") {
			message.error("GOscan is not a function");
			return;
		}

		this.setState({isLoadingRescan: true}, async () => {
			try {
				await window.GOscan();
				await data.getDictionaries();
				await data.getRules();
				await data.getMasks();
				await data.getHashes();
				await data.getAlgorithms();	
			} catch (e) {
				message.error(e.toString());
			}
			this.setState({isLoadingRescan: false});
		})
	}

	onClickRefreshTaskCounter() {
		if (typeof window.GOsettingsCurrentTaskCounter !== "function") {
			message.error("GOsettingsCurrentTaskCounter is not a function");
			return;
		}

		this.setState({isLoadingRefreshTaskCounter: true}, () => {
			window.GOsettingsCurrentTaskCounter().then(
				response => {
					this.setState({
						taskCounter: response,
						isLoadingRefreshTaskCounter: false
					});
				},
				error => {
					message.error(error);
					this.setState({isLoadingRefreshTaskCounter: false});
				}
			);
		})
	}

	onClickResetTaskCounter() {
		if (typeof window.GOsettingsResetTaskCounter !== "function") {
			message.error("GOsettingsResetTaskCounter is not a function");
			return;
		}

		this.setState({isLoadingResetTaskCounter: true}, () => {
			window.GOsettingsResetTaskCounter().then(
				response => {
					this.setState({
						taskCounter: response,
						isLoadingResetTaskCounter: false
					});
				},
				error => {
					message.error(error);
					this.setState({isLoadingResetTaskCounter: false});
				}
			);
		})
	}

	componentDidMount() {
		EventBus.on("dataUpdate", "Settings", () => {
			this.setState({
				_dictionaries: getDictionaries(),
				_rules: getRules(),
				_masks: getMasks(),
				_hashes: getHashes(),
				_algorithms: getAlgorithms()
			});
		});
	}

	componentWillUnmount() {
		EventBus.remove("dataUpdate", "Settings");
	}

	render() {
		return (
			<>
				<PageHeader
					title="Settings"
				/>
				<Content style={{ padding: '16px 24px' }}>
					<Row gutter={[16, 14]}>
						<Col span={12}>
							<Statistic title="Hashes" value={this.state._hashes.length} />
						</Col>
						<Col span={12}>
							<Statistic title="Algorithms" value={Object.keys(this.state._algorithms).length} />
						</Col>
						<Col span={12}>
							<Statistic title="Dictionaries" value={this.state._dictionaries.length} />
						</Col>
						<Col span={12}>
							<Statistic title="Rules" value={this.state._rules.length} />
						</Col>
						<Col span={12}>
							<Statistic title="Masks" value={this.state._masks.length} />
						</Col>
						<Col span={24}>
							<Button
								icon={<SyncOutlined />}
								type="primary"
								onClick={this.onClickRescan}
								loading={this.state.isLoadingRescan}
							>
								Rescan
							</Button>
						</Col>
					</Row>
					<Row style={{ marginTop: "2rem" }} gutter={[16, 14]}>
						<Col span={24}>
							<Statistic
								title="Task counter"
								value={this.state.taskCounter}
							/>
							<Space>
							<Button
								style={{ marginTop: 16 }}
								type="default"
								onClick={this.onClickRefreshTaskCounter}
								loading={this.state.isLoadingRefreshTaskCounter}
							>
								Refresh
							</Button>
							<Popconfirm
								placement="topRight"
								title="Are you sure you want to reset the task counter?"
								onConfirm={this.onClickResetTaskCounter}
								okText="Yes"
								cancelText="No"
							>
								<Button
									style={{ marginTop: 16 }}
									type="danger"
									loading={this.state.isLoadingResetTaskCounter}
								>
									Reset counter
								</Button>
							</Popconfirm>
							</Space>
						</Col>
					</Row>
				</Content>
			</>
		)
	}
}

export default Settings;
