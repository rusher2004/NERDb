import NavBar from "@/app/ui/Nav/NavBar";
import Hero from "@/app/ui/Hero/Hero";

export default function Home() {
  return (
    <main className="flex min-h-dvh flex-col items-center bg-base-100">
      <NavBar />
      <Hero />
    </main>
  );
}
