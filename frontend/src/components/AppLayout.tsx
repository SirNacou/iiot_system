import { AppShell, Divider, NavLink, Text } from "@mantine/core";
import { Link } from "@tanstack/react-router";
import { Activity, BarChart3, Settings, Zap } from "lucide-react";
import { ReactNode } from "react";

interface AppLayoutProps {
  children: ReactNode;
}

export default function AppLayout({ children }: AppLayoutProps) {
  return (
    <AppShell
      navbar={{
        width: 280,
        breakpoint: "sm",
      }}
      className="min-h-screen"
    >
      <AppShell.Navbar>
        {/* Header Section */}
        <AppShell.Section>
          <div className="p-6">
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
                <Zap className="w-6 h-6 text-white" />
              </div>
              <div>
                <Text size="xl" fw={700}>
                  IIoT System
                </Text>
                <Text size="sm" c="dimmed">
                  Industrial IoT Platform
                </Text>
              </div>
            </div>
          </div>
        </AppShell.Section>

        <Divider />

        {/* Navigation Section */}
        <AppShell.Section grow className="p-4">
          <div className="space-y-2">
            <Text size="xs" fw={600} c="dimmed" tt="uppercase" className="px-3 mb-4">
              Main Navigation
            </Text>

            <NavLink
              component={Link}
              to="/dashboard"
              label="Dashboard"
              description="Real-time overview and monitoring"
              leftSection={<Activity size={20} />}
            />

            <NavLink
              component={Link}
              to="/analytics"
              label="Analytics & Reports"
              description="Data insights and performance metrics"
              leftSection={<BarChart3 size={20} />}
            />

            <Divider className="my-4" />

            <Text size="xs" fw={600} c="dimmed" tt="uppercase" className="px-3 mb-4">
              System
            </Text>

            <NavLink
              component={Link}
              to="/settings"
              label="Settings"
              description="System configuration and preferences"
              leftSection={<Settings size={20} />}
            />
          </div>
        </AppShell.Section>

        <Divider />

        {/* Footer Section */}
        <AppShell.Section>
          <div className="p-4">
            <div className="flex items-center space-x-2">
              <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
              <Text size="sm" fw={500}>
                System Online
              </Text>
            </div>
            <Text size="sm" c="dimmed" className="mt-1">
              Last updated: {new Date().toLocaleTimeString()}
            </Text>
          </div>
        </AppShell.Section>
      </AppShell.Navbar>

      <AppShell.Main className="p-0">
        <div className="w-full h-full">{children}</div>
      </AppShell.Main>
    </AppShell>
  );
}
