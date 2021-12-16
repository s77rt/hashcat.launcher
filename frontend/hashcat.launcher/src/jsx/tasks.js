import React, { Component } from "react";
import { Layout, PageHeader, Popconfirm, Tag, List, InputNumber, Table, Modal, message, Progress, Badge, Descriptions, Tree, Row, Col, Card, Select, Typography, Upload, Button, Space, Input, Form, Radio, Divider, Collapse, Checkbox, Tabs, Steps } from 'antd';
import {
	FileDoneOutlined,
	FileOutlined,
	AimOutlined,
	ToolOutlined,
	ExportOutlined,
	ExperimentOutlined,
	ReloadOutlined,
	PushpinOutlined,
	PlayCircleOutlined,
	PauseOutlined,
	CaretRightOutlined,
	StepForwardOutlined,
	CloseOutlined,
	ControlOutlined,
	CodeOutlined,
	EnvironmentOutlined
} from '@ant-design/icons';

import EventBus from "./eventbus/EventBus";
import TasksStats from './stats/tasks';

const { Content } = Layout;
const { Dragger } = Upload;
const { Option } = Select;
const { Text, Title, Paragraph } = Typography;
const { Panel } = Collapse;
const { TabPane } = Tabs;
const { Step } = Steps;
const { TreeNode } = Tree;

// https://github.com/hashcat/hashcat/blob/master/include/types.h
const HASHCAT_STATUS_INIT             = 0,
	HASHCAT_STATUS_AUTOTUNE           = 1,
	HASHCAT_STATUS_SELFTEST           = 2,
	HASHCAT_STATUS_RUNNING            = 3,
	HASHCAT_STATUS_PAUSED             = 4,
	HASHCAT_STATUS_EXHAUSTED          = 5,
	HASHCAT_STATUS_CRACKED            = 6,
	HASHCAT_STATUS_ABORTED            = 7,
	HASHCAT_STATUS_QUIT               = 8,
	HASHCAT_STATUS_BYPASS             = 9,
	HASHCAT_STATUS_ABORTED_CHECKPOINT = 10,
	HASHCAT_STATUS_ABORTED_RUNTIME    = 11,
	HASHCAT_STATUS_ERROR              = 13,
	HASHCAT_STATUS_ABORTED_FINISH     = 14,
	HASHCAT_STATUS_AUTODETECT         = 16;

const HASHCAT_STATUS_MESSAGES = {
	[HASHCAT_STATUS_INIT]               : "Init",
	[HASHCAT_STATUS_AUTOTUNE]           : "Autotune",
	[HASHCAT_STATUS_SELFTEST]           : "Selftest",
	[HASHCAT_STATUS_RUNNING]            : "Running",
	[HASHCAT_STATUS_PAUSED]             : "Paused",
	[HASHCAT_STATUS_EXHAUSTED]          : "Exhausted",
	[HASHCAT_STATUS_CRACKED]            : "Cracked",
	[HASHCAT_STATUS_ABORTED]            : "Aborted",
	[HASHCAT_STATUS_QUIT]               : "Quit",
	[HASHCAT_STATUS_BYPASS]             : "Bypass",
	[HASHCAT_STATUS_ABORTED_CHECKPOINT] : "Aborted (Checkpoint)",
	[HASHCAT_STATUS_ABORTED_RUNTIME]    : "Aborted (Runtime)",
	[HASHCAT_STATUS_ERROR]              : "Error",
	[HASHCAT_STATUS_ABORTED_FINISH]     : "Aborted (Finish)",
	[HASHCAT_STATUS_AUTODETECT]         : "Autodetect"
};

const HASHCAT_STATUS_BADGE_WARNING = [HASHCAT_STATUS_PAUSED];
const HASHCAT_STATUS_BADGE_PROCESSING = [HASHCAT_STATUS_RUNNING];
const HASHCAT_STATUS_BADGE_ERROR = [HASHCAT_STATUS_ABORTED, HASHCAT_STATUS_ABORTED_CHECKPOINT, HASHCAT_STATUS_ABORTED_RUNTIME, HASHCAT_STATUS_ABORTED_FINISH, HASHCAT_STATUS_QUIT, HASHCAT_STATUS_ERROR];
const HASHCAT_STATUS_BADGE_SUCCESS = [HASHCAT_STATUS_CRACKED];

