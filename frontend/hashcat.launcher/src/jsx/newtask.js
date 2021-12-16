import React, { Component } from "react";
import { Layout, PageHeader, Descriptions, Table, InputNumber, message, Row, Col, Card, Select, Typography, Upload, Button, Space, Input, Form, Radio, Divider, Collapse, Checkbox, Tabs, Steps } from 'antd';
import {
	ImportOutlined,
	ExportOutlined,
	UndoOutlined,
	PlusOutlined,
	FileDoneOutlined,
	AimOutlined,
	ToolOutlined,
	ExperimentOutlined
} from '@ant-design/icons';

import EventBus from "./eventbus/EventBus";
import { getHashes } from './data/hashes';
import { getAlgorithms } from './data/algorithms';
import { getDictionaries } from './data/dictionaries';
import { getRules } from './data/rules';
import { getMasks } from './data/masks';

const { Content } = Layout;
const { Option } = Select;
const { Title } = Typography;
const { Panel } = Collapse;
const { TabPane } = Tabs;
const { Step } = Steps;

const maxRules = 4;
const maskIncrementMin = 1;
const maskIncrementMax = 16;

function filename(path) {
	return path.split('\\').pop().split('/').pop();
}

function maskLength(mask) {
	if (typeof(mask) !== "string")
		return 0;

	if (mask.length === 0)
		return 0;

	var length = 0;
	var skipNext = false;
	for (var i = 0; i < mask.length; i++) {
		if (skipNext === true) {
			skipNext = false;
			continue;
		}

		let char = mask.charAt(i);
		if (char === "?") {
			skipNext = true;
		}
		length++;
	}

	return length
}

const initialConfig = {
	step: 0,
	advancedOptionsCollapse: [],

	maskInputType: "text", // text or file

	attackMode: undefined,
	algorithm: undefined,
	dictionaries: undefined,
	rules: undefined,
	mask: undefined,
	maskFile: undefined,
	leftDictionary: undefined,
	leftRule: undefined,
	rightDictionary: undefined,
	rightRule: undefined,

	customCharset1: undefined,
	customCharset2: undefined,
	customCharset3: undefined,
	customCharset4: undefined,

	enableMaskIncrementMode: false,
	maskIncrementMin: maskIncrementMin,
	maskIncrementMax: maskIncrementMax,

	hash: undefined,

	enableOptimizedKernel: true,
	enableSlowerCandidateGenerators: false,
	removeFoundHashes: false,
	disablePotFile: false,
	ignoreUsernames: false,
	disableSelfTest: false,
	ignoreWarnings: false,

	devicesIDs: undefined,
	devicesTypes: undefined,
	workloadProfile: undefined,

	disableMonitor: false,
	tempAbort: undefined,

	markovDisable: false,
	markovClassic: false,
	markovThreshold: 0,

	extraArguments: ["--quiet", "--logfile-disable"],

	statusTimer: 20,

	outputFile: undefined,
	outputFormat: [1, 2]
}

class NewTask extends Component {
	constructor(props) {
		super(props);

		this.onChangeStep = this.onChangeStep.bind(this);
		this.onChangeAdvancedOptionsCollapse = this.onChangeAdvancedOptionsCollapse.bind(this);

		this.onChangeMaskInputType = this.onChangeMaskInputType.bind(this);

		this.onChangeAttackMode = this.onChangeAttackMode.bind(this);
		this.onChangeAlgorithm = this.onChangeAlgorithm.bind(this);
		this.onChangeDictionaries = this.onChangeDictionaries.bind(this);
		this.onChangeRules = this.onChangeRules.bind(this);
		this.onChangeMask = this.onChangeMask.bind(this);
		this.onChangeMaskFile = this.onChangeMaskFile.bind(this);
		this.onChangeLeftDictionary = this.onChangeLeftDictionary.bind(this);
		this.onChangeLeftRule = this.onChangeLeftRule.bind(this);
		this.onChangeRightDictionary = this.onChangeRightDictionary.bind(this);
		this.onChangeRightRule = this.onChangeRightRule.bind(this);

		this.onChangeCustomCharset1 = this.onChangeCustomCharset1.bind(this);
		this.onChangeCustomCharset2 = this.onChangeCustomCharset2.bind(this);
		this.onChangeCustomCharset3 = this.onChangeCustomCharset3.bind(this);
		this.onChangeCustomCharset4 = this.onChangeCustomCharset4.bind(this);

		this.onChangeEnableMaskIncrementMode = this.onChangeEnableMaskIncrementMode.bind(this);
		this.onChangeMaskIncrementMin = this.onChangeMaskIncrementMin.bind(this);
		this.onChangeMaskIncrementMax = this.onChangeMaskIncrementMax.bind(this);

		this.onChangeHash = this.onChangeHash.bind(this);

		this.onChangeEnableOptimizedKernel = this.onChangeEnableOptimizedKernel.bind(this);
		this.onChangeEnableSlowerCandidateGenerators = this.onChangeEnableSlowerCandidateGenerators.bind(this);
		this.onChangeRemoveFoundHashes = this.onChangeRemoveFoundHashes.bind(this);
		this.onChangeDisablePotFile = this.onChangeDisablePotFile.bind(this);
		this.onChangeIgnoreUsernames = this.onChangeIgnoreUsernames.bind(this);
		this.onChangeDisableSelfTest = this.onChangeDisableSelfTest.bind(this);
		this.onChangeIgnoreWarnings = this.onChangeIgnoreWarnings.bind(this);

		this.onChangeDevicesIDs = this.onChangeDevicesIDs.bind(this);
		this.onChangeDevicesTypes = this.onChangeDevicesTypes.bind(this);
		this.onChangeWorkloadProfile = this.onChangeWorkloadProfile.bind(this);

		this.onChangeDisableMonitor = this.onChangeDisableMonitor.bind(this);
		this.onChangeTempAbort = this.onChangeTempAbort.bind(this);

		this.onChangeMarkovDisable = this.onChangeMarkovDisable.bind(this);
		this.onChangeMarkovClassic = this.onChangeMarkovClassic.bind(this);
		this.onChangeMarkovThreshold = this.onChangeMarkovThreshold.bind(this);

		this.onChangeExtraArguments = this.onChangeExtraArguments.bind(this);

		this.onChangeStatusTimer = this.onChangeStatusTimer.bind(this);

		this.onChangeOutputFile = this.onChangeOutputFile.bind(this);
		this.onChangeOutputFormat = this.onChangeOutputFormat.bind(this);

		this.onClickCreateTask = this.onClickCreateTask.bind(this);

		this.onChangePreserveTaskConfig = this.onChangePreserveTaskConfig.bind(this);

		this.onClickImportConfig = this.onClickImportConfig.bind(this);
		this.onClickExportConfig = this.onClickExportConfig.bind(this);

		this.state = {
			...initialConfig,

			preserveTaskConfig: false,

			isLoadingImportConfig: false,
			isLoadingExportConfig: false,
			isLoadingSetOutputFile: false,
			isLoadingCreateTask: false,

			_dictionaries: getDictionaries(),
			_rules: getRules(),
			_masks: getMasks(),
			_hashes: getHashes(),
			_algorithms: getAlgorithms()
		}
	}

