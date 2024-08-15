import Image from "next/image";

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
    <Image
      className="size-8 md:size-11"
      src={src}
      alt={alt}
      width={size}
      height={size}
    />
  );
}