class Tasks extends Component {
	constructor(props) {
		super(props);

		this.onSelect = this.onSelect.bind(this);

		this.onClickStart = this.onClickStart.bind(this);
		this.onClickRefresh = this.onClickRefresh.bind(this);
		this.onClickPause = this.onClickPause.bind(this);
		this.onClickResume = this.onClickResume.bind(this);
		this.onClickCheckpoint = this.onClickCheckpoint.bind(this);
		this.onClickSkip = this.onClickSkip.bind(this);
		this.onClickQuit = this.onClickQuit.bind(this);

		this.onClickArguments = this.onClickArguments.bind(this);

		this.state = {
			data: [],

			taskKey: undefined,
			task: undefined,

			isLoadingStart: false,
			isLoadingRefresh: false,
			isLoadingPause: false,
			isLoadingResume: false,
			isLoadingCheckpoint: false,
			isLoadingSkip: false,
			isLoadingQuit: false
		};
	}

	onSelect(keys) {
		const taskKey = keys.shift();
		this.setState({
			taskKey: taskKey,
			task: TasksStats.tasks[taskKey]
		})
	}

	onClickStart() {
		const task = this.state.task;
		if (!task) {
			message.error("no task is selected");
			return;
		}

		if (typeof window.GOstartTask !== "function") {
			message.error("GOstartTask is not a function");
			return;
		}

		this.setState({isLoadingStart: true}, () => {
			window.GOstartTask(task.id).then(
				response => {
					this.setState({isLoadingStart: false});
				},
				error => {
					message.error(error);
					this.setState({isLoadingStart: false});
				}
			);
		})
	}

	onClickRefresh() {
		const task = this.state.task;
		if (!task) {
			message.error("no task is selected");
			return;
		}

		if (typeof window.GOrefreshTask !== "function") {
			message.error("GOrefreshTask is not a function");
			return;
		}

		this.setState({isLoadingRefresh: true}, () => {
			window.GOrefreshTask(task.id).then(
				response => {
					this.setState({isLoadingRefresh: false});
				},
				error => {
					message.error(error);
					this.setState({isLoadingRefresh: false});
				}
			);
		})
	}

	onClickPause() {
		const task = this.state.task;
		if (!task) {
			message.error("no task is selected");
			return;
		}

		if (typeof window.GOpauseTask !== "function") {
			message.error("GOpauseTask is not a function");
			return;
		}

		this.setState({isLoadingPause: true}, () => {
			window.GOpauseTask(task.id).then(
				response => {
					this.setState({isLoadingPause: false});
				},
				error => {
					message.error(error);
					this.setState({isLoadingPause: false});
				}
			);
		})
	}

	onClickResume() {
		const task = this.state.task;
		if (!task) {
			message.error("no task is selected");
			return;
		}

		if (typeof window.GOresumeTask !== "function") {
			message.error("GOresumeTask is not a function");
			return;
		}

		this.setState({isLoadingResume: true}, () => {
			window.GOresumeTask(task.id).then(
				response => {
					this.setState({isLoadingResume: false});
				},
				error => {
					message.error(error);
					this.setState({isLoadingResume: false});
				}
			);
		})
	}

	onClickCheckpoint() {
		const task = this.state.task;
		if (!task) {
			message.error("no task is selected");
			return;
		}

		if (typeof window.GOcheckpointTask !== "function") {
			message.error("GOcheckpointTask is not a function");
			return;
		}

		this.setState({isLoadingCheckpoint: true}, () => {
			window.GOcheckpointTask(task.id).then(
				response => {
					this.setState({isLoadingCheckpoint: false});
				},
				error => {
					message.error(error);
					this.setState({isLoadingCheckpoint: false});
				}
			);
		})
	}

	onClickSkip() {
		const task = this.state.task;
		if (!task) {
			message.error("no task is selected");
			return;
		}

		if (typeof window.GOskipTask !== "function") {
			message.error("GOskipTask is not a function");
			return;
		}

		this.setState({isLoadingSkip: true}, () => {
			window.GOskipTask(task.id).then(
				response => {
					this.setState({isLoadingSkip: false});
				},
				error => {
					message.error(error);
					this.setState({isLoadingSkip: false});
				}
			);
		})
	}

	onClickQuit() {
		const task = this.state.task;
		if (!task) {
			message.error("no task is selected");
			return;
		}

		if (typeof window.GOquitTask !== "function") {
			message.error("GOquitTask is not a function");
			return;
		}

		this.setState({isLoadingQuit: true}, () => {
			window.GOquitTask(task.id).then(
				response => {
					this.setState({isLoadingQuit: false});
				},
				error => {
					message.error(error);
					this.setState({isLoadingQuit: false});
				}
			);
		})
	}

	onClickArguments() {
		const task = this.state.task;
		if (!task) {
			message.error("no task is selected");
			return;
		}

		Modal.info({
			title: 'Arguments',
			content: (
				<div style={{ maxHeight: '300px' }}>
					<Text code copyable>
						{task.arguments.join(" ")}
					</Text>
				</div>
			),
		});
	}

