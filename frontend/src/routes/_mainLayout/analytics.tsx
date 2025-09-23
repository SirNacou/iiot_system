import { createFileRoute } from '@tanstack/react-router'
import { Title, Card, SimpleGrid, Text, List, ListItem } from '@mantine/core'

export const Route = createFileRoute('/_mainLayout/analytics')({
  component: Analytics,
})

function Analytics() {
  return (
    <>
      <Title order={1} className="mb-6 text-white">Analytics and Reports</Title>
      <Text className="mb-4 text-lg text-white">IIoT System Analytics: Telemetry trends, alert history, and production reports.</Text>
      <SimpleGrid cols={{ base: 1, md: 2 }} spacing="md">
        <Card shadow="sm" padding="lg" radius="md" withBorder className="bg-white/10 border-white/20 backdrop-blur-sm">
          <Title order={3} className="mb-2 text-white">Telemetry Trends</Title>
          <Text className="text-sm text-gray-300 mb-2">Chart placeholder for device metrics over time.</Text>
          <div className="h-32 bg-gray-800 rounded-md flex items-center justify-center">
            <Text className="text-gray-400">Line Chart</Text>
          </div>
        </Card>
        <Card shadow="sm" padding="lg" radius="md" withBorder className="bg-white/10 border-white/20 backdrop-blur-sm">
          <Title order={3} className="mb-2 text-white">Alert History</Title>
          <Text className="text-sm text-gray-300 mb-2">Recent alerts and resolutions.</Text>
          <List spacing="xs" size="sm" center className="text-sm">
            <ListItem className="text-red-400">High temperature alert - Resolved</ListItem>
            <ListItem className="text-yellow-400">Low pressure warning - Active</ListItem>
            <ListItem className="text-green-400">Status update - Normal</ListItem>
          </List>
        </Card>
      </SimpleGrid>
      <Card shadow="sm" padding="lg" radius="md" withBorder className="mt-6 bg-white/10 border-white/20 backdrop-blur-sm">
        <Title order={3} className="mb-2 text-white">Production Reports</Title>
        <Text className="text-sm text-gray-300 mb-4">Summary of production events and efficiency.</Text>
        <SimpleGrid cols={2} spacing="md">
          <div>
            <Text className="text-2xl font-bold text-white">95%</Text>
            <Text className="text-xs text-gray-400">Uptime</Text>
          </div>
          <div>
            <Text className="text-2xl font-bold text-white">1,200</Text>
            <Text className="text-xs text-gray-400">Units Produced</Text>
          </div>
        </SimpleGrid>
      </Card>
    </>
  )
}