import { Title, Card, Text } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_mainLayout/device/$deviceId")({
  component: DeviceDetail,
});

function DeviceDetail() {
  const { deviceId } = Route.useParams();

  return (
    <>
      <Title order={1} className="mb-6 text-white">
        Device Detail View
      </Title>
      <Text className="mb-4 text-lg text-white">
        Details for Device ID: {deviceId}
      </Text>
      <Card
        shadow="sm"
        padding="lg"
        radius="md"
        withBorder
        className="mb-4 bg-white/10 border-white/20 backdrop-blur-sm"
      >
        <Title order={3} className="mb-2 text-white">
          Status
        </Title>
        <Text className="text-green-400">Online</Text>
      </Card>
      <Card
        shadow="sm"
        padding="lg"
        radius="md"
        withBorder
        className="mb-4 bg-white/10 border-white/20 backdrop-blur-sm"
      >
        <Title order={3} className="mb-2 text-white">
          Last Telemetry
        </Title>
        <Text className="text-sm text-gray-300">2025-09-23 15:00:00</Text>
      </Card>
      <Card
        shadow="sm"
        padding="lg"
        radius="md"
        withBorder
        className="bg-white/10 border-white/20 backdrop-blur-sm"
      >
        <Title order={3} className="mb-2 text-white">
          Alerts
        </Title>
        <Text className="text-yellow-400">No active alerts</Text>
      </Card>
    </>
  );
}