	reBuildData() {
		var data = [];

		Object.values(TasksStats.tasks).forEach(task => {
			data.push({
				key: task.id,
				title: (
					task.stats.hasOwnProperty("progress") ? (
						task.id + " (" + Math.trunc((task.stats["progress"][0] / task.stats["progress"][1])*100) + "%)"
					) : (
						task.id
					)
				),
				icon: (
					task.stats.hasOwnProperty("status") ? (
						HASHCAT_STATUS_BADGE_WARNING.indexOf(task.stats["status"]) > -1 ? (
							<Badge status="warning" />
						) : HASHCAT_STATUS_BADGE_PROCESSING.indexOf(task.stats["status"]) > -1 ? (
							<Badge status="processing" />
						) : HASHCAT_STATUS_BADGE_ERROR.indexOf(task.stats["status"]) > -1 ? (
							<Badge status="error" />
						) : HASHCAT_STATUS_BADGE_SUCCESS.indexOf(task.stats["status"]) > -1 ? (
							<Badge status="success" />
						) : (
							<Badge status="default" />
						)
					) : (
						<Badge status="default" />
					)
				),
			});
		});

		this.setState({
			data: data
		});
	}

	componentDidMount() {
		EventBus.on("tasksUpdate", () => {
			this.reBuildData();
		});

		this.reBuildData();
	}

	componentWillUnmount() {
		EventBus.remove("tasksUpdate");
	}

