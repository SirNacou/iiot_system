import { AppShell, Container, NavLink, Text } from "@mantine/core";
import { Link } from "@tanstack/react-router";
import { ReactNode } from "react";

interface AppLayoutProps {
  children: ReactNode;
}

export default function AppLayout({ children }: AppLayoutProps) {
  return (
    <AppShell
      navbar={{
        width: 300,
        breakpoint: "sm",
      }}
      className="min-h-screen"
    >
      <AppShell.Navbar>
        <AppShell.Section>
          <Text className="p-4 text-xl font-bold text-white">IIoT System</Text>
        </AppShell.Section>
        <AppShell.Section grow className="p-md">
          <NavLink
            component={Link}
            to="/dashboard"
            label="Main Dashboard"
            description="Overview of devices and alerts"
            className="text-white hover:bg-white/10"
          />
          <NavLink
            component={Link}
            to="/analytics"
            label="Analytics & Reports"
            description="Telemetry trends and production summaries"
            className="text-white hover:bg-white/10"
          />
        </AppShell.Section>
      </AppShell.Navbar>
      <AppShell.Main>
        <div className="flex items-center justify-center min-h-screen p-4">
          <Container size="xl" className="w-full">
            {children}
          </Container>
        </div>
      </AppShell.Main>
    </AppShell>
  );
}
