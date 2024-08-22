"use client"; // Error boundaries must be Client Components

import * as Sentry from "@sentry/nextjs";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";

export default function Error({
  error,
}: {
  error: Error & { digest?: string };
}) {
  useEffect(() => {
    Sentry.captureException(error);
  }, [error]);

  return (
    <div className="flex flex-col justify-center items-center gap-5">
      <h1 className="text-2xl">Something went wrong!</h1>
      <p className="text-left">
        We've collected the error. It might even get fixed soon.
      </p>

      <img
        src="/oopsie.png"
        width={150}
        alt="error image"
        className="absolute bottom-0 right-10"
      />

      <Link href="/">
        <button className="btn btn-primary">Back to the home page</button>
      </Link>
    </div>
  );
}