	onChangeAdvancedOptionsCollapse(e) {
		this.setState({
			advancedOptionsCollapse: e
		});
	}

	onChangeMaskInputType(e) {
		this.setState({
			maskInputType: e.type
		});
	}

	onChangeEnableOptimizedKernel(e) {
		this.setState({
			enableOptimizedKernel: e.target.checked
		});
	}

	onChangeEnableSlowerCandidateGenerators(e) {
		this.setState({
			enableSlowerCandidateGenerators: e.target.checked
		});
	}

	onChangeRemoveFoundHashes(e) {
		this.setState({
			removeFoundHashes: e.target.checked
		});
	}

	onChangeDisablePotFile(e) {
		this.setState({
			disablePotFile: e.target.checked
		});
	}

	onChangeIgnoreUsernames(e) {
		this.setState({
			ignoreUsernames: e.target.checked
		});
	}

	onChangeDisableSelfTest(e) {
		this.setState({
			disableSelfTest: e.target.checked
		});
	}

	onChangeIgnoreWarnings(e) {
		this.setState({
			ignoreWarnings: e.target.checked
		});
	}

	onChangeStep(step) {
		this.setState({
			step: step
		});
	}

	onChangeAttackMode(e) {
		this.setState({
			attackMode: e.target.value
		});
	}

	onChangeAlgorithm(e) {
		this.setState({
			algorithm: e
		});
	}

	onChangeDictionaries(e) {
		this.setState({
			dictionaries: e
		});
	}

	onChangeRules(e) {
		this.setState({
			rules: e.slice(0, maxRules)
		});
	}

	onChangeMask(e) {
		this.setState({
			mask: e.target.value
		});
	}

	onChangeMaskFile(e) {
		this.setState({
			maskFile: e
		});
	}

	onChangeLeftDictionary(e) {
		this.setState({
			leftDictionary: e
		});
	}

	onChangeLeftRule(e) {
		this.setState({
			leftRule: e.target.value
		});
	}

	onChangeRightDictionary(e) {
		this.setState({
			rightDictionary: e
		});
	}

	onChangeRightRule(e) {
		this.setState({
			rightRule: e.target.value
		});
	}

	onChangeCustomCharset1(e) {
		this.setState({
			customCharset1: e.target.value
		});
	}

	onChangeCustomCharset2(e) {
		this.setState({
			customCharset2: e.target.value
		});
	}

	onChangeCustomCharset3(e) {
		this.setState({
			customCharset3: e.target.value
		});
	}

	onChangeCustomCharset4(e) {
		this.setState({
			customCharset4: e.target.value
		});
	}

	onChangeEnableMaskIncrementMode(e) {
		this.setState({
			enableMaskIncrementMode: e.target.checked
		});
	}

	onChangeMaskIncrementMin(e) {
		this.setState({
			maskIncrementMin: e
		});
	}

	onChangeMaskIncrementMax(e) {
		this.setState({
			maskIncrementMax: e
		});
	}

	onChangeHash(e) {
		this.setState({
			hash: e
		});
	}

	onChangeDevicesIDs(e) {
		this.setState({
			devicesIDs: e
		});
	}

	onChangeDevicesTypes(e) {
		this.setState({
			devicesTypes: e
		});
	}

	onChangeWorkloadProfile(e) {
		this.setState({
			workloadProfile: e
		});
	}

	onChangeDisableMonitor(e) {
		this.setState({
			disableMonitor: e.target.checked
		});
	}

	onChangeTempAbort(e) {
		this.setState({
			tempAbort: e
		});
	}

	onChangeMarkovDisable(e) {
		this.setState({
			markovDisable: e.target.checked
		});
	}

	onChangeMarkovClassic(e) {
		this.setState({
			markovClassic: e.target.checked
		});
	}

	onChangeMarkovThreshold(e) {
		this.setState({
			markovThreshold: e
		});
	}

	onChangeExtraArguments(e) {
		this.setState({
			extraArguments: e.target.value.split(" ")
		});
	}

	onChangeStatusTimer(e) {
		this.setState({
			statusTimer: e
		});
	}

