export interface PointSnapshot {
  x: number;
  y_true: number;
  y_pred: number;
  point_loss: number;
  point_grad: number;
}

export interface UpdateDetails {
  w_old: number;
  lr: number;
  grad_w: number;
  delta_w: number;
  w_new: number;
}

export interface Snapshot {
  step: number;
  w: number;
  grad_w: number;
  loss: number;
  point_details: PointSnapshot[];
  update_components: UpdateDetails;
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
