import { withTranslation } from 'react-i18next';
import i18n, { SupportedLanguages } from '../i18n';

import React, { Component } from "react";
import { Popconfirm, Layout, PageHeader, message, Statistic, Row, Col, Card, Select, Typography, Upload, Button, Space, Form, Radio, Divider, Collapse, Checkbox, Tabs, Steps } from 'antd';
import { TranslationOutlined, FileOutlined, AimOutlined, ToolOutlined, ExportOutlined, ExperimentOutlined, SyncOutlined } from '@ant-design/icons';

import EventBus from "./eventbus/EventBus";

import moment from "moment/min/moment-with-locales"

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
		this.onChangeLanguage = this.onChangeLanguage.bind(this);

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

	onChangeLanguage(e) {
		i18n.changeLanguage(e).then(
			() => {
				moment.locale(e);
				if (typeof window.GOsettingsChangeLanguage === "function") {
					window.GOsettingsChangeLanguage(e).then(
						() => null,
						error => {
							message.error(error);
						}
					);
				}
			},
			error => {
				message.error(error);
			}
		);
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
		const LANG = this.props.t;
		return (
			<>
				<PageHeader
					title={LANG('settings.title')}
					extra={[
						<Select
							suffixIcon={<TranslationOutlined />}
							style={{textTransform: 'uppercase'}}
							onChange={this.onChangeLanguage}
							value={i18n.language}
						>
							{SupportedLanguages.map(lang =>
								<Option style={{textTransform: 'uppercase'}} value={lang} key={lang}>{lang}</Option>
							)}
						</Select>					
					]}
				/>
				<Content style={{ padding: '16px 24px' }}>
					<Row gutter={[16, 14]}>
						<Col span={12}>
							<Statistic title={LANG('settings.hashes')} value={this.state._hashes.length} />
						</Col>
						<Col span={12}>
							<Statistic title={LANG('settings.algorithms')} value={Object.keys(this.state._algorithms).length} />
						</Col>
						<Col span={12}>
							<Statistic title={LANG('settings.dictionaries')} value={this.state._dictionaries.length} />
						</Col>
						<Col span={12}>
							<Statistic title={LANG('settings.rules')} value={this.state._rules.length} />
						</Col>
						<Col span={12}>
							<Statistic title={LANG('settings.masks')} value={this.state._masks.length} />
						</Col>
						<Col span={24}>
							<Button
								icon={<SyncOutlined />}
								type="primary"
								onClick={this.onClickRescan}
								loading={this.state.isLoadingRescan}
							>
								{LANG('settings.rescan')}
							</Button>
						</Col>
					</Row>
					<Row style={{ marginTop: "2rem" }} gutter={[16, 14]}>
						<Col span={24}>
							<Statistic
								title={LANG('settings.task_counter')}
								value={this.state.taskCounter}
							/>
							<Space>
							<Button
								style={{ marginTop: 16 }}
								type="default"
								onClick={this.onClickRefreshTaskCounter}
								loading={this.state.isLoadingRefreshTaskCounter}
							>
								{LANG('settings.refresh')}
							</Button>
							<Popconfirm
								placement="topRight"
								title={LANG('settings.reset_counter_confirm.message')}
								onConfirm={this.onClickResetTaskCounter}
								okText={LANG('settings.reset_counter_confirm.yes')}
								cancelText={LANG('settings.reset_counter_confirm.no')}
							>
								<Button
									style={{ marginTop: 16 }}
									type="danger"
									loading={this.state.isLoadingResetTaskCounter}
								>
									{LANG('settings.reset_counter')}
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

export default withTranslation()(Settings);
