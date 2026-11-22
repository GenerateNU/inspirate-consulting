const GRADE_NAMES: Record<number, string> = {
  9: "Freshman",
  10: "Sophomore",
  11: "Junior",
  12: "Senior",
};

export function currentClassYear(today = new Date()) {
  return today.getMonth() >= 6 ? today.getFullYear() + 1 : today.getFullYear();
}

export function gradeForGraduationYear(graduationYear: number) {
  const grade = 12 - (graduationYear - currentClassYear());
  return grade >= 9 && grade <= 12 ? grade : null;
}

export function graduationYearLabel(graduationYear: number) {
  const grade = gradeForGraduationYear(graduationYear);
  return grade === null
    ? `Class of ${graduationYear}`
    : `${GRADE_NAMES[grade]} (${graduationYear})`;
}

export function graduationYearOptions() {
  const seniors = currentClassYear();
  return [0, 1, 2, 3].map((offset) => {
    const value = seniors + offset;
    return { value, label: graduationYearLabel(value) };
  });
}
