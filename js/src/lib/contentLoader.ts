import type { Glossary, FAQData, TutorialContent } from '../types';

/**
 * Loads the glossary from external JSON file
 * @param phase - Optional phase identifier (phase1 or phase2)
 */
export async function loadGlossary(phase?: 'phase1' | 'phase2'): Promise<Glossary> {
  try {
    const suffix = phase === 'phase2' ? '-phase2' : '';
    const response = await fetch(`/content/glossary${suffix}.json`);
    if (!response.ok) {
      throw new Error(`Failed to load glossary: ${response.statusText}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Error loading glossary:', error);
    return {};
  }
}

/**
 * Loads FAQs from external JSON file
 * @param phase - Optional phase identifier (phase1 or phase2)
 */
export async function loadFAQs(phase?: 'phase1' | 'phase2'): Promise<FAQData> {
  try {
    const suffix = phase === 'phase2' ? '-phase2' : '';
    const response = await fetch(`/content/faqs${suffix}.json`);
    if (!response.ok) {
      throw new Error(`Failed to load FAQs: ${response.statusText}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Error loading FAQs:', error);
    return { categories: [] };
  }
}

/**
 * Loads the tutorial from external JSON file
 * @param phase - Optional phase identifier (phase1 or phase2)
 */
export async function loadTutorial(phase?: 'phase1' | 'phase2'): Promise<TutorialContent> {
  try {
    const suffix = phase === 'phase2' ? '-phase2' : '';
    const response = await fetch(`/content/tutorial${suffix}.json`);
    if (!response.ok) {
      throw new Error(`Failed to load tutorial: ${response.statusText}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Error loading tutorial:', error);
    return {
      meta: {
        title: '',
        description: '',
        estimatedReadTime: 0,
        version: '1.0'
      },
      chapters: []
    };
  }
}

/**
 * Loads all educational content at once
 * @param phase - Optional phase identifier (phase1 or phase2)
 */
export async function loadAllContent(phase?: 'phase1' | 'phase2') {
  const [glossary, faqs, tutorial] = await Promise.all([
    loadGlossary(phase),
    loadFAQs(phase),
    loadTutorial(phase)
  ]);

  return {
    glossary,
    faqs,
    tutorial
  };
}
