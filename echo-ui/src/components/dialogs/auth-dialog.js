"use client";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useState } from "react";
import { Button } from "@/components/ui/button";

export default function AuthDialog({ action, children }) {
  const [open, setOpen] = useState(false);
  const [mode, setMode] = useState(action);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>
            {mode === "login" ? "Login to Echo-UI" : "Create an Account"}
          </DialogTitle>
        </DialogHeader>
        <form className="space-y-4">
          {mode === "signup" && (
            <>
              <div>
                <Label htmlFor="firstName">Name</Label>
                <Input id="firstName" placeholder="Your first name" required />
              </div>
              <div>
                <Label htmlFor="lastName">Name</Label>
                <Input id="lastName" placeholder="Your last name" required />
              </div>
            </>
          )}
          <div>
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              placeholder="you@example.com"
              required
            />
          </div>
          <div>
            <Label htmlFor="password">Password</Label>
            <Input
              id="password"
              type="password"
              placeholder="••••••••"
              required
            />
          </div>
          <Button type="submit" className="w-full">
            {mode === "login" ? "Login" : "Sign Up"}
          </Button>
        </form>
        <p className="text-sm text-center text-gray-500 mt-2">
          {mode === "login"
            ? "Don't have an account?"
            : "Already have an account?"}{" "}
          <span
            onClick={() => setMode(mode === "login" ? "signup" : "login")}
            className="text-primary cursor-pointer hover:underline"
          >
            {mode === "login" ? "Sign up" : "Login"}
          </span>
        </p>
      </DialogContent>
    </Dialog>
  );
}
