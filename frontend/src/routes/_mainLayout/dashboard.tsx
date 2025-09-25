import { Badge, Card, Text, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { Activity, AlertTriangle, CheckCircle, Clock } from "lucide-react";
import {
    CartesianGrid,
    Legend,
    Line,
    PolarAngleAxis,
    RadialBar,
    RadialBarChart,
    LineChart as RechartsLineChart,
    ResponsiveContainer,
    Tooltip,
    XAxis,
    YAxis,
} from "recharts";

export const Route = createFileRoute("/_mainLayout/dashboard")({
    component: Dashboard,
});

function Dashboard() {
    // Mock data for demonstration
    const oee = 85;
    const availability = 90;
    const performance = 82;
    const quality = 98;

    const devices = [
        { id: "DEV001", status: "running" as const, name: "Production Line A" },
        { id: "DEV002", status: "idle" as const, name: "Assembly Unit B" },
        { id: "DEV003", status: "fault" as const, name: "Quality Station C" },
        { id: "DEV004", status: "running" as const, name: "Packaging Line D" },
    ];

    const getStatusColor = (status: string) => {
        switch (status) {
            case "running":
                return "green";
            case "idle":
                return "yellow";
            case "fault":
                return "red";
            default:
                return "gray";
        }
    };

    const getStatusIcon = (status: string) => {
        switch (status) {
            case "running":
                return <CheckCircle className="w-4 h-4 text-green-500" />;
            case "idle":
                return <Clock className="w-4 h-4 text-yellow-500" />;
            case "fault":
                return <AlertTriangle className="w-4 h-4 text-red-500" />;
            default:
                return <Activity className="w-4 h-4 text-gray-500" />;
        }
    };

    const telemetryData = [
        { time: "00:00", temperature: 25, humidity: 60, pressure: 1013 },
        { time: "00:10", temperature: 26, humidity: 62, pressure: 1015 },
        { time: "00:20", temperature: 24, humidity: 58, pressure: 1012 },
        { time: "00:30", temperature: 27, humidity: 65, pressure: 1018 },
        { time: "00:40", temperature: 25.5, humidity: 61, pressure: 1016 },
        { time: "00:50", temperature: 26.5, humidity: 63, pressure: 1014 },
    ];

    const alerts = [
        {
            time: "16:00",
            deviceId: "DEV001",
            type: "High Temperature",
            message: "Temperature exceeded 30°C",
            severity: "high" as const,
        },
        {
            time: "15:45",
            deviceId: "DEV002",
            type: "Idle Detected",
            message: "Device idle for over 10 minutes",
            severity: "medium" as const,
        },
        {
            time: "15:30",
            deviceId: "DEV003",
            type: "Fault",
            message: "Sensor failure detected",
            severity: "high" as const,
        },
    ];

    const getSeverityStyle = (severity: string) => {
        switch (severity) {
            case "high":
                return "bg-red-50 border-l-4 border-red-500";
            case "medium":
                return "bg-yellow-50 border-l-4 border-yellow-500";
            case "low":
                return "bg-green-50 border-l-4 border-green-500";
            default:
                return "bg-gray-50 border-l-4 border-gray-500";
        }
    };

    const getSeverityColor = (severity: string) => {
        switch (severity) {
            case "high":
                return "red";
            case "medium":
                return "yellow";
            case "low":
                return "green";
            default:
                return "gray";
        }
    };

    const gaugeData = (value: number, fillColor: string) => [{ value, fill: fillColor }];

    const getOEEColor = (value: number) => {
        if (value >= 85) return "#10b981"; // green
        if (value >= 75) return "#f59e0b"; // yellow
        return "#ef4444"; // red
    };

    return (
        <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50 p-6">
            <div className="max-w-7xl mx-auto">
                {/* Header */}
                <div className="mb-8">
                    <Title
                        order={1}
                        className="text-4xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent"
                    >
                        Production Dashboard
                    </Title>
                    <Text className="text-gray-600 mt-2 text-lg">Real-time monitoring and analytics</Text>
                </div>

                {/* KPI Gauges */}
                <div className="mb-8">
                    {/* Main OEE Gauge */}
                    <div className="flex justify-center mb-8">
                        <Card
                            shadow="xl"
                            padding="xl"
                            radius="xl"
                            className="bg-white/80 backdrop-blur-sm border border-white/20 shadow-2xl"
                        >
                            <div className="text-center">
                                <Text size="xl" fw={600} className="text-gray-700 mb-4">
                                    Overall Equipment Effectiveness
                                </Text>
                                <div className="relative">
                                    <ResponsiveContainer width={280} height={280}>
                                        <RadialBarChart
                                            innerRadius="60%"
                                            outerRadius="90%"
                                            barSize={25}
                                            startAngle={90}
                                            endAngle={-270}
                                            data={gaugeData(oee, getOEEColor(oee))}
                                        >
                                            <PolarAngleAxis
                                                type="number"
                                                domain={[0, 100]}
                                                angleAxisId={0}
                                                tick={false}
                                            />
                                            <RadialBar
                                                background={{ fill: "#f1f5f9" }}
                                                dataKey="value"
                                                cornerRadius={15}
                                                className="drop-shadow-sm"
                                            />
                                        </RadialBarChart>
                                    </ResponsiveContainer>
                                    <div className="absolute inset-0 flex items-center justify-center">
                                        <div className="text-center">
                                            <div className="text-5xl font-bold text-gray-800">{oee}%</div>
                                            <div className="text-sm text-gray-500 font-medium">OEE Score</div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </Card>
                    </div>

                    {/* Sub-metrics */}
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                        {[
                            { title: "Availability", value: availability, color: "#3b82f6" },
                            { title: "Performance", value: performance, color: "#f59e0b" },
                            { title: "Quality", value: quality, color: "#10b981" },
                        ].map((metric, index) => (
                            <Card
                                key={metric.title}
                                shadow="lg"
                                padding="lg"
                                radius="lg"
                                className="bg-white/80 backdrop-blur-sm border border-white/20 hover:shadow-xl transition-all duration-300"
                            >
                                <div className="text-center">
                                    <Text size="lg" fw={600} className="text-gray-700 mb-3">
                                        {metric.title}
                                    </Text>
                                    <ResponsiveContainer width="100%" height={140}>
                                        <RadialBarChart
                                            innerRadius="50%"
                                            outerRadius="80%"
                                            barSize={18}
                                            startAngle={90}
                                            endAngle={-270}
                                            data={gaugeData(metric.value, metric.color)}
                                        >
                                            <PolarAngleAxis
                                                type="number"
                                                domain={[0, 100]}
                                                angleAxisId={0}
                                                tick={false}
                                            />
                                            <RadialBar
                                                background={{ fill: "#f1f5f9" }}
                                                dataKey="value"
                                                cornerRadius={10}
                                            />
                                        </RadialBarChart>
                                    </ResponsiveContainer>
                                    <div className="absolute inset-0 flex items-center justify-center mt-6">
                                        <div className="text-2xl font-bold text-gray-800">{metric.value}%</div>
                                    </div>
                                </div>
                            </Card>
                        ))}
                    </div>
                </div>

                {/* Main Content Grid */}
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
                    {/* Device Status */}
                    <Card
                        shadow="lg"
                        padding="lg"
                        radius="lg"
                        className="bg-white/80 backdrop-blur-sm border border-white/20"
                    >
                        <div className="flex items-center mb-4">
                            <Activity className="w-6 h-6 text-blue-600 mr-2" />
                            <Title order={3} className="text-gray-800">
                                Device Status
                            </Title>
                        </div>
                        <div className="space-y-3">
                            {devices.map(device => (
                                <div
                                    key={device.id}
                                    className="flex items-center justify-between p-4 rounded-lg bg-gradient-to-r from-white to-gray-50 border border-gray-200 hover:shadow-md transition-all duration-200 cursor-pointer group"
                                >
                                    <div className="flex items-center space-x-3">
                                        {getStatusIcon(device.status)}
                                        <div>
                                            <Text
                                                fw={600}
                                                className="text-gray-800 group-hover:text-blue-600 transition-colors"
                                            >
                                                {device.id}
                                            </Text>
                                            <Text size="sm" className="text-gray-500">
                                                {device.name}
                                            </Text>
                                        </div>
                                    </div>
                                    <Badge
                                        color={getStatusColor(device.status)}
                                        variant="light"
                                        size="md"
                                        className="shadow-sm"
                                    >
                                        {device.status.toUpperCase()}
                                    </Badge>
                                </div>
                            ))}
                        </div>
                    </Card>

                    {/* Telemetry Chart */}
                    <Card
                        shadow="lg"
                        padding="lg"
                        radius="lg"
                        className="bg-white/80 backdrop-blur-sm border border-white/20"
                    >
                        <div className="flex items-center mb-4">
                            <Activity className="w-6 h-6 text-green-600 mr-2" />
                            <Title order={3} className="text-gray-800">
                                Real-time Telemetry
                            </Title>
                        </div>
                        <ResponsiveContainer width="100%" height={300}>
                            <RechartsLineChart data={telemetryData}>
                                <CartesianGrid strokeDasharray="3 3" stroke="#e2e8f0" />
                                <XAxis dataKey="time" stroke="#64748b" fontSize={12} />
                                <YAxis stroke="#64748b" fontSize={12} />
                                <Tooltip
                                    contentStyle={{
                                        backgroundColor: "white",
                                        border: "1px solid #e2e8f0",
                                        borderRadius: "8px",
                                        boxShadow: "0 4px 6px -1px rgb(0 0 0 / 0.1)",
                                    }}
                                />
                                <Legend />
                                <Line
                                    type="monotone"
                                    dataKey="temperature"
                                    stroke="#ef4444"
                                    strokeWidth={3}
                                    name="Temperature (°C)"
                                    dot={{ fill: "#ef4444", strokeWidth: 2, r: 4 }}
                                />
                                <Line
                                    type="monotone"
                                    dataKey="humidity"
                                    stroke="#3b82f6"
                                    strokeWidth={3}
                                    name="Humidity (%)"
                                    dot={{ fill: "#3b82f6", strokeWidth: 2, r: 4 }}
                                />
                            </RechartsLineChart>
                        </ResponsiveContainer>
                    </Card>
                </div>

                {/* Alerts Table */}
                <Card
                    shadow="lg"
                    padding="lg"
                    radius="lg"
                    className="bg-white/80 backdrop-blur-sm border border-white/20"
                >
                    <div className="flex items-center mb-6">
                        <AlertTriangle className="w-6 h-6 text-orange-600 mr-2" />
                        <Title order={3} className="text-gray-800">
                            Recent Alerts & Events
                        </Title>
                    </div>

                    <div className="space-y-3">
                        {alerts.map((alert, index) => (
                            <div
                                key={index}
                                className={`p-4 rounded-lg ${getSeverityStyle(alert.severity)} hover:shadow-md transition-all duration-200`}
                            >
                                <div className="flex items-center justify-between">
                                    <div className="flex-1">
                                        <div className="flex items-center space-x-3 mb-1">
                                            <Badge color={getSeverityColor(alert.severity)} variant="filled" size="sm">
                                                {alert.type}
                                            </Badge>
                                            <Text size="sm" className="text-gray-600">
                                                {alert.time}
                                            </Text>
                                        </div>
                                        <Text fw={500} className="text-gray-800 mb-1">
                                            {alert.deviceId}
                                        </Text>
                                        <Text size="sm" className="text-gray-600">
                                            {alert.message}
                                        </Text>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}
