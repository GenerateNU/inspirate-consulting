import { useState } from "react";
import type { ApiError } from "../api/client";
import { useCreateUser } from "../api/endpoints";
import { CreateUserInputBodyRole } from "../api/models";
import { useAuthContext } from "../contexts/use-auth-context";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  NativeSelect,
  NativeSelectOption,
} from "@/components/ui/native-select";
import { graduationYearOptions } from "@/lib/graduation-year";
import { ErrorMessage } from "@/components/ui/error-message";

export default function CreateUser() {
  const { counselorID } = useAuthContext();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [role, setRole] = useState<CreateUserInputBodyRole>(
    CreateUserInputBodyRole.student,
  );
  const [graduationYear, setGraduationYear] = useState(
    graduationYearOptions()[3].value,
  );
  const [gpa, setGpa] = useState("");
  const [status, setStatus] = useState("");
  // set when the account was created but the sign-in email couldn't be sent
  const [tempPassword, setTempPassword] = useState<string | null>(null);

  const { mutate, error, isPending, reset } = useCreateUser();

  // drop the last result's message and error
  // the temp password stays until the next submit, since it's the only copy
  const clearMessages = () => {
    if (!tempPassword) setStatus("");
    reset();
  };

  // return readable error message
  const backendMessage = (error as ApiError | null)?.info as
    { message?: unknown } | undefined;
  const errorText =
    typeof backendMessage?.message === "string"
      ? backendMessage.message
      : undefined;

  const isStudent = role === CreateUserInputBodyRole.student;

  const handleSubmit = (event: React.SubmitEvent<HTMLElement>) => {
    event.preventDefault();
    setStatus("");
    setTempPassword(null);
    mutate(
      {
        data: {
          name,
          email,
          role,
          ...(isStudent
            ? {
                graduation_year: graduationYear,
                gpa: Number(gpa),
                counselor_id: counselorID ?? undefined,
              }
            : {}),
        },
      },
      {
        onSuccess: (data) => {
          if (data.email_sent) {
            setStatus(
              `Created user ${name}. An email with sign-in instructions has been sent to ${email}.`,
            );
            return;
          }
          setStatus(
            `Created user ${name}, but the sign-in email to ${email} couldn't be sent. Give them this temporary password so they can sign in:`,
          );
          setTempPassword(data.temp_password ?? null);
        },
      },
    );
  };

  return (
    <div className="mx-auto w-full max-w-md px-4 py-6">
      <Card className="gap-4">
        <CardHeader>
          <CardTitle className="text-2xl">Invite a New User</CardTitle>
        </CardHeader>
        {/* change events bubble, so this runs when any field is edited */}
        <form onSubmit={handleSubmit} onChange={clearMessages}>
          <CardContent className="flex flex-col gap-3">
            <div className="flex flex-col gap-2">
              <Label htmlFor="name">Name</Label>
              <Input
                id="name"
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
              />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="role">Role</Label>
              <NativeSelect
                id="role"
                className="w-full"
                value={role}
                onChange={(e) =>
                  setRole(e.target.value as CreateUserInputBodyRole)
                }
              >
                <NativeSelectOption value={CreateUserInputBodyRole.student}>
                  Student
                </NativeSelectOption>
                <NativeSelectOption value={CreateUserInputBodyRole.counselor}>
                  Counselor
                </NativeSelectOption>
              </NativeSelect>
            </div>
            {isStudent && (
              <>
                <div className="flex flex-col gap-2">
                  <Label htmlFor="graduation-year">Graduation Year</Label>
                  <NativeSelect
                    id="graduation-year"
                    className="w-full"
                    value={graduationYear}
                    onChange={(e) => setGraduationYear(Number(e.target.value))}
                  >
                    {graduationYearOptions().map(({ value, label }) => (
                      <NativeSelectOption key={value} value={value}>
                        {label}
                      </NativeSelectOption>
                    ))}
                  </NativeSelect>
                </div>
                <div className="flex flex-col gap-2">
                  <Label htmlFor="gpa">GPA</Label>
                  <Input
                    id="gpa"
                    type="number"
                    min={0}
                    max={5}
                    step="0.01"
                    value={gpa}
                    onChange={(e) => setGpa(e.target.value)}
                    required
                  />
                </div>
              </>
            )}
          </CardContent>
          <CardFooter className="mt-4 flex flex-col items-stretch gap-2">
            <Button type="submit" disabled={isPending}>
              {isPending ? "Creating User..." : "Create User"}
            </Button>
            {status && !tempPassword && (
              <p className="text-sm text-muted-foreground">{status}</p>
            )}
            {tempPassword && (
              <ErrorMessage
                variant="banner"
                title="Email not sent"
                message={
                  <>
                    {status}{" "}
                    <code className="rounded bg-white px-1 font-mono select-all">
                      {tempPassword}
                    </code>
                  </>
                }
              />
            )}
            <ErrorMessage variant="banner" error={error} message={errorText} />
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}