	render() {
		const { taskKey, task } = this.state;

		return (
			<>
				<PageHeader
					title="Tasks"
				/>
				<Content style={{ padding: '16px 24px' }}>
					<Row gutter={16}>
						<Col span={5}>
							<Tree
								showIcon
								treeData={this.state.data}
								onSelect={this.onSelect}
								selectedKeys={[taskKey]}
								style={{
									height: 'calc(100vh - 195px)',
									paddingRight: '.5rem',
									overflow: 'auto',
									background: '#0a0a0a',
									border: '1px solid #303030'
								}}
							/>
						</Col>
						<Col span={19}>
							{task ? (
								<Row gutter={[16, 14]}>
									<Col span={24}>
										<PageHeader
											title={task.id}
											tags={
												task.stats.hasOwnProperty("status") ? (	
													HASHCAT_STATUS_BADGE_WARNING.indexOf(task.stats["status"]) > -1 ? (
														<Tag color="warning">{HASHCAT_STATUS_MESSAGES[task.stats["status"]]}</Tag>
													) : HASHCAT_STATUS_BADGE_PROCESSING.indexOf(task.stats["status"]) > -1 ? (
														<Tag color="processing">{HASHCAT_STATUS_MESSAGES[task.stats["status"]]}</Tag>
													) : HASHCAT_STATUS_BADGE_ERROR.indexOf(task.stats["status"]) > -1 ? (
														<Tag color="error">{HASHCAT_STATUS_MESSAGES[task.stats["status"]]}</Tag>
													) : HASHCAT_STATUS_BADGE_SUCCESS.indexOf(task.stats["status"]) > -1 ? (
														<Tag color="success">{HASHCAT_STATUS_MESSAGES[task.stats["status"]]}</Tag>
													) : (
														<Tag color="default">{HASHCAT_STATUS_MESSAGES[task.stats["status"]]}</Tag>
													)
												) : null
											}
											style={{ padding: 0 }}
											extra={[
												<Button
													type="dashed"
													icon={<ControlOutlined />}
													onClick={this.onClickArguments}
													style={{ marginRight: '1rem' }}
													key="arguments"
												>
													Arguments
												</Button>
											]}
										/>
									</Col>
									<Col span={24}>
										{task.stats.hasOwnProperty("progress") ? (
											<Progress type="line" percent={Math.trunc((task.stats["progress"][0] / task.stats["progress"][1])*100)} />
										) : (
											<Progress type="line" percent={0} />
										)}
									</Col>
									<Col span={24}>
										<Row gutter={[12, 10]}>
											<Col>
												<Button
													type="primary"
													icon={<PlayCircleOutlined />}
													onClick={this.onClickStart}
													loading={this.state.isLoadingStart}
												>
													Start
												</Button>
											</Col>
											<Col>
												<Button
													icon={<ReloadOutlined />}
													onClick={this.onClickRefresh}
													loading={this.state.isLoadingRefresh}
												>
													Refresh
												</Button>
											</Col>
											<Col>
												<Button
													icon={<PauseOutlined />}
													onClick={this.onClickPause}
													loading={this.state.isLoadingPause}
												>
													Pause
												</Button>
											</Col>
											<Col>
												<Button
													icon={<CaretRightOutlined />}
													onClick={this.onClickResume}
													loading={this.state.isLoadingResume}
												>
													Resume
												</Button>
											</Col>
											<Col>
												<Button
													icon={<EnvironmentOutlined />}
													onClick={this.onClickCheckpoint}
													loading={this.state.isLoadingCheckpoint}
												>
													Checkpoint
												</Button>
											</Col>
											<Col>
												<Button
													icon={<StepForwardOutlined />}
													onClick={this.onClickSkip}
													loading={this.state.isLoadingSkip}
												>
													Skip
												</Button>
											</Col>
											<Col>
												<Popconfirm
													placement="topRight"
													title="Are you sure to quit the task?"
													onConfirm={this.onClickQuit}
													okText="Yes"
													cancelText="No"
												>
													<Button
														type="danger"
														icon={<CloseOutlined />}
														
														loading={this.state.isLoadingQuit}
													>
														Quit
													</Button>
												</Popconfirm>
											</Col>
										</Row>
									</Col>
									<Col span={16}>
										<Descriptions
											column={2}
											layout="horizontal"
											bordered
										>
											{task.stats.hasOwnProperty("status") && (
												<Descriptions.Item label="Status" span={2}>
													{HASHCAT_STATUS_BADGE_WARNING.indexOf(task.stats["status"]) > -1 ? (
														<Badge status="warning" text={HASHCAT_STATUS_MESSAGES[task.stats["status"]]} />
													) : HASHCAT_STATUS_BADGE_PROCESSING.indexOf(task.stats["status"]) > -1 ? (
														<Badge status="processing" text={HASHCAT_STATUS_MESSAGES[task.stats["status"]]} />
													) : HASHCAT_STATUS_BADGE_ERROR.indexOf(task.stats["status"]) > -1 ? (
														<Badge status="error" text={HASHCAT_STATUS_MESSAGES[task.stats["status"]]} />
													) : HASHCAT_STATUS_BADGE_SUCCESS.indexOf(task.stats["status"]) > -1 ? (
														<Badge status="success" text={HASHCAT_STATUS_MESSAGES[task.stats["status"]]} />
													) : (
														<Badge status="default" text={HASHCAT_STATUS_MESSAGES[task.stats["status"]]} />
													)}
												</Descriptions.Item>
											)}
											{task.stats.hasOwnProperty("target") && (
												<Descriptions.Item label="Target" span={2}>
													{task.stats["target"]}
												</Descriptions.Item>
											)}
											{task.stats.hasOwnProperty("progress") && (
												<Descriptions.Item label="Progress" span={2}>
													{task.stats["progress"][0] + " / " + task.stats["progress"][1] + " (" + Math.trunc((task.stats["progress"][0] / task.stats["progress"][1])*100) + "%)"}
												</Descriptions.Item>
											)}
											{task.stats.hasOwnProperty("rejected") && (
												<Descriptions.Item label="Rejected" span={1}>
													{task.stats["rejected"]}
												</Descriptions.Item>
											)}
											{task.stats.hasOwnProperty("restore_point") && (
												<Descriptions.Item label="Restore point" span={1}>
													{task.stats["restore_point"]}
												</Descriptions.Item>
											)}
											{task.stats.hasOwnProperty("recovered_hashes") && (
												<Descriptions.Item label="Recovered hashes" span={1}>
													{task.stats["recovered_hashes"][0] + " / " + task.stats["recovered_hashes"][1] + " (" + Math.trunc((task.stats["recovered_hashes"][0] / task.stats["recovered_hashes"][1])*100) + "%)"}
												</Descriptions.Item>
											)}
											{task.stats.hasOwnProperty("recovered_salts") && (
												<Descriptions.Item label="Recovered salts" span={1}>
													{task.stats["recovered_salts"][0] + " / " + task.stats["recovered_salts"][1] + " (" + Math.trunc((task.stats["recovered_salts"][0] / task.stats["recovered_salts"][1])*100) + "%)"}
												</Descriptions.Item>
											)}
											{task.stats.hasOwnProperty("session") && (
												<Descriptions.Item label="Session" span={2}>
													{task.stats["session"]}
												</Descriptions.Item>
											)}
										</Descriptions>
									</Col>
									<Col span={8}>
										<CodeOutlined /> Terminal
										<pre style={{
											height: 'calc(100vh - 395px)',
											overflow: 'auto',
											padding: '.5rem',
											margin: '0',
											border: '1px solid #303030'
										}}>
											{task.journal.map(j => j.message + "\n")}
										</pre>
									</Col>
								</Row>
							) : (
								"No selected task."
							)}
						</Col>
					</Row>
				</Content>
			</>
		)
	}
}

export default Tasks;
