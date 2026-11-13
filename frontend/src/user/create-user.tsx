import { useState } from "react";
import { errorMessage } from "../api/client";
import { useCreateUser } from "../api/endpoints";
import { CreateUserInputBodyRole, CreateUserInputBodyYear } from "../api/models";
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

export default function CreateUser() {
  const { counselorID } = useAuthContext();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [role, setRole] = useState<CreateUserInputBodyRole>(
    CreateUserInputBodyRole.student,
  );
  const [year, setYear] = useState<CreateUserInputBodyYear>(
    CreateUserInputBodyYear.freshman,
  );
  const [gpa, setGpa] = useState("");
  const [status, setStatus] = useState("");

  const { mutateAsync, error, isPending } = useCreateUser();

  const isStudent = role === CreateUserInputBodyRole.student;

  const handleSubmit = (event: React.SubmitEvent<HTMLElement>) => {
    event.preventDefault();
    setStatus("");
    void mutateAsync(
      {
        data: {
          name,
          email,
          role,
          pfp_key: null,
          ...(isStudent
            ? {
                year,
                gpa: Number(gpa),
                counselor_id: counselorID ?? undefined,
              }
            : {}),
        },
      },
      {
        onSuccess: () =>
          setStatus(
            `Created user ${name}. An email with sign-in instructions has been sent to ${email}.`,
          ),
      },
    );
  };

  return (
    <div className="mx-auto w-full max-w-md px-4 py-6">
      <Card className="gap-4">
        <CardHeader>
          <CardTitle className="text-2xl">Invite a New User</CardTitle>
        </CardHeader>
        <form onSubmit={handleSubmit}>
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
                  <Label htmlFor="year">Graduation Year / Cohort</Label>
                  <NativeSelect
                    id="year"
                    className="w-full"
                    value={year}
                    onChange={(e) =>
                      setYear(e.target.value as CreateUserInputBodyYear)
                    }
                  >
                    <NativeSelectOption value={CreateUserInputBodyYear.freshman}>
                      Freshman
                    </NativeSelectOption>
                    <NativeSelectOption value={CreateUserInputBodyYear.sophomore}>
                      Sophomore
                    </NativeSelectOption>
                    <NativeSelectOption value={CreateUserInputBodyYear.junior}>
                      Junior
                    </NativeSelectOption>
                    <NativeSelectOption value={CreateUserInputBodyYear.senior}>
                      Senior
                    </NativeSelectOption>
                  </NativeSelect>
                </div>
                <div className="flex flex-col gap-2">
                  <Label htmlFor="gpa">GPA</Label>
                  <Input
                    id="gpa"
                    type="number"
                    min={0}
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
            {status && <p className="text-sm text-muted-foreground">{status}</p>}
            {error && (
              <p
                role="alert"
                className="rounded-md bg-destructive/10 p-4 text-destructive"
              >
                {errorMessage(error)}
              </p>
            )}
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}
