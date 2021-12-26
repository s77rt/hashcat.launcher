import React, { Component } from "react";
import { Tooltip, Layout, PageHeader, Popconfirm, Tag, List, InputNumber, Table, Modal, message, Progress, Badge, Descriptions, Tree, Row, Col, Card, Select, Typography, Upload, Button, Space, Input, Form, Radio, Divider, Collapse, Checkbox, Tabs, Steps } from 'antd';
import {
	InfoCircleOutlined,
	DeleteOutlined,
	LineOutlined,
	FireOutlined,
	ClockCircleOutlined,
	CheckCircleOutlined,
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

import moment from "moment"

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
const HASHCAT_STATUS_BADGE_PINK = [HASHCAT_STATUS_EXHAUSTED];

const PROCESS_STATUS_NOTSTARTED = 0,
	PROCESS_STATUS_RUNNING = 1,
	PROCESS_STATUS_FINISHED = 2;

function totalSpeed(devices) {
	var speed = 0;
	devices.forEach(device => {
		speed += device.speed;
	});
	return speed;
}

function humanizeSpeed(H) {
	let KH = 1000;
	let MH = KH * KH;
	let GH = MH * KH;
	let TH = GH * KH;

	if (H < KH*100) {
		return `${(H).toFixed(0)} H/s`
	} else if (H < MH) {
		return `${(H/KH).toFixed(1)} kH/s`
	} else if (H < GH) {
		return `${(H/MH).toFixed(1)} MH/s`
	} else if (H < TH) {
		return `${(H/GH).toFixed(1)} GH/s`
	} else {
		return `${(H/TH).toFixed(1)} TH/s`
	}
}

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

		this.onChangePriority = this.onChangePriority.bind(this);

		this.onClickArguments = this.onClickArguments.bind(this);
		this.onClickDelete = this.onClickDelete.bind(this);

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
			isLoadingQuit: false,

			isReadOnlyPriority: false,
			isLoadingDelete: false
		};
	}

	refresh() {
		this.setState({
			task: TasksStats.tasks[this.state.taskKey]
		});
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

	onChangePriority(priority) {
		if (typeof(priority) !== "number")
			return

		const task = this.state.task;
		if (!task) {
			message.error("no task is selected");
			return;
		}

		if (typeof window.GOpriorityTask !== "function") {
			message.error("GOpriorityTask is not a function");
			return;
		}

		this.setState({isReadOnlyPriority: true}, () => {
			window.GOpriorityTask(task.id, priority).then(
				response => {
					task.priority = priority;
					this.reBuildData();
					this.setState({isReadOnlyPriority: false});
					if (task.priority >= 0 && this.state.data[2].children.length === 0) {
						if (typeof window.GOstartNextTask === "function") {
							window.GOstartNextTask();
						}
					}
				},
				error => {
					message.error(error);
					this.setState({isReadOnlyPriority: false});
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
				<div style={{ maxHeight: '300px', overflow: 'auto' }}>
					<Text code copyable>
						{task.arguments.join(" ")}
					</Text>
				</div>
			),
		});
	}

	onClickDelete() {
		const task = this.state.task;
		if (!task) {
			message.error("no task is selected");
			return;
		}

		if (typeof window.GOdeleteTask !== "function") {
			message.error("GOdeleteTask is not a function");
			return;
		}

		this.setState({isLoadingDelete: true}, () => {
			window.GOdeleteTask(task.id).then(
				response => {
					this.refresh();
					this.setState({isLoadingDelete: false});
				},
				error => {
					message.error(error);
					this.setState({isLoadingDelete: false});
				}
			);
		})
	}

	reBuildData() {
		var data = [
			{
				key: "Idle",
				title: "Idle",
				selectable: false,
				icon: <LineOutlined />,
				children: []
			},
			{
				key: "Queued",
				title: "Queued",
				selectable: false,
				icon: <ClockCircleOutlined />,
				children: []
			},
			{
				key: "In Progress",
				title: "In Progress",
				selectable: false,
				icon: <FireOutlined />,
				children: []
			},
			{
				key: "Finished",
				title: "Finished",
				selectable: false,
				icon: <CheckCircleOutlined />,
				children: []
			}
		];

		Object.values(TasksStats.tasks).forEach(task => {
			var category;
			switch (task.process.status) {
				case PROCESS_STATUS_NOTSTARTED:
					if (task.priority >= 0)
						category = data[1];
					else
						category = data[0];
					break;
				case PROCESS_STATUS_RUNNING:
					category = data[2];
					break;
				case PROCESS_STATUS_FINISHED:
					category = data[3];
					break;
				default:
					category = data[0];
					message.warning("unrecognized process status");
			}
			category.children.push({
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
						) : HASHCAT_STATUS_BADGE_PINK.indexOf(task.stats["status"]) > -1 ? (
							<Badge color="pink" />
						) : (
							<Badge status="default" />
						)
					) : (
						<Badge color="#b7b7b7" />
					)
				),
			});
		});

		for (let i = 0; i < data.length; i++) {
			data[i].title = data[i].title + " (" + data[i].children.length + ")";
		}

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
					<Row gutter={16} className="height-100 tree-height-100">
						<Col span={5}>
							<Tree
								showIcon
								blockNode
								treeData={this.state.data}
								onSelect={this.onSelect}
								selectedKeys={[taskKey]}
								style={{
									height: '100%',
									paddingRight: '.5rem',
									overflow: 'auto',
									background: '#0a0a0a',
									border: '1px solid #303030'
								}}
							/>
						</Col>
						<Col span={19}>
							{task ? (
								<Row gutter={[16, 14]} className="height-100" style={{ flexDirection: "column" }}>
									<Col flex="0 0 auto">
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
															) : HASHCAT_STATUS_BADGE_PINK.indexOf(task.stats["status"]) > -1 ? (
																<Tag color="pink">{HASHCAT_STATUS_MESSAGES[task.stats["status"]]}</Tag>
															) : (
																<Tag color="default">{HASHCAT_STATUS_MESSAGES[task.stats["status"]]}</Tag>
															)
														) : null
													}
													style={{ padding: 0 }}
													extra={
														<Form layout="inline">
															<Form.Item
																label="Priority"
															>
																<InputNumber
																	min={-1}
																	max={999}
																	value={task.priority}
																	onChange={this.onChangePriority}
																	readOnly={this.state.isReadOnlyPriority}
																	bordered={false}
																/>
															</Form.Item>
															<Button
																icon={<ControlOutlined />}
																onClick={this.onClickArguments}
																style={{ marginRight: '1rem' }}
															>
																Arguments
															</Button>
															<Popconfirm
																placement="topRight"
																title="Are you sure you want to delete this task?"
																onConfirm={this.onClickDelete}
																okText="Yes"
																cancelText="No"
															>
																<Button
																	type="danger"
																	icon={<DeleteOutlined />}
																	loading={this.state.isLoadingDelete}
																>
																	Delete
																</Button>
															</Popconfirm>
														</Form>
													}
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
															title="Are you sure you want to quit this task?"
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
										</Row>
									</Col>
									<Col flex="1 1 auto">
										<Row gutter={[16, 14]} className="height-100">
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
															) : HASHCAT_STATUS_BADGE_PINK.indexOf(task.stats["status"]) > -1 ? (
																<Badge color="pink" text={HASHCAT_STATUS_MESSAGES[task.stats["status"]]} />
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
															{task.stats.hasOwnProperty("guess") && (
																<Tooltip title={
																	<Descriptions bordered size="small" column={1} layout="horizontal">
																		{task.stats.guess.guess_base !== null ? (
																			<Descriptions.Item label="Guess Base">{task.stats.guess.guess_base} ({task.stats.guess.guess_base_offset}/{task.stats.guess.guess_base_count})</Descriptions.Item>
																		) : (
																			<Descriptions.Item label="Guess Base">-</Descriptions.Item>
																		)}
																		{task.stats.guess.guess_mod !== null ? (
																			<Descriptions.Item label="Guess Mod">{task.stats.guess.guess_mod} ({task.stats.guess.guess_mod_offset}/{task.stats.guess.guess_mod_count})</Descriptions.Item>
																		) : (
																			<Descriptions.Item label="Guess Mod">-</Descriptions.Item>
																		)}
																	</Descriptions>
																}>
																	<InfoCircleOutlined style={{ marginLeft: ".5rem" }} />
																</Tooltip>
															)}
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
													{task.stats.hasOwnProperty("devices") && (
														<Descriptions.Item label="Speed" span={2}>
															{humanizeSpeed(totalSpeed(task.stats["devices"]))}
															<Tooltip title={
																<Table
																	columns={[
																		{
																			title: 'ID',
																			dataIndex: 'id',
																			key: 'ID'
																		},
																		{
																			title: 'Speed',
																			dataIndex: 'speed',
																			key: 'Speed'
																		},
																		{
																			title: 'Temp',
																			dataIndex: 'temp',
																			key: 'Temp'
																		},
																		{
																			title: 'Util',
																			dataIndex: 'util',
																			key: 'Util'
																		}
																	]}
																	dataSource={task.stats["devices"].map(device =>
																		({
																			key: device.device_id,
																			id: device.device_id,
																			speed: humanizeSpeed(device.speed),
																			temp: device.hasOwnProperty("temp") ? device.temp + " °C": "-",
																			util: device.util + "%",
																		})
																	)}
																	size="small"
																	pagination={false}
																	style={{ overflow: 'auto' }}
																/>
															}>
																<InfoCircleOutlined style={{ marginLeft: ".5rem" }} />
															</Tooltip>
														</Descriptions.Item>
													)}
													{task.stats.hasOwnProperty("time_start") && (
														<Descriptions.Item label="Started" span={1}>
															<Tooltip title={moment.unix(task.stats["time_start"]).format("MMMM Do YYYY, HH:mm")}>
																{moment.unix(task.stats["time_start"]).fromNow()}
															</Tooltip>
														</Descriptions.Item>
													)}
													{task.stats.hasOwnProperty("estimated_stop") && (
														<Descriptions.Item label="ETA" span={1}>
															<Tooltip title={moment.unix(task.stats["estimated_stop"]).format("MMMM Do YYYY, HH:mm")}>
																{moment.unix(task.stats["estimated_stop"]).fromNow()}
															</Tooltip>
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
												<div className="height-100" style={{ display: "flex", flexDirection: "column" }}>
												<span><CodeOutlined /> Terminal</span>
												<pre style={{
													flex: '1 0 auto',
													overflow: 'auto',
													padding: '.5rem',
													margin: '0',
													border: '1px solid #303030'
												}}>
													{task.journal.map(j => j.message + "\n")}
												</pre>
												</div>
											</Col>
										</Row>
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
