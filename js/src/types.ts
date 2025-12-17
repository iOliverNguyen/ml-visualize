export interface Snapshot {
  step: number;
  w: number;
  grad_w: number;
  loss: number;
}

export type LayoutMode = '1col' | '2col' | '3col' | 'auto';

export interface MetricsData {
  initial: {
    w: number;
    grad_w: number;
    loss: number;
  };
  current: {
    w: number;
    grad_w: number;
    loss: number;
  };
  final: {
    w: number;
    grad_w: number;
    loss: number;
  };
  deltaFromStart: {
    w: number;
    loss: number;
    progress: number;
  };
  deltaToTarget: {
    w: number;
    loss: number;
  };
}
