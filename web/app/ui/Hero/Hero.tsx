import Search from "@/app/ui/Search/Search";

export default function Hero() {
  return (
    <div className="hero size-full">
      <div className="hero-content text-center text-base-content size-full">
        <div>
          <h1 className="text-4xl font-bold">New Eden Rivalries Database</h1>
          <p className="py-6">{/* {put something nice here} */}</p>
          <Search />
        </div>
      </div>
    </div>
  );
}
