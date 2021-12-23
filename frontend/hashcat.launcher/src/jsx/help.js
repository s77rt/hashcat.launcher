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
		return (
			<>
				<PageHeader
					title="Help"
				/>
				<Content style={{ padding: '16px 24px' }}>
					<Typography>
						<Title level={5}>Frequently Asked Questions</Title>
						<Collapse bordered={false}>
							<Panel header="What hashcat versions are supported" key="1">
								hashcat.launcher supports hashcat v6.2.1 and higher.
							</Panel>
							<Panel header="How to add hashes, dictionaries, etc..." key="2">
								Files are expected to be in the following folders:
								<ul>
									<li>Hashes: <Text code>/hashcat/hashes</Text></li>
									<li>Dictionaries: <Text code>/hashcat/dictionaries</Text></li>
									<li>Rules: <Text code>/hashcat/rules</Text></li>
									<li>Masks: <Text code>/hashcat/masks</Text></li>
								</ul>
							</Panel>
							<Panel header="Algorithms list is empty" key="3">
								Algorithms get loaded automatically and depends on hashcat.
								Make sure hashcat exists then go to Settings and click on Rescan.
								<br />
								hashcat is expected to be in the same directory as hashcat.launcher
								inside a subfolder <Text code>/hashcat</Text>.
							</Panel>
							<Panel header="I added a file but it's not in the list" key="4">
								Try perform a manual scan.
								Go to Settings then click on Rescan.
							</Panel>
							<Panel header="What is the difference between an idle task and a queued task" key="5">
								An idle task have a priority that is less than zero
								<br />
								while, a queued task have a priority that is equal or greater than zero.
								<br />
								<br />
								An idle task can only be started manually
								<br />
								while, a queued task can be started manually or automatically.
							</Panel>
						</Collapse>
					</Typography>
				</Content>
			</>
		)
	}
}

export default Help;
