import React, { Component } from "react";
import { message, Statistic, Row, Col, Card, Select, Typography, Upload, Button, Space, Form, Radio, Divider, Collapse, Checkbox, Tabs, Steps } from 'antd';
import { FileOutlined, AimOutlined, ToolOutlined, ExportOutlined, ExperimentOutlined, SyncOutlined } from '@ant-design/icons';

const { Dragger } = Upload;
const { Option } = Select;
const { Title, Paragraph, Text, Link } = Typography;
const { Panel } = Collapse;
const { TabPane } = Tabs;
const { Step } = Steps;

class About extends Component {
	render() {
		return (
			<>
				<Typography>
						<Paragraph>
								hashcat.launcher is a cross-platform app that run and control hashcat
						</Paragraph>
						<Paragraph>
								it is designed to make it easier to use hashcat offering a friendly graphical user interface
						</Paragraph>
						<Title level={5}>Contribute</Title>
								<Link target="_blank" href="https://github.com/s77rt/hashcat.launcher/issues/new/">Report a bug / Request a feature</Link>
						<Title level={5}>License</Title>
						<Paragraph>
								hashcat.launcher is licensed under the MIT License
						</Paragraph>
						<Paragraph>
								Copyright &copy; 2021 Abdelhafidh Belalia (s77rt)
						</Paragraph>
				</Typography>
			</>
		)
	}
}

export default About;