	onChangeOutputFile() {
		if (typeof window.GOsaveDialog !== "function") {
			message.error("GOsaveDialog is not a function");
			return;
		}

		this.setState({isLoadingSetOutputFile: true}, async () => {
			try {
				let outputFile = await window.GOsaveDialog();
				this.setState({
					outputFile: outputFile,
					isLoadingSetOutputFile: false
				});
			} catch (e) {
				message.error(e.toString());
				this.setState({
					isLoadingSetOutputFile: false
				});
			}
		});
	}

	onChangeOutputFormat(e) {
		this.setState({
			outputFormat: e
		});
	}

	onChangePreserveTaskConfig(e) {
		this.setState({
			preserveTaskConfig: e.target.checked
		});
	}

	resetInitialConfig() {
		this.setState(initialConfig);
	}

	onClickCreateTask() {
		if (typeof window.GOcreateTask !== "function") {
			message.error("GOcreateTask is not a function");
			return;
		}

		this.setState({
			isLoadingCreateTask: true
		}, () => {
			window.GOcreateTask({
				attackMode: this.state.attackMode,
				hashMode: this.state.algorithm,

				dictionaries: this.state.dictionaries,
				rules: this.state.rules,
				mask: this.state.maskInputType === "text" ? this.state.mask : undefined,
				maskFile: this.state.maskInputType === "file" ? this.state.maskFile : undefined,
				leftDictionary: this.state.leftDictionary,
				leftRule: this.state.leftRule,
				rightDictionary: this.state.rightDictionary,
				rightRule: this.state.rightRule,

				customCharset1: this.state.customCharset1,
				customCharset2: this.state.customCharset2,
				customCharset3: this.state.customCharset3,
				customCharset4: this.state.customCharset4,

				enableMaskIncrementMode: this.state.enableMaskIncrementMode,
				maskIncrementMin: this.state.maskIncrementMin,
				maskIncrementMax: this.state.maskIncrementMax,

				hash: this.state.hash,

				enableOptimizedKernel: this.state.enableOptimizedKernel,
				enableSlowerCandidateGenerators: this.state.enableSlowerCandidateGenerators,
				removeFoundHashes: this.state.removeFoundHashes,
				disablePotFile: this.state.disablePotFile,
				ignoreUsernames: this.state.ignoreUsernames,
				disableSelfTest: this.state.disableSelfTest,
				ignoreWarnings: this.state.ignoreWarnings,

				devicesIDs: this.state.devicesIDs,
				devicesTypes: this.state.devicesTypes,
				workloadProfile: this.state.workloadProfile,

				disableMonitor: this.state.disableMonitor,
				tempAbort: this.state.tempAbort,

				markovDisable: this.state.markovDisable,
				markovClassic: this.state.markovClassic,
				markovThreshold: this.state.markovThreshold,

				extraArguments: this.state.extraArguments.filter(n => n),

				statusTimer: this.state.statusTimer,

				outputFile: this.state.outputFile,
				outputFormat: this.state.outputFormat,
			}).then(
				response => {
					message.success("Task has been created successfully!");
					this.setState({isLoadingCreateTask: false});
					if (!this.state.preserveTaskConfig)
						this.resetInitialConfig();
				},
				error => {
					message.error(error);
					this.setState({isLoadingCreateTask: false});
				}
			);
		});
	}

	importConfig(config) {
		const newConfig = {...initialConfig};
		Object.keys(newConfig).forEach(key => {
			if (config.hasOwnProperty(key))
				newConfig[key] = config[key];
		});
		this.setState({
			...newConfig
		});
	}

	onClickImportConfig(e) {
		var fileList = e.fileList;

		if (fileList.length === 0)
			return;

		this.setState({
			isLoadingImportConfig: true
		}, () => {
			var file = fileList[0].originFileObj;

			var reader = new FileReader();
			reader.onload = (e) => {
				try {
					const config = JSON.parse(e.target.result);
					this.importConfig(config);
					message.success("Imported!");
				} catch (e) {
					message.error(e.toString());
				}
				this.setState({
					isLoadingImportConfig: false
				});
			};
			reader.readAsText(file);
		});
	}

	exportConfig() {
		const config = {};
		Object.keys(initialConfig).forEach(key => { config[key] = this.state[key] } );
		return config;
	}

	onClickExportConfig() {
		this.setState({
			isLoadingExportConfig: true
		}, () => {
			var dataStr = "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(this.exportConfig(), null, '\t'));
			var downloadAnchorNode = document.createElement('a');
			downloadAnchorNode.style.display = "none";
			downloadAnchorNode.setAttribute("href", dataStr);
			downloadAnchorNode.setAttribute("download", "config.json");
			document.body.appendChild(downloadAnchorNode);
			downloadAnchorNode.click();
			downloadAnchorNode.remove();
			message.success("Exported!");
			this.setState({
				isLoadingExportConfig: false
			});
		});
	}

