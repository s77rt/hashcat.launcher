import React, { Component } from "react";
import { message, Statistic, Row, Col, Card, Select, Typography, Upload, Button, Space, Form, Radio, Divider, Collapse, Checkbox, Tabs, Steps } from 'antd';
import { FileOutlined, AimOutlined, ToolOutlined, ExportOutlined, ExperimentOutlined, SyncOutlined } from '@ant-design/icons';

import EventBus from "./eventbus/EventBus";

import data from "./data/data";
import { getHashes } from './data/hashes';
import { getAlgorithms } from './data/algorithms';
import { getDictionaries } from './data/dictionaries';
import { getRules } from './data/rules';
import { getMasks } from './data/masks';

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

		this.state = {
			_dictionaries: getDictionaries(),
			_rules: getRules(),
			_masks: getMasks(),
			_hashes: getHashes(),
			_algorithms: getAlgorithms(),

			isLoadingRescan: false,
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
				message.error(e);
			}
			this.setState({isLoadingRescan: false});
		})
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
			</>
		)
	}
}

export default Settings;
