import Search from "@/app/ui/Search/Search";

export default function Hero() {
  return (
    <div className="hero h-full">
      <div className="hero-content text-center text-base-content">
        <div className="max-w-md">
          <h1 className="text-5xl font-bold">Hello there</h1>
          <p className="py-6">
            Provident cupiditate voluptatem et in. Quaerat fugiat ut assumenda
            excepturi exercitationem quasi. In deleniti eaque aut repudiandae et
            a id nisi.
          </p>
          <Search />
        </div>
      </div>
    </div>
  );
}
