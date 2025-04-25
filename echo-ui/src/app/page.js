import Footer from "@/components/footer";
import { Button } from "@/components/ui/button";
import Link from "next/link";

export default function Home() {
  return (
    <div className="flex flex-col min-h-screen">
      <header className="w-full">
        <div className="container mx-auto px-4 lg:px-6 h-14 flex items-center max-w-7xl">
          <Link className="flex items-center justify-center" href="#">
            {/* <Eye className="h-6 w-6" /> */}
            <h1 className="ml-2 text-2xl font-bold">Echo Vision</h1>
          </Link>
          <nav className="ml-auto flex gap-4 sm:gap-6">
            <Link
              className="text-sm font-medium hover:underline underline-offset-4"
              href="/sign-in"
            >
              Sign in
            </Link>
            <Link
              className="text-sm font-medium hover:underline underline-offset-4"
              href="/sign-up"
            >
              Sign up
            </Link>
          </nav>
        </div>
      </header>
      <main className="flex-1">
        <section className="w-full py-12 md:py-24 lg:py-32 xl:py-48">
          <div className="container mx-auto px-4 md:px-6 max-w-6xl">
            <div className="flex flex-col items-center space-y-4 text-center">
              <div className="space-y-2">
                <h2 className="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl lg:text-6xl/none">
                  Unlock the Power of Visual AI
                </h2>
                <p className="mx-auto max-w-[700px] text-gray-500 md:text-xl dark:text-gray-400">
                  Echo Vision analyzes images and videos with cutting-edge AI,
                  providing instant insights and powerful visual understanding.
                </p>
              </div>
              <div className="space-x-4">
                <Button asChild>
                  <Link href="/dashboard">Get Started for Free</Link>
                </Button>
                <Button variant="outline" asChild>
                  <Link href="#">Learn More</Link>
                </Button>
              </div>
            </div>
          </div>
        </section>

        {/* Call to Action */}
        <section className="text-center space-y-4 bg-gray-100 py-10 md:py-20">
          <div className="container mx-auto px-4 md:px-6 max-w-6xl">
            <div className="flex flex-col items-center space-y-4 text-center">
              <h3 className="text-xl sm:text-2xl md:text-3xl lg:text-4xl xl:text-5xl font-semibold">
                Join the Echo Vision Experience
              </h3>
              <h2 className="text-gray-600">
                Start using Echo Vision today and transform your workflow.
              </h2>
              <Button size="lg" variant="outline">
                Learn More
              </Button>
            </div>
          </div>
        </section>
      </main>

      <Footer />
    </div>
  );
}
