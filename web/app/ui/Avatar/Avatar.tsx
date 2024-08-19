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
      <div className="max-w-12 max-h-12">
        <img src={src} alt={alt} />
      </div>
    </div>
  );
}
