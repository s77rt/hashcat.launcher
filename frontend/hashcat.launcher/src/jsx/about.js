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

class About extends Component {
	render() {
		const LANG = this.props.t;
		return (
			<>
				<PageHeader
					title={LANG('about.title')}
				/>
				<Content style={{ padding: '16px 24px' }}>
					<Typography>
							<Paragraph>
								{LANG('about.paragraph1')}
							</Paragraph>
							<Paragraph>
								{LANG('about.paragraph2')}
							</Paragraph>
							<Title level={5}>{LANG('about.contribute')}</Title>
									<Link target="_blank" href="https://github.com/s77rt/hashcat.launcher/issues/new/">{LANG('about.report_a_bug') + " / " + LANG('about.request_a_feature')}</Link>
							<Title level={5}>{LANG('about.license')}</Title>
							<Paragraph>
								{LANG('about.license_paragraph')}
							</Paragraph>
							<Title level={5}>{LANG('about.copyright')}</Title>
							<Paragraph>
								{LANG('about.copyright')} &copy; {new Date().getFullYear()} Abdelhafidh Belalia (s77rt)
							</Paragraph>
					</Typography>
				</Content>
			</>
		)
	}
}

export default withTranslation()(About);
