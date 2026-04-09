import type { Kid, Behavior, Prize, Record, DrawResult } from '../types';

export const API_BASE = (import.meta.env.VITE_API_BASE as string) || '';

// Kids
export async function getKids(): Promise<Kid[]> {
  const res = await fetch(`${API_BASE}/kids`);
  if (!res.ok) throw new Error('Failed to fetch kids');
  return res.json();
}

export async function createKid(data: { name: string; avatar: string; points?: number }): Promise<Kid> {
  const res = await fetch(`${API_BASE}/kids`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to create kid');
  return res.json();
}

export async function uploadAvatar(file: File): Promise<string> {
  const formData = new FormData();
  formData.append('avatar', file);
  const res = await fetch(`${API_BASE}/kids/upload-avatar`, {
    method: 'POST',
    body: formData,
  });
  if (!res.ok) throw new Error('Failed to upload avatar');
  const data = await res.json();
  return data.avatar;
}

export async function updateKid(id: number, data: { name?: string; avatar?: string }): Promise<Kid> {
  const res = await fetch(`${API_BASE}/kids/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to update kid');
  return res.json();
}

export async function deleteKid(id: number): Promise<void> {
  const res = await fetch(`${API_BASE}/kids/${id}`, { method: 'DELETE' });
  if (!res.ok) throw new Error('Failed to delete kid');
}

export async function getPoints(idOrName: number | string): Promise<number> {
  const res = await fetch(`${API_BASE}/kids/${idOrName}/points`);
  if (!res.ok) throw new Error('Failed to fetch points');
  return res.json();
}

export async function addPoints(idOrName: number | string, behaviorId: number): Promise<Kid> {
  const res = await fetch(`${API_BASE}/kids/points/add`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ id: idOrName, behavior_id: behaviorId }),
  });
  if (!res.ok) throw new Error('Failed to add points');
  return res.json();
}

export async function subPoints(idOrName: number | string, behaviorId: number): Promise<Kid> {
  const res = await fetch(`${API_BASE}/kids/points/sub`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ id: idOrName, behavior_id: behaviorId }),
  });
  if (!res.ok) throw new Error('Failed to subtract points');
  return res.json();
}

// Behaviors
export async function getBehaviors(): Promise<Behavior[]> {
  const res = await fetch(`${API_BASE}/behaviors`);
  if (!res.ok) throw new Error('Failed to fetch behaviors');
  return res.json();
}

export async function createBehavior(data: Omit<Behavior, 'id'>): Promise<Behavior> {
  const res = await fetch(`${API_BASE}/behaviors`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to create behavior');
  return res.json();
}

export async function updateBehavior(id: number, data: Partial<Omit<Behavior, 'id'>>): Promise<Behavior> {
  const res = await fetch(`${API_BASE}/behaviors/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to update behavior');
  return res.json();
}

export async function deleteBehavior(id: number): Promise<void> {
  const res = await fetch(`${API_BASE}/behaviors/${id}`, { method: 'DELETE' });
  if (!res.ok) throw new Error('Failed to delete behavior');
}

// Prizes
export async function getPrizes(): Promise<Prize[]> {
  const res = await fetch(`${API_BASE}/prizes`);
  if (!res.ok) throw new Error('Failed to fetch prizes');
  return res.json();
}

export async function createPrize(data: Omit<Prize, 'id'>): Promise<Prize> {
  const res = await fetch(`${API_BASE}/prizes`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to create prize');
  return res.json();
}

export async function updatePrize(id: number, data: Partial<Omit<Prize, 'id'>>): Promise<Prize> {
  const res = await fetch(`${API_BASE}/prizes/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to update prize');
  return res.json();
}

export async function deletePrize(id: number): Promise<void> {
  const res = await fetch(`${API_BASE}/prizes/${id}`, { method: 'DELETE' });
  if (!res.ok) throw new Error('Failed to delete prize');
}

// Draw
export async function draw(kidId: number): Promise<DrawResult> {
  const res = await fetch(`${API_BASE}/draw`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ kid_id: kidId }),
  });
  if (!res.ok) throw new Error('Failed to draw');
  return res.json();
}

// Records
export async function getRecords(kidId?: number): Promise<Record[]> {
  const url = kidId ? `${API_BASE}/records?kid_id=${kidId}` : `${API_BASE}/records`;
  const res = await fetch(url);
  if (!res.ok) throw new Error('Failed to fetch records');
  return res.json();
}

export async function reverseRecord(id: number, reason: string): Promise<Record> {
  const res = await fetch(`${API_BASE}/records/${id}/reverse`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ reason }),
  });
  if (!res.ok) throw new Error('Failed to reverse record');
  return res.json();
}
