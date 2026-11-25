import { useState } from "react";
import { useAuthContext } from "../contexts/use-auth-context";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { ErrorMessage } from "@/components/ui/error-message";

type Audience = "student" | "counselor";

interface LogInProps {
  audience?: Audience;
}

const emailLabel: Record<"default" | Audience, string> = {
  default: "Email / Username",
  student: "Student Email / Username",
  counselor: "Counselor Email / Username",
};

export default function LogIn({ audience }: LogInProps) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<unknown>(null);

  const { login } = useAuthContext();

  const handleSubmit = async (event: React.SubmitEvent<HTMLElement>) => {
    event.preventDefault();
    setError(null);
    try {
      await login(email, password);
    } catch (err) {
      setError(err);
    }
  };

  return (
    <div className="mx-auto w-full max-w-md p-8">
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl">Log In</CardTitle>
        </CardHeader>
        <form onSubmit={handleSubmit}>
          <CardContent className="flex flex-col gap-4">
            <div className="flex flex-col gap-2">
              <Label htmlFor="email">{emailLabel[audience ?? "default"]}</Label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
            <Button type="submit">Log In</Button>
            <ErrorMessage
              variant="banner"
              title="Couldn't log in"
              error={error}
            />
            <Separator />
            <p className="text-sm text-muted-foreground">
              Don&apos;t have an account yet? Contact an administrator.
            </p>
          </CardContent>
        </form>
      </Card>
    </div>
  );
}
