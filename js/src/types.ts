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

// Educational content types
export interface GlossaryEntry {
  term: string;
  brief: string;
  detailed: string;
  comprehensive?: string;
  level: string[];
  formula?: string;
  relatedTerms: string[];
  example?: string;
}

export interface Glossary {
  [key: string]: GlossaryEntry;
}

export interface FAQItem {
  id: string;
  question: string;
  answer: string;
  level: string;
  tags: string[];
}

export interface FAQCategory {
  name: string;
  questions: FAQItem[];
}

export interface FAQData {
  categories: FAQCategory[];
}

// Tutorial content types
export interface TutorialSection {
  type: 'text' | 'highlight' | 'code' | 'callout' | 'visual-reference';
  content: string;
  code?: {
    language: string;
    snippet: string;
  };
  calloutType?: 'info' | 'tip' | 'warning';
}

export interface TutorialChapter {
  id: string;
  number: number;
  title: string;
  sections: TutorialSection[];
  visualCues?: {
    highlightElements?: string[];
    focusMetrics?: string[];
    suggestedStep?: number;
  };
  estimatedReadTime?: number;
}

export interface TutorialContent {
  meta: {
    title: string;
    description: string;
    estimatedReadTime: number;
    version: string;
  };
  chapters: TutorialChapter[];
}

// Case Library types
export interface DataGenConfig {
  num_points: number;
  x_min: number;
  x_max: number;
  true_slope: number;
  noise_level: number;
  seed: number;
}

export interface TrainingConfig {
  w_init: number;
  lr: number;
  steps: number;
}

export interface CaseConfig {
  id: string;
  name: string;
  description: string;
  emoji: string;
  category: 'foundational' | 'learning-rate' | 'initialization';
  data_config: DataGenConfig;
  training_config: TrainingConfig;
  insights: string[];
}

export interface CaseManifest {
  version: string;
  cases: CaseConfig[];
}
