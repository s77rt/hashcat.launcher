import { withTranslation } from 'react-i18next';

import React, { Component } from "react";
import { Layout, PageHeader, message, Statistic, Row, Col, Card, Select, Typography, Upload, Button, Space, Form, Radio, Divider, Collapse, Checkbox, Tabs, Steps } from 'antd';
import { FileOutlined, AimOutlined, ToolOutlined, ExportOutlined, ExperimentOutlined, SyncOutlined } from '@ant-design/icons';

const { Content } = Layout;
const { Dragger } = Upload;
const { Option } = Select;
const { Title, Paragraph, Text, Link } = Typography;
const { Panel } = Collapse;
const { TabPane } = Tabs;
const { Step } = Steps;

class Help extends Component {
	render() {
		const LANG = this.props.t;
		return (
			<>
				<PageHeader
					title={LANG('help.title')}
				/>
				<Content style={{ padding: '16px 24px' }}>
					<Typography>
						<Title level={5}>{LANG('help.faq_long')}</Title>
						<Collapse bordered={false}>
							<Panel header={LANG('help.question_hashcat_supported_versions')} key="1">
								{LANG('help.answer_hashcat_supported_versions')}
							</Panel>
							<Panel header={LANG('help.question_add_files')} key="2">
								{LANG('help.answer_add_files.part1')}
								<ul>
									<li>{LANG('help.answer_add_files.hashes') + ": "} <Text code>/hashcat/hashes</Text></li>
									<li>{LANG('help.answer_add_files.dictionaries') + ": "} <Text code>/hashcat/dictionaries</Text></li>
									<li>{LANG('help.answer_add_files.rules') + ": "} <Text code>/hashcat/rules</Text></li>
									<li>{LANG('help.answer_add_files.masks') + ": "} <Text code>/hashcat/masks</Text></li>
								</ul>
							</Panel>
							<Panel header={LANG('help.question_empty_algorithms_list')} key="3">
								{LANG('help.answer_empty_algorithms_list.part1') + " " + LANG('help.answer_empty_algorithms_list.part2')}
								<br />
								{LANG('help.answer_empty_algorithms_list.part3') + " " + LANG('help.answer_empty_algorithms_list.part4')} <Text code>/hashcat</Text>.
							</Panel>
							<Panel header={LANG('help.question_added_file_but_not_listed')} key="4">
								{LANG('help.answer_added_file_but_not_listed.part1') + " " + LANG('help.answer_added_file_but_not_listed.part2')}
							</Panel>
							<Panel header={LANG('help.question_difference_idle_vs_queued')} key="5">
								{LANG('help.answer_difference_idle_vs_queued.part1')}
								<br />
								{LANG('help.answer_difference_idle_vs_queued.part2')}
								<br />
								<br />
								{LANG('help.answer_difference_idle_vs_queued.part3')}
								<br />
								{LANG('help.answer_difference_idle_vs_queued.part4')}
							</Panel>
						</Collapse>
					</Typography>
				</Content>
			</>
		)
	}
}

export default withTranslation()(Help);
