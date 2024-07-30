export default function Home() {
  return (
    <main
      data-theme="amarr"
      className="flex min-h-screen flex-col items-center gap-5 p-24"
    >
      <input
        type="text"
        placeholder="Type here"
        className="input input-bordered w-full max-w-xs"
      />
      <div className="flex flex-col items-center gap-2">
        <div className="flex items-center justify-between gap-2">
          <button className="btn btn-primary">Submit</button>
          <button className="btn btn-secondary">Cancel</button>
          <button className="btn btn-accent">Reset</button>
          <button className="btn btn-warning">Delete</button>
        </div>
        <div
          data-theme="caldari"
          className="flex items-center justify-between gap-2"
        >
          <button className="btn btn-primary">Submit</button>
          <button className="btn btn-secondary">Cancel</button>
          <button className="btn btn-accent">Reset</button>
          <button className="btn btn-warning">Delete</button>
        </div>
        <div
          data-theme="gallente"
          className="flex items-center justify-between gap-2"
        >
          <button className="btn btn-primary">Submit</button>
          <button className="btn btn-secondary">Cancel</button>
          <button className="btn btn-accent">Reset</button>
          <button className="btn btn-warning">Delete</button>
        </div>
        <div
          data-theme="minmatar"
          className="flex items-center justify-between gap-2"
        >
          <button className="btn btn-primary">Submit</button>
          <button className="btn btn-secondary">Cancel</button>
          <button className="btn btn-accent">Reset</button>
          <button className="btn btn-warning">Delete</button>
        </div>
        <div
          data-theme="ore"
          className="flex items-center justify-between gap-2"
        >
          <button className="btn btn-primary">Submit</button>
          <button className="btn btn-secondary">Cancel</button>
          <button className="btn btn-accent">Reset</button>
          <button className="btn btn-warning">Delete</button>
        </div>
        <div
          data-theme="sisters"
          className="flex items-center justify-between gap-2"
        >
          <button className="btn btn-primary">Submit</button>
          <button className="btn btn-secondary">Cancel</button>
          <button className="btn btn-accent">Reset</button>
          <button className="btn btn-warning">Delete</button>
        </div>
      </div>
    </main>
  );
}
