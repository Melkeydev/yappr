export function calculateTimeRemaining(expiresAt: string): string {
  const now = new Date();
  const expirationDate = new Date(expiresAt);
  
  const diffMs = expirationDate.getTime() - now.getTime();
  
  // If expired or will expire in less than a minute, show as expired
  if (diffMs <= 0) {
    return "Expired";
  }
  
  // Calculate total minutes remaining
  const totalMinutes = Math.floor(diffMs / (1000 * 60));
  
  if (totalMinutes < 60) {
    // Less than an hour, show just minutes
    return `${totalMinutes}m`;
  }
  
  // Calculate hours and remaining minutes
  const hours = Math.floor(totalMinutes / 60);
  const minutes = totalMinutes % 60;
  
  // Format as HH:MM
  return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`;
}

export function isRoomExpired(expiresAt: string): boolean {
  const now = new Date();
  const expirationDate = new Date(expiresAt);
  return expirationDate.getTime() <= now.getTime();
}