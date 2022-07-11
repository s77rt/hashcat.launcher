import { withTranslation } from 'react-i18next';

import React, { Component } from "react";
import { Modal, Layout, PageHeader, Descriptions, Table, InputNumber, message, Row, Col, Card, Select, Typography, Upload, Button, Space, Input, Form, Radio, Divider, Collapse, Checkbox, Tabs, Steps } from 'antd';
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

import filename from './lib/filename';

const { Content } = Layout;
const { Option } = Select;
const { Title } = Typography;
const { Panel } = Collapse;
const { TabPane } = Tabs;
const { Step } = Steps;

const maxRules = 4;
const maskIncrementMin = 1;
const maskIncrementMax = 16;

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

	quiet: true,
	disablePotFile: true,
	disableLogFile: true,
	enableOptimizedKernel: true,
	enableSlowerCandidateGenerators: false,
	removeFoundHashes: false,
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

	extraArguments: [],

	statusTimer: 20,

	outputFile: undefined,
	outputFormat: [1, 2],

	priority: 0
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

		this.onChangeQuiet = this.onChangeQuiet.bind(this);
		this.onChangeDisablePotFile = this.onChangeDisablePotFile.bind(this);
		this.onChangeDisableLogFile = this.onChangeDisableLogFile.bind(this);
		this.onChangeEnableOptimizedKernel = this.onChangeEnableOptimizedKernel.bind(this);
		this.onChangeEnableSlowerCandidateGenerators = this.onChangeEnableSlowerCandidateGenerators.bind(this);
		this.onChangeRemoveFoundHashes = this.onChangeRemoveFoundHashes.bind(this);
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

		this.onChangePriority = this.onChangePriority.bind(this);

		this.onClickCreateTask = this.onClickCreateTask.bind(this);

		this.onChangePreserveTaskConfig = this.onChangePreserveTaskConfig.bind(this);

		this.onClickImportConfig = this.onClickImportConfig.bind(this);
		this.onClickExportConfig = this.onClickExportConfig.bind(this);

		this.onClickDevicesInfo = this.onClickDevicesInfo.bind(this);
		this.onClickBenchmark = this.onClickBenchmark.bind(this);

		this.state = {
			...initialConfig,

			preserveTaskConfig: false,

			isLoadingImportConfig: false,
			isLoadingExportConfig: false,
			isLoadingSetOutputFile: false,
			isLoadingCreateTask: false,

			isLoadingDevicesInfo: false,
			isLoadingBenchmark: false,

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

	onChangeQuiet(e) {
		this.setState({
			quiet: e.target.checked
		});
	}

	onChangeDisablePotFile(e) {
		this.setState({
			disablePotFile: e.target.checked
		});
	}

	onChangeDisableLogFile(e) {
		this.setState({
			disableLogFile: e.target.checked
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

	onChangePriority(priority) {
		if (typeof(priority) !== "number")
			return

		this.setState({
			priority: priority
		});
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

				quiet: this.state.quiet,
				disablePotFile: this.state.disablePotFile,
				disableLogFile: this.state.disableLogFile,
				enableOptimizedKernel: this.state.enableOptimizedKernel,
				enableSlowerCandidateGenerators: this.state.enableSlowerCandidateGenerators,
				removeFoundHashes: this.state.removeFoundHashes,
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
			}, this.state.priority).then(
				response => {
					message.success(this.props.t('newtask.task_success'));
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

	onChangePreserveTaskConfig(e) {
		this.setState({
			preserveTaskConfig: e.target.checked
		});
	}

	resetInitialConfig() {
		this.setState(initialConfig);
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
					message.success(this.props.t('newtask.import_success') + " (" + filename(file.name) + ")");
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
		if (typeof window.GOexportConfig !== "function") {
			message.error("GOexportConfig is not a function");
			return;
		}

		this.setState({
			isLoadingExportConfig: true
		}, () => {
			window.GOexportConfig(this.exportConfig()).then(
				response => {
					message.success(this.props.t('newtask.export_success') + " (" + filename(response) + ")");
					this.setState({isLoadingExportConfig: false});
				},
				error => {
					message.error(error);
					this.setState({isLoadingExportConfig: false});
				}
			);
		});
	}

	onClickDevicesInfo() {
		if (typeof window.GOhashcatDevices !== "function") {
			message.error("GOhashcatDevices is not a function");
			return;
		}

		this.setState({
			isLoadingDevicesInfo: true
		}, () => {
			window.GOhashcatDevices().then(
				response => {
					Modal.info({
						title: this.props.t('newtask.devices'),
						okText: this.props.t('newtask.devices_modal_ok'),
						content: (
							<pre style={{
								maxHeight: '300px',
								overflow: 'auto',
								padding: '.5rem',
								margin: '0',
								border: '1px solid #303030'
							}}>
								{response}
							</pre>
						),
						width: 720
					});
					this.setState({isLoadingDevicesInfo: false});
				},
				error => {
					message.error(error);
					this.setState({isLoadingDevicesInfo: false});
				}
			);
		});
	}

	onClickBenchmark() {
		if (typeof window.GOhashcatBenchmark !== "function") {
			message.error("GOhashcatBenchmark is not a function");
			return;
		}

		const algorithm = this.state.algorithm;
		if (typeof(algorithm) !== "number") {
			message.error(this.props.t('newtask.no_algorithm_error'));
			return;
		}

		this.setState({
			isLoadingBenchmark: true
		}, () => {
			window.GOhashcatBenchmark(algorithm).then(
				response => {
					Modal.info({
						title: this.props.t('newtask.benchmark'),
						okText: this.props.t('newtask.benchmark_modal_ok'),
						content: (
							<pre style={{
								maxHeight: '300px',
								overflow: 'auto',
								padding: '.5rem',
								margin: '0',
								border: '1px solid #303030'
							}}>
								{response}
							</pre>
						),
						width: 720
					});
					this.setState({isLoadingBenchmark: false});
				},
				error => {
					message.error(error);
					this.setState({isLoadingBenchmark: false});
				}
			);
		});
	}

	componentDidMount() {
		EventBus.on("dataUpdate", "NewTask", () => {
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
		EventBus.remove("dataUpdate", "NewTask");
	}

	render() {
		const LANG = this.props.t;
		return (
			<>
				<PageHeader
					title={LANG('newtask.title')}
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
								{LANG('newtask.import_config')}
							</Button>
						</Upload>,
						<Button
							key="Export Config"
							type="text"
							icon={<ExportOutlined />}
							onClick={this.onClickExportConfig}
							loading={this.state.isLoadingExportConfig}
						>
							{LANG('newtask.export_config')}
						</Button>
					]}
				/>
				<Content style={{ padding: '16px 24px' }}>
					<Row gutter={[16]}>
						<Col span={4}>
							<Steps direction="vertical" current={this.state.step} onChange={this.onChangeStep}>
								<Step title={LANG('newtask.target')} icon={<AimOutlined />} description={LANG('newtask.select_target')} />
								<Step title={LANG('newtask.attack')} icon={<ToolOutlined />} description={LANG('newtask.configure_attack')} />
								<Step title={LANG('newtask.advanced')} icon={<ExperimentOutlined />} description={LANG('newtask.advanced_options')} />
								<Step title={LANG('newtask.output')} icon={<ExportOutlined />} description={LANG('newtask.set_output')} />
								<Step title={LANG('newtask.finalize')} icon={<FileDoneOutlined />} description={LANG('newtask.review_and_finalize')} />
							</Steps>
						</Col>
						<Col span={20}>
							<div className="steps-content">
								{this.state.step === 0 ? (
									<Form layout="vertical">
										<Form.Item
											label={LANG('newtask.hash')}
										>
											<Select
												showSearch
												style={{ width: "100%" }}
												size="large"
												placeholder={LANG('newtask.select_hash')}
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
											label={LANG('newtask.algorithm')}
										>
											<Select
												showSearch
												style={{ width: "100%" }}
												size="large"
												placeholder={LANG('newtask.select_algorithm')}
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
											label={LANG('newtask.attack_mode')}
											required
										>
											<Radio.Group value={this.state.attackMode} onChange={this.onChangeAttackMode}>
												<Radio value={0}>{LANG('newtask.dictionary_attack')}</Radio>
												<Radio value={1}>{LANG('newtask.combinator_attack')}</Radio>
												<Radio value={3}>{LANG('newtask.mask_attack')}</Radio>
												<Radio value={6}>{LANG('newtask.hybrid1_attack')}</Radio>
												<Radio value={7}>{LANG('newtask.Hybrid2_attack')}</Radio>
											</Radio.Group>
										</Form.Item>
											{this.state.attackMode === 0 ? (
												<>
													<Form.Item
														label={LANG('newtask.dictionaries')}
														required
													>
														<Select
															mode="multiple"
															allowClear
															style={{ width: '100%' }}
															placeholder={LANG('newtask.select_dictionaries')}
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
														label={LANG('newtask.rules')}
													>
														<Select
															mode="multiple"
															allowClear
															style={{ width: '100%' }}
															placeholder={LANG('newtask.select_rules') + " [" + LANG('newtask.max_short') + " " + maxRules + "]"}
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
																	label={LANG('newtask.left_dictionary')}
																	required
																>
																	<Select
																		showSearch
																		allowClear
																		style={{ width: '100%' }}
																		placeholder={LANG('newtask.select_left_dictionary')}
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
																	label={LANG('newtask.left_rule')}
																>
																	<Input
																		allowClear
																		style={{ width: '100%' }}
																		placeholder={LANG('newtask.set_left_rule')}
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
																	label={LANG('newtask.right_dictionary')}
																	required
																>
																	<Select
																		showSearch
																		allowClear
																		style={{ width: '100%' }}
																		placeholder={LANG('newtask.select_right_dictionary')}
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
																	label={LANG('newtask.right_rule')}
																>
																	<Input
																		allowClear
																		style={{ width: '100%' }}
																		placeholder={LANG('newtask.set_right_rule')}
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
																label={LANG('newtask.mask')}
																required
															>
																<Input
																	allowClear
																	style={{ width: '100%' }}
																	placeholder={LANG('newtask.set_mask')}
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
																	{LANG('newtask.use_hcmask_file_instead')}
																</Button>
															</Form.Item>
														) : this.state.maskInputType === "file" ? (
															<Form.Item
																label={LANG('newtask.mask')}
																required
															>
																<Select
																	showSearch
																	allowClear
																	style={{ width: '100%' }}
																	placeholder={LANG('newtask.select_mask')}
																	size="large"
																	onChange={this.onChangeMaskFile}
																	value={this.state.maskFile}
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
																	{LANG('newtask.use_mask_text_instead')}
																</Button>
															</Form.Item>
														) : LANG('newtask.unsupported_mask_input_type') }
														<Form.Item
															label={LANG('newtask.mask_increment_mode')}
														>
															<Checkbox
																checked={this.state.enableMaskIncrementMode}
																onChange={this.onChangeEnableMaskIncrementMode}
															>
																{LANG('newtask.enable')}
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
															label={LANG('newtask.custom_charset_1')}
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder={LANG('newtask.set_custom_charset_1')}
																size="large"
																onChange={this.onChangeCustomCharset1}
																value={this.state.customCharset1}
																disabled={!(this.state.maskInputType === "text")}
															/>
														</Form.Item>
														<Form.Item
															label={LANG('newtask.custom_charset_2')}
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder={LANG('newtask.set_custom_charset_2')}
																size="large"
																onChange={this.onChangeCustomCharset2}
																value={this.state.customCharset2}
																disabled={!(this.state.maskInputType === "text")}
															/>
														</Form.Item>
														<Form.Item
															label={LANG('newtask.custom_charset_3')}
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder={LANG('newtask.set_custom_charset_3')}
																size="large"
																onChange={this.onChangeCustomCharset3}
																value={this.state.customCharset3}
																disabled={!(this.state.maskInputType === "text")}
															/>
														</Form.Item>
														<Form.Item
															label={LANG('newtask.custom_charset_4')}
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder={LANG('newtask.set_custom_charset_4')}
																size="large"
																onChange={this.onChangeCustomCharset4}
																value={this.state.customCharset4}
																disabled={!(this.state.maskInputType === "text")}
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
																	label={LANG('newtask.dictionary')}
																	required
																>
																	<Select
																		showSearch
																		allowClear
																		style={{ width: '100%' }}
																		placeholder={LANG('newtask.select_dictionary')}
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
																	label={LANG('newtask.rule')}
																>
																	<Input
																		allowClear
																		style={{ width: '100%' }}
																		placeholder={LANG('newtask.set_rule')}
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
																label={LANG('newtask.mask')}
																required
															>
																<Input
																	allowClear
																	style={{ width: '100%' }}
																	placeholder={LANG('newtask.set_mask')}
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
																	{LANG('newtask.use_hcmask_file_instead')}
																</Button>
															</Form.Item>
														) : this.state.maskInputType === "file" ? (
															<Form.Item
																label={LANG('newtask.mask')}
																required
															>
																<Select
																	showSearch
																	allowClear
																	style={{ width: '100%' }}
																	placeholder={LANG('newtask.select_mask')}
																	size="large"
																	onChange={this.onChangeMaskFile}
																	value={this.state.maskFile}
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
																	{LANG('newtask.use_mask_text_instead')}
																</Button>
															</Form.Item>
														) : LANG('newtask.unsupported_mask_input_type') }
														<Form.Item
															label={LANG('newtask.mask_increment_mode')}
														>
															<Checkbox
																checked={this.state.enableMaskIncrementMode}
																onChange={this.onChangeEnableMaskIncrementMode}
															>
																{LANG('newtask.enable')}
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
															label={LANG('newtask.custom_charset_1')}
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder={LANG('newtask.set_custom_charset_1')}
																size="large"
																onChange={this.onChangeCustomCharset1}
																value={this.state.customCharset1}
																disabled={!(this.state.maskInputType === "text")}
															/>
														</Form.Item>
														<Form.Item
															label={LANG('newtask.custom_charset_2')}
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder={LANG('newtask.set_custom_charset_2')}
																size="large"
																onChange={this.onChangeCustomCharset2}
																value={this.state.customCharset2}
																disabled={!(this.state.maskInputType === "text")}
															/>
														</Form.Item>
														<Form.Item
															label={LANG('newtask.custom_charset_3')}
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder={LANG('newtask.set_custom_charset_3')}
																size="large"
																onChange={this.onChangeCustomCharset3}
																value={this.state.customCharset3}
																disabled={!(this.state.maskInputType === "text")}
															/>
														</Form.Item>
														<Form.Item
															label={LANG('newtask.custom_charset_4')}
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder={LANG('newtask.set_custom_charset_4')}
																size="large"
																onChange={this.onChangeCustomCharset4}
																value={this.state.customCharset4}
																disabled={!(this.state.maskInputType === "text")}
															/>
														</Form.Item>
													</Col>
												</Row>
											) : this.state.attackMode === 7 ? (
												<Row gutter={[18, 16]}>
													<Col span={12}>
														{this.state.maskInputType === "text" ? (
															<Form.Item
																label={LANG('newtask.mask')}
																required
															>
																<Input
																	allowClear
																	style={{ width: '100%' }}
																	placeholder={LANG('newtask.set_mask')}
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
																	{LANG('newtask.use_hcmask_file_instead')}
																</Button>
															</Form.Item>
														) : this.state.maskInputType === "file" ? (
															<Form.Item
																label={LANG('newtask.mask')}
																required
															>
																<Select
																	showSearch
																	allowClear
																	style={{ width: '100%' }}
																	placeholder={LANG('newtask.select_mask')}
																	size="large"
																	onChange={this.onChangeMaskFile}
																	value={this.state.maskFile}
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
																	{LANG('newtask.use_mask_text_instead')}
																</Button>
															</Form.Item>
														) : LANG('newtask.unsupported_mask_input_type') }
														<Form.Item
															label={LANG('newtask.mask_increment_mode')}
														>
															<Checkbox
																checked={this.state.enableMaskIncrementMode}
																onChange={this.onChangeEnableMaskIncrementMode}
															>
																{LANG('newtask.enable')}
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
															label={LANG('newtask.custom_charset_1')}
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder={LANG('newtask.set_custom_charset_1')}
																size="large"
																onChange={this.onChangeCustomCharset1}
																value={this.state.customCharset1}
																disabled={!(this.state.maskInputType === "text")}
															/>
														</Form.Item>
														<Form.Item
															label={LANG('newtask.custom_charset_2')}
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder={LANG('newtask.set_custom_charset_2')}
																size="large"
																onChange={this.onChangeCustomCharset2}
																value={this.state.customCharset2}
																disabled={!(this.state.maskInputType === "text")}
															/>
														</Form.Item>
														<Form.Item
															label={LANG('newtask.custom_charset_3')}
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder={LANG('newtask.set_custom_charset_3')}
																size="large"
																onChange={this.onChangeCustomCharset3}
																value={this.state.customCharset3}
																disabled={!(this.state.maskInputType === "text")}
															/>
														</Form.Item>
														<Form.Item
															label={LANG('newtask.custom_charset_4')}
														>
															<Input
																allowClear
																style={{ width: '100%' }}
																placeholder={LANG('newtask.set_custom_charset_4')}
																size="large"
																onChange={this.onChangeCustomCharset4}
																value={this.state.customCharset4}
																disabled={!(this.state.maskInputType === "text")}
															/>
														</Form.Item>
													</Col>
													<Col span={24}>
														<Row gutter={[18, 16]}>
															<Col span={12}>
																<Form.Item
																	label={LANG('newtask.dictionary')}
																	required
																>
																	<Select
																		showSearch
																		allowClear
																		style={{ width: '100%' }}
																		placeholder={LANG('newtask.select_dictionary')}
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
																	label={LANG('newtask.rule')}
																>
																	<Input
																		allowClear
																		style={{ width: '100%' }}
																		placeholder={LANG('newtask.set_rule')}
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
												LANG('newtask.select_attack_mode')
											)}
									</Form>
								) : this.state.step === 2 ? (
									<Form layout="vertical">
										<Collapse ghost onChange={this.onChangeAdvancedOptionsCollapse} activeKey={this.state.advancedOptionsCollapse}>
											<Panel header={LANG('newtask.general')} key="General">
												<Row gutter={[18, 16]}>
													<Col>
														<Checkbox
															checked={this.state.quiet}
															onChange={this.onChangeQuiet}
														>
															{LANG('newtask.quiet')}
														</Checkbox>
													</Col>
													<Col>
														<Checkbox
															checked={this.state.disablePotFile}
															onChange={this.onChangeDisablePotFile}
														>
															{LANG('newtask.disable_pot_file')}
														</Checkbox>
													</Col>
													<Col>
														<Checkbox
															checked={this.state.disableLogFile}
															onChange={this.onChangeDisableLogFile}
														>
															{LANG('newtask.disable_log_file')}
														</Checkbox>
													</Col>
													<Col>
														<Checkbox
															checked={this.state.enableOptimizedKernel}
															onChange={this.onChangeEnableOptimizedKernel}
														>
															{LANG('newtask.enable_optimized_kernel')}
														</Checkbox>
													</Col>
													<Col>
														<Checkbox
															checked={this.state.enableSlowerCandidateGenerators}
															onChange={this.onChangeEnableSlowerCandidateGenerators}
														>
															{LANG('newtask.enable_slower_candidate_generators')}
														</Checkbox>
													</Col>
													<Col>
														<Checkbox
															checked={this.state.removeFoundHashes}
															onChange={this.onChangeRemoveFoundHashes}
														>
															{LANG('newtask.remove_found_hashes')}
														</Checkbox>
													</Col>
													<Col>
														<Checkbox
															checked={this.state.ignoreUsernames}
															onChange={this.onChangeIgnoreUsernames}
														>
															{LANG('newtask.ignore_usernames')}
														</Checkbox>
													</Col>
													<Col>
														<Checkbox
															checked={this.state.disableSelfTest}
															onChange={this.onChangeDisableSelfTest}
														>
															{LANG('newtask.disable_self-test') + " (" + LANG('newtask.not_recommended')+ ")"}
														</Checkbox>
													</Col>
													<Col>
														<Checkbox
															checked={this.state.ignoreWarnings}
															onChange={this.onChangeIgnoreWarnings}
														>
															{LANG('newtask.ignore_warnings') + " (" + LANG('newtask.not_recommended')+ ")"}
														</Checkbox>
													</Col>
												</Row>
											</Panel>
											<Panel header={LANG('newtask.devices')} key="Devices">
												<Row gutter={[18, 16]}>
													<Col span={8}>
														<Form.Item
															label={LANG('newtask.devices_ids')}
														>
															<Select
																mode="multiple"
																allowClear
																style={{ width: '100%' }}
																placeholder={LANG('newtask.select_devices_ids')}
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
															label={LANG('newtask.devices_types')}
														>
															<Select
																mode="multiple"
																allowClear
																style={{ width: '100%' }}
																placeholder={LANG('newtask.select_devices_types')}
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
															label={LANG('newtask.workload_profile')}
															tooltip={
																<Table
																	columns={[
																		{
																			title: LANG('newtask.performance'),
																			dataIndex: 'performance',
																			key: 'Performance'
																		},
																		{
																			title: LANG('newtask.runtime'),
																			dataIndex: 'runtime',
																			key: 'Runtime'
																		},
																		{
																			title: LANG('newtask.power_consumption'),
																			dataIndex: 'powerConsumption',
																			key: 'Power Consumption'
																		},
																		{
																			title: LANG('newtask.desktop_impact'),
																			dataIndex: 'desktopImpact',
																			key: 'Desktop Impact'
																		}
																	]}
																	dataSource={[
																		{
																			key: '1',
																			performance: LANG('newtask.low'),
																			runtime: '2 ms',
																			powerConsumption: LANG('newtask.low'),
																			desktopImpact: LANG('newtask.minimal')
																		},
																		{
																			key: '2',
																			performance: LANG('newtask.default'),
																			runtime: '12 ms',
																			powerConsumption: LANG('newtask.economic'),
																			desktopImpact: LANG('newtask.noticeable')
																		},
																		{
																			key: '3',
																			performance: LANG('newtask.high'),
																			runtime: '96 ms',
																			powerConsumption: LANG('newtask.high'),
																			desktopImpact: LANG('newtask.unresponsive')
																		},
																		{
																			key: '4',
																			performance: LANG('newtask.nightmare'),
																			runtime: '480 ms',
																			powerConsumption: LANG('newtask.insane'),
																			desktopImpact: LANG('newtask.headless')
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
																placeholder={LANG('newtask.select_workload_profile')}
																size="large"
																onChange={this.onChangeWorkloadProfile}
																value={this.state.workloadProfile}
																filterOption={(input, option) =>
																	String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
																	String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
																}
															>
																<Option value={1} key={1}>{LANG('newtask.low')}</Option>
																<Option value={2} key={2}>{LANG('newtask.default')}</Option>
																<Option value={3} key={3}>{LANG('newtask.high')}</Option>
																<Option value={4} key={4}>{LANG('newtask.nightmare')}</Option>
															</Select>
														</Form.Item>
													</Col>
												</Row>
												<Row gutter={[18, 16]}>
													<Col>
														<Button
															loading={this.state.isLoadingDevicesInfo}
															onClick={this.onClickDevicesInfo}
														>
															{LANG('newtask.devices_info')}
														</Button>
													</Col>
													<Col>
														<Button
															loading={this.state.isLoadingBenchmark}
															onClick={this.onClickBenchmark}
														>
															{LANG('newtask.benchmark')}
														</Button>
													</Col>
												</Row>
											</Panel>
											<Panel header={LANG('newtask.markov')} key="Markov">
												<Row gutter={[18, 16]}>
													<Col>
														<Checkbox
															checked={this.state.markovDisable}
															onChange={this.onChangeMarkovDisable}
														>
															{LANG('newtask.disable_markov-chains')}
														</Checkbox>
													</Col>
													<Col>
														<Checkbox
															checked={this.state.markovClassic}
															onChange={this.onChangeMarkovClassic}
														>
															{LANG('newtask.enable_classic_markov-chains')}
														</Checkbox>
													</Col>
													<Col span={24}>
														<Row gutter={[18, 16]}>
															<Col span={8}>
																<Form.Item
																	label={LANG('newtask.markov_threshold')}
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
											<Panel header={LANG('newtask.monitor')} key="Monitor">
												<Row gutter={[18, 16]}>
													<Col>
														<Checkbox
															checked={this.state.disableMonitor}
															onChange={this.onChangeDisableMonitor}
														>
															{LANG('newtask.disable_monitor')}
														</Checkbox>
													</Col>
													<Col span={24}>
														<Row gutter={[18, 16]}>
															<Col span={8}>
																<Form.Item
																	label={LANG('newtask.temp_abort')}
																>
																	<Select
																		allowClear
																		style={{ width: '100%' }}
																		placeholder={LANG('newtask.select_temp_abort')}
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
											<Panel header={LANG('newtask.extra_arguments')} key="Extra Arguments">
												<Form.Item
													label={LANG('newtask.extra_arguments')}
												>
													<Input
														allowClear
														style={{ width: '100%' }}
														placeholder={LANG('newtask.set_extra_arguments')}
														size="large"
														onChange={this.onChangeExtraArguments}
														value={this.state.extraArguments.join(" ")}
													/>
												</Form.Item>
											</Panel>
											<Panel header={LANG('newtask.misc')} key="Misc">
												<Form.Item
													label={LANG('newtask.status_timer')}
												>
													<Select
														allowClear
														style={{ width: '100%' }}
														placeholder={LANG('newtask.select_status_timer')}
														size="large"
														onChange={this.onChangeStatusTimer}
														value={this.state.statusTimer}
														filterOption={(input, option) =>
															String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
															String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
														}
													>
														<Option value={10} key={10}>{"10 " + LANG('newtask.seconds')}</Option>
														<Option value={20} key={20}>{"20 " + LANG('newtask.seconds')}</Option>
														<Option value={30} key={30}>{"30 " + LANG('newtask.seconds')}</Option>
														<Option value={45} key={45}>{"45 " + LANG('newtask.seconds')}</Option>
														<Option value={60} key={60}>{"60 " + LANG('newtask.seconds')}</Option>
														<Option value={90} key={90}>{"90 " + LANG('newtask.seconds')}</Option>
														<Option value={120} key={120}>{"120 " + LANG('newtask.seconds')}</Option>
														<Option value={300} key={300}>{"300 " + LANG('newtask.seconds')}</Option>
													</Select>
												</Form.Item>
											</Panel>
										</Collapse>
									</Form>
								) : this.state.step === 3 ? (
									<Form layout="vertical">
										<Form.Item
											label={LANG('newtask.output_file')}
											extra={this.state.outputFile ? this.state.outputFile : null}
										>
											<Button
												type="primary"
												onClick={this.onChangeOutputFile}
												loading={this.state.isLoadingSetOutputFile}
											>
												{LANG('newtask.set_output_file')}
											</Button>
										</Form.Item>
										<Form.Item
											label={LANG('newtask.output_format')}
										>
											<Select
												mode="multiple"
												allowClear
												style={{ width: '100%' }}
												placeholder={LANG('newtask.select_output_format')}
												size="large"
												onChange={this.onChangeOutputFormat}
												value={this.state.outputFormat}
												filterOption={(input, option) =>
													String(option.value).toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
													String(option.children).toLowerCase().indexOf(input.toLowerCase()) >= 0
												}
											>
												<Option value={1} key={1}>{LANG('newtask.output_format_1')}</Option>
												<Option value={2} key={2}>{LANG('newtask.output_format_2')}</Option>
												<Option value={3} key={3}>{LANG('newtask.output_format_3')}</Option>
												<Option value={4} key={4}>{LANG('newtask.output_format_4')}</Option>
												<Option value={5} key={5}>{LANG('newtask.output_format_5')}</Option>
												<Option value={6} key={6}>{LANG('newtask.output_format_6')}</Option>
											</Select>
										</Form.Item>
									</Form>
								) : this.state.step === 4 ? (
									<Space direction="vertical">
										<Form layout="vertical">
											<Form.Item
												label={LANG('newtask.priority')}
												tooltip={
													<Typography>
														{LANG('newtask.priority_tooltip.part1')}
													<br />
														{LANG('newtask.priority_tooltip.part2')}
													</Typography>
												}
											>
												<InputNumber
													min={-1}
													max={999}
													value={this.state.priority}
													onChange={this.onChangePriority}
													bordered={true}
												/>
											</Form.Item>
										</Form>
										<Space size="large" direction="horizontal">
											<Button
												type="primary"
												icon={<PlusOutlined />}
												onClick={this.onClickCreateTask}
												loading={this.state.isLoadingCreateTask}
											>
												{LANG('newtask.create_task')}
											</Button>
												<Checkbox
													checked={this.state.preserveTaskConfig}
													onChange={this.onChangePreserveTaskConfig}
												>
													{LANG('newtask.preserve_task_config')}
												</Checkbox>
										</Space>
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

export default withTranslation()(NewTask);
