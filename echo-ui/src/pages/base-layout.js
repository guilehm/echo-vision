"use client";

import { Button } from "@/components/ui/button";
import { Menu, X } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useState } from "react";

export default function BaseLayout({ children }) {
  const [isOpen, setIsOpen] = useState(false);
  return (
    <div className="flex flex-col min-h-screen">
      {/* App Bar */}
      <header className="w-full p-4 border-b bg-white shadow-md">
        <div className="container mx-auto flex justify-between items-center">
          <h1 className="text-xl font-bold">Echo-UI</h1>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex space-x-4">
            <Button variant="ghost" className="cursor-pointer">
              Home
            </Button>
          </nav>

          {/* Mobile Navigation */}
          <div className="md:hidden">
            <DropdownMenu open={isOpen} onOpenChange={setIsOpen}>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon">
                  {isOpen ? (
                    <X className="w-6 h-6" />
                  ) : (
                    <Menu className="w-6 h-6" />
                  )}
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem onSelect={() => setIsOpen(false)}>
                  Home
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-grow container mx-auto p-6">{children}</main>

      {/* Footer */}
      <footer className="w-full p-4 border-t bg-white shadow-inner">
        <div className="container mx-auto text-center text-sm text-gray-500">
          © {new Date().getFullYear()} Echo-Vision. All rights reserved.
        </div>
      </footer>
    </div>
  );
}
