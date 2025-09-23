import { createFileRoute, useNavigate, Link } from "@tanstack/react-router";
import React from "react";
import { Title, Card, SimpleGrid, Text, Table, Badge } from "@mantine/core";
import {
  RadialBarChart,
  RadialBar,
  LineChart as RechartsLineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  PolarAngleAxis,
} from "recharts";

export const Route = createFileRoute("/_mainLayout/dashboard")({
  component: Dashboard,
});

function Dashboard() {
  const navigate = useNavigate();

  // Mock data for demonstration
  const oee = 85;
  const availability = 90;
  const performance = 82;
  const quality = 98;

  const devices = [
    { id: "DEV001", status: "running" as const },
    { id: "DEV002", status: "idle" as const },
    { id: "DEV003", status: "fault" as const },
    { id: "DEV004", status: "running" as const },
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

  const getStatusDotColor = (status: string) => {
    const colorMap: Record<string, string> = {
      running: "bg-green-500",
      idle: "bg-yellow-500",
      fault: "bg-red-500",
      default: "bg-gray-500",
    };
    return colorMap[status] || colorMap.default;
  };

  const getStatusBorderClass = (status: string) => {
    const colorMap: Record<string, string> = {
      running: "border-green-500",
      idle: "border-yellow-500",
      fault: "border-red-500",
      default: "border-gray-500",
    };
    return colorMap[status] || colorMap.default;
  };

  const telemetryData = [
    { time: "00:00", temperature: 25, humidity: 60 },
    { time: "00:10", temperature: 26, humidity: 62 },
    { time: "00:20", temperature: 24, humidity: 58 },
    { time: "00:30", temperature: 27, humidity: 65 },
    { time: "00:40", temperature: 25.5, humidity: 61 },
    { time: "00:50", temperature: 26.5, humidity: 63 },
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

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case "high":
        return "bg-red-100 text-red-800";
      case "medium":
        return "bg-yellow-100 text-yellow-800";
      case "low":
        return "bg-green-100 text-green-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  // Gauge data for Recharts RadialBarChart
  const gaugeData = (value: number, fillColor: string) => [
    { value, fill: fillColor },
  ];

  return (
    <>
      <Title order={1} className="mb-6">
        Main Dashboard
      </Title>

      {/* Top Gauges Section */}
      <div className="mb-6">
        {/* Large OEE Gauge in Center */}
        <div className="flex justify-center mb-4">
          <div className="w-full max-w-xs">
            <Text fz="28px" fw={"bold"} className="text-center mb-2">
              OEE Score
            </Text>
            <ResponsiveContainer width="100%" height={200}>
              <RadialBarChart
                innerRadius="50%"
                outerRadius="100%"
                barSize={20}
                startAngle={90}
                endAngle={-270}
                data={gaugeData(oee, "#10b981")}
              >
                <PolarAngleAxis
                  type="number"
                  domain={[0, 100]} // Define the domain from 0 to 100 for percentage
                  angleAxisId={0}
                  tick={false} // Hide the ticks if not needed
                />
                <RadialBar background dataKey="value" cornerRadius={50} />
                <text
                  x="50%"
                  y="50%"
                  textAnchor="middle"
                  dominantBaseline="middle"
                  style={{
                    fontSize: "1.5rem",
                    fontWeight: "bold",
                  }}
                >
                  {oee}%
                </text>
              </RadialBarChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Availability and Performance Gauges */}
        <div className="flex gap-4 mb-4 max-w-4xl mx-auto">
          <div className="flex-1">
            <Text fz="28px" fw={"bold"} className="text-center mb-2">
              Availability
            </Text>
            <ResponsiveContainer width="100%" height={150}>
              <RadialBarChart
                innerRadius="50%"
                outerRadius="80%"
                barSize={15}
                startAngle={90}
                endAngle={-270}
                data={gaugeData(availability, "#3b82f6")}
              >
                <PolarAngleAxis
                  type="number"
                  domain={[0, 100]} // Define the domain from 0 to 100 for percentage
                  angleAxisId={0}
                  tick={false} // Hide the ticks if not needed
                />
                <RadialBar background dataKey="value" cornerRadius={25} />
                <text
                  x="50%"
                  y="50%"
                  textAnchor="middle"
                  dominantBaseline="middle"
                  style={{
                    fontSize: "1.2rem",
                    fontWeight: "bold",
                  }}
                >
                  {availability}%
                </text>
              </RadialBarChart>
            </ResponsiveContainer>
          </div>
          <div className="flex-1">
            <Text fz="28px" fw={"bold"} className="text-center mb-2">
              Performance
            </Text>
            <ResponsiveContainer width="100%" height={150}>
              <RadialBarChart
                innerRadius="50%"
                outerRadius="80%"
                barSize={15}
                startAngle={90}
                endAngle={-270}
                data={gaugeData(performance, "#f59e0b")}
              >
                <PolarAngleAxis
                  type="number"
                  domain={[0, 100]} // Define the domain from 0 to 100 for percentage
                  angleAxisId={0}
                  tick={false} // Hide the ticks if not needed
                />

                <RadialBar background dataKey="value" cornerRadius={25} />
                <text
                  x="50%"
                  y="50%"
                  textAnchor="middle"
                  dominantBaseline="middle"
                  style={{
                    fontSize: "1.2rem",
                    fontWeight: "bold",
                  }}
                >
                  {performance}%
                </text>
              </RadialBarChart>
            </ResponsiveContainer>
          </div>
          <div className="flex-1">
            <Text fz="28px" fw={"bold"} className="text-center mb-2">
              Quality
            </Text>
            <ResponsiveContainer width="100%" height={150}>
              <RadialBarChart
                innerRadius="50%"
                outerRadius="80%"
                barSize={15}
                startAngle={90}
                endAngle={-270}
                data={gaugeData(quality, "#10b981")}
              >
                <PolarAngleAxis
                  type="number"
                  domain={[0, 100]} // Define the domain from 0 to 100 for percentage
                  angleAxisId={0}
                  tick={false} // Hide the ticks if not needed
                />
                <RadialBar background dataKey="value" cornerRadius={25} />
                <text
                  x="50%"
                  y="50%"
                  textAnchor="middle"
                  dominantBaseline="middle"
                  style={{
                    fontSize: "1.2rem",
                    fontWeight: "bold",
                  }}
                >
                  {quality}%
                </text>
              </RadialBarChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>

      {/* Central Area: Device Status and Telemetry Chart */}
      <SimpleGrid cols={{ base: 1, md: 2 }} spacing="md" className="mb-6">
        {/* Device Status Panel */}
        <Card
          shadow="sm"
          padding="lg"
          radius="md"
          withBorder
          className="bg-white/10 border-white/20 backdrop-blur-sm overflow-hidden"
        >
          <Title order={3} mb={"sm"}>
            Device Status
          </Title>
          <div className="-mx-6 space-y-1">
            {devices.map((device) => (
              <Link
                key={device.id}
                to="/device/$deviceId"
                params={{ deviceId: device.id }}
                className={`w-full flex items-center justify-between cursor-pointer hover:bg-white/20 px-6 py-2 rounded-none transition-all duration-200 border-b-0 hover:border-b-2 ${getStatusBorderClass(device.status)} block`}
              >
                <div className="flex items-center space-x-3 flex-1 min-w-0">
                  <div
                    className={`w-3 h-3 rounded-full ${getStatusDotColor(device.status)}`}
                  ></div>
                  <Text fw="bold" className="text-lg truncate">
                    {device.id}
                  </Text>
                </div>
                <Badge
                  color={getStatusColor(device.status)}
                  variant="light"
                  size="sm"
                >
                  {device.status.toUpperCase()}
                </Badge>
              </Link>
            ))}
          </div>
        </Card>

        {/* Real-time Telemetry Chart */}
        <Card
          shadow="sm"
          padding="lg"
          radius="md"
          withBorder
          className="bg-white/10 border-white/20 backdrop-blur-sm"
        >
          <Title order={3} mb={"sm"}>
            Real-time Telemetry (Last 60 Minutes)
          </Title>
          <ResponsiveContainer width="100%" height={300}>
            <RechartsLineChart data={telemetryData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="time" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Line
                type="monotone"
                dataKey="temperature"
                stroke="#ef4444"
                strokeWidth={2}
                name="Temperature (°C)"
              />
              <Line
                type="monotone"
                dataKey="humidity"
                stroke="#3b82f6"
                strokeWidth={2}
                name="Humidity (%)"
              />
            </RechartsLineChart>
          </ResponsiveContainer>
        </Card>
      </SimpleGrid>

      {/* Bottom: Alerts and Events Table */}
      <Card
        shadow="sm"
        padding="lg"
        radius="md"
        withBorder
        className="bg-white/10 border-white/20 backdrop-blur-sm"
      >
        <Title order={3} mb={"sm"}>
          Recent Alerts and Events
        </Title>
        <Table striped highlightOnHover>
          <Table.Thead>
            <Table.Tr>
              <Table.Th>Time</Table.Th>
              <Table.Th>Device ID</Table.Th>
              <Table.Th>Alert Type</Table.Th>
              <Table.Th>Message</Table.Th>
            </Table.Tr>
          </Table.Thead>
          <Table.Tbody>
            {alerts.map((alert, index) => (
              <Table.Tr
                key={index}
                className={getSeverityColor(alert.severity)}
              >
                <Table.Td>{alert.time}</Table.Td>
                <Table.Td>{alert.deviceId}</Table.Td>
                <Table.Td>{alert.type}</Table.Td>
                <Table.Td>{alert.message}</Table.Td>
              </Table.Tr>
            ))}
          </Table.Tbody>
        </Table>
      </Card>
    </>
  );
}
