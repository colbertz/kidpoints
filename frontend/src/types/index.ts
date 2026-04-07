export interface Kid {
  id: number;
  name: string;
  avatar: string;
  points: number;
}

export interface Behavior {
  id: number;
  name: string;
  type: 'add' | 'sub';
  points: number;
}

export interface Prize {
  id: number;
  name: string;
  icon: string;
  probability: number;
}

export interface Record {
  id: number;
  kid_id: number;
  behavior_id: number;
  behavior_name: string;
  points: number;
  created_at: string;
  reversed: boolean;
  reversed_reason?: string;
  reversed_at?: string;
}

export interface DrawResult {
  prize: Prize;
  kid: Kid;
}
