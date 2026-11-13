import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useQueryClient } from "@tanstack/react-query";
import { errorMessage } from "../api/client";
import {
  getGetCurrentUserQueryKey,
  useGetCurrentUser,
  useLoginUser,
  useResetUserPassword,
  useUpdateStudent,
} from "../api/endpoints";
import { useAuthContext } from "../contexts/use-auth-context";
import { CurrentUserRole, UpdateStudentBodyYear } from "../api/models";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  NativeSelect,
  NativeSelectOption,
} from "@/components/ui/native-select";

// Reached only via the auth-context redirect gate once a just-invited user
// logs in with their temp password.
export default function ClaimAccount() {
  const { data: currentUser } = useGetCurrentUser();
  const { studentID } = useAuthContext();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const [step, setStep] = useState<"password" | "profile">("password");
  const [newPassword, setNewPassword] = useState("");
  const [passwordError, setPasswordError] = useState("");

  const { mutateAsync: resetPassword, isPending: isResettingPassword } =
    useResetUserPassword();
  const { mutateAsync: loginAgain } = useLoginUser();
  const {
    mutateAsync: updateStudent,
    isPending: isSavingProfile,
    error: profileError,
  } = useUpdateStudent();

  const isStudent = currentUser?.role === CurrentUserRole.student;

  const handlePasswordSubmit = async (
    event: React.SubmitEvent<HTMLElement>,
  ) => {
    event.preventDefault();
    setPasswordError("");
    if (!currentUser?.email) return;
    try {
      await resetPassword({ data: { new_password: newPassword } });
      // Resetting the password revokes the current Supabase session, so the
      // cookies from the original (temp-password) login no longer work.
      // Re-authenticate with the new password to get a fresh session before
      // continuing - the user never sees a login prompt for this.
      await loginAgain({
        data: { email: currentUser.email, password: newPassword },
      });
      // The AuthProvider's own currentUser query must see reset_time=null
      // before we navigate, or its redirect gate bounces us right back here.
      await queryClient.invalidateQueries({
        queryKey: getGetCurrentUserQueryKey(),
      });
      if (isStudent) {
        setStep("profile");
      } else {
        void navigate("/counselor", { replace: true });
      }
    } catch (err) {
      console.error(err);
      setPasswordError("Something went wrong setting your password.");
    }
  };

  const [year, setYear] = useState<UpdateStudentBodyYear>(
    UpdateStudentBodyYear.freshman,
  );
  const [highSchool, setHighSchool] = useState("");
  const [gpa, setGpa] = useState("");
  const [testScore, setTestScore] = useState("");
  const [major, setMajor] = useState("");
  const [minor, setMinor] = useState("");
  const [geoPrefs, setGeoPrefs] = useState("");

  const handleProfileSubmit = async (
    event: React.SubmitEvent<HTMLElement>,
  ) => {
    event.preventDefault();
    if (!studentID) return;
    try {
      await updateStudent({
        id: studentID,
        data: {
          year,
          high_school: highSchool,
          gpa: Number(gpa),
          standardized_test_score: testScore || undefined,
          intended_major: major || undefined,
          intended_minor: minor || undefined,
          geographic_preferences: geoPrefs || undefined,
        },
      });
      void navigate("/student", { replace: true });
    } catch (err) {
      console.error(err);
    }
  };

  if (step === "password") {
    return (
      <div className="mx-auto w-full max-w-md p-8">
        <Card>
          <CardHeader>
            <CardTitle className="text-2xl">Create Account</CardTitle>
            <CardDescription>Enter in your information.</CardDescription>
          </CardHeader>
          <form onSubmit={handlePasswordSubmit}>
            <CardContent className="flex flex-col gap-4">
              <div className="flex flex-col gap-2">
                <Label>Name</Label>
                <Input value={currentUser?.name ?? ""} disabled readOnly />
              </div>
              <div className="flex flex-col gap-2">
                <Label>Email</Label>
                <Input value={currentUser?.email ?? ""} disabled readOnly />
              </div>
              <div className="flex flex-col gap-2">
                <Label htmlFor="new-password">Password</Label>
                <Input
                  id="new-password"
                  type="password"
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  required
                />
              </div>
              <Button type="submit" disabled={isResettingPassword}>
                {isResettingPassword ? "Saving..." : "Continue"}
              </Button>
              {passwordError && (
                <p
                  role="alert"
                  className="rounded-md bg-destructive/10 p-4 text-destructive"
                >
                  {passwordError}
                </p>
              )}
            </CardContent>
          </form>
        </Card>
      </div>
    );
  }

  return (
    <div className="mx-auto w-full max-w-md p-8">
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl">Profile Details</CardTitle>
          <CardDescription>
            Add your details, you can always go back to edit via your profile
            page
          </CardDescription>
        </CardHeader>
        <form onSubmit={handleProfileSubmit}>
          <CardContent className="flex flex-col gap-4">
            <div className="flex flex-col gap-2">
              <Label htmlFor="year">Graduation Year / Cohort</Label>
              <NativeSelect
                id="year"
                className="w-full"
                value={year}
                onChange={(e) =>
                  setYear(e.target.value as UpdateStudentBodyYear)
                }
              >
                <NativeSelectOption value={UpdateStudentBodyYear.freshman}>
                  Freshman
                </NativeSelectOption>
                <NativeSelectOption value={UpdateStudentBodyYear.sophomore}>
                  Sophomore
                </NativeSelectOption>
                <NativeSelectOption value={UpdateStudentBodyYear.junior}>
                  Junior
                </NativeSelectOption>
                <NativeSelectOption value={UpdateStudentBodyYear.senior}>
                  Senior
                </NativeSelectOption>
              </NativeSelect>
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="high-school">High School</Label>
              <Input
                id="high-school"
                value={highSchool}
                onChange={(e) => setHighSchool(e.target.value)}
                required
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
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
              <div className="flex flex-col gap-2">
                <Label htmlFor="test-score">SAT/ACT</Label>
                <Input
                  id="test-score"
                  value={testScore}
                  onChange={(e) => setTestScore(e.target.value)}
                />
              </div>
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="major" className="justify-between">
                Intended Major
                <span className="font-normal text-muted-foreground">
                  Optional
                </span>
              </Label>
              <Input
                id="major"
                value={major}
                onChange={(e) => setMajor(e.target.value)}
              />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="minor" className="justify-between">
                Intended Minor
                <span className="font-normal text-muted-foreground">
                  Optional
                </span>
              </Label>
              <Input
                id="minor"
                value={minor}
                onChange={(e) => setMinor(e.target.value)}
              />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="geo" className="justify-between">
                Geographic Preferences
                <span className="font-normal text-muted-foreground">
                  Optional
                </span>
              </Label>
              <Input
                id="geo"
                value={geoPrefs}
                onChange={(e) => setGeoPrefs(e.target.value)}
              />
            </div>
            <Button type="submit" disabled={isSavingProfile}>
              {isSavingProfile ? "Saving..." : "Proceed to Dashboard"}
            </Button>
            {profileError && (
              <p
                role="alert"
                className="rounded-md bg-destructive/10 p-4 text-destructive"
              >
                {errorMessage(profileError)}
              </p>
            )}
          </CardContent>
        </form>
      </Card>
    </div>
  );
}
