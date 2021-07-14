export function formatDate(timestamp: number) {
  const today = new Date();
  const date = new Date(timestamp);

  if (
    // is same day?
    date.getFullYear() === today.getFullYear() &&
    date.getMonth() === today.getMonth() &&
    date.getDay() === today.getDay()
  ) {
    return date.toLocaleTimeString();
  } else {
    return date.toLocaleString();
  }
}