	componentDidMount() {
		EventBus.on("dataUpdate", () => {
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
		EventBus.remove("dataUpdate");
	}

	render() {
		return (
			<>
				<PageHeader
					title="New Task"
					extra={[
						<Upload
							key="Import Config"
							accept=".json"
							maxCount={1}
							showUploadList={false}
							onChange={this.onClickImportConfig}
							beforeUpload={() => {return false;}}
						>
							<Button
								type="text"
								icon={<ImportOutlined />}
								loading={this.state.isLoadingImportConfig}
							>
								Import Config
							</Button>
						</Upload>,
						<Button
							key="Export Config"
							type="text"
							icon={<ExportOutlined />}
							onClick={this.onClickExportConfig}
							loading={this.state.isLoadingExportConfig}
						>
							Export Config
						</Button>
					]}
				/>
				<Content style={{ padding: '16px 24px' }}>
					<Row gutter={[16]}>
						<Col span={4}>
							<Steps direction="vertical" current={this.state.step} onChange={this.onChangeStep}>
								<Step title="Target" icon={<AimOutlined />} description="Select Target" />
								<Step title="Attack" icon={<ToolOutlined />} description="Configure Attack" />
								<Step title="Advanced" icon={<ExperimentOutlined />} description="Advanced Options" />
								<Step title="Output" icon={<ExportOutlined />} description="Set Output" />
								<Step title="Finalize" icon={<FileDoneOutlined />} description="Review and Finalize" />
							</Steps>
						</Col>
						<Col span={20}>
							<div className="steps-content">
								{this.state.step === 0 ? (
									<Form layout="vertical">
										<Form.Item
											label="Hash"
										>
											<Select
												showSearch
												style={{ width: "100%" }}
												size="large"
												placeholder="Select Hash"
												value={this.state.hash}
												onChange={this.onChangeHash}
												filterOption={(input, option) =>
													String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
													String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
												}
											>
												{this.state._hashes.map(hash =>
													<Option value={hash} key={hash} title={hash}>{filename(hash)}</Option>
												)}
											</Select>
										</Form.Item>
										<Form.Item
											label="Algorithm"
										>
											<Select
												showSearch
												style={{ width: "100%" }}
												size="large"
												placeholder="Select Algorithm"
												value={this.state.algorithm}
												onChange={this.onChangeAlgorithm}
												filterOption={(input, option) =>
													String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
													String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
												}
											>
												{Object.keys(this.state._algorithms).map(key =>
													<Option value={Number(key)} key={Number(key)}>{key + " - " + this.state._algorithms[key]}</Option>
												)}
											</Select>
										</Form.Item>
									</Form>
								) : this.state.step === 1 ? (
									<Form layout="vertical" requiredMark="optional">
										<Form.Item
											label="Attack Mode"
											required
										>
											<Radio.Group value={this.state.attackMode} onChange={this.onChangeAttackMode}>
												<Radio value={0}>Dictionary Attack</Radio>
												<Radio value={1}>Combinator Attack</Radio>
												<Radio value={3}>Mask Attack</Radio>
												<Radio value={6}>Hybrid1 (Dictionary + Mask)</Radio>
												<Radio value={7}>Hybrid2 (Mask + Dictionary)</Radio>
											</Radio.Group>
										</Form.Item>
											{this.state.attackMode === 0 ? (
												<>
													<Form.Item
														label="Dictionaries"
														required
													>
														<Select
															mode="multiple"
															allowClear
															style={{ width: '100%' }}
															placeholder="Select Dictionaries"
															size="large"
															onChange={this.onChangeDictionaries}
															value={this.state.dictionaries}
															filterOption={(input, option) =>
																String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
																String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
															}
														>
															{this.state._dictionaries.map(dictionary =>
																<Option value={dictionary} key={dictionary} title={dictionary}>{filename(dictionary)}</Option>
															)}
														</Select>
													</Form.Item>
													<Form.Item
														label="Rules"
													>
														<Select
															mode="multiple"
															allowClear
															style={{ width: '100%' }}
															placeholder={"Select Rules [max. "+maxRules+"]"}
															size="large"
															onChange={this.onChangeRules}
															value={this.state.rules}
															filterOption={(input, option) =>
																String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
																String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
															}
														>
															{this.state._rules.map(rule =>
																<Option value={rule} key={rule} title={rule}>{filename(rule)}</Option>
															)}
														</Select>
													</Form.Item>
												</>
											) : this.state.attackMode === 1 ? (
												<Row gutter={[18, 16]}>
													<Col span={12}>
														<Row>
															<Col span={24}>
																<Form.Item
																	label="Left Dictionary"
																	required
																>
																	<Select
																		showSearch
																		allowClear
																		style={{ width: '100%' }}
																		placeholder="Select Left Dictionary"
																		size="large"
																		onChange={this.onChangeLeftDictionary}
																		value={this.state.leftDictionary}
																		filterOption={(input, option) =>
																			String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
																			String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
																		}
																	>
																		{this.state._dictionaries.map(dictionary =>
																			<Option value={dictionary} key={dictionary} title={dictionary}>{filename(dictionary)}</Option>
																		)}
																	</Select>
																</Form.Item>
															</Col>
															<Col span={24}>
																<Form.Item
																	label="Left Rule"
																>
																	<Input
																		allowClear
																		style={{ width: '100%' }}
																		placeholder="Set Left Rule"
																		size="large"
																		onChange={this.onChangeLeftRule}
																		value={this.state.leftRule}
																	/>
																</Form.Item>
															</Col>
														</Row>
													</Col>
													<Col span={12}>
														<Row>
															<Col span={24}>
																<Form.Item
																	label="Right Dictionary"
																	required
																>
																	<Select
																		showSearch
																		allowClear
																		style={{ width: '100%' }}
																		placeholder="Select Right Dictionary"
																		size="large"
																		onChange={this.onChangeRightDictionary}
																		value={this.state.rightDictionary}
																		filterOption={(input, option) =>
																			String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
																			String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
																		}
																	>
																		{this.state._dictionaries.map(dictionary =>
																			<Option value={dictionary} key={dictionary} title={dictionary}>{filename(dictionary)}</Option>
																		)}
																	</Select>
																</Form.Item>
															</Col>
															<Col span={24}>
																<Form.Item
																	label="Right Rule"
																>
																	<Input
																		allowClear
																		style={{ width: '100%' }}
																		placeholder="Set Right Rule"
																		size="large"
																		onChange={this.onChangeRightRule}
																		value={this.state.rightRule}
																	/>
																</Form.Item>
															</Col>
														</Row>
													</Col>
												</Row>
											) : this.state.attackMode === 3 ? (
												<Row gutter={[18, 16]}>
													<Col span={12}>
														{this.state.maskInputType === "text" ? (
															<Form.Item
																label="Mask"
																required
															>
																<Input
																	allowClear
																	style={{ width: '100%' }}
																	placeholder="Set Mask"
																	size="large"
																	onChange={this.onChangeMask}
																	value={this.state.mask}
																	suffix={
																		this.state.mask ? maskLength(this.state.mask) : undefined
																	}
																/>
																<Button
																	type="link"
																	style={{ padding: '0' }}
																	onClick={() => this.onChangeMaskInputType({type: "file"})}
																>
																	Use .hcmask file instead
																</Button>
															</Form.Item>
														) : this.state.maskInputType === "file" ? (
															<Form.Item
																label="Mask"
																required
															>
																<Select
																	allowClear
																	style={{ width: '100%' }}
																	placeholder="Select Mask"
																	size="large"
																	onChange={this.onChangeMaskFile}
																	value={this.state.onChangeMaskFile}
																	filterOption={(input, option) =>
																		String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
																		String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
																	}
																>
																	{this.state._masks.map(mask =>
																		<Option value={mask} key={mask} title={mask}>{filename(mask)}</Option>
																	)}
																</Select>
																<Button
																	type="link"
																	style={{ padding: '0' }}
																	onClick={() => this.onChangeMaskInputType({type: "text"})}
																>
																	Use mask text instead
																</Button>
															</Form.Item>
														) : "unsupported mask input type" }
														<Form.Item
															label="Mask increment mode"
														>
															<Checkbox
																checked={this.state.enableMaskIncrementMode}
																onChange={this.onChangeEnableMaskIncrementMode}
															>
																Enable
															</Checkbox>
															<InputNumber
																disabled={!this.state.enableMaskIncrementMode}
																min={maskIncrementMin}
																max={maskIncrementMax}
																value={this.state.maskIncrementMin}
																onChange={this.onChangeMaskIncrementMin}
															/>
															<InputNumber
																disabled={!this.state.enableMaskIncrementMode}
																min={maskIncrementMin}
																max={maskIncrementMax}
																value={this.state.maskIncrementMax}
																onChange={this.onChangeMaskIncrementMax}
															/>
														</Form.Item>
													</Col>
													<Col span={12}>
														<Form.Item
															label="Custom charset 1"
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder="Set Custom charset 1"
																size="large"
																onChange={this.onChangeCustomCharset1}
																value={this.state.customCharset1}
															/>
														</Form.Item>
														<Form.Item
															label="Custom charset 2"
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder="Set Custom charset 2"
																size="large"
																onChange={this.onChangeCustomCharset2}
																value={this.state.customCharset2}
															/>
														</Form.Item>
														<Form.Item
															label="Custom charset 3"
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder="Set Custom charset 3"
																size="large"
																onChange={this.onChangeCustomCharset3}
																value={this.state.customCharset3}
															/>
														</Form.Item>
														<Form.Item
															label="Custom charset 4"
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder="Set Custom charset 4"
																size="large"
																onChange={this.onChangeCustomCharset4}
																value={this.state.customCharset4}
															/>
														</Form.Item>
													</Col>
												</Row>
											) : this.state.attackMode === 6 ? (
												<Row gutter={[18, 16]}>
													<Col span={24}>
														<Row gutter={[18, 16]}>
															<Col span={12}>
																<Form.Item
																	label="Dictionary"
																	required
																>
																	<Select
																		showSearch
																		allowClear
																		style={{ width: '100%' }}
																		placeholder="Select Dictionary"
																		size="large"
																		onChange={this.onChangeLeftDictionary}
																		value={this.state.leftDictionary}
																		filterOption={(input, option) =>
																			String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
																			String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
																		}
																	>
																		{this.state._dictionaries.map(dictionary =>
																			<Option value={dictionary} key={dictionary} title={dictionary}>{filename(dictionary)}</Option>
																		)}
																	</Select>
																</Form.Item>
															</Col>
															<Col span={12}>
																<Form.Item
																	label="Rule"
																>
																	<Input
																		allowClear
																		style={{ width: '100%' }}
																		placeholder="Set Rule"
																		size="large"
																		onChange={this.onChangeLeftRule}
																		value={this.state.leftRule}
																	/>
																</Form.Item>
															</Col>
														</Row>
													</Col>
													<Col span={12}>
														{this.state.maskInputType === "text" ? (
															<Form.Item
																label="Mask"
																required
															>
																<Input
																	allowClear
																	style={{ width: '100%' }}
																	placeholder="Set Mask"
																	size="large"
																	onChange={this.onChangeMask}
																	value={this.state.mask}
																	suffix={
																		this.state.mask ? maskLength(this.state.mask) : undefined
																	}
																/>
																<Button
																	type="link"
																	style={{ padding: '0' }}
																	onClick={() => this.onChangeMaskInputType({type: "file"})}
																>
																	Use .hcmask file instead
																</Button>
															</Form.Item>
														) : this.state.maskInputType === "file" ? (
															<Form.Item
																label="Mask"
																required
															>
																<Select
																	allowClear
																	style={{ width: '100%' }}
																	placeholder="Select Mask"
																	size="large"
																	onChange={this.onChangeMaskFile}
																	value={this.state.onChangeMaskFile}
																	filterOption={(input, option) =>
																		String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
																		String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
																	}
																>
																	{this.state._masks.map(mask =>
																		<Option value={mask} key={mask} title={mask}>{filename(mask)}</Option>
																	)}
																</Select>
																<Button
																	type="link"
																	style={{ padding: '0' }}
																	onClick={() => this.onChangeMaskInputType({type: "text"})}
																>
																	Use mask text instead
																</Button>
															</Form.Item>
														) : "unsupported mask input type" }
														<Form.Item
															label="Mask increment mode"
														>
															<Checkbox
																checked={this.state.enableMaskIncrementMode}
																onChange={this.onChangeEnableMaskIncrementMode}
															>
																Enable
															</Checkbox>
															<InputNumber
																disabled={!this.state.enableMaskIncrementMode}
																min={maskIncrementMin}
																max={maskIncrementMax}
																value={this.state.maskIncrementMin}
																onChange={this.onChangeMaskIncrementMin}
															/>
															<InputNumber
																disabled={!this.state.enableMaskIncrementMode}
																min={maskIncrementMin}
																max={maskIncrementMax}
																value={this.state.maskIncrementMax}
																onChange={this.onChangeMaskIncrementMax}
															/>
														</Form.Item>
													</Col>
													<Col span={12}>
														<Form.Item
															label="Custom charset 1"
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder="Set Custom charset 1"
																size="large"
																onChange={this.onChangeCustomCharset1}
																value={this.state.customCharset1}
															/>
														</Form.Item>
														<Form.Item
															label="Custom charset 2"
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder="Set Custom charset 2"
																size="large"
																onChange={this.onChangeCustomCharset2}
																value={this.state.customCharset2}
															/>
														</Form.Item>
														<Form.Item
															label="Custom charset 3"
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder="Set Custom charset 3"
																size="large"
																onChange={this.onChangeCustomCharset3}
																value={this.state.customCharset3}
															/>
														</Form.Item>
														<Form.Item
															label="Custom charset 4"
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder="Set Custom charset 4"
																size="large"
																onChange={this.onChangeCustomCharset4}
																value={this.state.customCharset4}
															/>
														</Form.Item>
													</Col>
												</Row>
											) : this.state.attackMode === 7 ? (
												<Row gutter={[18, 16]}>
													<Col span={12}>
														{this.state.maskInputType === "text" ? (
															<Form.Item
																label="Mask"
																required
															>
																<Input
																	allowClear
																	style={{ width: '100%' }}
																	placeholder="Set Mask"
																	size="large"
																	onChange={this.onChangeMask}
																	value={this.state.mask}
																	suffix={
																		this.state.mask ? maskLength(this.state.mask) : undefined
																	}
																/>
																<Button
																	type="link"
																	style={{ padding: '0' }}
																	onClick={() => this.onChangeMaskInputType({type: "file"})}
																>
																	Use .hcmask file instead
																</Button>
															</Form.Item>
														) : this.state.maskInputType === "file" ? (
															<Form.Item
																label="Mask"
																required
															>
																<Select
																	allowClear
																	style={{ width: '100%' }}
																	placeholder="Select Mask"
																	size="large"
																	onChange={this.onChangeMaskFile}
																	value={this.state.onChangeMaskFile}
																	filterOption={(input, option) =>
																		String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
																		String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
																	}
																>
																	{this.state._masks.map(mask =>
																		<Option value={mask} key={mask} title={mask}>{filename(mask)}</Option>
																	)}
																</Select>
																<Button
																	type="link"
																	style={{ padding: '0' }}
																	onClick={() => this.onChangeMaskInputType({type: "text"})}
																>
																	Use mask text instead
																</Button>
															</Form.Item>
														) : "unsupported mask input type" }
														<Form.Item
															label="Mask increment mode"
														>
															<Checkbox
																checked={this.state.enableMaskIncrementMode}
																onChange={this.onChangeEnableMaskIncrementMode}
															>
																Enable
															</Checkbox>
															<InputNumber
																disabled={!this.state.enableMaskIncrementMode}
																min={maskIncrementMin}
																max={maskIncrementMax}
																value={this.state.maskIncrementMin}
																onChange={this.onChangeMaskIncrementMin}
															/>
															<InputNumber
																disabled={!this.state.enableMaskIncrementMode}
																min={maskIncrementMin}
																max={maskIncrementMax}
																value={this.state.maskIncrementMax}
																onChange={this.onChangeMaskIncrementMax}
															/>
														</Form.Item>
													</Col>
													<Col span={12}>
														<Form.Item
															label="Custom charset 1"
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder="Set Custom charset 1"
																size="large"
																onChange={this.onChangeCustomCharset1}
																value={this.state.customCharset1}
															/>
														</Form.Item>
														<Form.Item
															label="Custom charset 2"
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder="Set Custom charset 2"
																size="large"
																onChange={this.onChangeCustomCharset2}
																value={this.state.customCharset2}
															/>
														</Form.Item>
														<Form.Item
															label="Custom charset 3"
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder="Set Custom charset 3"
																size="large"
																onChange={this.onChangeCustomCharset3}
																value={this.state.customCharset3}
															/>
														</Form.Item>
														<Form.Item
															label="Custom charset 4"
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder="Set Custom charset 4"
																size="large"
																onChange={this.onChangeCustomCharset4}
																value={this.state.customCharset4}
															/>
														</Form.Item>
													</Col>
													<Col span={24}>
														<Row gutter={[18, 16]}>
															<Col span={12}>
																<Form.Item
																	label="Dictionary"
																	required
																>
																	<Select
																		showSearch
																		allowClear
																		style={{ width: '100%' }}
																		placeholder="Select Dictionary"
																		size="large"
																		onChange={this.onChangeRightDictionary}
																		value={this.state.rightDictionary}
																		filterOption={(input, option) =>
																			String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
																			String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
																		}
																	>
																		{this.state._dictionaries.map(dictionary =>
																			<Option value={dictionary} key={dictionary} title={dictionary}>{filename(dictionary)}</Option>
																		)}
																	</Select>
																</Form.Item>
															</Col>
															<Col span={12}>
																<Form.Item
																	label="Rule"
																>
																	<Input
																		allowClear
																		style={{ width: '100%' }}
																		placeholder="Set Rule"
																		size="large"
																		onChange={this.onChangeRightRule}
																		value={this.state.rightRule}
																	/>
																</Form.Item>
															</Col>
														</Row>
													</Col>
												</Row>
											) : (
												"Select Attack Mode"
											)}
									</Form>
								) : this.state.step === 2 ? (
									<Form layout="vertical">
										<Collapse ghost onChange={this.onChangeAdvancedOptionsCollapse} activeKey={this.state.advancedOptionsCollapse}>
											<Panel header="General" key="General">
												<Row gutter={[18, 16]}>
													<Col>
														<Checkbox
															checked={this.state.enableOptimizedKernel}
															onChange={this.onChangeEnableOptimizedKernel}
														>
															Enable optimized kernel
														</Checkbox>
													</Col>
													<Col>
														<Checkbox
															checked={this.state.enableSlowerCandidateGenerators}
															onChange={this.onChangeEnableSlowerCandidateGenerators}
														>
															Enable slower candidate generators
														</Checkbox>
													</Col>
													<Col>
														<Checkbox
															checked={this.state.removeFoundHashes}
															onChange={this.onChangeRemoveFoundHashes}
														>
															Remove found hashes
														</Checkbox>
													</Col>
													<Col>
														<Checkbox
															checked={this.state.disablePotFile}
															onChange={this.onChangeDisablePotFile}
														>
															Disable Pot File
														</Checkbox>
													</Col>
													<Col>
														<Checkbox
															checked={this.state.ignoreUsernames}
															onChange={this.onChangeIgnoreUsernames}
														>
															Ignore Usernames
														</Checkbox>
													</Col>
													<Col>
														<Checkbox
															checked={this.state.disableSelfTest}
															onChange={this.onChangeDisableSelfTest}
														>
															Disable self-test (Not Recommended)
														</Checkbox>
													</Col>
													<Col>
														<Checkbox
															checked={this.state.ignoreWarnings}
															onChange={this.onChangeIgnoreWarnings}
														>
															Ignore warnings (Not Recommended)
														</Checkbox>
													</Col>
												</Row>
											</Panel>
											<Panel header="Devices" key="Devices">
												<Row gutter={[18, 16]}>
													<Col span={8}>
														<Form.Item
															label="Devices IDs"
														>
															<Select
																mode="multiple"
																allowClear
																style={{ width: '100%' }}
																placeholder="Select Devices IDs"
																size="large"
																onChange={this.onChangeDevicesIDs}
																value={this.state.devicesIDs}
																filterOption={(input, option) =>
																	String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
																	String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
																}
															>
																{Array.from(Array(16)).map((x,i) =>
																	<Option value={i+1} key={i+1}>{"Device #"+(i+1)}</Option>
																)}
															</Select>
														</Form.Item>
													</Col>
													<Col span={8}>
														<Form.Item
															label="Devices Types"
														>
															<Select
																mode="multiple"
																allowClear
																style={{ width: '100%' }}
																placeholder="Select Devices Types"
																size="large"
																onChange={this.onChangeDevicesTypes}
																value={this.state.devicesTypes}
																filterOption={(input, option) =>
																	String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
																	String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
																}
															>
																<Option value={1} key={1}>CPU</Option>
																<Option value={2} key={2}>GPU</Option>
																<Option value={3} key={3}>FPGA, DSP, Co-Processor</Option>
															</Select>
														</Form.Item>
													</Col>
													<Col span={8}>
														<Form.Item
															label="Workload Profile"
															tooltip={
																<Table
																	columns={[
																		{
																			title: 'Performance',
																			dataIndex: 'performance',
																			key: 'Performance'
																		},
																		{
																			title: 'Runtime',
																			dataIndex: 'runtime',
																			key: 'Runtime'
																		},
																		{
																			title: 'Power Consumption',
																			dataIndex: 'powerConsumption',
																			key: 'Power Consumption'
																		},
																		{
																			title: 'Desktop Impact',
																			dataIndex: 'desktopImpact',
																			key: 'Desktop Impact'
																		}
																	]}
																	dataSource={[
																		{
																			key: '1',
																			performance: 'Low',
																			runtime: '2 ms',
																			powerConsumption: 'Low',
																			desktopImpact: 'Minimal'
																		},
																		{
																			key: '2',
																			performance: 'Default',
																			runtime: '12 ms',
																			powerConsumption: 'Economic',
																			desktopImpact: 'Noticeable'
																		},
																		{
																			key: '3',
																			performance: 'High',
																			runtime: '96 ms',
																			powerConsumption: 'High',
																			desktopImpact: 'Unresponsive'
																		},
																		{
																			key: '4',
																			performance: 'Nightmare',
																			runtime: '480 ms',
																			powerConsumption: 'Insane',
																			desktopImpact: 'Headless'
																		}
																	]}
																	size="small"
																	pagination={false}
																	style={{ overflow: 'auto' }}
																/>
															}
														>
															<Select
																allowClear
																style={{ width: '100%' }}
																placeholder="Select Workload Profile"
																size="large"
																onChange={this.onChangeWorkloadProfile}
																value={this.state.workloadProfile}
																filterOption={(input, option) =>
																	String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
																	String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
																}
															>
																<Option value={1} key={1}>Low</Option>
																<Option value={2} key={2}>Default</Option>
																<Option value={3} key={3}>High</Option>
																<Option value={4} key={4}>Nightmare</Option>
															</Select>
														</Form.Item>
													</Col>
												</Row>
											</Panel>
											<Panel header="Markov" key="Markov">
												<Row gutter={[18, 16]}>
													<Col>
														<Checkbox
															checked={this.state.markovDisable}
															onChange={this.onChangeMarkovDisable}
														>
															Disables markov-chains, emulates classic brute-force
														</Checkbox>
													</Col>
													<Col>
														<Checkbox
															checked={this.state.markovClassic}
															onChange={this.onChangeMarkovClassic}
														>
															Enables classic markov-chains, no per-position
														</Checkbox>
													</Col>
													<Col span={24}>
														<Row gutter={[18, 16]}>
															<Col span={8}>
																<Form.Item
																	label="Threshold X when to stop accepting new markov-chains"
																>
																	<InputNumber
																		value={this.state.markovThreshold}
																		onChange={this.onChangeMarkovThreshold}
																	/>
																</Form.Item>
															</Col>
														</Row>
													</Col>
												</Row>
											</Panel>
											<Panel header="Monitor" key="Monitor">
												<Row gutter={[18, 16]}>
													<Col>
														<Checkbox
															checked={this.state.disableMonitor}
															onChange={this.onChangeDisableMonitor}
														>
															Disable Monitor
														</Checkbox>
													</Col>
													<Col span={24}>
														<Row gutter={[18, 16]}>
															<Col span={8}>
																<Form.Item
																	label="Temp Abort (°C)"
																>
																	<Select
																		allowClear
																		style={{ width: '100%' }}
																		placeholder="Select Temp Abort (°C)"
																		size="large"
																		onChange={this.onChangeTempAbort}
																		value={this.state.tempAbort}
																		disabled={this.state.disableMonitor}
																		filterOption={(input, option) =>
																			String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
																			String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
																		}
																	>
																		<Option value={60} key={60}>60 °C</Option>
																		<Option value={65} key={65}>65 °C</Option>
																		<Option value={70} key={70}>70 °C</Option>
																		<Option value={75} key={75}>75 °C</Option>
																		<Option value={80} key={80}>80 °C</Option>
																		<Option value={85} key={85}>85 °C</Option>
																		<Option value={90} key={90}>90 °C</Option>
																		<Option value={95} key={95}>95 °C</Option>
																		<Option value={100} key={100}>100 °C</Option>
																	</Select>
																</Form.Item>
															</Col>
														</Row>
													</Col>
												</Row>
											</Panel>
											<Panel header="Extra Arguments" key="Extra Arguments">
												<Form.Item
													label="Extra Arguments"
												>
													<Input
														allowClear
														style={{ width: '100%' }}
														placeholder="Set Extra Arguments"
														size="large"
														onChange={this.onChangeExtraArguments}
														value={this.state.extraArguments.join(" ")}
													/>
												</Form.Item>
											</Panel>
											<Panel header="Misc" key="Misc">
												<Form.Item
													label="Status Timer"
												>
													<Select
														allowClear
														style={{ width: '100%' }}
														placeholder="Select Status Timer"
														size="large"
														onChange={this.onChangeStatusTimer}
														value={this.state.statusTimer}
														filterOption={(input, option) =>
															String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
															String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
														}
													>
														<Option value={10} key={10}>10 Seconds</Option>
														<Option value={20} key={20}>20 Seconds</Option>
														<Option value={30} key={30}>30 Seconds</Option>
														<Option value={45} key={45}>45 Seconds</Option>
														<Option value={60} key={60}>60 Seconds</Option>
														<Option value={90} key={90}>90 Seconds</Option>
														<Option value={120} key={120}>120 Seconds</Option>
														<Option value={300} key={300}>300 Seconds</Option>
													</Select>
												</Form.Item>
											</Panel>
										</Collapse>
									</Form>
								) : this.state.step === 3 ? (
									<Form layout="vertical">
										<Form.Item
											label="Output File"
											extra={this.state.outputFile ? this.state.outputFile : null}
										>
											<Button
												type="primary"
												onClick={this.onChangeOutputFile}
												loading={this.state.isLoadingSetOutputFile}
											>
												Set Output File
											</Button>
										</Form.Item>
										<Form.Item
											label="Output Format"
										>
											<Select
												mode="multiple"
												allowClear
												style={{ width: '100%' }}
												placeholder="Select Output Format"
												size="large"
												onChange={this.onChangeOutputFormat}
												value={this.state.outputFormat}
												filterOption={(input, option) =>
													String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
													String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
												}
											>
												<Option value={1} key={1}>hash[:salt]</Option>
												<Option value={2} key={2}>plain</Option>
												<Option value={3} key={3}>hex_plain</Option>
												<Option value={4} key={4}>crack_pos</Option>
												<Option value={5} key={5}>timestamp absolute</Option>
												<Option value={6} key={6}>timestamp relative</Option>
											</Select>
										</Form.Item>
									</Form>
								) : this.state.step === 4 ? (
									<Space size="large">
										<Button
											type="primary"
											icon={<PlusOutlined />}
											onClick={this.onClickCreateTask}
											loading={this.state.isLoadingCreateTask}
										>
											Create Task
										</Button>
										<Checkbox
											checked={this.state.preserveTaskConfig}
											onChange={this.onChangePreserveTaskConfig}
										>
											Preserve task config
										</Checkbox>
									</Space>
								) : null }
							</div>
						</Col>
					</Row>
				</Content>
			</>
		)
	}
}

export default NewTask;
