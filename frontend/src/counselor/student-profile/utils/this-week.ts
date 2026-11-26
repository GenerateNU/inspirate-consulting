// Sunday 00:00 through the next Sunday 00:00, in local time
export function thisWeekRange(today = new Date()) {
  const start = new Date(today);
  start.setHours(0, 0, 0, 0);
  start.setDate(start.getDate() - start.getDay());
  const end = new Date(start);
  end.setDate(end.getDate() + 7);
  return { start, end };
}

export function isThisWeek(value: string, today = new Date()) {
  const { start, end } = thisWeekRange(today);
  const time = new Date(value).getTime();
  return time >= start.getTime() && time < end.getTime();
}

// due before the week ends, including overdue
export function isDueByEndOfWeek(value: string, today = new Date()) {
  return new Date(value).getTime() < thisWeekRange(today).end.getTime();
}
