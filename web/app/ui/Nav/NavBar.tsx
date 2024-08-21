import Link from "next/link";
import Image from "next/image";
import { useEffect, useState } from "react";
import ThemeSelector from "@/app/ui/Theme/ThemeSelector";

export default function NavBar() {
  return (
    <nav className="navbar bg-base-200">
      <div className="flex-1">
        <Link href="/" className="btn btn-ghost text-4xl text-primary">
          NERDb
        </Link>
      </div>
      <div className="flex-none gap-2">
        {/* {pathname !== "/" && (
          // <div className="form-control">
          //   <input
          //     type="text"
          //     placeholder="Search"
          //     className="input input-bordered w-24 md:w-auto"
          //   />
          // </div>
          // <Search />
        )} */}
        <ThemeSelector />
      </div>
    </nav>
  );
}
