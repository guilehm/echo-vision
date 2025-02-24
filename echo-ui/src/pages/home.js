import { Button } from "@/components/ui/button";
import { Rocket, Eye, ShieldCheck } from "lucide-react";
import FeatureCard from "@/components/cards/feature-card";
import AuthDialog from "@/components/dialogs/auth-dialog";

export default function Home() {
  return (
    <div className="container mx-auto px-4 py-12 space-y-16">
      {/* Hero Section */}
      <section className="text-center space-y-4">
        <h1 className="text-4xl font-bold tracking-tight sm:text-5xl">
          Welcome to Echo-UI
        </h1>
        <p className="text-lg text-gray-600">
          The next-generation interface for seamless user interaction.
        </p>

        {/* Authentication  */}
        <div className="mt-4 space-x-4">
          <AuthDialog action="signup">
            <Button size="lg">Sign Up</Button>
          </AuthDialog>
          <AuthDialog action="login">
            <Button size="lg" variant="outline">
              Login
            </Button>
          </AuthDialog>
        </div>
      </section>

      {/* Features Section */}
      <section className="grid grid-cols-1 sm:grid-cols-3 gap-6">
        <FeatureCard
          icon={<Rocket className="w-10 h-10 text-primary" />}
          title="High Performance"
          description="Built with efficiency and speed in mind."
        />
        <FeatureCard
          icon={<Eye className="w-10 h-10 text-primary" />}
          title="Intuitive UI"
          description="Designed for a seamless and engaging user experience."
        />
        <FeatureCard
          icon={<ShieldCheck className="w-10 h-10 text-primary" />}
          title="Secure & Reliable"
          description="Advanced security to keep your data safe."
        />
      </section>

      {/* Call to Action */}
      <section className="text-center space-y-4">
        <h2 className="text-3xl font-semibold">
          Join the Echo-Vision Experience
        </h2>
        <p className="text-gray-600">
          Start using Echo-UI today and transform your workflow.
        </p>
        <Button size="lg" variant="outline">
          Learn More
        </Button>
      </section>
    </div>
  );
}
