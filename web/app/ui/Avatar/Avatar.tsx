export default function Avatar({
  alt,
  size,
  src,
}: {
  alt: string;
  size: number;
  src: string;
}) {
  return (
    <div className="avatar">
      <div className="max-w-24 max-h-24">
        <img src={src} alt={alt} />
      </div>
    </div>
  );
}
